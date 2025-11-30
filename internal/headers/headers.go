package headers

import (
	"bytes"
	"fmt"
	"http/internal/constants"
	"log"
	"strings"
)

var VALID_FIELD_NAME_CHARACTERS = func() map[rune]bool {
	m := make(map[rune]bool)

	// A–Z
	for c := 'A'; c <= 'Z'; c++ {
		m[c] = true
	}

	// a–z
	for c := 'a'; c <= 'z'; c++ {
		m[c] = true
	}

	// 0–9
	for c := '0'; c <= '9'; c++ {
		m[c] = true
	}

	// Special characters
	specials := "!#$%&'*+-.^_`|~"
	for _, c := range specials {
		m[c] = true
	}

	return m
}()

func isValidFieldName(fieldName string) bool {
	if len(fieldName) == 0 {
		return false
	}

	if strings.HasSuffix(fieldName, " ") {
		return false
	}

	for _, char := range fieldName {
		if _, exists := VALID_FIELD_NAME_CHARACTERS[char]; !exists {
			return false
		}
	}

	return true
}

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(name string) string {
	return h[strings.ToLower(name)]
}

func (h Headers) set(name, value string) {
	fieldName := strings.ToLower(name)

	fieldValue, exists := h[fieldName]

	if exists {
		h[fieldName] = fmt.Sprintf("%s,%s", fieldValue, value)
	} else {
		h[fieldName] = value
	}
}

func (h Headers) parseFieldLine(fieldLineData string) error {
	//1. find the first : and separate the field name from field value

	splitIndex := strings.Index(fieldLineData, ":")
	if splitIndex == -1 {
		log.Println(": missing")
		return fmt.Errorf(": is missing")
	}

	fieldName := fieldLineData[:splitIndex]
	fieldValue := fieldLineData[splitIndex+1:]

	//2. trim the necessary white space in field name

	fieldName = strings.TrimLeft(fieldName, " ")

	//3. trim the necessary white space in field value

	fieldValue = strings.TrimSpace(fieldValue)

	//4. validate the field name

	if !isValidFieldName(fieldName) {
		log.Println("not a valid field name: ", fieldName)
		return fmt.Errorf("not a valid field name: %s", fieldName)
	}

	//5. set the header with field name set to all lowercase

	h.set(fieldName, fieldValue)

	return nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {

	bytesConsumed := 0
	done := false

	for {
		idx := bytes.Index(data, []byte(constants.REQUEST_SEPARATOR))
		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			bytesConsumed += len(constants.REQUEST_SEPARATOR)
			break
		}

		if err := h.parseFieldLine(string(data[:idx])); err != nil {
			return 0, false, err
		}

		bytesConsumed += len(data[:idx]) + len(constants.REQUEST_SEPARATOR)
		data = data[idx+len(constants.REQUEST_SEPARATOR):]
	}

	return bytesConsumed, done, nil

}
