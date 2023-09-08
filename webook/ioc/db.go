package ioc

import (
	"awesomeProject/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("indigo:indigotest@tcp(10.1.80.122:3306)/go_test"))
	//db, err := gorm.Open(mysql.Open("indigo:indigotest@tcp(10.1.90.122:3306)/go_test"))
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
