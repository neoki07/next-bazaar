package test_util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gofiber/fiber/v2"
)

type RequestParams struct {
	Method    string
	URL       string
	Body      fiber.Map
	SetupAuth func(request *http.Request)
}

func NewRequest(
	t *testing.T,
	params RequestParams,
) *http.Request {
	body, err := json.Marshal(params.Body)
	require.NoError(t, err)

	request, err := http.NewRequest(params.Method, params.URL, bytes.NewReader(body))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")

	if params.SetupAuth != nil {
		params.SetupAuth(request)
	}

	return request
}

func SendRequest(
	t *testing.T,
	app *fiber.App,
	request *http.Request,
) *http.Response {
	response, err := app.Test(request)
	require.NoError(t, err)

	return response
}
