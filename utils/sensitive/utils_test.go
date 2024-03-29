package sensitive

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindImagesInMarkdown(t *testing.T) {

	type wantStruct struct {
		clearContent string
		imageUrls    []string
	}
	tests := []struct {
		text string
		want wantStruct
	}{
		{
			text: `![image1](https://example.com/image1)`,
			want: wantStruct{
				clearContent: `image1`,
				imageUrls:    []string{"https://example.com/image1"},
			},
		},
		{
			text: `![image1](https://example.com/image1) ![image2](https://example.com/image2)`,
			want: wantStruct{
				clearContent: `image1 image2`,
				imageUrls:    []string{"https://example.com/image1", "https://example.com/image2"},
			},
		},
		{
			text: `![image1](https://example.com/image1 "title1") ![image2](https://example.com/image2 "title2")`,
			want: wantStruct{
				clearContent: `image1 title1 image2 title2`,
				imageUrls:    []string{"https://example.com/image1", "https://example.com/image2"},
			},
		},
	}

	for _, tt := range tests {
		imageUrls, cleanText := findImagesInMarkdownContent(tt.text)
		assert.EqualValues(t, tt.want.clearContent, cleanText, "cleanText should be equal")
		assert.EqualValues(t, tt.want.imageUrls, imageUrls, "imageUrls should be equal")
	}
}
