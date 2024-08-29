package parser

import (
	"os"
)

func ReadFile(filePath string) (string, error) {
	filebyt, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(filebyt), nil
}

func Separtor(program string) []string {
	return []string{""}
}
