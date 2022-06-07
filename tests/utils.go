package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"treehole_next/bootstrap"
)

var App = bootstrap.Init()

// testCommon tests status code and returns response body in bytes
func testCommon(t *testing.T, method string, route string, statusCode int, data ...fiber.Map) []byte {
	var requestData []byte
	var err error

	if len(data) > 0 { // data[0] is request data
		requestData, err = json.Marshal(data[0])
		assert.Nilf(t, err, "encode request body")
	}
	req, err := http.NewRequest(
		strings.ToUpper(method),
		route,
		bytes.NewBuffer(requestData),
	)
	assert.Nilf(t, err, "constructs http request")

	res, err := App.Test(req)
	assert.Nilf(t, err, "perform request")
	assert.Equalf(t, statusCode, res.StatusCode, "status code")

	responseBody, err := ioutil.ReadAll(res.Body)
	assert.Nilf(t, err, "decode response")

	return responseBody
}

// testAPI inherits testCommon, decodes response body to json, tests whether it's expected and returns json
func testAPI(t *testing.T, method string, route string, statusCode int, data ...fiber.Map) fiber.Map {
	responseBody := testCommon(t, method, route, statusCode, data...)
	var responseData fiber.Map
	err := json.Unmarshal(responseBody, &responseData)
	assert.Nilf(t, err, "decode response")

	if len(data) > 1 { // data[1] is response data
		assert.Equalf(t, data[1], responseData, "response data")
	}

	return responseData
}
