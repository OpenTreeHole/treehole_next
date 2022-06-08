package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"treehole_next/config"
	"treehole_next/middlewares"
	. "treehole_next/models"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

func TestDivision(t *testing.T) {

	config.InitConfig()
	InitDB()

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.MyErrorHandler,
	})
	middlewares.RegisterMiddlewares(app)
	app.Post("/divisions", AddDivision)
	app.Get("/divisions", ListDivisions)
	app.Get("/divisions/:id", GetDivision)

	divisions := []Division{
		{
			BaseModel: BaseModel{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:        "dqd",
			Description: "123",
			Pinned:      IntArray{1, 3, 5},
		},
		{
			BaseModel: BaseModel{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:        "vvs",
			Description: "svsv",
			Pinned:      IntArray{2, 3, 5},
		},
		{
			BaseModel: BaseModel{
				ID:        2,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:        "dqdtt",
			Description: "335",
			Pinned:      IntArray{5, 7, 5},
		},
	}
	route := "/divisions"
	tests := []struct {
		name         string
		method       string
		id           int
		expectedCode int
	}{
		{
			name:         "add dqd",
			method:       fiber.MethodPost,
			expectedCode: 200,
		},
		{
			name:         "add vvs",
			method:       fiber.MethodPost,
			expectedCode: 201,
		},
		{
			name:         "add dqdtt",
			method:       fiber.MethodPost,
			expectedCode: 200,
		},
		{
			name:         "get id1",
			method:       fiber.MethodGet,
			id:           1,
			expectedCode: 200,
		},
		{
			name:         "get id2",
			method:       fiber.MethodGet,
			id:           2,
			expectedCode: 404,
		},
		{
			name:         "get id3",
			method:       fiber.MethodGet,
			id:           3,
			expectedCode: 200,
		},
	}
	for i, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			var req *http.Request
			if test.method == fiber.MethodPost {
				requestBody, err := json.Marshal(divisions[i])
				if err != nil {
					fmt.Printf("err: %v\n", err)
					return
				}
				req = httptest.NewRequest(
					http.MethodPost,
					route,
					bytes.NewReader(requestBody),
				)
				requestString := string(requestBody)
				fmt.Printf("requestString: %v\n", requestString)
			} else if test.method == fiber.MethodGet {
				req = httptest.NewRequest(
					http.MethodGet,
					route,
					nil,
				)
			}

			req.Header.Add("Content-Type", fiber.MIMEApplicationJSON)
			resp, _ := app.Test(req, 1000)
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("respbody: %s\n", string(body))
		})
	}
	DB.Exec("DELETE from division")
	DB.Find(&divisions)
	if len(divisions) != 0 {
		fmt.Println(json.Marshal(divisions[0]))
	}else {
		fmt.Println("Data cleaned...")
	}
}
