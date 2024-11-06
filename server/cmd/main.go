package main

import (
	"log"
	"server/db"
	"server/internal/user"
	"server/internal/ws"
	"server/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	userRep := user.NewRepository(dbConn.GetDB())
	userService := user.NewUserService(userRep)
	userHandler := user.NewHandler(userService)

	hub := ws.NewWebsocketHub()
	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	err = router.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
