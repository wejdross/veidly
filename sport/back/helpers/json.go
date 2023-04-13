package helpers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
)

func JsonMustSerialize(obj interface{}) []byte {
	var b []byte
	var err error
	if b, err = json.Marshal(obj); err != nil {
		panic(err)
	}
	return b
}

func JsonMustSerializeStr(obj interface{}) string {
	return string(JsonMustSerialize(obj))
}

func JsonMustSerializeFormatStr(obj interface{}) string {
	var b []byte
	var err error
	if b, err = json.MarshalIndent(obj, "", "    "); err != nil {
		panic(err)
	}
	return string(b)
}

func JsonMustDeserialize(v []byte, obj interface{}) {
	err := json.Unmarshal(v, obj)
	if err != nil {
		panic(err)
	}
}

func JsonMustSerializeReader(obj interface{}) *strings.Reader {
	s := JsonMustSerializeStr(obj)
	return strings.NewReader(s)
}

func ParseAndValidateJson(j []byte, val interface{}, valFunc func() error) error {

	var err error

	if err = json.Unmarshal(j, val); err != nil {
		return err
	}

	if valFunc != nil {
		if err = valFunc(); err != nil {
			return err
		}
	}

	return nil
}

/*
	Unmarshals (deserializes) value from io.Reader
	Note that if val implements Validate interface (present in this package)
	return value will be validated
	If mustValidate flag is specified then
		this function panics if val does not implement said interface
*/
func ReadJsonBodyFromReader(
	r io.Reader, val interface{}, valFunc func() error,
) error {

	var jb []byte
	var err error

	if jb, err = ioutil.ReadAll(r); err != nil {
		return err
	}

	return ParseAndValidateJson(jb, val, valFunc)
}
