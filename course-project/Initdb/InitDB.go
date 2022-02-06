package Initdb

import (
	"course-project/types"
	"fmt"
	"github.com/jinzhu/gorm"
)

func InitDB() *gorm.DB {

	driverName := "mysql"
	host := "127.0.0.1"
	port := "3306"
	database := "course"
	username := "root"
	password := "12345"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset,
	)

	DB, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}

	//自动创建数据表
	DB.AutoMigrate(&types.TMember{})
	DB.AutoMigrate(&types.TCourse{})

	return DB
}
