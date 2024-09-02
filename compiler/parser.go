package compiler

import (
	"regexp"
	"strings"
)

var minifyExp *regexp.Regexp 

func init() {
	minifyExp,err = regexp.Compile(` *\n(?:\s*)+`)
	if err!=nil {
		log.Fatal(err)
	}
}

func ParseProgram(program string) {
	cname, program := ParseCollectionName(program)
	
}

func minifyProgram(program string) string {
	
}

func ParseCollectionName(program string) (string, string) {
	start, end := -1, -1
	for i, letter := range program {
		if letter == rune('#') {
			start = i
		}
		if start != -1 && letter == rune('\n') {
			end = i
			break
		}
	}
	return strings.TrimSpace(program[start:end]), program[end:]
}

func SetVaraiblles(program string) map[string]string {
	start, end := -1, -1
	for 
}
