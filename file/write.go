package file

import (
	"fmt"
	"io"
	"os"
)

type Encodeable interface {
	Encode() string
}

func WriteTo[T Encodeable](path string, perms int, item T) error {
	file, err := os.OpenFile(path, perms, 0644)
	if err != nil {
		return fmt.Errorf(`can not open "%s": %w`, path, err)
	}
	defer file.Close()
	return Write(file, item)

}

func Write[T Encodeable](write io.Writer, item T) error {
	_, err := fmt.Fprintln(write, item.Encode())
	return err

}
