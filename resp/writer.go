package resp

import "io"

type RespWriter struct {
	writer io.Writer
}

func NewRespWriter(w io.Writer) *RespWriter {
	return &RespWriter{writer: w}
}

func (w *RespWriter) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
