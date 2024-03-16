package resp

import (
	"reflect"
	"testing"
)

type MockWriter struct {
	WrittenData []byte
}

func (w *MockWriter) Write(p []byte) (n int, err error) {
	w.WrittenData = append(w.WrittenData, p...)
	return len(p), nil
}

func TestWriter(t *testing.T) {
	mockWriter := &MockWriter{}
	writer := NewRespWriter(mockWriter)

	v := Value{Type: "null"}

	err := writer.Write(v)
	if err != nil {
		t.Errorf("Error: %q", err)
	}

	expected := v.Marshal()

	if !reflect.DeepEqual(mockWriter.WrittenData, expected) {
		t.Errorf("returned %q, expected %q", mockWriter.WrittenData, expected)
	}
}
