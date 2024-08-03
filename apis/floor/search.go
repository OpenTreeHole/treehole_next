package floor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"

	. "treehole_next/config"
	. "treehole_next/models"
	. "treehole_next/utils"
)

type SearchQuery struct {
	Search   string `json:"search" query:"search" validate:"required"`
	Size     int    `json:"size" query:"size" validate:"min=0" default:"10"`
	Offset   int    `json:"offset" query:"offset" validate:"min=0" default:"0"`
	Accurate bool   `json:"accurate" query:"accurate" default:"false"`
}

// SearchFloors
//
// @Summary SearchFloors In ElasticSearch
// @Tags Search
// @Produce application/json
// @Router /floors/search [get]
// @Router /floors/search [post]
// @Param object query SearchQuery true "search_query"
// @Success 200 {array} models.Floor
func SearchFloors(c *fiber.Ctx) error {
	var query SearchQuery
	err := common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	floors, err := Search(c, query.Search, query.Size, query.Offset, query.Accurate)
	if err != nil {
		return err
	}

	return Serialize(c, floors)
}

// SearchConfig
//
// @Summary change search config
// @Tags Search
// @Produce application/json
// @Router /config/search [post]
// @Param json body SearchConfigModel true "json"
// @Success 200 {object} Map
func SearchConfig(c *fiber.Ctx) error {
	var body SearchConfigModel
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	user, err := GetUser(c)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return common.Forbidden()
	}
	if DynamicConfig.OpenSearch.Load() == body.Open {
		return c.Status(200).JSON(Map{"message": "已经被修改"})
	} else {
		DynamicConfig.OpenSearch.Store(body.Open)
		return c.Status(201).JSON(Map{"message": "修改成功"})
	}
}

func SearchFloorsOld(c *fiber.Ctx, query *ListOldModel) error {
	if DynamicConfig.OpenSearch.Load() == false {
		return common.Forbidden("茶楼流量激增，搜索功能暂缓开放")
	}

	floors, err := Search(c, query.Search, query.Size, query.Offset, false)
	if err != nil {
		return err
	}

	return Serialize(c, &floors)
}
