package files

import (
	"os"
)

func Write(path string, data []byte) error {
	return os.WriteFile(path, data, 0755)
}
