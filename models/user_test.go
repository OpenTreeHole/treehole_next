package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJWT(t *testing.T) {
	var user User
	jwt := "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1aWQiOjE2LCJpc3MiOiJEU2lSa2NvWDJZV3dta3VqM3FFdFVxSE1uUnNvMjZQYiIsImlhdCI6MTY2MjUyNzg5OSwiaWQiOjE2LCJpc19hZG1pbiI6ZmFsc2UsIm5pY2tuYW1lIjoidXNlciIsIm9mZmVuc2VfY291bnQiOjAsInJvbGVzIjpbXSwidHlwZSI6ImFjY2VzcyIsImV4cCI6MTY2MjUyOTY5OX0.Ov_8cJay-Ta0jsPYUx1D-XDc_D1WK1iTdjnuEKAelaM"
	err := user.parseJWT(jwt)
	assert.Nilf(t, err, "parseJWT failed: %v", err)
}
