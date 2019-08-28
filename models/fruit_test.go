package models_test

import (
	"go-api-twodb/models"
	"testing"

	"github.com/pangpanglabs/goutils/test"
)

func TestFruitCRUD(t *testing.T) {
	var ID int64
	t.Run("Create", func(t *testing.T) {
		code := "apple"
		f := &models.Fruit{
			Code: code,
		}
		affectedRow, err := f.Create(ctx)
		test.Ok(t, err)
		test.Equals(t, int64(1), affectedRow)
		ID = f.Id
		test.Assert(t, ID != int64(0), "create failure")
	})
	t.Run("Get", func(t *testing.T) {
		has, v, err := models.Fruit{}.GetById(ctx, ID)
		test.Ok(t, err)
		test.Equals(t, true, has)
		test.Equals(t, "apple", v.Code)
	})
	t.Run("Update", func(t *testing.T) {
		var price int64 = 10
		f := &models.Fruit{
			Price: price,
		}
		affectedRow, err := f.Update(ctx, ID)
		test.Ok(t, err)
		test.Equals(t, int64(1), affectedRow)

		has, v, err := models.Fruit{}.GetById(ctx, ID)
		test.Ok(t, err)
		test.Equals(t, true, has)
		test.Equals(t, int64(10), v.Price)

	})
	t.Run("Delete", func(t *testing.T) {
		affectedRow, err := models.Fruit{}.Delete(ctx, ID)
		test.Ok(t, err)
		test.Equals(t, int64(1), affectedRow)
		has, v, err := models.Fruit{}.GetById(ctx, ID)
		test.Ok(t, err)
		test.Equals(t, false, has)
		test.Equals(t, int64(0), v.Id)
	})
}
