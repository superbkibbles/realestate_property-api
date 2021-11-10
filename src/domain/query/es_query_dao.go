package query

import "github.com/olivere/elastic"

func (q *EsQuery) Build() elastic.Query {
	query := elastic.NewBoolQuery()
	equalsQuery := make([]elastic.Query, 0)
	for _, eq := range q.Equals {
		equalsQuery = append(equalsQuery, elastic.NewMatchQuery(eq.Field, eq.Value))
	}
	query.Must(equalsQuery...)
	return query
}
