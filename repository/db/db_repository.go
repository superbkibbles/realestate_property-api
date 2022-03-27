package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_property-api/clients/elasticsearch"
	"github.com/superbkibbles/realestate_property-api/domain/property"
	"github.com/superbkibbles/realestate_property-api/domain/query"
	"github.com/superbkibbles/realestate_property-api/helpers"
)

const (
	indexProperties        = "property"
	indexTranslateProperty = "translate_property"
	typeProperty           = "_doc"
)

type DbRepository interface {
	Create(property.Property) (*property.Property, rest_errors.RestErr)
	Get(sort string, asc bool) (property.Properties, rest_errors.RestErr)
	GetByID(string) (*property.Property, rest_errors.RestErr)
	Search(query query.EsQuery, sort string, asc bool) (property.Properties, rest_errors.RestErr)
	Update(id string, updateRequest property.EsUpdate) (*property.Property, rest_errors.RestErr)
	UploadMedia(visuals []property.Visual, videos []property.Video, propertyID string) rest_errors.RestErr
	GetActive(sort string, asc bool) (property.Properties, rest_errors.RestErr)
	GetDeactive(sort string, asc bool) (property.Properties, rest_errors.RestErr)
	Translate(translateProperty property.TranslateProperty) rest_errors.RestErr
	GetTranslateById(id string, local string) (*property.TranslateProperty, rest_errors.RestErr)
	UpdateTranslate(id string, translateProperty property.TranslateProperty) (*property.TranslateProperty, rest_errors.RestErr)
	GetAllTranslated(local string) (property.TranslateProperties, rest_errors.RestErr)
	GetActiveLocal(local string) (property.TranslateProperties, rest_errors.RestErr)
	GetDeactiveLocal(local string) (property.TranslateProperties, rest_errors.RestErr)
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (db *dbRepository) UploadMedia(visuals []property.Visual, videos []property.Video, propertyID string) rest_errors.RestErr {
	var esUpdate property.EsUpdate
	updated := property.UpdatePropertyRequest{
		Field: "visuals",
		Value: visuals,
	}
	videosUpdate := property.UpdatePropertyRequest{
		Field: "videos",
		Value: videos,
	}
	esUpdate.Fields = append(esUpdate.Fields, updated)
	esUpdate.Fields = append(esUpdate.Fields, videosUpdate)
	_, err := db.Update(propertyID, esUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (db *dbRepository) Update(id string, updateRequest property.EsUpdate) (*property.Property, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Update(indexProperties, typeProperty, id, updateRequest)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, rest_errors.NewNotFoundErr(fmt.Sprintf("no Property was found with id %s", id))
		}
		return nil, rest_errors.NewInternalServerErr("error when trying to Update Property", errors.New("databse error"))
	}

	var property property.Property

	bytes, _ := result.GetResult.Source.MarshalJSON()
	// if err != nil {
	// 	return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	// }
	// if err := json.Unmarshal(bytes, &property); err != nil {
	// 	return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	// }
	json.Unmarshal(bytes, &property)

	property.ID = result.Id
	return &property, nil
}

func (db *dbRepository) Translate(translateProperty property.TranslateProperty) rest_errors.RestErr {
	if err := elasticsearch.Client.Translate(indexTranslateProperty, typeProperty, translateProperty); err != nil {
		return rest_errors.NewInternalServerErr("Internal server error", err)
	}
	return nil
}

func (db *dbRepository) GetAllTranslated(local string) (property.TranslateProperties, rest_errors.RestErr) {
	results, err := elasticsearch.Client.GetAllTranslated(indexTranslateProperty, typeProperty, local)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, rest_errors.NewInternalServerErr("error when trying to Update Property", errors.New("databse error"))
	}

	tp, restErr := helpers.SearchResultToTranslatedProperties(results)

	if restErr != nil {
		return nil, restErr
	}

	return tp, nil
}

func (db *dbRepository) UpdateTranslate(id string, translateProperty property.TranslateProperty) (*property.TranslateProperty, rest_errors.RestErr) {
	var esUpdate property.EsUpdate

	description := property.UpdatePropertyRequest{
		Field: "description",
		Value: translateProperty.Description,
	}
	title := property.UpdatePropertyRequest{
		Field: "title",
		Value: translateProperty.Title,
	}
	directionFace := property.UpdatePropertyRequest{
		Field: "direction_face",
		Value: translateProperty.DirectionFace,
	}
	propertyType := property.UpdatePropertyRequest{
		Field: "property_type",
		Value: translateProperty.PropertyType,
	}
	category := property.UpdatePropertyRequest{
		Field: "category",
		Value: translateProperty.Category,
	}
	location := property.UpdatePropertyRequest{
		Field: "location",
		Value: translateProperty.Location,
	}
	city := property.UpdatePropertyRequest{
		Field: "city",
		Value: translateProperty.City,
	}

	esUpdate.Fields = append(esUpdate.Fields, description)
	esUpdate.Fields = append(esUpdate.Fields, city)
	esUpdate.Fields = append(esUpdate.Fields, location)
	esUpdate.Fields = append(esUpdate.Fields, category)
	esUpdate.Fields = append(esUpdate.Fields, propertyType)
	esUpdate.Fields = append(esUpdate.Fields, directionFace)
	esUpdate.Fields = append(esUpdate.Fields, title)

	result, err := elasticsearch.Client.UpdateTranslate(indexTranslateProperty, typeProperty, id, esUpdate)
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to Update Property", errors.New("databse error"))
	}

	var tp property.TranslateProperty
	bytes, _ := result.GetResult.Source.MarshalJSON()
	json.Unmarshal(bytes, &tp)
	return &tp, nil
}

func (db *dbRepository) GetTranslateById(id string, local string) (*property.TranslateProperty, rest_errors.RestErr) {
	results, err := elasticsearch.Client.GetTranslatByID(indexTranslateProperty, typeProperty, id, local)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, rest_errors.NewNotFoundErr("no Property was found with")
		}
		return nil, rest_errors.NewInternalServerErr("error when trying to Update Property", errors.New("databse error"))
	}

	var p property.TranslateProperty
	if len(results.Hits.Hits) > 0 {
		bytes, _ := results.Hits.Hits[0].Source.MarshalJSON()
		json.Unmarshal(bytes, &p)

		p.ID = results.Hits.Hits[0].Id
		return &p, nil
	}

	return &property.TranslateProperty{}, nil
}

func (db *dbRepository) Create(property property.Property) (*property.Property, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Save(indexProperties, typeProperty, property)
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to save Property", errors.New("databse error"))
	}
	property.ID = result.Id
	return &property, nil
}

func (db *dbRepository) Get(sort string, asc bool) (property.Properties, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Get(indexProperties, typeProperty, sort, asc)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, rest_errors.NewNotFoundErr("no Property was found with")
		}
		return nil, rest_errors.NewInternalServerErr("error when trying to Get all Properties", errors.New("databse error"))
	}

	// properties := helpers.SearchResultToProperties(result)
	// properties := utils.SearchResultToProperties(result)
	properties, restErr := helpers.SearchResultToProperties(result)
	if restErr != nil {
		return nil, restErr
	}
	// fmt.Println("************")
	// fmt.Println(len(properties))

	return properties, nil
}

func (db *dbRepository) GetByID(id string) (*property.Property, rest_errors.RestErr) {
	result, err := elasticsearch.Client.GetByID(indexProperties, typeProperty, id)
	var property property.Property
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, rest_errors.NewNotFoundErr(fmt.Sprintf("no Property was found with id %s", id))
		}
		return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to id %s", id), errors.New("database error"))
	}
	bytes, _ := result.Source.MarshalJSON()
	// if err != nil {
	// 	return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	// }
	// if err := json.Unmarshal(bytes, &property); err != nil {
	// 	return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	// }
	json.Unmarshal(bytes, &property)
	property.ID = result.Id

	return &property, nil
}

func (db *dbRepository) Search(query query.EsQuery, sort string, asc bool) (property.Properties, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Search(indexProperties, query.Build(), sort, asc)
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to search documents", errors.New("database error"))
	}

	properties, restErr := helpers.SearchResultToProperties(result)
	if restErr != nil {
		return nil, restErr
	}

	return properties, nil
}

func (db *dbRepository) GetActive(sort string, asc bool) (property.Properties, rest_errors.RestErr) {
	result, err := elasticsearch.Client.GetActive(indexProperties, typeProperty, sort, asc)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, rest_errors.NewNotFoundErr("no Property was found with Status Active")
		}
		return nil, rest_errors.NewInternalServerErr("error when trying to Get Status Active", errors.New("database error"))
	}
	properties, restErr := helpers.SearchResultToProperties(result)
	if restErr != nil {
		return nil, restErr
	}
	return properties, nil
}

func (db *dbRepository) GetDeactive(sort string, asc bool) (property.Properties, rest_errors.RestErr) {
	result, err := elasticsearch.Client.GetDeactive(indexProperties, typeProperty, sort, asc)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, rest_errors.NewNotFoundErr("no Property was found with Status Deactive")
		}
		return nil, rest_errors.NewInternalServerErr("error when trying to Get Status Deactive", errors.New("database error"))
	}

	properties, restErr := helpers.SearchResultToProperties(result)
	if restErr != nil {
		return nil, restErr
	}

	return properties, nil
}

func (db *dbRepository) GetActiveLocal(local string) (property.TranslateProperties, rest_errors.RestErr) {
	result, err := elasticsearch.Client.GetAllTranslated(indexTranslateProperty, typeProperty, local)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, rest_errors.NewInternalServerErr("error when trying to Get Status Active", errors.New("database error"))
	}
	properties, restErr := helpers.SearchResultToTranslatedProperties(result)
	if restErr != nil {
		return nil, restErr
	}
	return properties, nil
}

func (db *dbRepository) GetDeactiveLocal(local string) (property.TranslateProperties, rest_errors.RestErr) {
	result, err := elasticsearch.Client.GetAllTranslated(indexTranslateProperty, typeProperty, local)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			fmt.Println("Nothing was found!!!!")
			fmt.Println("Nothing was found!!!!")
			fmt.Println("Nothing was found!!!!")
			fmt.Println("Nothing was found!!!!")
			fmt.Println("Nothing was found!!!!")
			return nil, nil
		}
		return nil, rest_errors.NewInternalServerErr("error when trying to Get Status Deactive", errors.New("database error"))
	}

	properties, restErr := helpers.SearchResultToTranslatedProperties(result)
	if restErr != nil {
		return nil, restErr
	}

	return properties, nil
}
