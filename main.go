package main

import (
	"fmt"
	"go-api-twodb/config"
	"go-api-twodb/controllers"
	"go-api-twodb/factory"
	"go-api-twodb/models"
	"log"
	"net/http"
	"nomni/utils/validator"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pangpanglabs/echoswagger"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/sirupsen/logrus"
)

func main() {
	c := config.Init(os.Getenv("APP_ENV"))

	fmt.Println("Config===", c)
	db, err := models.InitDB(c.Database.Fruit.Driver, c.Database.Fruit.Connection)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := models.InitTable(db); err != nil {
		panic(err)
	}

	db2, err := models.InitDB2(c.Database.Fruit2.Driver, c.Database.Fruit2.Connection)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := models.InitTable2(db2); err != nil {
		panic(err)
	}

	e := echo.New()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	r := echoswagger.New(e, "docs", &echoswagger.Info{
		Title:       "Sample Fruit API",
		Description: "This is docs for fruit service",
		Version:     "1.0.0",
	})
	r.AddSecurityAPIKey("Authorization", "JWT token", echoswagger.SecurityInHeader)
	r.SetUI(echoswagger.UISetting{
		HideTop: true,
	})
	controllers.FruitApiController{}.Init(r.Group("fruits", "v1/fruits"))

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(middleware.RequestID())
	e.Use(echomiddleware.ContextLogger())
	e.Use(echomiddleware.ContextDBWithName(c.ServiceName, factory.ContextDBName, db, c.Database.Logger.Kafka))
	e.Use(echomiddleware.ContextDB(c.ServiceName, db, c.Database.Logger.Kafka))
	e.Use(echomiddleware.BehaviorLogger(c.ServiceName, c.BehaviorLog.Kafka))
	if !strings.HasSuffix(c.Appenv, "production") {
		behaviorlog.SetLogLevel(logrus.InfoLevel)
	}

	e.Validator = validator.New()

	e.Debug = c.Debug

	if err := e.Start(":" + c.HttpPort); err != nil {
		log.Println(err)
	}
}
