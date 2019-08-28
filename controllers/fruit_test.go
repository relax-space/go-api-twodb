package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-api-twodb/controllers"
	"go-api-twodb/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/test"
)

func TestFruitCRUD(t *testing.T) {
	inputs := []models.Fruit{
		models.Fruit{
			Code:  "fruit#1",
			Color: "red",
		},
		models.Fruit{
			Code:  "fruit#2",
			Color: "green",
		},
	}
	var fruits []models.Fruit

	for i, p := range inputs {
		pb, _ := json.Marshal(p)
		t.Run(fmt.Sprint("Create#", i+1), func(t *testing.T) {
			req := httptest.NewRequest(echo.POST, "/v1/fruits", bytes.NewReader(pb))
			setHeader(req)
			rec := httptest.NewRecorder()
			test.Ok(t, handleWithFilter(controllers.FruitApiController{}.Create, echoApp.NewContext(req, rec)))
			test.Equals(t, http.StatusCreated, rec.Code)
		})
	}

	t.Run("GetAll", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/products", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(controllers.FruitApiController{}.GetAll, echoApp.NewContext(req, rec)))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result struct {
				TotalCount int            `json:"totalCount"`
				Items      []models.Fruit `json:"items"`
			} `json:"result"`
			Success bool `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result.TotalCount, 2)
		test.Equals(t, v.Result.Items[0].Code, "fruit#1")
		test.Equals(t, v.Result.Items[1].Code, "fruit#2")
		fruits = v.Result.Items
	})

	var fruit models.Fruit

	t.Run("GetOne", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/fruits/1", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetPath("/v1/fruit/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		test.Ok(t, handleWithFilter(controllers.FruitApiController{}.GetOne, c))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result.Code, "fruit#1")
		fruit = v.Result
	})

	t.Run("Update", func(t *testing.T) {
		fruit.Code = "fruit#updated"
		pb, _ := json.Marshal(fruit)
		req := httptest.NewRequest(echo.PUT, "/", bytes.NewReader(pb))
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetPath("/v1/fruits/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%v", fruit.Id))
		test.Ok(t, handleWithFilter(controllers.FruitApiController{}.Update, c))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result  models.Fruit `json:"result"`
			Success bool         `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result.Code, "fruit#updated")
	})
	for i, p := range fruits {
		t.Run(fmt.Sprint("Delete#", i+1), func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", nil)
			setHeader(req)
			rec := httptest.NewRecorder()
			c := echoApp.NewContext(req, rec)
			c.SetPath("/v1/fruits/:id")
			c.SetParamNames("id")
			c.SetParamValues(fmt.Sprintf("%v", p.Id))
			test.Ok(t, handleWithFilter(controllers.FruitApiController{}.Delete, c))
			test.Equals(t, http.StatusOK, rec.Code)
		})
	}
}
