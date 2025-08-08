package main

import (
	"blogplatform/conf"
	"blogplatform/internal/models"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main(){
	err := models.Connect()
	if err != nil {
		log.Fatalf("Error from main while tryin to connect")
	}
	defer models.Db.Close()
	
	log.Printf("Starting Server on %s:", conf.Cfg.ServerPort)
	err = http.ListenAndServe(conf.Cfg.ServerPort, routes())
	log.Fatalf("Error while trying to up server %s", err) 
}