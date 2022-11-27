package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"around/model"
	"around/service"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
)

// add a literal map for different media types
var (
	mediaTypes = map[string]string{
		".jpeg": "image",
		".jpg":  "image",
		".gif":  "image",
		".png":  "image",
		".mov":  "video",
		".mp4":  "video",
		".avi":  "video",
		".flv":  "video",
		".wmv":  "video",
	}
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	/* * -> 要求指向request的对象的指针（地址）
	需要地址的原因： pass by value like java, 保证更改之后，也可以更新；
	也因为是shallow copy快很多； 能传ptr时候尽量传ptr
	Parse from body of request to get a json object.
	fmt.Println("Received one post request")
	decoder := json.NewDecoder(r.Body)
	把request的body变成一个post类型的对象，出错抛出异常panic */
	// var p model.Post
	// if err := decoder.Decode(&p); err != nil {
	// 	panic(err)
	// }

	// fmt.Fprintf(w, "Post received: %s\n", p.Message)
	fmt.Println("Received one upload request")

	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]

	p := model.Post{
		Id:      uuid.New(),
		User:    username.(string),
		Message: r.FormValue("message"),
	}

	file, header, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		//error code 400 客户端问题
		fmt.Printf("Media file is not available %v\n", err)
		return
	}

	suffix := filepath.Ext(header.Filename)
	if t, ok := mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
		//error handle http.Error(w, "File type is not support", http.StatusBadRequest)
	}
	//根据title返回type

	err = service.SavePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
		//error code 500 服务器问题
		fmt.Printf("Failed to save post to backend %v\n", err)
		return
	}

	fmt.Println("Post is saved successfully.")
}

/* 1. req := http.Request(...)
resp := http.ResponseWriter(...)
uploadHandler(resp, &req)
2. req := &http.Request(...)
resp := http.ResponseWriter(...)
uploadHandler(resp, req) */

// add a new handler function to handler search-related requests.
func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for search")
	w.Header().Set("Content-Type", "application/json")

	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	var posts []model.Post
	var err error
	if user != "" {
		posts, err = service.SearchPostsByUser(user)
	} else {
		posts, err = service.SearchPostsByKeywords(keywords)
	}

	if err != nil {
		http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
		fmt.Printf("Failed to read post from backend %v.\n", err)
		return
	}

	js, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
		return
	}
	w.Write(js)
}
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for delete")

	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)
	id := mux.Vars(r)["id"]

	if err := service.DeletePost(id, username); err != nil {
		http.Error(w, "Failed to delete post from backend", http.StatusInternalServerError)
		fmt.Printf("Failed to delete post from backend %v\n", err)
		return
	}
	fmt.Println("Post is deleted successfully")
}
