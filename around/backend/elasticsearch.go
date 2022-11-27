package backend

import (
	"context"
	"fmt"

	"around/constants"
	"around/util"

	"github.com/olivere/elastic/v7"
)

// https://github.com/olivere/elastic/blob/release-branch.v7/example_test.go
// 用于与ES数据库交互，相当于DAO in online order
var (
	ESBackend *ElasticsearchBackend
	//大写的，有且只有一个，globlal的，谁需要谁调用他；只初始化一次；
	//singlton
)

type ElasticsearchBackend struct { //这个strut想到于customerDAO
	client *elastic.Client
	//new a client, like new a sessionFactory in Java
}

func InitElasticsearchBackend(config *util.ElasticsearchInfo) {
	//初始化，new ElasticsearchBackend object
	//hibernate里面已经写好了，这里要自己写
	client, err := elastic.NewClient(
		elastic.SetURL(config.Address),
		elastic.SetBasicAuth(config.Username, config.Password))
	if err != nil {
		panic(err)
	}

	exists, err := client.IndexExists(constants.POST_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}
	//创建post index
	if !exists {
		mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "message":  { "type": "text" }, 
                    //text 模糊匹配
                    "url":      { "type": "keyword", "index": false },
                    "type":     { "type": "keyword", "index": false }
                }
            }
        }`
		_, err := client.CreateIndex(constants.POST_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
	//context -> 可以放一个deadline，避免傻等,特定时间，任务不完成就结束
	if err != nil {
		panic(err)
	}
	//判断user在不在，不在创建user index
	if !exists {
		mapping := `{
                        "mappings": {
                                "properties": {
                                        "username": {"type": "keyword"},
										//term query,精准匹配
                                        "password": {"type": "keyword"},
                                        "age":      {"type": "long", "index": false},
                                        "gender":   {"type": "keyword", "index": false}
                                        //without index, searching linear, compared slow
                                }
                        }
                }`
		_, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Indexes are created.")

	ESBackend = &ElasticsearchBackend{client: client}
	//type client : func init client
	//创建好的sessonfactory 保存到es backend
	//类似于v = &Vertex(X:x)
}

// go 把方法拿出来，backend ElasticsearchBackend相当于告诉是哪个类的方法；
func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	searchResult, err := backend.client.Search().
		Index(index).            //search in index
		Query(query).            //specify the query
		Pretty(true).            //pretty print request and response Json
		Do(context.Background()) //execute
	if err != nil {
		return nil, err
	}

	return searchResult, nil
}

// save a post data to ElasticSearch
func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
	//i interface{} 相当于object class in Java
	_, err := backend.client.Index().
		Index(index).
		Id(id).
		BodyJson(i).
		Do(context.Background())
	return err
}

func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
	_, err := backend.client.DeleteByQuery().
		Index(index).
		Query(query).
		Pretty(true).
		Do(context.Background())

	return err
}
