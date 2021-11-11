package csv

import (
	"bytes"
	"fmt"
	"testing"
)

type Something struct {
	Name string `csv:"name"`
	ID   int    `csv:"id"`
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
)

func Test_Empty(t *testing.T) {
	p := []*Something{{"a", 1}}
	body := bytes.Buffer{}

	if err := NewEncoder(&body).Encode(p); err != nil {
		t.Error(err)
	}
}

func Test_E2EPtrPtr(t *testing.T) {
	body := bytes.Buffer{}
	err := NewEncoder(&body).Encode(&payloadPtr)

	if err != ErrNotAnSlice {
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

	if err != ErrNotAnSlice {
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
