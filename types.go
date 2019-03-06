package sszgen

import (
	"fmt"
	"github.com/holiman/saszy/ssz"
	"math/big"
)

type Alpha struct {
	Auint32 uint32   `ssz:"uint32"`
	Abytes  []byte   `ssz:"bytes"`
	Abigint *big.Int `ssz:"uint256"`
	Alist   []uint32 `ssz:"list:uint32"`
}
type Beta struct {
	Auint32    uint32 `ssz:"uint32"`
	AContainer Alpha  `ssz:"container"`
}
type BetaPointer struct {
	Auint32           uint32 `ssz:"uint32"`
	AContainerPointer *Alpha `ssz:"container"`
}

func (f *Alpha) SszSize() uint32 {
	return uint32(4 + // uint32
		4 + len(f.Abytes) + // byteslen + bytesdata
		32 + // uint256
		4 + len(f.Alist)*4) // length + list data (uint32)
}

func (f *Alpha) DecodeSSZ(buf []byte) error {
	var (
		err    error
		offset uint32
		min    = uint32(4 + 4 + 0 + 4)
	)
	if uint32(len(buf)) < min {
		return fmt.Errorf("decode err: byte slice too small for type %s, has %d require %d", "Foobaz", len(buf), min)
	}

	f.Auint32, offset = ssz.DecodeUint32(buf, offset)
	f.Abytes, offset, err = ssz.DecodeBytesX(buf, offset)
	if err != nil {
		return err
	}
	f.Abigint, offset, err = ssz.DecodeUintN(buf, offset, 256, nil)
	if err != nil {
		return err
	}
	f.Alist, offset, err = ssz.DecodeListUint32(buf, offset, 4)
	return err
}

func (f *Alpha) EncodeSSZ(buf []byte) (uint32, error) {
	var (
		err    error
		offset uint32
	)
	var MINLEN = 4 + // uint32
		4 + // bytearray length
		len(f.Abytes) // bytearray data

	if len(buf) < MINLEN {
		return offset, fmt.Errorf("encode err: byte slice too small for type %s, has %d require %d", "Foobaz", len(buf), MINLEN)
	}
	offset, err = ssz.EncodeUint32(buf, offset, f.Auint32)
	offset, err = ssz.EncodeBytesWithLengthPrefix(buf, offset, f.Abytes)
	offset, err = ssz.EncodeUintN(buf, offset, 256, f.Abigint)
	//offset, err = EncodeListUint32s(buf, offset, f.Listy)
	return offset, err
}

func (obj *Alpha) EnxcodeSSZ(buf []byte) (offset uint32 , err error) {
	offset, err = ssz.EncodeUint32(buf, offset, obj.Auint32)
	if err != nil {
		return
	}
	offset, err = ssz.EncodeBytesWithLengthPrefix(buf, offset, obj.Abytes)
	if err != nil {
		return
	}
	offset, err = ssz.EncodeUintN(buf, offset, 256, obj.Abigint)
	if err != nil {
		return
	}
	return
}

func (b *Beta) SszSize() uint32 {
	// Size
	// 4      : Pre [uint32]
	// + 4    : container size [uint32]
	// + len(ssz_encoded_container)
	return uint32(4 + 4 + ssz.SszSize(&(b.AContainer)))
}

func (obj *Beta) EncodeSSZ(buf []byte) (offset uint32, err error) {
	offset, err = ssz.EncodeUint32(buf, offset, obj.Auint32)

	if offset, err = ssz.SszEncode(buf[offset:], &obj.AContainer); err != nil {
		return
	}
	return
}

func (b *Beta) DecodeSSZ(buf []byte) error {
	var (
		err    error
		offset uint32
		min    = uint32(4 + 4)
	)
	if uint32(len(buf)) < min {
		return fmt.Errorf("decode err: byte slice too small for type %s, has %d require %d", "Foobaz", len(buf), min)
	}
	b.Auint32, offset = ssz.DecodeUint32(buf, offset)
	containerByteSize, offset := ssz.DecodeUint32(buf, offset)
	min = min + containerByteSize
	if uint32(len(buf)) < min {
		return fmt.Errorf("decode err: byte slice too small for type %s, has %d require %d", "Foobaz", len(buf), min)
	}
	err = b.AContainer.DecodeSSZ(buf[offset:])
	if err != nil {
		return fmt.Errorf("decode err: %v", err)
	}
	return nil
}

func (b *BetaPointer) SszSize() uint32 {
	// Size
	// 4      : Pre [uint32]
	// + 4    : container size [uint32]
	// + len(ssz_encoded_container)
	return uint32(4 + 4 + ssz.SszSize(b.AContainerPointer))
}

func (b *BetaPointer) EncodeSSZ(buf []byte) (uint32, error) {
	var err error
	var offset uint32

	offset, err = ssz.EncodeUint32(buf, 0, b.Auint32)
	offset, err = ssz.SszEncode(buf[offset:], b.AContainerPointer)
	return offset, err
}

func (b *BetaPointer) DecodeSSZ(buf []byte) error {
	var (
		err    error
		offset uint32
		min    = uint32(4 + 4)
	)
	if uint32(len(buf)) < min {
		return fmt.Errorf("decode err: byte slice too small for type %s, has %d require %d", "Foobaz", len(buf), min)
	}
	b.Auint32, offset = ssz.DecodeUint32(buf, offset)
	containerByteSize, offset := ssz.DecodeUint32(buf, offset)
	min = min + containerByteSize
	b.AContainerPointer = new(Alpha)
	if err = b.AContainerPointer.DecodeSSZ(buf[offset:]); err != nil {
		return fmt.Errorf("decode err: %v", err)
	}
	return err
}
