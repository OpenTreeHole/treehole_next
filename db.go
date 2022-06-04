package main

import (
	"treehole_next/apis/division"
	"treehole_next/db"
)

func initDB() {
	db.InitDB()
	err := db.DB.AutoMigrate(&division.Division{})
	if err != nil {
		panic(err)
	}
}
