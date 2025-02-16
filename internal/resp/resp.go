package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	ErrInvalidRESP = errors.New("invalid RESP format")
)

func ProcessRESP(reader *bufio.Reader) (cmd string, args []interface{}, err error) {
	parsedData, err := ParseRESP(reader)
	if err != nil {
		return "", nil, err
	}

	arr, ok := parsedData.([]interface{})
	if !ok || len(arr) == 0 {
		return "", nil, ErrInvalidRESP
	}

	cmd, ok = arr[0].(string)
	if !ok {
		return "", nil, fmt.Errorf("command must be a string")
	}

	args = arr[1:]
	return cmd, args, nil
}

func ParseRESP(reader *bufio.Reader) (interface{}, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch prefix {
	case '+':
		return parseString(reader)
	case '-':
		return parseError(reader)
	case ':':
		return parseInteger(reader)
	case '$':
		return parseBulkString(reader)
	case '*':
		return parseArray(reader)
	default:
		return nil, ErrInvalidRESP
	}
}

func parseString(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(line, "\r\n"), nil
}

func parseError(reader *bufio.Reader) (string, error) {
	return parseString(reader)
}

func parseInteger(reader *bufio.Reader) (int64, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSuffix(line, "\r\n"), 10, 64)
}

func parseBulkString(reader *bufio.Reader) (string, error) {
	lengthLine, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	length, err := strconv.Atoi(strings.TrimSuffix(lengthLine, "\r\n"))
	if err != nil || length < -1 {
		return "", ErrInvalidRESP
	}
	if length == -1 {
		return "", nil
	}

	data := make([]byte, length+2)
	if _, err := io.ReadFull(reader, data); err != nil {
		return "", err
	}
	return string(data[:length]), nil
}

func parseArray(reader *bufio.Reader) ([]interface{}, error) {
	lengthLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(strings.TrimSuffix(lengthLine, "\r\n"))
	if err != nil || length < -1 {
		return nil, ErrInvalidRESP
	}
	if length == -1 {
		return nil, nil
	}

	array := make([]interface{}, length)
	for i := 0; i < length; i++ {
		element, err := ParseRESP(reader)
		if err != nil {
			return nil, err
		}
		array[i] = element
	}
	return array, nil
}
