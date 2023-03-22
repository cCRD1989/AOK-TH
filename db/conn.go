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

}
func Migrate() {
	err := Conn.AutoMigrate(
		&model.LogWeb{},
		// &model.Category{},
		// &model.Product{},
		// &model.Order{},
		// &model.OrderItem{},
	)
	if err != nil {
		log.Fatal("Connot ")
		return
	}
}
