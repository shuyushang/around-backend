package service

import (
	"mime/multipart"
	"reflect"

	"around/backend"
	"around/constants"
	"around/model"

	"github.com/olivere/elastic/v7"
)

//support user-based search and keyword-based search, and call backend.ESBackend.ReadFromES

func SearchPostsByUser(user string) ([]model.Post, error) { //返回slice of posts，不一定是一个post返回
	query := elastic.NewTermQuery("user", user)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func SearchPostsByKeywords(keywords string) ([]model.Post, error) {
	query := elastic.NewMatchQuery("message", keywords)
	query.Operator("AND") //所有单词必须全部包含 hello + world
	//corner case，空字符keyword
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

// 模块话 filter部分
func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
	var ptype model.Post
	var posts []model.Post

	//filter
	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		p := item.(model.Post)
		posts = append(posts, p)
	}
	return posts
}

// help save uploaded data to ES and GCS
func SavePost(post *model.Post, file multipart.File) error {
	//save to GCS
	medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
	if err != nil {
		return err
	}
	//save to ES
	post.Url = medialink

	return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)

	/*
		err = backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
		if err != nil {
			1. rollback: delete from GCS
			2. retry: call SaveToES again
			3. offline service:
		}
	*/
}
func DeletePost(id string, user string) error {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", id))
	query.Must(elastic.NewTermQuery("user", user))

	return backend.ESBackend.DeleteFromES(query, constants.POST_INDEX)
}
