package property

type EsUpdate struct {
	Fields []UpdatePropertyRequest `json:"fields"`
}

type UpdatePropertyRequest struct {
	Field string      `json:"field"`
	Value interface{} `json:"Value"`
}
