package opusfile

import (
	"errors"
	"io"
)

func readBytes(in io.Reader, c int) ([]byte, error) {
	buf := make([]byte, c)
	n, err := in.Read(buf)
	if err != nil {
		return nil, err
	}
	if n != c {
		return nil, errors.New("unexpected number of bytes read")
	}
	return buf, nil
}
