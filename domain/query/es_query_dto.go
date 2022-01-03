package query

type EsQuery struct {
	Equals []FieldValue  `json:"equals"`
	Gt     []GtValue     `json:"gt"`
	Range  []RangeStruct `json:"range"`
}

type FieldValue struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type GtValue struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type RangeStruct struct {
	Field string `json:"field"`
	From  int64  `json:"from"`
	To    int64  `json:"to"`
}
