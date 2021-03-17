package es

import (
	"bytes"
	"encoding/json"
	"io"
)

func (body SearchBody) encodeSearchBody() (io.Reader, error) {
	buf := &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}
	return buf, nil
}
