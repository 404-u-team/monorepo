package main

import (
	db "DevSpace/DB"
	net "DevSpace/Net"
	system "DevSpace/System"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	config, err := system.LoadConfig()

	if err != nil {
		fmt.Println("Ошибка создания конфига: ", err.Error())
		return
	}

	db, err := db.InitDB(config)
	if err != nil {
		fmt.Println("Ошибка подключения к БД: ", err.Error())
		return
	}

	err = db.AutoMigrate()

	if err != nil {
		fmt.Println("Ошибка миграции: ", err.Error())
		return
	}

	rest := net.Rest{Config: config, DB: &db}
	//	fmt.Printf("Загружен конфиг: %+v\n", config)
	//  fmt.Printf("Порт из конфига: %d\n", config.APIPort)

	port := fmt.Sprintf(":%d", config.APIPort)

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/users/create", rest.RegisterUser)
	router.GET("/users/auth", rest.AuthUser)

	err = router.Run(port)

	if err != nil {
		fmt.Println("Слушаю %s", port)
	}
}
