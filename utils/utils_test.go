package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStripContent(t *testing.T) {
	var str string
	str = "愿中国青年都摆脱冷气，只是向上走，不必听自暴自弃者流的话。能做事的做事，能发声的发声。有一分热，发一分光。就令萤火一般，也可以在黑暗里发一点光，不必等候炬火。"
	println(len(str))
	println(len([]rune(str)))
	assert.Equal(t, "愿中国青年都摆脱冷气", StripContent(str, 10))
	assert.Equal(t, str, StripContent(str, 100))
}
