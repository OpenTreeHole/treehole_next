package sensitive

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"treehole_next/config"
)

func TestFindImagesInMarkdown(t *testing.T) {
	config.Config.ValidImageUrl = []string{"example.com"}

	type wantStruct struct {
		clearContent string
		imageUrls    []string
		err          error
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
				err:          nil,
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
		{
			text: `![image1](123) ![image2](456)`,
			want: wantStruct{
				clearContent: `image1 123 image2 456`,
				imageUrls:    nil,
			},
		},
		{
			text: `![](123) ![](456)`,
			want: wantStruct{
				clearContent: `123 456`,
				imageUrls:    nil,
			},
		},
		{
			text: "![](https://example2.com/image1)",
			want: wantStruct{
				clearContent: "",
				imageUrls:    nil,
				err:          ErrInvalidImageHost,
			},
		},
	}

	for _, tt := range tests {
		imageUrls, cleanText, err := findImagesInMarkdownContent(tt.text)
		assert.EqualValues(t, tt.want.clearContent, cleanText, "cleanText should be equal")
		assert.EqualValues(t, tt.want.imageUrls, imageUrls, "imageUrls should be equal")
		assert.EqualValues(t, tt.want.err, err, "err should be equal")
	}
}

func TestCheckValidUrl(t *testing.T) {
	config.Config.ValidImageUrl = []string{"example.com"}
	type wantStruct struct {
		err error
	}
	tests := []struct {
		url  string
		want wantStruct
	}{
		{
			url: "https://example.com/image1",
			want: wantStruct{
				err: nil,
			},
		},
		{
			url: "https://example.com/image2",
			want: wantStruct{
				err: nil,
			},
		},
		{
			url: "123456",
			want: wantStruct{
				err: ErrImageLinkTextOnly,
			},
		},
		{
			url: "https://example2.com",
			want: wantStruct{
				err: ErrInvalidImageHost,
			},
		},
	}

	for _, tt := range tests {
		err := checkValidUrl(tt.url)
		assert.EqualValues(t, tt.want.err, err, "err should be equal")
	}
}
