package resp

import (
	"reflect"
	"testing"
)

func TestMarshallError(t *testing.T) {
	v1 := Value{Type: "error", String: "something went wrong"}
	expected := []byte("-something went wrong\r\n")

	b := v1.marshallError()

	if !reflect.DeepEqual(b, expected) {
		t.Errorf("returned %q, expected %q", b, expected)
	}
}

func TestMarshallNull(t *testing.T) {
	v1 := Value{Type: "null"}
	expected := []byte("$-1\r\n")

	b := v1.marshallNull()

	if !reflect.DeepEqual(b, expected) {
		t.Errorf("returned %q, expected %q", b, expected)
	}
}

func TestMarshallBulk(t *testing.T) {
	v1 := Value{Type: "bulk", String: "foobar"}
	expected := []byte("$6\r\nfoobar\r\n")

	b := v1.marshalBulk()

	if !reflect.DeepEqual(b, expected) {
		t.Errorf("returned %q, expected %q", b, expected)
	}
}

func TestMarshallString(t *testing.T) {
	v1 := Value{Type: "string", String: "OK"}
	expected := []byte("+OK\r\n")

	b := v1.marshalString()

	if !reflect.DeepEqual(b, expected) {
		t.Errorf("returned %q, expected %q", b, expected)
	}
}

func TestMarshallArray(t *testing.T) {
	v1 := Value{Type: "array", Array: []Value{
		{Type: "string", String: "OK"},
		{Type: "null"},
	}}
	expected := [][]byte{[]byte("+OK\r\n"), []byte("$-1\r\n")}

	for i := range expected {
		b := v1.Array[i].Marshal()
		if !reflect.DeepEqual(b, expected[i]) {
			t.Errorf("returned %q, expected %q", b, expected)
		}
	}
}
