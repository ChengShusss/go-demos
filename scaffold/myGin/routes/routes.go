package routes

//routes/routes.go

import (
	"myGin/controller"

	"github.com/gin-gonic/gin"
)

func Load(r *gin.Engine) {

	r.GET("/", controller.Index)

}
