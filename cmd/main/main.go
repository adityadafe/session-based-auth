package main

import (
	"fmt"
	"log"
	"session-based-auth/internal/db"
	"session-based-auth/internal/routes"
)

func main() {

	sqliteStore, err := db.NewSqliteStore()

	if err != nil {
		log.Fatal("Error opening db ", err)
	}

	if err = sqliteStore.Init(); err != nil {
		log.Fatal("Error initing db ", err)
	}

	port := 8080
	addr := fmt.Sprintf(":%d", port)

	router := routes.NewApiServer(addr, sqliteStore)
	router.Run()

}
