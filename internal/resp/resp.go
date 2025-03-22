package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	SimpleString = '+'
	Error        = '-'
	Integer      = ':'
	BulkString   = '$'
	Array        = '*'
)

type Value struct {
	Type  byte
	Data  interface{}
	Error error
}

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(reader),
	}
}

func (p *Parser) Parse() (*Value, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading line: %v", err)
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil, fmt.Errorf("empty line")
	}

	switch line[0] {
	case SimpleString:
		return p.parseSimpleString(line[1:])
	case Error:
		return p.parseError(line[1:])
	case Integer:
		return p.parseInteger(line[1:])
	case BulkString:
		return p.parseBulkString(line[1:])
	case Array:
		return p.parseArray(line[1:])
	default:
		return nil, fmt.Errorf("unknown type: %c", line[0])
	}
}

func (p *Parser) parseSimpleString(data string) (*Value, error) {
	return &Value{
		Type: SimpleString,
		Data: data,
	}, nil
}

func (p *Parser) parseError(data string) (*Value, error) {
	return &Value{
		Type: Error,
		Data: data,
	}, nil
}

func (p *Parser) parseInteger(data string) (*Value, error) {
	val, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid integer: %v", err)
	}
	return &Value{
		Type: Integer,
		Data: val,
	}, nil
}

func (p *Parser) parseBulkString(data string) (*Value, error) {
	length, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid bulk string length: %v", err)
	}

	if length == -1 {
		return &Value{
			Type: BulkString,
			Data: nil,
		}, nil
	}

	// Read the actual string
	str := make([]byte, length)
	_, err = io.ReadFull(p.reader, str)
	if err != nil {
		return nil, fmt.Errorf("error reading bulk string: %v", err)
	}

	// Read the trailing \r\n
	_, err = p.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading bulk string terminator: %v", err)
	}

	return &Value{
		Type: BulkString,
		Data: string(str),
	}, nil
}

func (p *Parser) parseArray(data string) (*Value, error) {
	length, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid array length: %v", err)
	}

	if length == -1 {
		return &Value{
			Type: Array,
			Data: nil,
		}, nil
	}

	array := make([]*Value, length)
	for i := int64(0); i < length; i++ {
		val, err := p.Parse()
		if err != nil {
			return nil, fmt.Errorf("error parsing array element: %v", err)
		}
		array[i] = val
	}

	return &Value{
		Type: Array,
		Data: array,
	}, nil
}

// WriteValue writes a RESP value to the given writer
func WriteValue(w io.Writer, v *Value) error {
	switch v.Type {
	case SimpleString:
		_, err := fmt.Fprintf(w, "+%s\r\n", v.Data.(string))
		return err
	case Error:
		_, err := fmt.Fprintf(w, "-%s\r\n", v.Data.(string))
		return err
	case Integer:
		_, err := fmt.Fprintf(w, ":%d\r\n", v.Data.(int64))
		return err
	case BulkString:
		if v.Data == nil {
			_, err := fmt.Fprintf(w, "$-1\r\n")
			return err
		}
		str := v.Data.(string)
		_, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len(str), str)
		return err
	case Array:
		if v.Data == nil {
			_, err := fmt.Fprintf(w, "*-1\r\n")
			return err
		}
		array := v.Data.([]*Value)
		_, err := fmt.Fprintf(w, "*%d\r\n", len(array))
		if err != nil {
			return err
		}
		for _, elem := range array {
			if err := WriteValue(w, elem); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown type: %c", v.Type)
	}
}
