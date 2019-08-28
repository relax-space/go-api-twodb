package controllers

import (
	"fmt"
	"go-api-twodb/models"
	"net/http"
	"nomni/utils/api"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type FruitApiController struct {
}

// localhost:8080/docs
func (d FruitApiController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")
	g.GET("", d.GetAll).
		AddParamQueryNested(SearchInput{})
	g.GET("/:id", d.GetOne).
		AddParamPath("", "id", "id").AddParamQuery("", "with_store", "with_store", false)
	g.PUT("/:id", d.Update).
		AddParamPath("", "id", "id").
		AddParamBody(models.Fruit{}, "fruit", "only can modify name,color,price", true)
	g.POST("", d.Create).
		AddParamBody(models.Fruit{}, "fruit", "new fruit", true)
	g.DELETE("/:id", d.Delete).
		AddParamPath("", "id", "id")
}

/*
localhost:8080/fruits
localhost:8080/fruits?name=apple
localhost:8080/fruits?skipCount=0&maxResultCount=2
localhost:8080/fruits?skipCount=0&maxResultCount=2&sortby=store_code&order=desc
*/
func (FruitApiController) GetAll(c echo.Context) error {
	var v SearchInput
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	if v.MaxResultCount == 0 {
		v.MaxResultCount = DefaultMaxResultCount
	}
	totalCount, items, err := models.Fruit{}.GetAll(c.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, api.NotFoundError(), err)
	}
	if len(items) == 0 {
		return ReturnApiFail(c, http.StatusBadRequest, api.NotFoundError())
	}
	return ReturnApiSucc(c, http.StatusOK, api.ArrayResult{
		TotalCount: totalCount,
		Items:      items,
	})
}

/*
localhost:8080/fruits/1?with_store=true
localhost:8080/fruits/1
*/
func (d FruitApiController) GetOne(c echo.Context) error {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.InvalidParamError("id", c.Param("id"), err))
	}

	var withStore bool
	if len(c.QueryParam("with_store")) != 0 {
		withStore, err = strconv.ParseBool(c.QueryParam("with_store"))
		if err != nil {
			return ReturnApiFail(c, http.StatusBadRequest, api.InvalidParamError("with_store", c.Param("with_store"), err))
		}
	}
	if withStore == true {
		has, fruit, err := models.Fruit{}.GetWithStoreById(c.Request().Context(), id)
		if err != nil {
			return ReturnApiFail(c, http.StatusInternalServerError, api.NotFoundError(), err)
		}
		if has == false {
			param := fmt.Sprintf("?id=%v&with_store=true", id)
			return ReturnApiFail(c, http.StatusBadRequest, api.NotFoundError(), param)
		}
		return ReturnApiSucc(c, http.StatusOK, fruit)
	}

	has, fruit, err := models.Fruit{}.GetById(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, api.NotFoundError(), err)
	}
	if !has {
		param := fmt.Sprintf("?id=%v", id)
		return ReturnApiFail(c, http.StatusBadRequest, api.NotFoundError(), param)
	}
	return ReturnApiSucc(c, http.StatusOK, fruit)
}

/*
localhost:8080/fruits
 {
        "code": "AA01",
        "name": "Apple",
        "color": "",
        "price": 2,
        "store_code": ""
    }
*/
func (d FruitApiController) Create(c echo.Context) error {
	var v models.Fruit
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	if err := c.Validate(v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	has, _, err := models.Fruit{}.GetByCode(c.Request().Context(), v.Code)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, api.NotCreatedError(), err)
	}
	if has {
		return ReturnApiFail(c, http.StatusBadRequest, api.NotCreatedError(), "code has exist")
	}
	affectedRow, err := v.Create(c.Request().Context())
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, api.NotCreatedError(), err)
	}
	if affectedRow == int64(0) {
		return ReturnApiFail(c, http.StatusBadRequest, api.NotCreatedError())
	}
	return ReturnApiSucc(c, http.StatusCreated, v)
}

/*
localhost:8080/fruits
 {
        "price": 21,
    }
*/
func (d FruitApiController) Update(c echo.Context) error {
	var v models.Fruit
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	if err := c.Validate(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.ParameterParsingError(err))
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.InvalidParamError("id", c.Param("id"), err))
	}
	has, _, err := models.Fruit{}.GetById(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, api.NotUpdatedError(), err)

	}
	if has == false {
		return ReturnApiFail(c, http.StatusBadRequest, api.NotUpdatedError(), "id has not found")
	}
	affectedRow, err := v.Update(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, api.NotUpdatedError(), err)
	}
	if affectedRow == int64(0) {
		return ReturnApiFail(c, http.StatusBadRequest, api.NotUpdatedError())
	}
	return ReturnApiSucc(c, http.StatusOK, v)
}

/*
localhost:8080/fruits/45
*/
func (d FruitApiController) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, api.InvalidParamError("id", c.Param("id"), err))
	}
	has, v, err := models.Fruit{}.GetById(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, api.NotDeletedError(), err)
	}
	if has == false {
		return ReturnApiFail(c, http.StatusBadRequest, api.NotDeletedError(), "id has not found")
	}
	affectedRow, err := models.Fruit{}.Delete(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, api.NotDeletedError(), err)
	}
	if affectedRow == int64(0) {
		return ReturnApiFail(c, http.StatusBadRequest, api.NotDeletedError())
	}
	return ReturnApiSucc(c, http.StatusOK, v)
}
