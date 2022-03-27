package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/bookstore_utils-go/logger"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	domainProperty "github.com/superbkibbles/realestate_property-api/domain/property"
	"github.com/superbkibbles/realestate_property-api/domain/query"
	"github.com/superbkibbles/realestate_property-api/services/property"
)

type Propertyhandler interface {
	Create(*gin.Context)
	Get(*gin.Context)
	GetByID(*gin.Context)
	Search(*gin.Context)
	Update(*gin.Context)
	UploadMedia(*gin.Context)
	DeleteMedia(*gin.Context)
	UploadPropertyPic(c *gin.Context)
	GetActive(*gin.Context)
	GetDeactive(*gin.Context)
	Translate(*gin.Context)
	GetTranslated(*gin.Context)
}

type propertyHandler struct {
	service property.Service
}

func NewPropertyHandler(serv property.Service) Propertyhandler {
	return &propertyHandler{
		service: serv,
	}
}

func (ph *propertyHandler) UploadMedia(c *gin.Context) {
	propertyID := strings.TrimSpace(c.Param("id"))

	form, err := c.MultipartForm()
	if err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid JSON Body")
		c.JSON(restErr.Status(), restErr)
		return
	}
	files := form.File["files"]
	if len(files) < 1 {
		restErr := rest_errors.NewBadRequestErr("Invalid JSON Body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	if err := ph.service.UploadMedia(files, propertyID); err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.String(200, "uploaded")
}

func (ph *propertyHandler) Update(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	var updateRequest domainProperty.EsUpdate

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}

	property, err := ph.service.Update(id, updateRequest)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, property)
}

func (ph *propertyHandler) Translate(c *gin.Context) {
	local := c.GetHeader("local")
	id := strings.TrimSpace(c.Param("id"))
	var translateProperty domainProperty.TranslateProperty

	if err := c.ShouldBindJSON(&translateProperty); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}
	property, err := ph.service.Translate(id, translateProperty, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, property)
}

func (ph *propertyHandler) Create(c *gin.Context) {
	var property domainProperty.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		logger.Error("error when trying to fetch the body", err)
		restErr := rest_errors.NewBadRequestErr("Invalid JSON body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	newProperty, resultErr := ph.service.Create(property)
	if resultErr != nil {
		logger.Error("error when trying to create service property", nil)
		c.JSON(resultErr.Status(), resultErr)
		return
	}

	c.JSON(http.StatusCreated, newProperty)
}

func (ph *propertyHandler) Get(c *gin.Context) {
	sort := c.Query("sort")
	asc := c.Query("asc") == "true"
	local := c.GetHeader("local")
	properties, err := ph.service.Get(sort, asc, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, properties)
}

func (ph *propertyHandler) GetActive(c *gin.Context) {
	sort := c.Query("sort")
	asc := c.Query("asc") == "true"
	local := c.GetHeader("local")

	p, err := ph.service.GetActive(sort, asc, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, p)
}

func (ph *propertyHandler) GetDeactive(c *gin.Context) {
	sort := c.Query("sort")
	asc := c.Query("asc") == "true"
	local := c.GetHeader("local")

	p, err := ph.service.GetDeactive(sort, asc, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, p)
}

func (ph *propertyHandler) GetByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	local := c.GetHeader("local")
	if len(id) == 0 {
		restErr := rest_errors.NewBadRequestErr("Invalid ID")
		c.JSON(restErr.Status(), restErr)
		return
	}

	property, err := ph.service.GetByID(id, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, property)
}

func (ph *propertyHandler) GetTranslated(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	local := c.GetHeader("local")
	if len(id) == 0 {
		restErr := rest_errors.NewBadRequestErr("Invalid ID")
		c.JSON(restErr.Status(), restErr)
		return
	}

	property, err := ph.service.GetTranslated(id, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, property)
}

func (ph *propertyHandler) Search(c *gin.Context) {
	var q query.EsQuery
	sort := c.Query("sort")
	asc := c.Query("asc") == "true"
	local := c.GetHeader("local")

	if err := c.ShouldBindJSON(&q); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid Body JSON")
		c.JSON(restErr.Status(), restErr)
		return
	}

	properties, err := ph.service.Search(q, sort, asc, local)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, properties)
}

func (ph *propertyHandler) DeleteMedia(c *gin.Context) {
	propertyID := strings.TrimSpace(c.Param("id"))
	mediaID := strings.TrimSpace(c.Param("media_id"))

	if err := ph.service.DeleteMedia(propertyID, mediaID); err != nil {
		c.JSON(err.Status(), err)
	}

	c.String(http.StatusOK, "Deleted")
}

func (ph *propertyHandler) UploadPropertyPic(c *gin.Context) {
	agencyID := strings.TrimSpace(c.Param("id"))

	file, err := c.FormFile("property_pic")
	if err != nil {
		logger.Info(err.Error())
		restErr := rest_errors.NewBadRequestErr("Bad Request")
		c.JSON(restErr.Status(), restErr)
		return
	}

	p, uploadErr := ph.service.UploadProperyPic(agencyID, file)
	if err != nil {
		c.JSON(uploadErr.Status(), uploadErr)
		return
	}

	c.JSON(http.StatusOK, p)
}
