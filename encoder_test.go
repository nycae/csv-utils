package csv

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

type Something struct {
	Name string `csv:"name"`
	ID   int    `csv:"id"`
}

type Other struct {
	Name string
	ID   int
}

type Payload []Something
type PayloadPtr []*Something

var (
	payload = Payload{
		{"Harry Sheldon", 1},
		{"Arkady", 2},
		{"Charly Johns", 3},
		{"a name, with, 'comas\" and quotes", 4},
	}
	payloadPtr = PayloadPtr{
		{"Harry Sheldon", 1},
		{"Arkady", 2},
		{"Charly Johns", 3},
		{"a name, with, 'comas\" and quotes", 4},
	}
	payloadUntagged = []Other{
		{"Harry Sheldon", 1},
		{"Arkady", 2},
		{"Charly Johns", 3},
		{"a name, with, 'comas\" and quotes", 4},
	}
)

var (
	ErrPlanned  = errors.New("planned error")
	ErrExpected = errors.New("error was expected but not returned")
)

type writer struct{}

func (writer) Write([]byte) (int, error) {
	return 42, ErrPlanned
}

func Test_WriterErr(t *testing.T) {

}

func Test_TypeErr(t *testing.T) {
	b := bytes.Buffer{}
	err := NewEncoder(&b).Encode([]func(){func() {}})
	if err == nil {
		t.Error(ErrExpected)
	}

	if err != ErrNotAnStructSlice {
		t.Error(err)
	}

}

func Test_Untagged(t *testing.T) {
	body := bytes.Buffer{}
	expected := "Name,ID\nHarry Sheldon,1\nArkady,2\nCharly Johns,3\n\"a name, with, 'comas\"\" and quotes\",4\n"
	if err := NewEncoder(&body).Encode(payloadUntagged); err != nil {
		t.Error(err)
	}

	if expected != body.String() {
		t.Error(fmt.Sprintf("Expected: %s\nGot: %s", expected, body.String()))
	}
}

func Test_Empty(t *testing.T) {
	var p []Something
	var pPtr []*Something
	var body bytes.Buffer

	if err := NewEncoder(&body).Encode(p); err != nil {
		t.Error(err)
	}

	if err := NewEncoder(&body).Encode(pPtr); err != nil {
		t.Error(err)
	}
}

func Test_E2EPtrPtr(t *testing.T) {
	body := bytes.Buffer{}
	err := NewEncoder(&body).Encode(&payloadPtr)

	if err != ErrNotAnStructSlice {
		t.Error(err)
	}
}

func Test_E2ESlicePtr(t *testing.T) {
	body := bytes.Buffer{}
	expected := "name,id\nHarry Sheldon,1\nArkady,2\nCharly Johns,3\n\"a name, with, 'comas\"\" and quotes\",4\n"
	if err := NewEncoder(&body).Encode(payloadPtr); err != nil {
		t.Error(err)
	}

	if expected != body.String() {
		t.Error(fmt.Sprintf("Expected: %s\nGot: %s", expected, body.String()))
	}
}

func Test_E2EPtrStruct(t *testing.T) {
	body := bytes.Buffer{}
	err := NewEncoder(&body).Encode(&payload)

	if err != ErrNotAnStructSlice {
		t.Error(err)
	}
}

func Test_E2ESliceStruct(t *testing.T) {
	body := bytes.Buffer{}
	expected := "name,id\nHarry Sheldon,1\nArkady,2\nCharly Johns,3\n\"a name, with, 'comas\"\" and quotes\",4\n"
	if err := NewEncoder(&body).Encode(payload); err != nil {
		t.Error(err)
	}

	if expected != body.String() {
		t.Error(fmt.Sprintf("Expected: %s\nGot: %s", expected, body.String()))
	}
}
