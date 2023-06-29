package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ronilsonalves/go-wallet-watcher/cmd/server/handler"
	"github.com/ronilsonalves/go-wallet-watcher/docs"
	"github.com/ronilsonalves/go-wallet-watcher/internal/wallet"
	"github.com/ronilsonalves/go-wallet-watcher/internal/watcher"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
)

// @title Go Wallet Watcher API
// @version 1.0
// @description This API handle query to check crypto wallet info such as balance, transactions...
// @termsOfService https://github.com/ronilsonalves/go-wallet-watcher/blob/main/LICENSE.md
// @contact.name Ronilson Alves
// @contact.url https://www.linkedin.com/in/ronilsonalves

// @license.name MIT
// @license.url https://github.com/ronilsonalves/go-wallet-watcher/blob/main/LICENSE.md
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file", err.Error())
	}

	wService := wallet.NewService()
	wHandler := handler.NewWalletHandler(wService)

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	docs.SwaggerInfo.Host = os.Getenv("WATCHER_API_DOMAIN")
	docs.SwaggerInfo.BasePath = "/api/v1/"
	r.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Everything is okay here",
		})
	})

	api := r.Group("/api/v1")
	{
		ethNet := api.Group("/eth/wallets")
		{
			ethNet.GET(":address", wHandler.GetWalletByAddress())
			ethNet.GET(":address/transactions", wHandler.GetTransactionsByAddress())
		}
	}

	// Start the Gin server in a goroutine
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Println("ERROR IN GONIC: ", err.Error())
		}
	}()

	// Start our watcher
	go watcher.StartWatcherService()

	// Wait for the server and the watcher service to finish
	select {}
}
