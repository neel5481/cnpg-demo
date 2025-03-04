package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"strings"
	"os/signal"
    	"syscall"

	"github.com/jackc/pgx/v5"
)

func main() {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	c := make(chan os.Signal)
    	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    	go func() {
            <-c
    	    fmt.Printf("cleanup the connection")
            conn.Close(context.Background())
            os.Exit(1)
        }()

	for {
            time.Sleep(500 * time.Millisecond)
 	    rows, err := conn.Query(context.Background(), "SELECT version();")
	    if err != nil {
    		fmt.Printf("Error while executing select version() query: %v", err);
	    }
	    defer rows.Close()

	    var pgVersion string
	    for rows.Next() {
	        err := rows.Scan(&pgVersion)
	        if err != nil {
	            fmt.Printf("error while fetching row: %v", err);
	        }
	        versions := strings.Fields(pgVersion)
	        fmt.Printf("Time --> [%s], pgVersion --> [%s:%s]\n", time.Now().Format("2006-01-02 3:4:5.00000"), versions[0], versions[1])
	    }
	    if err := rows.Err(); err != nil {
	        fmt.Printf("error while iterating row: %v", err);
	    }
        }
}
