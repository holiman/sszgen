package sszgen

import (
	"fmt"
	"testing"
)

//func TestContainer(t *testing.T) {
//
//	u := Baxoz{
//		Pre: 0xffff,
//		Containy: Foobaz{0x1337,
//			[]byte{1, 2, 3},
//			big.NewInt(1337),
//			[]uint32{6, 7, 8, 9},
//		},
//	}
//	buf := make([]byte, u.SszSize())
//	u.EncodeSSZ(buf)
//	fmt.Printf("encoded %x\n", buf)
//
//	var v Baxoz
//	(&v).DecodeSSZ(buf)
//	fmt.Printf("decoded %v\n", v)
//
//	buf2 := make([]byte, v.SszSize())
//	v.EncodeSSZ(buf2)
//	fmt.Printf("encoded %x\n", buf2)
//
//}

//func TestContainerWithPointer(t *testing.T) {
//
//	u := BaxozPtr{
//		Pre: 0xffff,
//		Containy: &Foobaz{0x1337,
//			[]byte{1, 2, 3},
//			big.NewInt(1337),
//			[]uint32{6, 7, 8, 9},
//		},
//	}
//	buf := make([]byte, SszSize(&u))
//	u.EncodeSSZ(buf)
//	fmt.Printf("encoded %x\n", buf)
//
//	var v BaxozPtr
//	err := (&v).DecodeSSZ(buf[:len(buf)])
//	fmt.Printf("decoded %v\n", v)
//	if err != nil {
//		fmt.Printf("err: %v\n", err)
//	}
//
//	buf2 := make([]byte, SszSize(&v))
//	v.EncodeSSZ(buf2)
//	fmt.Printf("encoded %x\n", buf2)
//}

func TestCodeGen(t *testing.T) {
	var code string
	code, _ = genCode(Alpha{})
	fmt.Printf("code\n%v", code)

	//var exp string
	//	code, _ = genCode(Beta{})
	//	exp = `func (obj *Baxoz) EncodeSSZ(buf []byte) (offset uint32 , err error) {
	//	offset = EncodeUint32(buf, offset, obj.Pre)
	//	offset, err = SszEncode(buf[offset:], &obj.Containy)
	//	if err != nil{
	//		return offset, err
	//	}
	//}`
	//	if code != exp {
	//		fmt.Printf("code\n%v", code)
	//		t.Logf("got \n%v\n exp \n%v\n", code, exp)
	//		t.Errorf("Wrong code")
	//	}
	//	code, _ = genCode(BetaPointer{})
	//	exp = `func (obj *BaxozPtr) EncodeSSZ(buf []byte) (offset uint32 , err error) {
	//	offset = EncodeUint32(buf, offset, obj.Pre)
	//	offset, err = SszEncode(buf[offset:], obj.Containy)
	//	if err != nil{
	//		return offset, err
	//	}
	//}`
	//	if code != exp {
	//		fmt.Println(code)
	//		t.Logf("got \n%v\n exp \n%v\n", code, exp)
	//		t.Errorf("Wrong code")
	//	}
	//fmt.Printf("code\n%v", code)

}
