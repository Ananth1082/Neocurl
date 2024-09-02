package utils

import "os"

func ReadFile(filename string) (string, error) {
	program, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(program), nil
}
