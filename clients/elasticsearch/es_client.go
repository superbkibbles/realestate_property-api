package elasticsearch

import (
	"context"
	"fmt"
	"time"

	elastic "github.com/olivere/elastic/v7"
	"github.com/superbkibbles/bookstore_utils-go/logger"
	"github.com/superbkibbles/realestate_property-api/domain/property"
)

var (
	Client EsClientInterface = &esClient{}
)

type EsClientInterface interface {
	Init()
	setClient(*elastic.Client)
	Save(string, string, interface{}) (*elastic.IndexResponse, error)
	Get(string) (*elastic.SearchResult, error)
	GetByID(string, string, string) (*elastic.GetResult, error)
	Search(index string, query elastic.Query) (*elastic.SearchResult, error)
	Update(indexProperties string, typeProperty string, id string, updateRequest property.EsUpdate) (*elastic.UpdateResponse, error)
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
		elastic.SetURL("http://127.0.0.1:9200"),
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

func (c *esClient) Get(index string) (*elastic.SearchResult, error) {
	ctx := context.Background()
	query := elastic.MatchAllQuery{}

	result, err := c.client.Search(index).Query(query).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
		return nil, err
	}
	return result, nil
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

func (c *esClient) Search(index string, query elastic.Query) (*elastic.SearchResult, error) {
	ctx := context.Background()
	result, err := c.client.Search(index).Query(query).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
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
