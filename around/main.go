package main

import (
	"fmt"
	"log"
	"net/http"

	"around/backend"
	"around/handler"
	"around/util"
)

func main() {
	fmt.Println("started-service")

	config, err := util.LoadApplicationConfig("conf", "deploy.yml")
	if err != nil {
		panic(err)
	}

	backend.InitElasticsearchBackend(config.ElasticsearchConfig)
	backend.InitGCSBackend(config.GCSConfig)

	// backend.InitElasticsearchBackend()
	// backend.InitGCSBackend() //initialize GCS client when it starts

	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter(config.TokenConfig)))
	//log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
	//启动一个go里自带的标准库的server监听8080端口，并使用之前的router来分发
	//如果错误，打印到log里
}
