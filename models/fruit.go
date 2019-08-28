package models

import (
	"context"
	"time"

	"github.com/go-xorm/xorm"

	"go-api-twodb/factory"
)

type Fruit struct {
	Id        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Price     int64     `json:"price"`
	StoreCode string    `json:"storeCode"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated"`
}

type Store struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type FruitStoreDto struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Price     int64  `json:"price"`
	StoreName string `json:"storeName" xorm:"store_name"` // note: xorm:"store_name" ==== b.name as store_name
}

func (d *Fruit) Create(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Insert(d)
}

func (Fruit) GetById(ctx context.Context, id int64) (bool, Fruit, error) {
	fruit := Fruit{}
	has, err := factory.DB(ctx).Where("id=?", id).Get(&fruit)
	return has, fruit, err
}

func (Fruit) GetByCode(ctx context.Context, code string) (bool, Fruit, error) {
	fruit := Fruit{}
	has, err := factory.DB(ctx).Where("code=?", code).Get(&fruit)
	return has, fruit, err
}

func (Fruit) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (int64, []Fruit, error) {
	queryBuilder := func() xorm.Interface {
		q := factory.DB(ctx)
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}
	var items []Fruit
	totalCount, err := queryBuilder().Limit(limit, offset).FindAndCount(&items)
	if err != nil {
		return totalCount, items, err
	}
	return totalCount, items, nil
}

func (d *Fruit) Update(ctx context.Context, id int64) (int64, error) {
	return factory.DB(ctx).Where("id=?", id).Update(d)

}

func (Fruit) Delete(ctx context.Context, id int64) (int64, error) {
	return factory.DB(ctx).Where("id=?", id).Delete(&Fruit{})
}

func (Fruit) GetWithStoreById(ctx context.Context, id int64) (bool, FruitStoreDto, error) {
	var dto FruitStoreDto
	has, err := factory.DB(ctx).Table("fruit").Alias("a").
		Join("inner", []string{"store", "b"}, "a.store_code = b.code").
		Select(`a.id,a.name,a.color,a.price,b.name as store_name`).
		Where("a.id=?", id).Get(&dto)
	return has, dto, err
}
