package resp

import (
	"bufio"
	"io"
	"reflect"
	"testing"
)

// MockReader is a mock implementation of bufio.Reader.
type MockReader struct {
	data []byte
	pos  int
}

// NewMockReader creates a new MockReader with the given data.
func NewMockReader(data []byte) *MockReader {
	return &MockReader{data: data}
}

// Read reads data from the buffer.
func (r *MockReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func TestReadBulk(t *testing.T) {
	input := []byte("5\r\nhello\r\n")
	expected := Value{Type: "bulk", String: "hello"}

	reader := bufio.NewReader(NewMockReader(input))
	respReader := NewRespReader(reader)

	val, err := respReader.readBulk()

	if err != nil {
		t.Errorf("Error: %q", err)
	}

	if !reflect.DeepEqual(val, expected) {
		t.Errorf("returned %q, expected %q", val, expected)
	}
}

func TestReadLine(t *testing.T) {
	input := []byte("hello\r\n")
	expected := input[:len(input)-2]

	reader := bufio.NewReader(NewMockReader(input))
	respReader := NewRespReader(reader)

	line, _, err := respReader.readLine()

	if err != nil {
		t.Errorf("Error: %q", err)
	}

	if !reflect.DeepEqual(line, expected) {
		t.Errorf("returned %q, expected %q", line, expected)
	}
}

func TestReadInt(t *testing.T) {
	input := []byte("5\r\n")
	expected := 5

	reader := bufio.NewReader(NewMockReader(input))
	respReader := NewRespReader(reader)

	x, _, err := respReader.readInteger()

	if err != nil {
		t.Errorf("Error: %q", err)
	}

	if !reflect.DeepEqual(x, expected) {
		t.Errorf("returned %q, expected %q", x, expected)
	}
}

func TestReadArray(t *testing.T) {
	expected := Value{Type: "array", Array: []Value{
		{Type: "bulk", String: "hgetall"},
		{Type: "bulk", String: "users"},
	}}
	input := []byte("2\r\n$7\r\nhgetall\r\n$5\r\nusers\r\n")

	reader := bufio.NewReader(NewMockReader(input))
	respReader := NewRespReader(reader)

	val, err := respReader.readArray()

	if err != nil {
		t.Errorf("Error: %q", err)
	}

	if !reflect.DeepEqual(val, expected) {
		t.Errorf("returned %q, expected %q", val, expected)
	}
}

func TestRead(t *testing.T) {
	expected := Value{Type: "array", Array: []Value{
		{Type: "bulk", String: "hgetall"},
		{Type: "bulk", String: "users"},
	}}
	input := []byte("*2\r\n$7\r\nhgetall\r\n$5\r\nusers\r\n")

	reader := bufio.NewReader(NewMockReader(input))
	respReader := NewRespReader(reader)

	val, err := respReader.Read()

	if err != nil {
		t.Errorf("Error: %q", err)
	}

	if !reflect.DeepEqual(val, expected) {
		t.Errorf("returned %q, expected %q", val, expected)
	}
}
