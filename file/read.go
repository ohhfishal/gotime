package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

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
