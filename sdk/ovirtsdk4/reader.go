package ovirtsdk4

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

// XMLTagNotMatchError indicates the error of XML tag
// not matched when unmarshaling XML
type XMLTagNotMatchError struct {
	ActualTag   string
	ExpectedTag string
}

func (err XMLTagNotMatchError) Error() string {
	return fmt.Sprintf("Tag not matched: expect <%v> but got <%v>", err.ExpectedTag, err.ActualTag)
}

// CanForward indicates if Decoder has been finished
func CanForward(tok xml.Token) (bool, error) {
	switch tok.(type) {
	case xml.StartElement:
		return true, nil
	case xml.EndElement:
		return false, nil
	default:
		return true, nil
	}
}

// XMLReader unmarshalizes the xml to struct
type XMLReader struct {
	dec *xml.Decoder
}

// NewXMLReader creates a XMLReader instance
func NewXMLReader(b []byte) *XMLReader {
	return &XMLReader{
		dec: xml.NewDecoder(bytes.NewReader(b)),
	}
}

// Skip calls xml.Decoder.Skip to skip current XML element
func (reader *XMLReader) Skip() error {
	return reader.dec.Skip()
}

// FindStartElement finds the right next StartElement
func (reader *XMLReader) FindStartElement() (*xml.StartElement, error) {
	// Find start element if we need it.
	for {
		tok, err := reader.Next()
		if err != nil {
			fmt.Printf("err is %v\n", err)
			break
		}
		tok = xml.CopyToken(tok)
		if tok, ok := tok.(xml.StartElement); ok {
			return &tok, nil
		}
	}
	return nil, errors.New("Failed to find StartElement")
}

// Next calls xml.Decoder.Token() to get the next xml.Token
func (reader *XMLReader) Next() (xml.Token, error) {
	return reader.dec.Token()
}

// ReadString reads the xml.CharData as a string after xml.StartElement
func (reader *XMLReader) ReadString(start *xml.StartElement) (string, error) {
	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			return "", err
		}
		start = st
	}
	var buf []byte
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		switch t := t.(type) {
		case xml.CharData:
			if depth == 1 {
				buf = append(buf, t...)
			}
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}
	}

	return string(buf), nil
}

// ReadStrings reads the xml.CharData of all subelements with a slice of string returned
func (reader *XMLReader) ReadStrings(start *xml.StartElement) ([]string, error) {
	var strings []string

	if start == nil {
		st, err := reader.FindStartElement()
		if err != nil {
			return nil, err
		}
		start = st
	}
	depth := 1
	for depth > 0 {
		t, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch t := t.(type) {
		case xml.StartElement:
			str, err := reader.ReadString(&t)
			if err != nil {
				return nil, err
			}
			strings = append(strings, str)
			depth++
		case xml.EndElement:
			depth--
		}
	}

	return strings, nil
}

// ReadBool reads the xml.CharData as bool
func (reader *XMLReader) ReadBool(start *xml.StartElement) (bool, error) {
	str, err := reader.ReadString(start)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(str)
}

func (reader *XMLReader) ReadBools(start *xml.StartElement) ([]bool, error) {
	strs, err := reader.ReadStrings(start)
	if err != nil {
		return nil, err
	}
	var bools []bool
	for _, sv := range strs {
		bv, err := strconv.ParseBool(sv)
		if err != nil {
			return nil, err
		}
		bools = append(bools, bv)
	}
	return bools, nil
}

// ReadInt64 reads the xml.CharData as int64
func (reader *XMLReader) ReadInt64(start *xml.StartElement) (int64, error) {
	str, err := reader.ReadString(start)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(str, 10, 64)
}

func (reader *XMLReader) ReadInt64s(start *xml.StartElement) ([]int64, error) {
	strs, err := reader.ReadStrings(start)
	if err != nil {
		return nil, err
	}
	var int64s []int64
	for _, sv := range strs {
		iv, err := strconv.ParseInt(sv, 10, 64)
		if err != nil {
			return nil, err
		}
		int64s = append(int64s, iv)
	}
	return int64s, nil
}

func (reader *XMLReader) ReadFloat64(start *xml.StartElement) (float64, error) {
	str, err := reader.ReadString(start)
	if err != nil {
		return 0.0, err
	}
	return strconv.ParseFloat(str, 64)
}

func (reader *XMLReader) ReadFloat64s(start *xml.StartElement) ([]float64, error) {
	strs, err := reader.ReadStrings(start)
	if err != nil {
		return nil, err
	}
	var float64s []float64
	for _, sv := range strs {
		fv, err := strconv.ParseFloat(sv, 64)
		if err != nil {
			return nil, err
		}
		float64s = append(float64s, fv)
	}
	return float64s, nil
}

// ReadTime reads the xml.CharData as time.Time
func (reader *XMLReader) ReadTime(start *xml.StartElement) (time.Time, error) {
	str, err := reader.ReadString(start)
	if err != nil {
		var t time.Time
		return t, err
	}
	return time.Parse("2006-01-02T15:04:05.999999-07:00", str)
}

func (reader *XMLReader) ReadTimes(start *xml.StartElement) ([]time.Time, error) {
	strs, err := reader.ReadStrings(start)
	if err != nil {
		return nil, err
	}
	var times []time.Time
	for _, sv := range strs {
		tv, err := time.Parse("2006-01-02T15:04:05.999999-07:00", sv)
		if err != nil {
			return nil, err
		}
		times = append(times, tv)
	}
	return times, nil
}
