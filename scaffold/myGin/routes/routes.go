package routes

//routes/routes.go

import (
	"myGin/controller"
	"myGin/utils"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

const (
	Static_Relative_Path = "/ui"
)

func Load(r *gin.Engine) {

	// r.GET("/", controller.Index)

	r.GET("/insert", controller.Create)
	r.GET("/get", controller.GetAllInfo)

	// Static files  --> current for ui
	LoadStaticFiles(r)
}

func LoadStaticFiles(r *gin.Engine) {
	staticPath := path.Join(utils.GetExeDir(), Static_Relative_Path)
	entries, err := os.ReadDir(staticPath)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			r.Static("/"+entry.Name(), path.Join(staticPath, entry.Name()))
		} else {
			r.StaticFile("/"+entry.Name(), path.Join(staticPath, entry.Name()))
		}
	}

	r.StaticFile("/", path.Join(staticPath, "index.html"))
}
