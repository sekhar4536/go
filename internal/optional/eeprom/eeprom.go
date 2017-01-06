// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package eeprom

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	OnieDataSz = 8
	LenOffset  = OnieDataSz + 1
	HeaderSz   = LenOffset + 2
)

type Eeprom struct {
	Onie struct {
		Data    *OnieData
		Version *Hex8
	}
	Tlv TlvMap
}

func (p *Eeprom) Bytes() []byte {
	buf := new(bytes.Buffer)
	buf.Write(p.Onie.Data[:])
	buf.WriteByte(byte(*p.Onie.Version))
	tlvbytes := p.Tlv.Bytes()
	binary.Write(buf, binary.BigEndian, uint16(len(tlvbytes)))
	buf.Write(tlvbytes)
	return buf.Bytes()
}

func (p *Eeprom) Clone() (*Eeprom, error) {
	clone := new(Eeprom)
	_, err := clone.Write(p.Bytes())
	return clone, err
}

func (p *Eeprom) Del(name string) {
	delete(p.Tlv, TypesByName[name])
}

func (p *Eeprom) Equal(clone *Eeprom) error {
	if p.Onie.Data != clone.Onie.Data {
		return fmt.Errorf("Onie.Data: [% x] vs. [% x]",
			p.Onie.Data, clone.Onie.Data)
	}
	if p.Onie.Version != clone.Onie.Version {
		return fmt.Errorf("Onie.Version: %x vs. %x",
			p.Onie.Version, clone.Onie.Version)
	}
	return p.Tlv.Equal(clone.Tlv)
}

func (p *Eeprom) Set(name, s string) (err error) {
	switch name {
	case "Onie.Data":
		err = p.Onie.Data.Scan(s)
	case "Onie.Version":
		err = p.Onie.Version.Scan(s)
	default:
		t := TypesByName[name]
		v := p.Tlv[t]
		if v == nil {
			v = p.Tlv.Add(t)
		}
		b, isBytesBuffer := v.(*bytes.Buffer)
		if isBytesBuffer {
			b.Reset()
			_, err = b.Write([]byte(s))
		} else {
			err = v.(Scanner).Scan(s)
		}
	}
	return
}

func (p *Eeprom) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, "eeprom.Onie.Data:", p.Onie.Data)
	fmt.Fprintln(buf, "eeprom.Onie.Version:", p.Onie.Version)
	buf.WriteString(p.Tlv.String())
	return buf.String()
}

func (p *Eeprom) Write(buf []byte) (n int, err error) {
	if p.Onie.Data == nil {
		p.Onie.Data = new(OnieData)
	}
	if p.Onie.Version == nil {
		p.Onie.Version = new(Hex8)
	}
	if p.Tlv == nil {
		p.Tlv = make(TlvMap)
	}
	i, err := p.Onie.Data.Write(buf)
	if err != nil {
		return
	}
	n += i
	i, err = p.Onie.Version.Write(buf[OnieDataSz:])
	if err != nil {
		return
	}
	n += i
	i, err = p.Tlv.Write(buf[OnieDataSz+1+2:])
	n += i
	return
}
