package main

import (
	system "DevSpace/System"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	config, err := system.LoadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Загружен конфиг: %+v\n", config)
	fmt.Printf("Порт из конфига: %d\n", config.APIPort)

	port := fmt.Sprintf(":%d", config.APIPort)
	router := gin.Default()
	router.Use(cors.Default())

	err = router.Run(port)

	if err != nil {
		fmt.Println("Слушаю %s", port)
	}
}
