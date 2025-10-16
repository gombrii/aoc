package files

import (
	"fmt"
	"os"
)

const (
	Runner = "runner.go"
	Lock   = "lock"
	Res    = "res"
	Dur    = "dur"
)

func ReadAll(files map[string]string) (map[string]string, error) {
	data := make(map[string]string, len(files))

	for name, path := range files {
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading file: %v", err)
		}
		data[name] = string(bytes)
	}

	return data, nil
}
