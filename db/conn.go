package db

import (
	"ccrd/model"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Conn *gorm.DB
var AOK_DB *gorm.DB

func ConnectDB() {
	//user:pass@tcp
	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)},
	)
	if err != nil {

		log.Fatal("Connot Connect to The Database")
		return
	}
	Conn = db
	fmt.Println("Database Connect Dons.")

	// AOK DB
	db1, err1 := gorm.Open(
		mysql.Open(os.Getenv("DATABASEAOK_DSN")),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)},
	)
	if err1 != nil {

		log.Fatal(" Connect to The AOK_DB")
		return
	}
	AOK_DB = db1
	fmt.Println("AOK_DB Connect Dons.")

}
func Migrate() {

	err := Conn.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.LogWeb{}, &model.LogTopup{})

	// err := Conn.AutoMigrate(
	// 	&model.LogWeb{},
	// 	// &model.Category{},
	// 	// &model.Product{},
	// 	// &model.Order{},
	// 	// &model.OrderItem{},
	// )
	if err != nil {
		log.Fatal("Connot AutoMigrate Error.")
		return
	}
}
