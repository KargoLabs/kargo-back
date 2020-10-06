package phoneValidation

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"time"
)

// TruoraValidation is the struct handler for TruoraValidation responses
type TruoraValidation struct {
	ValidationID     string    `json:"validation_id"`
	AccountID        string    `json:"account_id"`
	Type             string    `json:"type"`
	ValidationStatus string    `json:"validation_status"`
	FailureStatus    string    `json:"failure_status"`
	DeclinedReason   string    `json:"declined_reason"`
	CreationDate     time.Time `json:"creation_date"`
}

func getTruoraTruoraValidationFromBody(body io.ReadCloser) (*TruoraValidation, error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	truoraValidation := TruoraValidation{}

	err = json.Unmarshal(bodyBytes, &truoraValidation)
	if err != nil {
		return nil, err
	}

	return &truoraValidation, nil
}
