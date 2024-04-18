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
		imageLink := convertImageURLForModeration(submatch[2])
		imageUrls = append(imageUrls, submatch[2])
		if len(submatch) > 3 && submatch[3] != "" {
			// If there is a title, return it along with the alt text
			title := strings.Trim(submatch[3], "\"")
			return altText + imageLink + title
		}
		// If there is no title, return the alt text
		return altText + imageLink
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

func convertImageURLForModeration(imageLink string) string {
	url, err := imageUrl.Parse(imageLink)
	if err != nil {
		return ""
	}
	if url.Scheme == "" && url.Host == "" {
		return imageLink
	}
	return " "
}

func checkValidUrl(input string) (bool, error) {
	url, err := imageUrl.Parse(input)
	if err != nil {
		return false, err
	}
	// if the url is a sticker, skip check
	if url.Scheme == "" && url.Host == "" {
		return true, nil
	}
	if !slices.Contains(config.Config.ValidImageUrl, url.Hostname()) {
		return false, nil
	}
	return true, nil
}

var reHole = regexp.MustCompile(`[^#]#(\d+)`)
var reFloor = regexp.MustCompile(`##(\d+)`)

func removeIDReprInContent(content string) string {
	content = " " + content
	content = reHole.ReplaceAllString(content, "")
	content = reFloor.ReplaceAllString(content, "")
	return strings.TrimSpace(content)
}
