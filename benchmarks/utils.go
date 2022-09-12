package benchmarks

import (
	"bytes"
	"encoding/json"
	"github.com/hetiansu5/urlquery"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"treehole_next/bootstrap"
	. "treehole_next/models"
)

var App = bootstrap.Init()

var _ Map

const (
	REQUEST_BODY = iota
	REQUEST_QUERY
)

func benchmarkCommon(b *testing.B, method string, route string, requestType int, data ...Map) []byte {
	var requestData []byte
	var err error
	var req *http.Request

	b.StopTimer()
	switch requestType {
	case REQUEST_BODY:
		if len(data) > 0 && data[0] != nil { // data[0] is request data
			requestData, err = json.Marshal(data[0])
			assert.Nilf(b, err, "encode request body")
		}
		req, err = http.NewRequest(
			strings.ToUpper(method),
			route,
			bytes.NewBuffer(requestData),
		)
	case REQUEST_QUERY:
		req, err = http.NewRequest(
			strings.ToUpper(method),
			route,
			nil,
		)
		if len(data) > 0 && data[0] != nil { // data[0] is query data
			queryData, err := urlquery.Marshal(data[0])
			req.URL.RawQuery = string(queryData)
			assert.Nilf(b, err, "encode request body")
		}
	}

	req.Header.Add("Content-Type", "application/json")
	assert.Nilf(b, err, "constructs http request")

	b.StartTimer()
	res, err := App.Test(req, -1)
	b.StopTimer()
	assert.Nilf(b, err, "perform request")

	responseBody, err := ioutil.ReadAll(res.Body)
	assert.Nilf(b, err, "decode response")

	if res.StatusCode != 200 && res.StatusCode != 201 {
		assert.Fail(b, string(responseBody))
	}
	return responseBody
}
