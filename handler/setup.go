package handler

import (
	"github.com/GGXXLL/core-processor-demo/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupDB() (*gorm.DB, func()) {
	db, err := gorm.Open(mysql.Open(":memory:"))
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
