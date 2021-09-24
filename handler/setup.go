package handler

import (
	"github.com/ggxxll/core-processor-demo/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&entity.Example{})
	if err != nil {
		panic(err)
	}

	return db, func() {
		if rdb, err := db.DB(); err == nil {
			rdb.Close()
		}
	}
}
