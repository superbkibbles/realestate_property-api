package app

import (
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

	mapURLS()
	router.Run(":3030")
}
