package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_property-api/src/clients/elasticsearch"
	"github.com/superbkibbles/realestate_property-api/src/domain/property"
	"github.com/superbkibbles/realestate_property-api/src/domain/query"
)

const (
	indexProperties = "property"
	typeProperty    = "_doc"
)

type DbRepository interface {
	Create(property.Property) (*property.Property, rest_errors.RestErr)
	Get() (property.Properties, rest_errors.RestErr)
	GetByID(string) (*property.Property, rest_errors.RestErr)
	Search(query query.EsQuery) (property.Properties, rest_errors.RestErr)
	Update(id string, updateRequest property.EsUpdate) (*property.Property, rest_errors.RestErr)
	UploadMedia(visuals []property.Visual, propertyID string) rest_errors.RestErr
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (db *dbRepository) UploadMedia(visuals []property.Visual, propertyID string) rest_errors.RestErr {
	var esUpdate property.EsUpdate
	updated := property.UpdatePropertyRequest{
		Field: "visuals",
		Value: visuals,
	}
	esUpdate.Fields = append(esUpdate.Fields, updated)
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

	bytes, err := result.GetResult.Source.MarshalJSON()
	if err != nil {
		return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	}
	if err := json.Unmarshal(bytes, &property); err != nil {
		return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	}

	property.ID = result.Id
	return &property, nil
}

func (db *dbRepository) Create(property property.Property) (*property.Property, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Save(indexProperties, typeProperty, property)
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to save Property", errors.New("databse error"))
	}
	property.ID = result.Id
	return &property, nil
}

func (db *dbRepository) Get() (property.Properties, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Get(indexProperties)
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to Get all Properties", errors.New("databse error"))
	}

	properties := make([]property.Property, result.TotalHits())
	for i, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var property property.Property
		if err := json.Unmarshal(bytes, &property); err != nil {
			return nil, rest_errors.NewInternalServerErr("error when trying to parse response", errors.New("database error"))
		}
		property.ID = hit.Id
		properties[i] = property
	}

	if len(properties) == 0 {
		return nil, rest_errors.NewNotFoundErr("no items found matching given critirial")
	}

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
	bytes, err := result.Source.MarshalJSON()
	if err != nil {
		return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	}
	if err := json.Unmarshal(bytes, &property); err != nil {
		return nil, rest_errors.NewInternalServerErr(fmt.Sprintf("error when trying to parse database response"), errors.New("database error"))
	}
	property.ID = result.Id

	return &property, nil
}

func (db *dbRepository) Search(query query.EsQuery) (property.Properties, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Search(indexProperties, query.Build())
	if err != nil {
		return nil, rest_errors.NewInternalServerErr("error when trying to search documents", errors.New("database error"))
	}

	properties := make([]property.Property, result.TotalHits())
	for i, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var property property.Property
		if err := json.Unmarshal(bytes, &property); err != nil {
			return nil, rest_errors.NewInternalServerErr("error when trying to parse response", errors.New("database error"))
		}
		property.ID = hit.Id
		properties[i] = property
	}

	if len(properties) == 0 {
		return nil, rest_errors.NewNotFoundErr("no items found matching given critirial")
	}

	return properties, nil
}
