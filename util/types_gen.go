package util

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *AlphabeticTerm) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Title":
			z.Title, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Title")
				return
			}
		case "Code":
			z.Code, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Code")
				return
			}
		case "Manif":
			z.Manif, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Manif")
				return
			}
		case "Attrs":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "Attrs")
				return
			}
			if cap(z.Attrs) >= int(zb0002) {
				z.Attrs = (z.Attrs)[:zb0002]
			} else {
				z.Attrs = make([]Attr, zb0002)
			}
			for za0001 := range z.Attrs {
				var zb0003 uint32
				zb0003, err = dc.ReadArrayHeader()
				if err != nil {
					err = msgp.WrapError(err, "Attrs", za0001)
					return
				}
				if zb0003 != 2 {
					err = msgp.ArrayError{Wanted: 2, Got: zb0003}
					return
				}
				z.Attrs[za0001].Attr, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Attrs", za0001, "Attr")
					return
				}
				z.Attrs[za0001].Value, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Attrs", za0001, "Value")
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *AlphabeticTerm) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "Title"
	err = en.Append(0x84, 0xa5, 0x54, 0x69, 0x74, 0x6c, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Title)
	if err != nil {
		err = msgp.WrapError(err, "Title")
		return
	}
	// write "Code"
	err = en.Append(0xa4, 0x43, 0x6f, 0x64, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Code)
	if err != nil {
		err = msgp.WrapError(err, "Code")
		return
	}
	// write "Manif"
	err = en.Append(0xa5, 0x4d, 0x61, 0x6e, 0x69, 0x66)
	if err != nil {
		return
	}
	err = en.WriteString(z.Manif)
	if err != nil {
		err = msgp.WrapError(err, "Manif")
		return
	}
	// write "Attrs"
	err = en.Append(0xa5, 0x41, 0x74, 0x74, 0x72, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Attrs)))
	if err != nil {
		err = msgp.WrapError(err, "Attrs")
		return
	}
	for za0001 := range z.Attrs {
		// array header, size 2
		err = en.Append(0x92)
		if err != nil {
			return
		}
		err = en.WriteString(z.Attrs[za0001].Attr)
		if err != nil {
			err = msgp.WrapError(err, "Attrs", za0001, "Attr")
			return
		}
		err = en.WriteString(z.Attrs[za0001].Value)
		if err != nil {
			err = msgp.WrapError(err, "Attrs", za0001, "Value")
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *AlphabeticTerm) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Title"
	o = append(o, 0x84, 0xa5, 0x54, 0x69, 0x74, 0x6c, 0x65)
	o = msgp.AppendString(o, z.Title)
	// string "Code"
	o = append(o, 0xa4, 0x43, 0x6f, 0x64, 0x65)
	o = msgp.AppendString(o, z.Code)
	// string "Manif"
	o = append(o, 0xa5, 0x4d, 0x61, 0x6e, 0x69, 0x66)
	o = msgp.AppendString(o, z.Manif)
	// string "Attrs"
	o = append(o, 0xa5, 0x41, 0x74, 0x74, 0x72, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Attrs)))
	for za0001 := range z.Attrs {
		// array header, size 2
		o = append(o, 0x92)
		o = msgp.AppendString(o, z.Attrs[za0001].Attr)
		o = msgp.AppendString(o, z.Attrs[za0001].Value)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *AlphabeticTerm) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Title":
			z.Title, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Title")
				return
			}
		case "Code":
			z.Code, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Code")
				return
			}
		case "Manif":
			z.Manif, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Manif")
				return
			}
		case "Attrs":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Attrs")
				return
			}
			if cap(z.Attrs) >= int(zb0002) {
				z.Attrs = (z.Attrs)[:zb0002]
			} else {
				z.Attrs = make([]Attr, zb0002)
			}
			for za0001 := range z.Attrs {
				var zb0003 uint32
				zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Attrs", za0001)
					return
				}
				if zb0003 != 2 {
					err = msgp.ArrayError{Wanted: 2, Got: zb0003}
					return
				}
				z.Attrs[za0001].Attr, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Attrs", za0001, "Attr")
					return
				}
				z.Attrs[za0001].Value, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Attrs", za0001, "Value")
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *AlphabeticTerm) Msgsize() (s int) {
	s = 1 + 6 + msgp.StringPrefixSize + len(z.Title) + 5 + msgp.StringPrefixSize + len(z.Code) + 6 + msgp.StringPrefixSize + len(z.Manif) + 6 + msgp.ArrayHeaderSize
	for za0001 := range z.Attrs {
		s += 1 + msgp.StringPrefixSize + len(z.Attrs[za0001].Attr) + msgp.StringPrefixSize + len(z.Attrs[za0001].Value)
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Attr) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zb0001}
		return
	}
	z.Attr, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Attr")
		return
	}
	z.Value, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Value")
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Attr) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 2
	err = en.Append(0x92)
	if err != nil {
		return
	}
	err = en.WriteString(z.Attr)
	if err != nil {
		err = msgp.WrapError(err, "Attr")
		return
	}
	err = en.WriteString(z.Value)
	if err != nil {
		err = msgp.WrapError(err, "Value")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Attr) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	o = msgp.AppendString(o, z.Attr)
	o = msgp.AppendString(o, z.Value)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Attr) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zb0001}
		return
	}
	z.Attr, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Attr")
		return
	}
	z.Value, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Value")
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Attr) Msgsize() (s int) {
	s = 1 + msgp.StringPrefixSize + len(z.Attr) + msgp.StringPrefixSize + len(z.Value)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Diag) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 3 {
		err = msgp.ArrayError{Wanted: 3, Got: zb0001}
		return
	}
	z.Code, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Code")
		return
	}
	z.Desc, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Desc")
		return
	}
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err, "Notes")
		return
	}
	if cap(z.Notes) >= int(zb0002) {
		z.Notes = (z.Notes)[:zb0002]
	} else {
		z.Notes = make([]Note, zb0002)
	}
	for za0001 := range z.Notes {
		var zb0003 uint32
		zb0003, err = dc.ReadArrayHeader()
		if err != nil {
			err = msgp.WrapError(err, "Notes", za0001)
			return
		}
		if zb0003 != 2 {
			err = msgp.ArrayError{Wanted: 2, Got: zb0003}
			return
		}
		z.Notes[za0001].Kind, err = dc.ReadString()
		if err != nil {
			err = msgp.WrapError(err, "Notes", za0001, "Kind")
			return
		}
		var zb0004 uint32
		zb0004, err = dc.ReadArrayHeader()
		if err != nil {
			err = msgp.WrapError(err, "Notes", za0001, "Notes")
			return
		}
		if cap(z.Notes[za0001].Notes) >= int(zb0004) {
			z.Notes[za0001].Notes = (z.Notes[za0001].Notes)[:zb0004]
		} else {
			z.Notes[za0001].Notes = make([]string, zb0004)
		}
		for za0002 := range z.Notes[za0001].Notes {
			z.Notes[za0001].Notes[za0002], err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Notes", za0001, "Notes", za0002)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Diag) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 3
	err = en.Append(0x93)
	if err != nil {
		return
	}
	err = en.WriteString(z.Code)
	if err != nil {
		err = msgp.WrapError(err, "Code")
		return
	}
	err = en.WriteString(z.Desc)
	if err != nil {
		err = msgp.WrapError(err, "Desc")
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Notes)))
	if err != nil {
		err = msgp.WrapError(err, "Notes")
		return
	}
	for za0001 := range z.Notes {
		// array header, size 2
		err = en.Append(0x92)
		if err != nil {
			return
		}
		err = en.WriteString(z.Notes[za0001].Kind)
		if err != nil {
			err = msgp.WrapError(err, "Notes", za0001, "Kind")
			return
		}
		err = en.WriteArrayHeader(uint32(len(z.Notes[za0001].Notes)))
		if err != nil {
			err = msgp.WrapError(err, "Notes", za0001, "Notes")
			return
		}
		for za0002 := range z.Notes[za0001].Notes {
			err = en.WriteString(z.Notes[za0001].Notes[za0002])
			if err != nil {
				err = msgp.WrapError(err, "Notes", za0001, "Notes", za0002)
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Diag) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 3
	o = append(o, 0x93)
	o = msgp.AppendString(o, z.Code)
	o = msgp.AppendString(o, z.Desc)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Notes)))
	for za0001 := range z.Notes {
		// array header, size 2
		o = append(o, 0x92)
		o = msgp.AppendString(o, z.Notes[za0001].Kind)
		o = msgp.AppendArrayHeader(o, uint32(len(z.Notes[za0001].Notes)))
		for za0002 := range z.Notes[za0001].Notes {
			o = msgp.AppendString(o, z.Notes[za0001].Notes[za0002])
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Diag) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 3 {
		err = msgp.ArrayError{Wanted: 3, Got: zb0001}
		return
	}
	z.Code, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Code")
		return
	}
	z.Desc, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Desc")
		return
	}
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Notes")
		return
	}
	if cap(z.Notes) >= int(zb0002) {
		z.Notes = (z.Notes)[:zb0002]
	} else {
		z.Notes = make([]Note, zb0002)
	}
	for za0001 := range z.Notes {
		var zb0003 uint32
		zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Notes", za0001)
			return
		}
		if zb0003 != 2 {
			err = msgp.ArrayError{Wanted: 2, Got: zb0003}
			return
		}
		z.Notes[za0001].Kind, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Notes", za0001, "Kind")
			return
		}
		var zb0004 uint32
		zb0004, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Notes", za0001, "Notes")
			return
		}
		if cap(z.Notes[za0001].Notes) >= int(zb0004) {
			z.Notes[za0001].Notes = (z.Notes[za0001].Notes)[:zb0004]
		} else {
			z.Notes[za0001].Notes = make([]string, zb0004)
		}
		for za0002 := range z.Notes[za0001].Notes {
			z.Notes[za0001].Notes[za0002], bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Notes", za0001, "Notes", za0002)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Diag) Msgsize() (s int) {
	s = 1 + msgp.StringPrefixSize + len(z.Code) + msgp.StringPrefixSize + len(z.Desc) + msgp.ArrayHeaderSize
	for za0001 := range z.Notes {
		s += 1 + msgp.StringPrefixSize + len(z.Notes[za0001].Kind) + msgp.ArrayHeaderSize
		for za0002 := range z.Notes[za0001].Notes {
			s += msgp.StringPrefixSize + len(z.Notes[za0001].Notes[za0002])
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *DocIDMap) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0003 uint32
	zb0003, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if (*z) == nil {
		(*z) = make(DocIDMap, zb0003)
	} else if len((*z)) > 0 {
		for key := range *z {
			delete((*z), key)
		}
	}
	for zb0003 > 0 {
		zb0003--
		var zb0001 string
		var zb0002 bool
		zb0001, err = dc.ReadString()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		zb0002, err = dc.ReadBool()
		if err != nil {
			err = msgp.WrapError(err, zb0001)
			return
		}
		(*z)[zb0001] = zb0002
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z DocIDMap) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(uint32(len(z)))
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0004, zb0005 := range z {
		err = en.WriteString(zb0004)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		err = en.WriteBool(zb0005)
		if err != nil {
			err = msgp.WrapError(err, zb0004)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z DocIDMap) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, uint32(len(z)))
	for zb0004, zb0005 := range z {
		o = msgp.AppendString(o, zb0004)
		o = msgp.AppendBool(o, zb0005)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DocIDMap) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0003 uint32
	zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if (*z) == nil {
		(*z) = make(DocIDMap, zb0003)
	} else if len((*z)) > 0 {
		for key := range *z {
			delete((*z), key)
		}
	}
	for zb0003 > 0 {
		var zb0001 string
		var zb0002 bool
		zb0003--
		zb0001, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		zb0002, bts, err = msgp.ReadBoolBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, zb0001)
			return
		}
		(*z)[zb0001] = zb0002
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z DocIDMap) Msgsize() (s int) {
	s = msgp.MapHeaderSize
	if z != nil {
		for zb0004, zb0005 := range z {
			_ = zb0005
			s += msgp.StringPrefixSize + len(zb0004) + msgp.BoolSize
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *DrugTerm) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 4 {
		err = msgp.ArrayError{Wanted: 4, Got: zb0001}
		return
	}
	z.Title, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Title")
		return
	}
	z.See, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "See")
		return
	}
	z.SeeAlso, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "SeeAlso")
		return
	}
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err, "Codes")
		return
	}
	if cap(z.Codes) >= int(zb0002) {
		z.Codes = (z.Codes)[:zb0002]
	} else {
		z.Codes = make([]string, zb0002)
	}
	for za0001 := range z.Codes {
		z.Codes[za0001], err = dc.ReadString()
		if err != nil {
			err = msgp.WrapError(err, "Codes", za0001)
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *DrugTerm) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 4
	err = en.Append(0x94)
	if err != nil {
		return
	}
	err = en.WriteString(z.Title)
	if err != nil {
		err = msgp.WrapError(err, "Title")
		return
	}
	err = en.WriteString(z.See)
	if err != nil {
		err = msgp.WrapError(err, "See")
		return
	}
	err = en.WriteString(z.SeeAlso)
	if err != nil {
		err = msgp.WrapError(err, "SeeAlso")
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Codes)))
	if err != nil {
		err = msgp.WrapError(err, "Codes")
		return
	}
	for za0001 := range z.Codes {
		err = en.WriteString(z.Codes[za0001])
		if err != nil {
			err = msgp.WrapError(err, "Codes", za0001)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *DrugTerm) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 4
	o = append(o, 0x94)
	o = msgp.AppendString(o, z.Title)
	o = msgp.AppendString(o, z.See)
	o = msgp.AppendString(o, z.SeeAlso)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Codes)))
	for za0001 := range z.Codes {
		o = msgp.AppendString(o, z.Codes[za0001])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DrugTerm) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 4 {
		err = msgp.ArrayError{Wanted: 4, Got: zb0001}
		return
	}
	z.Title, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Title")
		return
	}
	z.See, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "See")
		return
	}
	z.SeeAlso, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "SeeAlso")
		return
	}
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Codes")
		return
	}
	if cap(z.Codes) >= int(zb0002) {
		z.Codes = (z.Codes)[:zb0002]
	} else {
		z.Codes = make([]string, zb0002)
	}
	for za0001 := range z.Codes {
		z.Codes[za0001], bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Codes", za0001)
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *DrugTerm) Msgsize() (s int) {
	s = 1 + msgp.StringPrefixSize + len(z.Title) + msgp.StringPrefixSize + len(z.See) + msgp.StringPrefixSize + len(z.SeeAlso) + msgp.ArrayHeaderSize
	for za0001 := range z.Codes {
		s += msgp.StringPrefixSize + len(z.Codes[za0001])
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Note) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zb0001}
		return
	}
	z.Kind, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Kind")
		return
	}
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err, "Notes")
		return
	}
	if cap(z.Notes) >= int(zb0002) {
		z.Notes = (z.Notes)[:zb0002]
	} else {
		z.Notes = make([]string, zb0002)
	}
	for za0001 := range z.Notes {
		z.Notes[za0001], err = dc.ReadString()
		if err != nil {
			err = msgp.WrapError(err, "Notes", za0001)
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Note) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 2
	err = en.Append(0x92)
	if err != nil {
		return
	}
	err = en.WriteString(z.Kind)
	if err != nil {
		err = msgp.WrapError(err, "Kind")
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Notes)))
	if err != nil {
		err = msgp.WrapError(err, "Notes")
		return
	}
	for za0001 := range z.Notes {
		err = en.WriteString(z.Notes[za0001])
		if err != nil {
			err = msgp.WrapError(err, "Notes", za0001)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Note) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	o = msgp.AppendString(o, z.Kind)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Notes)))
	for za0001 := range z.Notes {
		o = msgp.AppendString(o, z.Notes[za0001])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Note) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zb0001}
		return
	}
	z.Kind, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Kind")
		return
	}
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Notes")
		return
	}
	if cap(z.Notes) >= int(zb0002) {
		z.Notes = (z.Notes)[:zb0002]
	} else {
		z.Notes = make([]string, zb0002)
	}
	for za0001 := range z.Notes {
		z.Notes[za0001], bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Notes", za0001)
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Note) Msgsize() (s int) {
	s = 1 + msgp.StringPrefixSize + len(z.Kind) + msgp.ArrayHeaderSize
	for za0001 := range z.Notes {
		s += msgp.StringPrefixSize + len(z.Notes[za0001])
	}
	return
}
