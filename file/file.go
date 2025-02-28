package file

import (
	"bufio"
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

func ReadAllFrom[T any](path string, decode func(string) (T, error)) ([]T, error) {
	file, err := os.Open(path)
	if err != nil {
		return []T{}, fmt.Errorf(`can not open "%s": %w`, path, err)
	}
	defer file.Close()
	return ReadAll(file, decode)
}

func ReadAll[T any](reader io.Reader, decode func(string) (T, error)) ([]T, error) {
	read := []T{}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		next, err := decode(scanner.Text())
		if err != nil {
			return []T{}, fmt.Errorf("reading line: %w", err)
		}
		read = append(read, next)
	}

	if err := scanner.Err(); err != nil {
		return []T{}, fmt.Errorf("reading: %w", err)
	}
	return read, nil
}
