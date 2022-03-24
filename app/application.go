package app

import (
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/realestate_property-api/clients/elasticsearch"
	"github.com/superbkibbles/realestate_property-api/constants"
	"github.com/superbkibbles/realestate_property-api/http"
	cloudstorage "github.com/superbkibbles/realestate_property-api/repository/cloudStorage"
	"github.com/superbkibbles/realestate_property-api/repository/db"
	"github.com/superbkibbles/realestate_property-api/services/property"
)

var (
	router  = gin.Default()
	handler http.Propertyhandler
)

func StartApplication() {
	elasticsearch.Client.Init()
	cld, err := cloudinary.NewFromParams(os.Getenv(constants.CLOUD_STORAGE_NAME), os.Getenv(constants.CLOUD_STORAGE_API_KEY), os.Getenv(constants.CLOUD_STORAGE_API_SECRET))
	if err != nil {
		panic(err)
	}

	handler = http.NewPropertyHandler(property.NewService(db.NewRepository(), cloudstorage.NewRepository(cld)))
	router.Use(cors.Default())
	mapURLS()
	// router.Static("assets", "clients/visuals")
	router.Run(os.Getenv(constants.PORT))
}
