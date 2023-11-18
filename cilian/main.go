package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/", base)

	router.Run("localhost:8020")
}

func base(c *gin.Context) {
	fmt.Println("hello, world!")
	c.JSON(http.StatusOK, "hello, world!")
}
