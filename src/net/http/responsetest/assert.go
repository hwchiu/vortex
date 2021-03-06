package responsetest

import (
	"net/http/httptest"
	"testing"

	response "github.com/hwchiu/vortex/src/net/http"

	"encoding/json"
	"encoding/xml"

	"github.com/stretchr/testify/assert"
)

// AssertStatusEqual will assert equal status
func AssertStatusEqual(t *testing.T, resp *httptest.ResponseRecorder, status int) {
	assert.Equal(t, status, resp.Code)
}

// AssertErrorMessage will assert error message
func AssertErrorMessage(t *testing.T, resp *httptest.ResponseRecorder, msg string) (err error) {
	var payload = response.ErrorPayload{}
	var out = resp.Body.Bytes()
	var contentType = resp.Header().Get("content-type")

	switch contentType {
	case "application/json", "text/json":
		err = json.Unmarshal(out, &payload)
	case "application/xml", "text/xml":
		err = xml.Unmarshal(out, &payload)
	}
	assert.Equal(t, msg, payload.Message)
	return err
}

// AssertError will assert error
func AssertError(t *testing.T, resp *httptest.ResponseRecorder) (err error) {
	var payload = response.ErrorPayload{}
	var out = resp.Body.Bytes()
	var contentType = resp.Header().Get("content-type")

	switch contentType {
	case "application/json", "text/json":
		err = json.Unmarshal(out, &payload)
	case "application/xml", "text/xml":
		err = xml.Unmarshal(out, &payload)
	}
	assert.NoError(t, err)
	return err
}
