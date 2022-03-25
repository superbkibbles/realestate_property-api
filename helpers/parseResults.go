package helpers

import (
	"encoding/json"

	elastic "github.com/olivere/elastic/v7"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_property-api/domain/property"
)

func SearchResultToProperties(result *elastic.SearchResult) (property.Properties, rest_errors.RestErr) {
	var properties []property.Property
	for _, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var property property.Property
		// if err := json.Unmarshal(bytes, &property); err != nil {
		// 	return nil, rest_errors.NewInternalServerErr("error when trying to parse response", errors.New("database error"))
		// }
		json.Unmarshal(bytes, &property)
		property.ID = hit.Id
		properties = append(properties, property)
	}

	if len(properties) == 0 {
		return nil, rest_errors.NewNotFoundErr("no Property was found with Status Deactive")
	}

	return properties, nil
}

func SearchResultToTranslatedProperties(result *elastic.SearchResult) (property.TranslateProperties, rest_errors.RestErr) {
	var properties property.TranslateProperties
	for i, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var property property.TranslateProperty
		// if err := json.Unmarshal(bytes, &property); err != nil {
		// 	return nil, rest_errors.NewInternalServerErr("error when trying to parse response", errors.New("database error"))
		// }
		json.Unmarshal(bytes, &property)
		property.ID = hit.Id
		properties[i] = property
	}

	if len(properties) == 0 {
		return nil, rest_errors.NewNotFoundErr("no Property was found with")
	}

	return properties, nil
}
