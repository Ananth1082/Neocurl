package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Ananth1082/m-v0.0/utils"
)

var (
	aliasValidator *regexp.Regexp = regexp.MustCompile(`^#([[:word:]]+)$`)
	varValidator   *regexp.Regexp = regexp.MustCompile(`([[:word:]]+)="([[:word:]]+)"`)
	verbValidator  *regexp.Regexp = regexp.MustCompile(`(?P<verb>GET|POST|PUT|DELETE) +"(?P<url>.+)"`)
	headValidator  *regexp.Regexp = regexp.MustCompile(`^HEAD .+$`)
	headExtractor  *regexp.Regexp = regexp.MustCompile(`^HEAD ([\s\S]*)$`)
	bodyExtractor  *regexp.Regexp = regexp.MustCompile(`^BODY ([\s\S]*)$`)
	bodyValidator  *regexp.Regexp = regexp.MustCompile(`^BODY .+$`)
	// urlValidator   *regexp.Regexp = regexp.MustCompile(`(?:http|https):\/\/(?:.*\/)*(?:.*)*`)
	// identifierValidator *regexp.Regexp = regexp.MustCompile(`[[:word:]]`)
)

var (
	ErrInvalidCollectionName = errors.New("invalid collection name")
	ErrInvalidVariableName   = errors.New("invalid variable name")
	ErrInvalidVerbName       = errors.New("invalid verb name")
)

type Request struct {
	name string
	verb string
	url  string
	head map[string]string
	body []byte
}

func main() {
	script, _ := utils.ReadFile("request.ncurl")
	parseScript(script)
}

func parseScript(script string) {
	lines := strings.Split(script, "\n")
	fmt.Println("lines ", lines)
	fmt.Println(ParseCollectionName(lines[0]))
	boundaries := FindEndOfBlocks(lines)
	bos := FindSetBlock(lines[boundaries[0]:boundaries[1]])
	fmt.Println(ParseVariables(lines[bos:boundaries[1]]))
	requests := processRequests(lines, boundaries)
	fmt.Println(requests)
}

func processRequests(lines []string, boundaries []int) []Request {
	requests := make([]Request, 0, len(boundaries)-1)
	for i := 1; i < len(boundaries)-1; i++ {
		r := ParseRequest(lines[boundaries[i]:boundaries[i+1]])
		requests = append(requests, r)
		r.run()
	}
	return requests
}

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

func ParseRequest(requestBlock []string) Request {
	end := len(requestBlock)
	r := Request{}
	r.setRequestName(requestBlock[0])
	r.setVerbAndURL(requestBlock[1])
	h, b := findHeadAndBody(requestBlock)
	fmt.Println("h: ", h, " b: ", b)
	if h != -1 {
		if b == -1 {
			b = end
		}
		r.setHead(strings.Join(requestBlock[h:b], ""))
	}
	if b != -1 {
		r.setBody(strings.Join(requestBlock[b:], "\n"))
	}
	return r
}

func FindEndOfBlocks(lines []string) []int {
	var boundaries []int
	for i, line := range lines {
		if aliasValidator.MatchString(line) {
			boundaries = append(boundaries, i)
		}
	}
	boundaries = append(boundaries, len(lines))
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

func (r *Request) setRequestName(aliasLine string) error {
	rname := aliasValidator.FindStringSubmatch(aliasLine)
	if rname == nil {
		return ErrInvalidCollectionName
	}
	r.name = rname[1]
	return nil
}

func (r *Request) setVerbAndURL(verbLine string) error {
	matches := verbValidator.FindStringSubmatch(verbLine)
	if matches == nil {
		return ErrInvalidVerbName
	}
	r.verb, r.url = matches[1], matches[2]
	return nil
}

func findHeadAndBody(req []string) (int, int) {
	h, b := -1, -1
	for i, line := range req {
		if h == -1 && headValidator.MatchString(line) {
			h = i
		}
		if b == -1 && bodyValidator.MatchString(line) {
			b = i
		}
	}
	return h, b
}

func (r *Request) setHead(headSubBlock string) {
	matches := headExtractor.FindSubmatch([]byte(headSubBlock))
	if matches == nil {
		log.Fatalln("No match")
	}
	err := json.Unmarshal(matches[1], &r.head)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("head ", r.head)
}

func (r *Request) setBody(bodySubBlock string) {
	matches := bodyExtractor.FindSubmatch([]byte(bodySubBlock))
	if matches == nil {
		log.Fatalln("No match")
	}
	r.body = matches[1]
	fmt.Println("body ", string(r.body))
}

func (r Request) run() {
	client := &http.Client{}
	req, err := http.NewRequest(r.verb, r.url, bytes.NewReader(r.body))
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range r.head {
		req.Header.Add(key, value)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("res: ", res)
}
