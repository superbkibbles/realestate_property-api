package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/realestate_property-api/src/clients/elasticsearch"
	"github.com/superbkibbles/realestate_property-api/src/http"
	"github.com/superbkibbles/realestate_property-api/src/repository/db"
	"github.com/superbkibbles/realestate_property-api/src/services/property"
)

var (
	router  = gin.Default()
	handler http.Propertyhandler
)

func StartApplication() {
	elasticsearch.Client.Init()
	handler = http.NewPropertyHandler(property.NewService(db.NewRepository()))
	router.Use(cors.Default())
	mapURLS()
	router.Static("assets", "clients/visuals")
	router.Run(":3030")
}
