package floor

import (
	"bytes"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"log"
	. "treehole_next/config"
	. "treehole_next/models"
	. "treehole_next/utils"
)

var ES *elasticsearch.Client

func InitSearch() {
	if Config.Mode == "test" || Config.Mode == "bench" || Config.ElasticsearchUrl == "" {
		return
	}

	// export ELASTICSEARCH_URL environment variable to set the ElasticSearch URL
	// example: http://user:pass@127.0.0.1:9200
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{Config.ElasticsearchUrl},
	})
	if err != nil {
		panic(err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
	ES = es
}

type SearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
		Skipped    int `json:"skipped"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		Hits []struct {
			Index  string              `json:"_index"`
			ID     string              `json:"_id"`
			Score  float64             `json:"_score"`
			Source SearchFloorResponse `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type SearchFloorResponse struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

// SearchFloors
//
//	@Summary	SearchFloors In ElasticSearch
//	@Tags		Search
//	@Produce	application/json
//	@Router		/floors/search [post]
//	@Param		json	body	any	true	"json"
//	@Success	200		{array}	models.Floor
func SearchFloors(c *fiber.Ctx) error {
	// forwarding
	var body bytes.Buffer
	body.Write(c.Body())
	return search(c, body)
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
		return Forbidden()
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
		return Forbidden("树洞流量激增，搜索功能暂缓开放")
	}
	floors := Floors{}
	result := DB.
		Where("content like ?", "%"+query.Search+"%").
		Where("hole_id in (?)", DB.Table("hole").Select("id").Where("hidden = false")).
		Offset(query.Offset).Limit(query.Size).Order("id desc").
		Preload("Mention").Find(&floors)
	if result.Error != nil {
		return result.Error
	}
	return Serialize(c, &floors)
}

func search(c *fiber.Ctx, body bytes.Buffer) error {
	res, err := ES.Search(
		ES.Search.WithIndex("floor"),
		ES.Search.WithBody(&body),
	)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	if res.IsError() {
		e := Map{}
		err = json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return err
		} else {
			return c.Status(502).JSON(&e)
		}
	}

	var response SearchResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}

	floorIDs := make([]int, len(response.Hits.Hits))
	for i, hit := range response.Hits.Hits {
		floorIDs[i] = hit.Source.ID
	}
	Logger.Debug("search response", zap.Ints("floorIDs", floorIDs))

	// get floors
	floors := Floors{}
	if len(floorIDs) > 0 {

		result := DB.Preload("Mention").Find(&floors, floorIDs)
		if result.Error != nil {
			return result.Error
		}

		// order
		floors = OrderInGivenOrder(floors, floorIDs)
	}

	return Serialize(c, &floors)
}
