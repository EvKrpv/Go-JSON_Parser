package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type JSONParser struct {
	input string
	pos   int
}

func NewJSONParser(input string) *JSONParser {
	return &JSONParser{
		input: input,
		pos:   0,
	}
}

func (p *JSONParser) current() byte {
	if p.pos >= len(p.input) {
		return 0
	}
	return p.input[p.pos]
}

func (p *JSONParser) skipSpace() {
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			p.pos++
		} else {
			break
		}
	}
}

func (p *JSONParser) parseString() (string, error) {
	if p.current() != '"' {
		return "", fmt.Errorf("expected '\"' but got '%c'", p.current())
	}

	p.pos++
	start := p.pos

	for p.pos < len(p.input) {
		if p.input[p.pos] == '"' {
			str := p.input[start:p.pos]
			p.pos++
			return str, nil
		}
		p.pos++
	}
	return "", fmt.Errorf("unterminated string")
}

func (p *JSONParser) parseNumber() (interface{}, error) {
	start := p.pos

	if p.current() == '-' {
		p.pos++
	}

	for p.pos < len(p.input) {
		ch := p.current()
		if ch >= '0' && ch <= '9' {
			p.pos++
		} else {
			break
		}
	}

	if p.pos == start || (p.input[start] == '-' && p.pos == start+1) {
		return nil, fmt.Errorf("Invalid number")
	}

	numberStr := p.input[start:p.pos]
	return numberStr, nil
}
func (p *JSONParser) parseObject() (map[string]interface{}, error) {
	if p.current() != '{' {
		return nil, fmt.Errorf("Expected '{'")
	}

	p.pos++
	p.skipSpace()

	obj := make(map[string]interface{})

	for {
		if p.current() == '}' {
			p.pos++
			return obj, nil
		}

		key, err := p.parseString()
		if err != nil {
			return nil, err
		}

		p.skipSpace()

		if p.current() != ':' {
			return nil, fmt.Errorf("expected ':'")
		}
		p.pos++
		p.skipSpace()

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		obj[key] = value

		p.skipSpace()

		if p.current() == ',' {
			p.pos++
			p.skipSpace()
		} else if p.current() == '}' {
			p.pos++
			return obj, nil
		} else {
			return nil, fmt.Errorf("expected ',' or '}'")
		}
	}
}

func (p *JSONParser) parseArray() ([]interface{}, error) {
	if p.current() != '[' {
		return nil, fmt.Errorf("Expected '['")
	}

	p.pos++
	p.skipSpace()

	arr := []interface{}{}

	for {
		if p.current() == ']' {
			p.pos++
			return arr, nil
		}

		value, err := p.parseValue()

		if err != nil {
			return nil, err
		}

		arr = append(arr, value)

		p.skipSpace()

		if p.current() == ',' {
			p.pos++
			p.skipSpace()
		} else if p.current() == ']' {
			p.pos++
			return arr, nil
		} else {
			return nil, fmt.Errorf("expected ',' or ']'")
		}
	}
}
func (p *JSONParser) parseTrue() (bool, error) {
	if p.pos+4 <= len(p.input) && p.input[p.pos:p.pos+4] == "true" {
		p.pos += 4
		return true, nil
	}
	return false, fmt.Errorf("Expecte 'true")
}

func (p *JSONParser) parseFalse() (bool, error) {
	if p.pos+5 <= len(p.input) && p.input[p.pos:p.pos+5] == "false" {
		p.pos += 5
		return false, nil
	}
	return false, fmt.Errorf("Expected 'false")
}

func (p *JSONParser) parseNull() (interface{}, error) {
	if p.pos+4 <= len(p.input) && p.input[p.pos:p.pos+4] == "null" {
		p.pos += 4
		return nil, nil
	}
	return nil, fmt.Errorf("Expected 'null")

}

func (p *JSONParser) parseValue() (interface{}, error) {
	p.skipSpace()

	switch p.current() {
	case '"':
		return p.parseString()
	case '[':
		return p.parseArray()
	case '{':
		return p.parseObject()
	case 't':
		return p.parseTrue()
	case 'f':
		return p.parseFalse()
	case 'n':
		return p.parseNull()
	default:
		return p.parseNumber()
	}
}
func PathIterator(data interface{}, path ...string) (func(func(string, interface{}) bool), error) {
	current := data

	if len(path) == 0 {
		return createIterator(current)
	}

	for i, key := range path {
		switch v := current.(type) {
		case map[string]interface{}:
			next, exists := v[key]
			if !exists {
				return nil, fmt.Errorf("key '%s' not found in object", key)
			}
			current = next

		case []interface{}:
			index, err := parseArrayIndex(key)
			if err != nil {
				return nil, fmt.Errorf("path element '%s' is not an index", key)
			}
			if index < 0 || index >= len(v) {
				return nil, fmt.Errorf("index '%s' is out of range", key)
			}
			current = v[index]

		default:
			if i == len(path)-1 {
				return createIterator(current)
			}
			return nil, fmt.Errorf("path does not lead to an object")
		}
	}

	return createIterator(current)
}

func createIterator(data interface{}) (func(func(string, interface{}) bool), error) {
	switch v := data.(type) {
	case map[string]interface{}:
		return func(callback func(string, interface{}) bool) {
			for key, value := range v {
				if !callback(key, value) {
					break
				}
			}
		}, nil

	case []interface{}:
		return func(callback func(string, interface{}) bool) {
			for i, value := range v {
				if !callback(fmt.Sprintf("%d", i), value) {
					break
				}
			}
		}, nil

	default:
		return nil, fmt.Errorf("path does not lead to an object")
	}
}

func parseArrayIndex(s string) (int, error) {
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0, fmt.Errorf("not a valid index")
		}
	}

	index := 0
	for _, ch := range s {
		index = index*10 + int(ch-'0')
	}
	return index, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var currentData interface{}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)

		switch parts[0] {
		case "PARSE":
			if len(parts) < 2 {
				currentData = nil
				continue
			}
			parser := NewJSONParser(parts[1])
			data, err := parser.parseValue()
			if err != nil {
				currentData = nil
			} else {
				currentData = data
			}

		case "ITERATE":
			var path []string
			if len(parts) > 1 {
				path = strings.Split(parts[1], ".")
			}

			if currentData == nil {
				fmt.Println("Error: invalid JSON")
				continue
			}

			iter, err := PathIterator(currentData, path...)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			} else {
				iter(func(key string, value interface{}) bool {
					fmt.Printf("%s: %v\n", key, value)
					return true
				})
			}

		default:
			currentData = nil
		}
	}
}
