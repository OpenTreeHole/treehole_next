package sensitive

import (
	"fmt"
	"strconv"
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
)

type ParamsForCheck struct {
	Content  string
	Id       int
	TypeName string
}

type ResponseForCheck struct {
	Pass   bool
	Labels []int
}

func CheckSensitive(params ParamsForCheck) (resp *ResponseForCheck, err error) {
	request := single.NewTextCheckRequest(config.Config.YiDunBusinessIdText)
	textCheckClient := v5.NewTextClientWithAccessKey(config.Config.YiDunSecretId, config.Config.YiDunSecretKey)

	request.SetDataID(strconv.Itoa(params.Id) + "_" + params.TypeName)
	request.SetContent(params.Content)
	request.SetTimestamp(time.Now().UnixMilli())

	response, err := textCheckClient.SyncCheckText(request)
	if err != nil {
		// 处理错误并打印日志
		utils.RequestLog(fmt.Sprintf("sync request error:%+v", err.Error()), params, false)
		resp = nil
	}

	resp = &ResponseForCheck{}
	if response.GetCode() == 200 {

		if *response.Result.Antispam.Suggestion == 0 {
			utils.RequestLog("Sensitive Check response code is 200", params, true)
			resp.Pass = true
			return
		}

		utils.RequestLog("Sensitive Check response code is 200", params, false)
		resp.Pass = false
		for _, label := range response.Result.Antispam.Labels {
			resp.Labels = append(resp.Labels, *label.Label)
		}
		return
	}

	utils.RequestLog("http response code is not 200", params, false)
	resp.Pass = false
	return
}
