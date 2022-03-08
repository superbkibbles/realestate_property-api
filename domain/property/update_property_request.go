package property

import "github.com/superbkibbles/bookstore_utils-go/rest_errors"

type EsUpdate struct {
	Fields []UpdatePropertyRequest `json:"fields"`
}

type UpdatePropertyRequest struct {
	Field string      `json:"field"`
	Value interface{} `json:"Value"`
}

func (u EsUpdate) Validate() rest_errors.RestErr {
	for _, field := range u.Fields {
		switch field.Field {
		case "category":	
			switch field.Value {
			case "apartment", "house", "villa", "land", "farm":
				return nil
			default:
				return rest_errors.NewBadRequestErr("invalid JSON BODY category")
			}
		}
	}
	return nil
}
