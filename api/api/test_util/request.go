package test_util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gofiber/fiber/v2"
)

type Query = map[string]string

type Body = map[string]interface{}

type RequestParams struct {
	Method string
	URL    string
	Query  fiber.Map
	Body   fiber.Map
}

func NewRequest(
	t *testing.T,
	params RequestParams,
) *http.Request {
	var bodyReader io.Reader

	if params.Body != nil {
		body, err := json.Marshal(params.Body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(body)
	}

	request, err := http.NewRequest(params.Method, params.URL, bodyReader)
	require.NoError(t, err)

	if params.Query != nil {
		query := request.URL.Query()
		for key, value := range params.Query {
			query.Add(key, value.(string))
		}
		request.URL.RawQuery = query.Encode()
	}

	request.Header.Set("Content-Type", "application/json")

	return request
}

type RequestParams2 struct {
	Method string
	URL    string
	Query  Query
	Body   Body
}

func NewRequest2(
	t *testing.T,
	params RequestParams2,
) *http.Request {
	var bodyReader io.Reader

	if params.Body != nil {
		body, err := json.Marshal(params.Body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(body)
	}

	request, err := http.NewRequest(params.Method, params.URL, bodyReader)
	require.NoError(t, err)

	if params.Query != nil {
		query := request.URL.Query()
		for key, value := range params.Query {
			query.Add(key, value)
		}
		request.URL.RawQuery = query.Encode()
	}

	request.Header.Set("Content-Type", "application/json")

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
