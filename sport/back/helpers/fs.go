package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

func MkdirAllIfNotExists(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0700); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a dir")
	}
	return nil
}

func JsonDeserializeFile(path string, v interface{}) error {
	fc, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(fc, &v)
}
