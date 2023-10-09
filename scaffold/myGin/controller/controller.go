package controller

import (
	"myGin/services"

	"github.com/gin-gonic/gin"
)

func Index(context *gin.Context) {
	context.JSON(200, gin.H{"msg": "hello world"})
}

func Create(context *gin.Context) {

	err := services.InsertInfo(services.GetInfoModel("title-1", "default/1233", []string{"123", "456"}))

	if err != nil {
		context.JSON(200, gin.H{"msg": "failed, err: " + err.Error()})
		return
	}

	context.JSON(200, gin.H{"msg": "succeed"})
}

func GetAllInfo(context *gin.Context) {
	res := services.GetAllInfo()
	context.JSON(200, res)
}
