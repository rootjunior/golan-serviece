package main

import (
	"log"

	"github.com/glebarez/sqlite" // pure-Go драйвер для GORM
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	log.Println("✅ SQLite connected")

	err = db.AutoMigrate(&PostDB{})
	if err != nil {
		log.Fatal("migration failed:", err)
	}

	return db
}
