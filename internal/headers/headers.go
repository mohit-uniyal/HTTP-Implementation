package headers

import (
	"bytes"
	"fmt"
	"http/internal/request"
	"log"
	"strings"
)

func isValidFieldName(fieldName string) bool {
	if len(fieldName) == 0 {
		return false
	}

	if strings.HasSuffix(fieldName, " ") {
		return false
	}

	return true
}

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
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

	h[fieldName] = fieldValue

	return nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {

	bytesConsumed := 0
	done := false

	for {
		idx := bytes.Index(data, []byte(request.REQUEST_SEPARATOR))
		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			bytesConsumed += len(request.REQUEST_SEPARATOR)
			break
		}

		if err := h.parseFieldLine(string(data[:idx])); err != nil {
			return 0, false, err
		}

		bytesConsumed += len(data[:idx]) + len(request.REQUEST_SEPARATOR)
		data = data[idx+len(request.REQUEST_SEPARATOR):]
	}

	return bytesConsumed, done, nil

}
