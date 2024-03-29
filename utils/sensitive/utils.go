package sensitive

import (
	"golang.org/x/exp/slices"
	"mvdan.cc/xurls/v2"
	imageUrl "net/url"
	"regexp"
	"strings"
	"treehole_next/config"
)

var imageRegex = regexp.MustCompile(
	`!\[(.*?)]\(([^" )]*?)\s*(".*?")?\)`,
)

// findImagesInMarkdown 从Markdown文本中查找所有图片链接，并且返回清除链接之后的文本
func findImagesInMarkdownContent(content string) (imageUrls []string, clearContent string) {

	clearContent = imageRegex.ReplaceAllStringFunc(content, func(s string) string {
		submatch := imageRegex.FindStringSubmatch(s)
		altText := submatch[1]
		imageUrls = append(imageUrls, submatch[2])
		if len(submatch) > 3 && submatch[3] != "" {
			// If there is a title, return it along with the alt text
			title := strings.Trim(submatch[3], "\"")
			return altText + " " + title
		}
		// If there is no title, return the alt text
		return altText
	})
	return
}

func checkType(params ParamsForCheck) bool {
	return slices.Contains(checkTypes, params.TypeName)
}

func hasTextUrl(content string) bool {
	xurlsRelaxed := xurls.Relaxed()
	output := xurlsRelaxed.FindAllString(content, -1)
	if len(output) == 0 {
		return false
	}
	return true
}

func checkValidUrl(input string) (bool, error) {
	url, err := imageUrl.Parse(input)
	if err != nil {
		return false, err
	}
	if !slices.Contains(config.Config.ValidImageUrl, url.Hostname()) {
		return false, nil
	}
	return true, nil
}
