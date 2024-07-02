package sensitive

import (
	"fmt"
	"github.com/yidun/yidun-golang-sdk/yidun/service/antispam/image/v5"
	"github.com/yidun/yidun-golang-sdk/yidun/service/antispam/image/v5/check"
	"net/url"
	"strconv"
	"strings"
	"time"
	"treehole_next/config"
	"treehole_next/utils"

	v5 "github.com/yidun/yidun-golang-sdk/yidun/service/antispam/text"
	"github.com/yidun/yidun-golang-sdk/yidun/service/antispam/text/v5/check/sync/single"
)

const (
	TypeHole  = "Hole"
	TypeFloor = "Floor"
	TypeTag   = "Tag"
	TypeImage = "Image"
)

var checkTypes = []string{TypeHole, TypeFloor, TypeTag, TypeImage}

type ParamsForCheck struct {
	Content  string
	Id       int64
	TypeName string
}

type ResponseForCheck struct {
	Pass   bool
	Labels []int
	Detail string
}

func CheckSensitive(params ParamsForCheck) (resp *ResponseForCheck, err error) {
	images, clearContent, err := findImagesInMarkdownContent(params.Content)
	if err != nil {
		return nil, err
	}
	if len(images) != 0 {
		for _, img := range images {
			imgUrl, err := url.Parse(img)
			if err != nil {
				return nil, err
			}
			if config.Config.ExternalImageHost != "" {
				imgUrl.Host = config.Config.ExternalImageHost
			}
			ret, err := checkSensitiveImage(ParamsForCheck{
				Content:  imgUrl.String(),
				Id:       time.Now().UnixNano(),
				TypeName: TypeImage,
			})
			if err != nil {
				return nil, err
			}
			if !ret.Pass {
				return ret, nil
			}
		}
	}

	contained, reason := containsUnsafeURL(clearContent)
	if contained {
		return &ResponseForCheck{
			Pass:   false,
			Labels: nil,
			Detail: "不允许使用外部链接" + reason,
		}, nil
	}
	params.Content = strings.TrimSpace(removeIDReprInContent(clearContent))
	if params.Content == "" {
		return &ResponseForCheck{
			Pass:   true,
			Labels: nil,
			Detail: "",
		}, nil
	}

	return CheckSensitiveText(params)
}

func CheckSensitiveText(params ParamsForCheck) (resp *ResponseForCheck, err error) {
	if !checkType(params) {
		return nil, fmt.Errorf("invalid type for sensitive check")
	}

	request := single.NewTextCheckRequest(config.Config.YiDunBusinessIdText)
	textCheckClient := v5.NewTextClientWithAccessKey(config.Config.YiDunSecretId, config.Config.YiDunSecretKey)

	request.SetDataID(strconv.FormatInt(params.Id, 10) + "_" + params.TypeName)
	request.SetContent(params.Content)
	request.SetTimestamp(time.Now().UnixMilli())

	response, err := textCheckClient.SyncCheckText(request)
	if err != nil {
		// 处理错误并打印日志
		utils.RequestLog(fmt.Sprintf("sync request error:%+v", err.Error()), params.TypeName, params.Id, false)
		return &ResponseForCheck{Pass: false}, nil
	}

	resp = &ResponseForCheck{}
	if response.GetCode() == 200 {

		if *response.Result.Antispam.Suggestion == 0 {
			utils.RequestLog("Sensitive text check response code is 200", params.TypeName, params.Id, true)
			resp.Pass = true
			return
		}

		utils.RequestLog("Sensitive text check response code is 200", params.TypeName, params.Id, false)
		resp.Pass = false
		var str string
		for _, label := range response.Result.Antispam.Labels {
			resp.Labels = append(resp.Labels, *label.Label)
			// response != nil && response.Result != nil && response.Result.Antispam != nil &&
			//if response.Result.Antispam.SecondLabel != nil && response.Result.Antispam.ThirdLabel != nil {
			//	str := *response.Result.Antispam.SecondLabel + " " + *response.Result.Antispam.ThirdLabel
			//}
			if label.SubLabels != nil {
				for _, subLabel := range label.SubLabels {
					if subLabel.Details != nil && subLabel.Details.HitInfos != nil {
						for _, hitInfo := range subLabel.Details.HitInfos {
							if str == "" {
								str = *hitInfo.Value
								continue
							}
							str += "\n" + *hitInfo.Value
						}
					}
				}
			}
		}
		if str == "" {
			str = "文本敏感，未知原因"
		}
		resp.Detail = str
		return
	}

	utils.RequestLog("Sensitive text check http response code is not 200", params.TypeName, params.Id, false)
	resp.Pass = false
	return
}

func checkSensitiveImage(params ParamsForCheck) (resp *ResponseForCheck, err error) {
	// 设置易盾内容安全分配的businessId
	imgUrl := params.Content

	request := check.NewImageV5CheckRequest(config.Config.YiDunBusinessIdImage)

	// 实例化一个textClient，入参需要传入易盾内容安全分配的secretId，secretKey
	imageCheckClient := image.NewImageClientWithAccessKey(config.Config.YiDunSecretId, config.Config.YiDunSecretKey)

	imageInst := check.NewImageBeanRequest()
	imageInst.SetData(imgUrl)
	imageInst.SetName(strconv.FormatInt(params.Id, 10) + "_" + params.TypeName)
	// 设置图片数据的类型，1：图片URL，2:图片BASE64
	imageInst.SetType(1)

	imageBeans := []check.ImageBeanRequest{*imageInst}
	request.SetImages(imageBeans)

	response, err := imageCheckClient.ImageSyncCheck(request)
	if err != nil {
		// 处理错误并打印日志
		utils.RequestLog(fmt.Sprintf("sync request error:%+v", err.Error()), params.TypeName, params.Id, false)
		// TODO: 通知管理员
		return &ResponseForCheck{Pass: false}, nil
	}

	resp = &ResponseForCheck{}
	if response.GetCode() == 200 {
		if len(*response.Result) == 0 {
			return nil, fmt.Errorf("sensitive image check returns empty response")
		}

		if *((*response.Result)[0].Antispam.Suggestion) == 0 {
			utils.RequestLog("Sensitive image check response code is 200", params.TypeName, params.Id, true)
			resp.Pass = true
			return
		}

		utils.RequestLog("Sensitive image check response code is 200", params.TypeName, params.Id, false)
		resp.Pass = false
		for _, label := range *((*response.Result)[0].Antispam.Labels) {
			resp.Labels = append(resp.Labels, *label.Label)
		}
		var str string
		for _, result := range *response.Result {
			if result.Ocr != nil {
				if result.Ocr.Details != nil {
					for _, detail := range *result.Ocr.Details {
						if str == "" {
							str = *detail.Content
							continue
						}
						str += "\n" + *detail.Content
					}
				}
			}
			if result.Face != nil {
				if result.Face.Details != nil {
					for _, detail := range *result.Face.Details {
						if detail.FaceContents != nil {
							for _, faceContent := range *detail.FaceContents {
								if str == "" {
									str = *faceContent.Name
									continue
								}
								str += "\n" + *faceContent.Name
							}
						}
					}
				}
			}
		}
		if str == "" {
			str = "图片敏感，未知原因"
		}
		resp.Detail = str
		return
	}

	utils.RequestLog("Sensitive image check http response code is not 200", params.TypeName, params.Id, false)
	resp.Pass = false
	return
}
