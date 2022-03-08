package property

import (
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
)

type TranslateProperty struct {
	ID            string `json:"id"`
	PropertyID    string `json:"property_id"`
	Description   string `json:"description"`
	Title         string `json:"title"`
	DirectionFace string `json:"direction_face"`
	PropertyType  string `json:"property_type"`
	Category      string `json:"category"`
	Location      string `json:"location"`
	City          string `json:"city"`
	Local         string `json:"local"`
}

type TranslateProperties []TranslateProperty

func (t TranslateProperty) Validate() rest_errors.RestErr {
	switch t.Category {
	case "apartment", "house", "villa", "land", "farm", "بيت", "شفة":
		break
	default:
		return rest_errors.NewBadRequestErr("invalid category")
	}
	switch t.Local {
	case "en", "ar", "kur":
		break
	default:
		return rest_errors.NewBadRequestErr("invalid Language")
	}

	return nil
}

func (t *TranslateProperties) Marshal(ps Properties) Properties {
	results := ps
	for index, p := range ps {
		for _, tp := range *t {
			if p.ID == tp.PropertyID {
				t := tp.Marshal(&p)
				results[index] = *t
			}
		}
	}
	return results
}

func (t *TranslateProperty) Marshal(p *Property) *Property {
	// translatedJson, _ := json.Marshal(t)
	// property := *p
	// json.Unmarshal(translatedJson, &property)
	// property.ID = p.ID
	if t.Category != "" {
		p.Category = t.Category
	}
	if t.City != "" {
		p.City = t.City
	}
	if t.Description != "" {
		p.City = t.Description
	}
	if t.Location != "" {
		p.Location = t.Location
	}
	if t.PropertyType != "" {
		p.PropertyType = t.PropertyType
	}
	if t.Title != "" {
		p.Title = t.Title
	}

	return p
}
