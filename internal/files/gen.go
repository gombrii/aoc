package files

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func Gen(structure map[string]string, data map[string]string) error {
	for fPath, tmpl := range structure {
		if _, err := os.Stat(fPath); err == nil {
			fmt.Printf("skipping %s, already exists\n", fPath)
			continue
		}
		fmt.Printf("creating %s\n", fPath)

		if err := os.MkdirAll(filepath.Dir(fPath), 0755); err != nil {
			return fmt.Errorf("creating dir %s: %v", filepath.Dir(fPath), err)
		}

		file, err := os.Create(fPath)
		if err != nil {
			return fmt.Errorf("creating file %s: %v", fPath, err)
		}
		defer file.Close()

		if err = template.Must(template.New(filepath.Base(fPath)).Parse(tmpl)).Execute(file, data); err != nil {
			os.Remove(file.Name())
			return fmt.Errorf("compiling file %s: %v", fPath, err)
		}
	}

	return nil
}

func GenTemp(files map[string]string, data map[string]string) (map[string]string, error) {
	tempFiles := make(map[string]string)

	for fName, tmpl := range files {
		file, err := os.CreateTemp("", "*")
		if err != nil {
			return nil, fmt.Errorf("creating temp file: %v", err)
		}
		defer file.Close()

		if err = template.Must(template.New(fName).Parse(tmpl)).Execute(file, data); err != nil {
			os.Remove(file.Name())
			return nil, fmt.Errorf("compiling: %v", err)
		}

		tempFiles[fName] = file.Name()
	}

	return tempFiles, nil
}
