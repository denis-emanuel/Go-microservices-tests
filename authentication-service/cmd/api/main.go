package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "8080"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to PG")
	}

	//set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	//ENV variable set in docker-compose.yml in the 'project' folder
	dsn := os.Getenv("DSN")

	//try to connect until successful
	for {
		connection, err := openDB(dsn)

		if err != nil {
			log.Println("PG not yet ready...")
			counts++
		} else {
			log.Println("Connected to PG")
			return connection
		}

		//limited number of retries
		if counts > 10 {
			log.Println(err)
			return nil
		}

		//retry with a 2 seconds pause
		log.Println("2s pause")
		time.Sleep(2 * time.Second)
		continue
	}
}
