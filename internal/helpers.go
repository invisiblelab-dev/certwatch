package certwatch

import (
	"errors"
	"io/fs"
	"os"
)

func MissingFile(fileName string, data []byte, err error) ([]byte, error) {
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err := os.WriteFile(fileName, nil, 0644)
			if err != nil {
				return nil, err
			}
			data, err := os.ReadFile(fileName)
			if err != nil {
				return nil, err
			}
			return data, nil
		} else {
			return nil, err
		}
	}
	return data, err
}
