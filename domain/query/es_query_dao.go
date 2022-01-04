package query

import "github.com/olivere/elastic/v7"

func (q *EsQuery) Build() elastic.Query {
	query := elastic.NewBoolQuery()
	equalsQuery := make([]elastic.Query, 0)
	for _, eq := range q.Equals {
		equalsQuery = append(equalsQuery, elastic.NewMatchQuery(eq.Field, eq.Value))
	}

	for _, gtFilter := range q.Gt {
		query.Filter(elastic.NewRangeQuery(gtFilter.Field).Gt(gtFilter.Value))
	}

	for _, fRange := range q.Range {
		query.Filter(elastic.NewRangeQuery(fRange.Field).From(fRange.From).To(fRange.To))
	}

	query.Must(equalsQuery...)
	return query
}
