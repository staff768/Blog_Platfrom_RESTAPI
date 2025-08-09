package main

import (
	"blogplatform/conf"
	"blogplatform/internal/models"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main(){

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	err := models.Connect()
	if err != nil {
		errorLog.Fatalf("Error from main while tryin to connect")
	}
	defer models.Db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	app.infoLog.Printf("Starting Server on %s:", conf.Cfg.ServerPort)
	err = http.ListenAndServe(conf.Cfg.ServerPort, app.routes())
	app.errorLog.Fatalf("Error while trying to up server %s", err) 
}