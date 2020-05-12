// NXP Ultra Secured Digital Host Controller (uSDHC) driver
// https://github.com/f-secure-foundry/tamago
//
// IP: https://www.mobiveil.com/esdhc/
//
// Copyright (c) F-Secure Corporation
// https://foundry.f-secure.com
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.
//
// +build tamago,arm

package usdhc

import (
	"bytes"
	"encoding/binary"
)

const (
	ATTR_VALID = 0
	ATTR_END   = 1
	ATTR_INT   = 2
	ATTR_ACT   = 4

	ACT_TRANSFER = 0b10
	ACT_LINK     = 0b11

	DMASEL_NONE  = 0b00
	DMASEL_ADMA1 = 0b01
	DMASEL_ADMA2 = 0b10

	ADMA_BD_MAX_LENGTH = 65532

	BLOCK_SIZE uint32 = 512
)

// p3964 58.4.2.4.1 ADMA Concept and Descriptor Format, IMX6ULLRM
type ADMABufferDescriptor struct {
	Attribute uint8
	res       uint8
	Length    uint16
	Address   uint32

	next *ADMABufferDescriptor
}

// Init initializes an ADMA2 buffer descriptor.
func (bd *ADMABufferDescriptor) Init(addr uint32, size int) {
	b := bd

	for size > 0 {
		if size <= ADMA_BD_MAX_LENGTH {
			b.Attribute = ACT_TRANSFER<<ATTR_ACT | 1<<ATTR_END | 1<<ATTR_VALID
			b.Length = uint16(size)
			b.Address = addr
			break
		} else {
			b.Attribute = ACT_TRANSFER<<ATTR_ACT | 1<<ATTR_VALID
			b.Length = uint16(ADMA_BD_MAX_LENGTH)
			b.Address = addr

			addr += ADMA_BD_MAX_LENGTH
			size -= ADMA_BD_MAX_LENGTH

			b.next = &ADMABufferDescriptor{}
			b = b.next
		}
	}
}

// Bytes converts the descriptor structure to byte array format.
func (bd *ADMABufferDescriptor) Bytes() []byte {
	buf := new(bytes.Buffer)

	b := bd

	for b != nil {
		binary.Write(buf, binary.LittleEndian, b.Attribute)
		binary.Write(buf, binary.LittleEndian, b.res)
		binary.Write(buf, binary.LittleEndian, b.Length)
		binary.Write(buf, binary.LittleEndian, b.Address)

		b = b.next
	}

	return buf.Bytes()
}
