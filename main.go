package main

import (
	"fast-filestore-server/route"
	"fmt"
)

func main() {
	router := route.Router()

	//启动服务并监听
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Printf("Failed to start server, err:%s\n", err.Error())
	}
}
