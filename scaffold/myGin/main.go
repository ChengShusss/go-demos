package main

import (
	"myGin/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.Load(r)

	r.Run()
}
