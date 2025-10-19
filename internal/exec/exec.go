package exec

import (
	"errors"
	"os"
	"os/exec"
)

func BinaryAndPrint(path string) error {
	cmd := exec.Command("go", "run", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var exitErr *exec.ExitError
	if err := cmd.Run(); err != nil && !errors.As(err, &exitErr) {
		return err
	}

	return nil
}

func BinaryAndCapture(path string) ([]byte, error) {
	return exec.Command("go", "run", path).Output()
}

func CommandAndCapture(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}
