package sensitive

import (
	"fmt"
	"github.com/opentreehole/go-common"
	"github.com/yidun/yidun-golang-sdk/yidun/service/antispam/image/v5"
	"github.com/yidun/yidun-golang-sdk/yidun/service/antispam/image/v5/check"
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
}

func CheckSensitive(params ParamsForCheck) (resp *ResponseForCheck, err error) {
	images, clearContent := findImagesInMarkdownContent(params.Content)
	if len(images) != 0 {
		for _, img := range images {
			pass, err := checkValidUrl(img)
			if err != nil {
				return nil, err
			}
			if !pass {
				return nil, common.BadRequest("不允许使用外部图片链接")
			}

			ret, err := checkSensitiveImage(ParamsForCheck{
				Content:  img,
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

	if hasTextUrl(clearContent) {
		return &ResponseForCheck{
			Pass:   false,
			Labels: nil,
		}, nil
	}
	params.Content = strings.TrimSpace(removeIDReprInContent(clearContent))
	if params.Content == "" {
		return &ResponseForCheck{
			Pass:   true,
			Labels: nil,
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
		for _, label := range response.Result.Antispam.Labels {
			resp.Labels = append(resp.Labels, *label.Label)
		}
		return
	}

	utils.RequestLog("Sensitive text check http response code is not 200", params.TypeName, params.Id, false)
	resp.Pass = false
	return
}

func checkSensitiveImage(params ParamsForCheck) (resp *ResponseForCheck, err error) {
	// 设置易盾内容安全分配的businessId
	url := params.Content

	request := check.NewImageV5CheckRequest(config.Config.YiDunBusinessIdImage)

	// 实例化一个textClient，入参需要传入易盾内容安全分配的secretId，secretKey
	imageCheckClient := image.NewImageClientWithAccessKey(config.Config.YiDunSecretId, config.Config.YiDunSecretKey)

	imageInst := check.NewImageBeanRequest()
	imageInst.SetData(url)
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
		return
	}

	utils.RequestLog("Sensitive image check http response code is not 200", params.TypeName, params.Id, false)
	resp.Pass = false
	return
}
