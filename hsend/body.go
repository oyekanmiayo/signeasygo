package hsend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// BodyProvider provides Body content for http.Request attachment.
type BodyProvider interface {
	ContentType() string
	Body() (io.Reader, error)
}

type multipartFormProvider struct {
	payload     interface{}
	contentType string
}

func (p multipartFormProvider) ContentType() string {
	return p.contentType
}

func (p multipartFormProvider) Body() (io.Reader, error) {
	pl, ok := p.payload.(io.Reader)
	if ok {
		return pl, nil
	}
	return nil, errors.New(fmt.Sprintf("invalid payload: %v. should be of type io.Reader", p.payload))
}

type jsonBodyProvider struct {
	payload interface{}
}

func (p jsonBodyProvider) ContentType() string {
	return "application/json"
}

func (p jsonBodyProvider) Body() (io.Reader, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(p.payload)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
