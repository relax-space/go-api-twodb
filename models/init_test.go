package models_test

import (
	"context"
	"go-api-twodb/config"
	"go-api-twodb/factory"
	"go-api-twodb/models"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

var ctx context.Context

func TestMain(m *testing.M) {
	db := enterTest()
	code := m.Run()
	exitTest(db)
	os.Exit(code)
}

func enterTest() *xorm.Engine {
	configutil.SetConfigPath("../")
	c := config.Init(os.Getenv("APP_ENV"))
	xormEngine, err := xorm.NewEngine(c.Database.Fruit.Driver, c.Database.Fruit.Connection)
	if err != nil {
		panic(err)
	}
	// xormEngine.ShowSQL(true)
	if err = models.DropTables(xormEngine); err != nil {
		panic(err)
	}
	if err = models.InitTable(xormEngine); err != nil {
		panic(err)
	}
	xormEngine2, err := xorm.NewEngine(c.Database.Fruit2.Driver, c.Database.Fruit2.Connection)
	if err != nil {
		panic(err)
	}
	// xormEngine2.ShowSQL(true)
	if err = models.DropTables2(xormEngine2); err != nil {
		panic(err)
	}
	if err = models.InitTable2(xormEngine2); err != nil {
		panic(err)
	}

	ctx = context.WithValue(context.Background(), echomiddleware.ContextDBName, xormEngine)
	ctx = context.WithValue(ctx, factory.ContextDBName, xormEngine2)
	return xormEngine
}

func exitTest(db *xorm.Engine) {
	//db.Close()
}
