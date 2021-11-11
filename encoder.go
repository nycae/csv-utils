package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
)

var (
	ErrNotAnStructSlice = errors.New("provided interface must be an slice of structs")
	ErrUnsupportedType  = errors.New("unsupported type, remember struct nesting is not supported")
)

type Encoder struct {
	out *csv.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{csv.NewWriter(w)}
}

func headers(t reflect.Type) ([]string, error) {
	t = safeType(t)
	if t.Kind() != reflect.Struct {
		return nil, ErrNotAnStructSlice
	}

	h := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		name, ok := t.Field(i).Tag.Lookup("csv")
		if !ok {
			name = t.Field(i).Name
		}
		h[i] = name
	}

	return h, nil
}

func values(obj interface{}) ([]string, error) {
	v := safeValue(reflect.ValueOf(obj))
	data := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		switch f := v.Field(i); f.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			data[i] = fmt.Sprintf("%d", reflect.ValueOf(f.Interface()).Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			data[i] = fmt.Sprintf("%d", reflect.ValueOf(f.Interface()).Uint())
		case reflect.String:
			data[i] = fmt.Sprintf("%s", reflect.ValueOf(f.Interface()).String())
		case reflect.Float32, reflect.Float64:
			data[i] = fmt.Sprintf("%f", reflect.ValueOf(f.Interface()).Float())
		case reflect.Bool:
			data[i] = fmt.Sprint(reflect.ValueOf(f.Interface()).Bool())
		default:
			return nil, ErrUnsupportedType
		}
	}

	return data, nil
}

func castToSlice(obj interface{}) []interface{} {
	v := reflect.ValueOf(obj)
	oo := make([]interface{}, v.Len())
	for i := 0; i < len(oo); i++ {
		oo[i] = v.Index(i).Interface()
	}

	return oo
}

func safeType(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Ptr {
		return t
	}

	return t.Elem()
}

func safeValue(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr {
		return v
	}

	return v.Elem()
}

func (e *Encoder) Encode(obj interface{}) error {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	if t.Kind() != reflect.Slice {
		return ErrNotAnStructSlice
	}

	if v.IsNil() || v.IsZero() || v.Len() == 0 {
		return nil
	}

	hh, err := headers(t.Elem())
	if err != nil {
		return err
	}

	err = e.out.Write(hh)
	if err != nil {
		return err
	}

	e.out.Flush()

	oo := castToSlice(obj)
	for i := 0; i < len(oo); i++ {
		vv, err := values(oo[i])
		if err != nil {
			return err
		}
		err = e.out.Write(vv)
		if err != nil {
			return err
		}
		e.out.Flush()
	}

	e.out.Flush()

	return nil
}
