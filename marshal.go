package maker

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type WriteSeeker interface {
	io.WriteSeeker

	Truncate(size int64) error
	Sync() error
}

func unmarshalInto(fd io.Reader, out interface{}) error {
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}

func marshalInto(in interface{}, fd io.Writer) error {
	marshalBytes, err := yaml.Marshal(in)
	if err != nil {
		return err
	}

	_, err = fd.Write(marshalBytes)
	if err != nil {
		return err
	}

	return nil
}
