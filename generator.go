package sszgen

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

var (
	tmplSimple              *template.Template
	tmplUintN               *template.Template
	tmplContainer           *template.Template
	tmplMethod              *template.Template
	tmplContainerNonPointer *template.Template
	tmplList                *template.Template
)

func init() {
	tmplSimple, _ = template.New("a").Parse(`	offset, err = {{ .Encoder }}(buf, offset, obj.{{.FieldName}})
	if err != nil{
		return
	}`)
	tmplUintN, _ = template.New("b").Parse(`	offset, err = {{ .Encoder }}(buf, offset, 256, obj.{{.FieldName}})
	if err != nil{
		return
	}`)

	tmplContainer, _ = template.New("b").Parse(`	offset, err = {{ .Encoder }}(buf[offset:], obj.{{.FieldName}})
	if err != nil{
		return offset, err
	}`)
	tmplContainerNonPointer, _ = template.New("b").Parse(`	offset, err = {{ .Encoder }}(buf[offset:], &obj.{{.FieldName}})
	if err != nil{
		return offset, err
	}`)

}

type sszInfo struct {
	fixedSize bool
	size      int
	statement *template.Template
	encoder   string
}

func sszTypeInfo(sszType string) (sszInfo, error) {
	var (
		info sszInfo
		err  error
	)
	// Small integers
	if sszType == "uint8" ||
		sszType == "uint16" ||
		sszType == "uint32" ||
		sszType == "uint64" {
		info.fixedSize = true
		info.size, err = strconv.Atoi(strings.TrimPrefix(sszType, "uint"))
		if err != nil {
			return info, fmt.Errorf("invalid uintn format: %s", sszType)
		}
		info.statement = tmplSimple
		info.encoder = fmt.Sprintf("ssz.EncodeUint%d", info.size)
		return info, err
	}
	// Large integers
	if strings.HasPrefix(sszType, "uint") {
		n, err := strconv.Atoi(strings.TrimPrefix(sszType, "uint"))
		if err != nil {
			return info, err
		}
		if n%8 != 0 {
			return info, fmt.Errorf("invalid 'uintn' format: uint%d not a multiple of '8'", n)
		}
		if n < 64 {
			return info, fmt.Errorf("uintn format uint%d not implemented", n)
		}
		// We're dealing with a bigint
		info.fixedSize = true
		info.size = n / 8 // byte size
		info.statement = tmplUintN
		info.encoder = "ssz.EncodeUintN"
		return info, nil
	}
	// Dynamic size bytes
	if sszType == "bytes" {
		info.fixedSize = false
		info.statement = tmplSimple
		info.encoder = "ssz.EncodeBytesX"
		return info, nil
	}

	// Static size bytes
	if strings.HasPrefix(sszType, "bytes") {
		n, err := strconv.Atoi(strings.TrimPrefix(sszType, "uint"))
		if err != nil {
			return info, err
		}
		info.fixedSize = true
		info.size = n
		info.statement = tmplSimple
		info.encoder = "ssz.EncodeBytesN"
		return info, nil
	}
	return info, fmt.Errorf("no typeinfo found for %s", sszType)

}
func generateStatement(sszType string, field reflect.StructField, obj reflect.Type) (string, error) {
	var (
		//tmpl *template.Template
		err error
	)
	type Meta struct {
		Encoder   string
		FieldName string
	}
	info, err := sszTypeInfo(sszType)
	if err != nil {
		return "", err
	}
	meta := Meta{
		Encoder:   info.encoder,
		FieldName: field.Name,
	}

	buf := new(bytes.Buffer)

	err = info.statement.Execute(buf, meta)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func generateListEncoder(sszType string, field reflect.StructField, obj reflect.Type) (string, error) {
	elemStmt, err := generateStatement(sszType, field, obj)

	type listFiller struct {
		ElemStmt  string
		FieldName string
	}

	tmpl, _ := template.New("c").Parse(`
	startIndex := offset
	offset += 4 // Fill this later
	for _, elem := range obj.{{.FieldName}} {
		{{ .ElemStmt }}
	}
	// Write the actual bytesize
	ssz.EncodeUint32(buf, startIndex, (offset - startIndex))
	return offset, err
	`)
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, listFiller{
		ElemStmt:  elemStmt,
		FieldName: field.Name,
	})

	return buf.String(), err
}


func genCode(obj interface{}) (string, error) {
	var body []string

	fmt.Printf("type %v\n", obj)

	t := reflect.TypeOf(obj)
	fmt.Printf("Type: %v\n", t)
	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("incorrect type used for encoding")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		sszType := field.Tag.Get("ssz")
		fmt.Printf("\tssz type: %q\n", sszType)
		// Might be a list:xxx
		if strings.HasPrefix(sszType, "list:") {
			listType := strings.TrimPrefix(sszType, "list:")
			if txt, err := generateListEncoder(listType, field, t); err != nil {
				fmt.Printf("Err : %v\n", err)
			} else {
				//fmt.Println(txt)
				body = append(body, txt)
			}
		} else {
			if txt, err := generateStatement(sszType, field, t); err != nil {
				fmt.Printf("Err : %v\n", err)
			} else {
				//fmt.Println(txt)
				body = append(body, txt)
			}
		}

	}
	methodBody := strings.Join(body, "\n")
	method := fmt.Sprintf("func (obj *%v) EncodeSSZ(buf []byte) (offset uint32 , err error) {\n%v\n\treturn}\n", t.Name(), methodBody)

	return method, nil
}
