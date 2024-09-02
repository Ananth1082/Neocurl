package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	aliasValidator *regexp.Regexp = regexp.MustCompile(`^#([[:word:]]+)$`)
	varValidator   *regexp.Regexp = regexp.MustCompile(`([[:word:]]+)="([[:word:]]+)"`)
	verbValidator  *regexp.Regexp = regexp.MustCompile(`(?P<verb>GET|POST|PUT|DELETE) +"(?P<url>.+)"`)
	// urlValidator   *regexp.Regexp = regexp.MustCompile(`(?:http|https):\/\/(?:.*\/)*(?:.*)*`)
	// identifierValidator *regexp.Regexp = regexp.MustCompile(`[[:word:]]`)
)

var (
	ErrInvalidCollectionName = errors.New("invalid collection name")
	ErrInvalidVariableName   = errors.New("invalid variable name")
	ErrInvalidVerbName       = errors.New("invalid verb name")
)

type Request struct {
	verb string
	url  string
	head map[string]string
	body []byte
}

func main() {
	script := `#Collection_name
SET key1="value1"
key2="value2"
#Request1
GET "https;//google.com"
HEAD "This is the head
*** *** ***
"
BODY "This is the body
*** *** ***
#Request2
GET "https;//google.com"
HEAD "This is the head
*** *** ***
"
BODY "This is the body
*** *** ***
"`
	lines := strings.Split(script, "\n")
	fmt.Println("lines ", lines)
	fmt.Println(ParseCollectionName(lines[0]))

	boundaries := FindEndOfBlocks(lines)
	bos := FindSetBlock(lines[boundaries[0]:boundaries[1]])
	fmt.Println(ParseVariables(lines[bos:boundaries[1]]))
	prev := 0
	for ind := range boundaries {
		fmt.Println(lines[prev:ind])
		prev = ind
	}

}

// func processRequest() []Request {
// 	requests := make([]Request, 0, n)

// 	for

// 	return requests
// }

func ParseCollectionName(collectionLine string) (string, error) {
	cname := aliasValidator.FindStringSubmatch(collectionLine)
	if cname == nil {
		return "", ErrInvalidCollectionName
	}
	return cname[1], nil
}

func ParseVariables(setBlock []string) (map[string]string, error) {
	varLookUpTable := make(map[string]string)
	for _, line := range setBlock {
		matches := varValidator.FindStringSubmatch(line)
		if matches == nil {
			return nil, ErrInvalidVariableName
		}
		key, value := matches[1], matches[2]
		fmt.Println(key, ":", value)
		varLookUpTable[key] = value
	}
	return varLookUpTable, nil
}

func ParseRequest() {

}

func FindEndOfBlocks(lines []string) []int {
	var boundaries []int
	for i, line := range lines {
		if aliasValidator.MatchString(line) {
			boundaries = append(boundaries, i)
		}
	}
	return boundaries
}

func FindSetBlock(metaBlock []string) int {
	for i, line := range metaBlock {
		if strings.HasPrefix(line, "SET") {
			return i
		}
	}
	return -1
}

func (r *Request) setVerbAndURL(verbLine string) error {
	matches := verbValidator.FindStringSubmatch(verbLine)
	if matches == nil {
		return ErrInvalidVerbName
	}
	r.verb, r.url = matches[1], matches[2]
	return nil
}

func setHead() {

}

func setBody() {

}
