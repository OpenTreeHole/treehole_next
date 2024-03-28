package sensitive

import (
	"golang.org/x/exp/slices"
	"mvdan.cc/xurls/v2"
	url2 "net/url"
	"regexp"
	"treehole_next/config"
)

// findImagesInMarkdown 从Markdown文本中查找所有图片链接
func findImagesInMarkdown(markdown string) []string {
	imageRegex := regexp.MustCompile(
		`!\[.*?\]\(([^" )]*)`,
	)

	matches := imageRegex.FindAllStringSubmatch(markdown, -1)
	images := make([]string, 0, len(matches))
	for _, match := range matches {
		images = append(images, match[1])
	}
	return images
}

func detect(markdownText string) []string {

	var ret []string
	images := findImagesInMarkdown(markdownText)
	for _, image := range images {
		ret = append(ret, image)
	}

	return ret
}

func checkType(params ParamsForCheck) bool {
	if params.TypeName != TypeTag && params.TypeName != TypeHole && params.TypeName != TypeFloor {
		return true
	}
	return false
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
	url, err := url2.Parse(input)
	if err != nil {
		return false, err
	}
	if slices.Contains(config.Config.ValidImageUrl, url.Hostname()) {
		return false, nil
	}
	return true, nil
}

func deleteImagesInMarkdown(markdown string) string {
	imageRegex := regexp.MustCompile(
		`!\[(.*?)\]\(([^" ]*)( ".*")?\)`,
	)
	return imageRegex.ReplaceAllStringFunc(markdown, func(s string) string {
		submatches := imageRegex.FindStringSubmatch(s)
		altText := submatches[1]
		if len(submatches) > 3 && submatches[3] != "" {
			// If there is a title, return it along with the alt text
			return altText + " " + submatches[3][2:len(submatches[3])-1]
		}
		// If there is no title, return the alt text
		return altText
	})
}
