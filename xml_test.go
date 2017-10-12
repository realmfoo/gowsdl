package gowsdl

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"testing"
)

type ValueList []int32

func (v *ValueList) UnmarshalText(text []byte) error {
	if text == nil {
		return nil
	}
	var s string
	i := 0
	for i != -1 {
		i = bytes.IndexByte(text, ' ')
		if i == -1 {
			s = string(text)
		} else {
			s = string(text[:i])
			text = text[i+1:]
		}
		val, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*v = append(*v, int32(val))
	}

	return nil
}

func (v ValueList) MarshalText() (text []byte, err error) {
	buf := new(bytes.Buffer)
	i := 0
	n := len(v)
	for i < n {
		val := []byte(strconv.Itoa((int)(v[i])))
		buf.Write(val)
		i++
		if i == n {
			break
		}
		buf.Write([]byte{' '})
	}
	return buf.Bytes(), nil
}

type ValueListElm struct {
	XMLName xml.Name  `xml:"elm"`
	Attr    ValueList `xml:"attr,attr"`
}

func TestValueListMarshalling(t *testing.T) {
	elm := ValueListElm{Attr: ValueList{10, 12}}

	data, err := xml.Marshal(elm)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<elm attr="10 12"></elm>`
	if string(data) != expected {
		t.Error("got `" + string(data) + "` wanted `" + expected + "`")
	}
}

func TestValueListUnmarshalling(t *testing.T) {
	elm := ValueListElm{}

	err := xml.Unmarshal([]byte(`<elm attr="10 12"></elm>`), &elm)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{{ elm} [10 12]}`
	if fmt.Sprintf("%v", elm) != expected {
		t.Error("got `" + fmt.Sprintf("%v", elm) + "` wanted `" + expected + "`")
	}
}
