package phoneValidation

import (
	kargoClient "kargo-back/shared/client"
	"kargo-back/shared/environment"
	"net/http"
	"net/url"
)

var (
	truoraAPIKey              = environment.GetString("TRUORA_API_KEY", "")
	truoraValidationsURL      = "https://api.validations.truora.com"
	truoraEnrollmentsEndpoint = "/v1/enrollments"
	truoraValidationsEnpoint  = "/v1/validations"

	client *kargoClient.Client
)

// TruoraResponse is the struct wrapper of both Truora responses
type TruoraResponse struct {
	*TruoraValidation
	*TruoraError
}

func init() {
	client = kargoClient.NewClient()
}

// EnrollPhone does enroll request to Truora phone validation and returns response if it is not succesful
func EnrollPhone(phoneNumber, truoraAccountID string) (*TruoraResponse, error) {
	params := url.Values{}

	params.Set("type", "phone-verification")
	params.Set("phone_number", phoneNumber)
	params.Set("phone_type", "home")
	params.Set("user_authorized", "true")
	params.Set("account_id", truoraAccountID)

	headers := http.Header{}

	headers.Set("Truora-API-Key", truoraAPIKey)

	response, err := client.PostWithURLEncodedParams(truoraValidationsURL+truoraEnrollmentsEndpoint, params, headers)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		truoraError, err := getTruoraErrorFromBody(response.Body)
		if err != nil {
			return nil, err
		}

		// Need to overwrite because some times Truora Response does not have it
		truoraError.HTTPCode = response.StatusCode

		return &TruoraResponse{
			TruoraError: truoraError,
		}, nil
	}

	// Enroll Response just matters to use if it is not succesful
	return nil, nil
}

// CreatePhoneValidation creates Truora phone validation and returns TruoraResponse
func CreatePhoneValidation(truoraAccountID string) (*TruoraResponse, error) {
	params := url.Values{}

	params.Set("type", "phone-verification")
	params.Set("verify_channel", "sms")
	params.Set("phone_locale", "es")
	params.Set("phone_type", "home")
	params.Set("account_id", truoraAccountID)

	headers := http.Header{}

	headers.Set("Truora-API-Key", truoraAPIKey)

	response, err := client.PostWithURLEncodedParams(truoraValidationsURL+truoraValidationsEnpoint, params, headers)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		truoraError, err := getTruoraErrorFromBody(response.Body)
		if err != nil {
			return nil, err
		}

		// Need to overwrite because some times Truora Response does not have it
		truoraError.HTTPCode = response.StatusCode

		return &TruoraResponse{
			TruoraError: truoraError,
		}, nil
	}

	truoraValidation, err := getTruoraTruoraValidationFromBody(response.Body)
	if err != nil {
		return nil, err
	}

	return &TruoraResponse{
		TruoraValidation: truoraValidation,
	}, nil
}

// PerformPhoneValidation performs Truora phone validation and returns TruoraResponse
func PerformPhoneValidation(code, validationID string) (*TruoraResponse, error) {
	params := url.Values{}

	params.Set("token", code)

	headers := http.Header{}

	headers.Set("Truora-API-Key", truoraAPIKey)

	response, err := client.PostWithURLEncodedParams(truoraValidationsURL+truoraValidationsEnpoint+"/"+validationID, params, headers)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		truoraError, err := getTruoraErrorFromBody(response.Body)
		if err != nil {
			return nil, err
		}

		// Need to overwrite because some times Truora Response does not have it
		truoraError.HTTPCode = response.StatusCode

		return &TruoraResponse{
			TruoraError: truoraError,
		}, nil
	}

	truoraValidation, err := getTruoraTruoraValidationFromBody(response.Body)
	if err != nil {
		return nil, err
	}

	return &TruoraResponse{
		TruoraValidation: truoraValidation,
	}, nil
}
