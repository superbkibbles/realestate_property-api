package elasticsearch

import (
	"context"
	"fmt"
	"os"
	"time"

	elastic "github.com/olivere/elastic/v7"
	"github.com/superbkibbles/bookstore_utils-go/logger"
	"github.com/superbkibbles/realestate_property-api/constants"
	"github.com/superbkibbles/realestate_property-api/domain/property"
)

var (
	Client EsClientInterface = &esClient{}
)

type EsClientInterface interface {
	Init()
	setClient(*elastic.Client)
	Save(string, string, interface{}) (*elastic.IndexResponse, error)
	Get(string, string, string, bool) (*elastic.SearchResult, error)
	GetByID(string, string, string) (*elastic.GetResult, error)
	Search(index string, query elastic.Query, sort string, asc bool) (*elastic.SearchResult, error)
	Update(indexProperties string, typeProperty string, id string, updateRequest property.EsUpdate) (*elastic.UpdateResponse, error)
	GetActive(index string, propertyType string, sort string, asc bool) (*elastic.SearchResult, error)
	GetDeactive(index string, propertyType string, sort string, asc bool) (*elastic.SearchResult, error)
	Translate(indexTranslateProperty string, docType string, doc interface{}) error
	GetTranslatByID(indexTranslateProperty string, docType string, propertyID string, local string) (*elastic.SearchResult, error)
	UpdateTranslate(indexTranslateProperty string, typeProperty string, id string, updateRequest property.EsUpdate) (*elastic.UpdateResponse, error)
	GetAllTranslated(indexTranslateProperty string, typeProperty string, local string) (*elastic.SearchResult, error)
	// GetDeactiveLocal(indexTranslateProperty string, typeProperty string, local string) (*elastic.SearchResult, error)
	// GetActiveLocal(indexTranslateProperty string, typeProperty string, local string) (*elastic.SearchResult, error)
}

type esClient struct {
	client *elastic.Client
}

func (c *esClient) setClient(client *elastic.Client) {
	c.client = client
}

func (c *esClient) Init() {
	log := logger.Getlogger()
	client, err := elastic.NewClient(
		elastic.SetURL(os.Getenv(constants.ELASTIC_URL)),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log),
		elastic.SetInfoLog(log),
	)
	if err != nil {
		panic(err)
	}
	Client.setClient(client)
}

func (c *esClient) Save(index string, docType string, doc interface{}) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	result, err := c.client.Index().
		Index(index).
		Type(docType).
		BodyJson(doc).
		Do(ctx)
	if err != nil {
		logger.Error(
			fmt.Sprintf("error while trying to index document in index %s", index), err)
		return nil, err
	}

	return result, nil
}

func (c *esClient) Translate(indexTranslateProperty string, docType string, doc interface{}) error {
	ctx := context.Background()

	if _, err := c.client.Index().Index(indexTranslateProperty).Type(docType).BodyJson(doc).Do(ctx); err != nil {
		return err
	}

	return nil
}

func (c *esClient) GetTranslatByID(indexTranslateProperty string, docType string, propertyID string, local string) (*elastic.SearchResult, error) {
	ctx := context.Background()
	query := elastic.NewBoolQuery()
	equalsQuery := make([]elastic.Query, 0)
	equalsQuery = append(equalsQuery, elastic.NewMatchQuery("property_id", propertyID))
	equalsQuery = append(equalsQuery, elastic.NewMatchQuery("local", local))

	query.Must(equalsQuery...)
	result, err := c.client.Search().Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// func (c *esClient) GetActiveLocal(indexTranslateProperty string, typeProperty string, local string) (*elastic.SearchResult, error) {
// 	ctx := context.Background()
// 	query := elastic.NewMatchQuery("local", local)
// 	equalsQuery := make([]elastic.Query, 0)
// 	equalsQuery = append(equalsQuery, elastic.NewMatchQuery("local", local))

// 	query.Must(equalsQuery...)
// 	result, err := c.client.Search().Query(query).Do(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func (c *esClient) GetDeactiveLocal(indexTranslateProperty string, typeProperty string, local string) (*elastic.SearchResult, error) {
// 	ctx := context.Background()
// 	query := elastic.NewBoolQuery()
// 	equalsQuery := make([]elastic.Query, 0)
// 	equalsQuery = append(equalsQuery, elastic.NewMatchQuery("local", local))

// 	query.Must(equalsQuery...)
// 	result, err := c.client.Search().Query(query).Do(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

func (c *esClient) GetAllTranslated(indexTranslateProperty string, typeProperty string, local string) (*elastic.SearchResult, error) {
	ctx := context.Background()
	query := elastic.NewMatchQuery("local", local)
	result, err := c.client.Search().Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *esClient) GetByID(index string, docType string, id string) (*elastic.GetResult, error) {
	ctx := context.Background()
	result, err := c.client.Get().
		Index(index).
		Type(docType).
		Id(id).
		Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when tring to get id %s", id), err)
		return nil, err
	}
	return result, nil
}

func (c *esClient) Get(index string, propertyType string, sort string, asc bool) (*elastic.SearchResult, error) {
	ctx := context.Background()
	query := elastic.NewMatchAllQuery()
	var result *elastic.SearchResult
	var err error
	// result, err := c.client.Index(index).Type(propertyType).Query(query).Do(ctx)
	if sort == "" {
		result, err = c.client.Search().Index(index).Type(propertyType).Query(query).Size(3).Do(ctx)
	} else {
		result, err = c.client.Search().Sort(sort, asc).Index(index).Type(propertyType).Query(query).Size(3).Do(ctx)
	}

	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
		return nil, err
	}
	return result, nil
}

func (c *esClient) GetActive(index string, propertyType string, sort string, asc bool) (*elastic.SearchResult, error) {
	ctx := context.Background()
	query := elastic.NewMatchQuery("status", "active")
	var results *elastic.SearchResult
	var err error
	if sort == "" {
		results, err = c.client.Search().Index(index).Type(propertyType).Query(query).Do(ctx)
	} else {
		results, err = c.client.Search().Sort(sort, asc).Index(index).Type(propertyType).Query(query).Do(ctx)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
		return nil, err
	}

	return results, nil
}

func (c *esClient) GetDeactive(index string, propertyType string, sort string, asc bool) (*elastic.SearchResult, error) {
	ctx := context.Background()
	query := elastic.NewMatchQuery("status", "deactive")
	var results *elastic.SearchResult
	var err error

	if sort == "" {
		results, err = c.client.Search().Index(index).Type(propertyType).Query(query).Do(ctx)
	} else {
		results, err = c.client.Search().Sort(sort, asc).Index(index).Type(propertyType).Query(query).Do(ctx)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
		return nil, err
	}
	return results, nil
}

func (c *esClient) Search(index string, query elastic.Query, sort string, asc bool) (*elastic.SearchResult, error) {
	ctx := context.Background()
	var results *elastic.SearchResult
	var err error
	if sort == "" {
		results, err = c.client.Search(index).Query(query).Do(ctx)
	} else {
		results, err = c.client.Search(index).Sort(sort, asc).Query(query).Do(ctx)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
		return nil, err
	}

	return results, nil
}

func (c *esClient) UpdateTranslate(indexTranslateProperty string, typeProperty string, id string, updateRequest property.EsUpdate) (*elastic.UpdateResponse, error) {
	ctx := context.Background()
	arr := make(map[string]interface{})
	for _, value := range updateRequest.Fields {
		arr[value.Field] = value.Value
	}
	result, err := c.client.Update().Index(indexTranslateProperty).Type(typeProperty).Id(id).Doc(arr).FetchSource(true).Do(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *esClient) Update(indexProperties string, typeProperty string, id string, updateRequest property.EsUpdate) (*elastic.UpdateResponse, error) {
	ctx := context.Background()
	arr := make(map[string]interface{})
	for _, value := range updateRequest.Fields {
		arr[value.Field] = value.Value
	}

	result, err := c.client.Update().Index(indexProperties).Type(typeProperty).Id(id).Doc(arr).FetchSource(true).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to Update documents in index %s", indexProperties), err)
		return nil, err
	}

	return result, nil
}
