package main

import (
	"log"
	"myGin/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	log.Printf("Start Webservice\n")
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	// r.StaticFS("/static", http.Dir("./"))

	routes.Load(r)

	r.Run()
}
