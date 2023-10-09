package services

import (
	"log"
	"myGin/model"
	"strings"
	"time"
)

func InsertInfo(info *model.Info) error {
	dbMap, err := model.GetDbMap("")
	if err != nil {
		return err
	}

	return dbMap.Insert(info)
}

func GetAllInfo() []model.Info {
	var infos []model.Info
	dbMap, err := model.GetDbMap("")
	if err != nil {
		return nil
	}

	_, err = dbMap.Select(&infos, "select * from infos order by id")
	if err != nil {
		log.Printf("Failed to get infos, err: %v\n", err)
		return nil
	}
	for i := range infos {
		infos[i].NewTags = strings.Split(infos[i].Tags, ",")
		infos[i].Tags = ""
	}
	return infos
}

func GetInfoModel(title, category string, tags []string) *model.Info {
	return &model.Info{
		Created:  time.Now().Unix(),
		Title:    title,
		Category: category,
		Tags:     strings.Join(tags, ","),
	}
}
