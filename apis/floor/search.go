package floor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	. "treehole_next/config"
	. "treehole_next/models"
	. "treehole_next/utils"
)

type SearchQuery struct {
	Search string `json:"search" query:"search" validate:"required"`
	Size   int    `json:"size" query:"size" validate:"min=0" default:"10"`
	Offset int    `json:"offset" query:"offset" validate:"min=0" default:"0"`
}

// SearchFloors
//
//	@Summary	SearchFloors In ElasticSearch
//	@Tags		Search
//	@Produce	application/json
//	@Router		/floors/search [get]
//	@Param		object	query	SearchQuery	true	"search_query"
//	@Success	200		{array}	models.Floor
func SearchFloors(c *fiber.Ctx) error {
	query, err := common.ValidateQuery[SearchQuery](c)
	if err != nil {
		return err
	}

	floors, err := Search(query.Search, query.Size, query.Offset)
	if err != nil {
		return err
	}

	return Serialize(c, floors)
}

// SearchConfig
//
//	@Summary	change search config
//	@Tags		Search
//	@Produce	application/json
//	@Router		/config/search [post]
//	@Param		json	body		SearchConfigModel	true	"json"
//	@Success	200		{object}	Map
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
		return common.Forbidden("树洞流量激增，搜索功能暂缓开放")
	}

	floors, err := Search(query.Search, query.Size, query.Offset)
	if err != nil {
		return err
	}

	return Serialize(c, &floors)
}
