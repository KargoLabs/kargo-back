package phoneValidation

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// TruoraError is the struct handler for Truora Error responses
type TruoraError struct {
	Code     int    `json:"code"`
	HTTPCode int    `json:"http_code"`
	Message  string `json:"message"`
}

func getTruoraErrorFromBody(body io.ReadCloser) (*TruoraError, error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	truoraError := TruoraError{}

	err = json.Unmarshal(bodyBytes, &truoraError)
	if err != nil {
		return nil, err
	}

	return &truoraError, nil
}
