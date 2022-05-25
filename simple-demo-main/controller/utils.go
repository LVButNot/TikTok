package controller

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GLOBAL_DB *gorm.DB

func ConnectionSQL() {
	dsn := "root:wulingwei@tcp(175.178.126.39:3306)/tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:   false, //是否跳过默认事务
		DisableNestedTransaction: true,  //在 AutoMigrate 或 CreateTable 时，GORM 会自动创建外键约束，若要禁用该特性，可将其设置为 true
	})
	if err != nil {
		fmt.Println(db)
	}
	GLOBAL_DB = db
}
