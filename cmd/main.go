package main

import (
	"log"
	"github.com/rirachii/golivechat/db"
	// "github.com/rirachii/golivechat/users"
)

func main() {
	_, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize postgres db connection: %s", err)
	}

	// userRep := users.NewRepository(dbConn.GetDB())
	// userSvc := users.NewService(userRep)
	// userHandler := users.NewHandler(userSvc)
}
