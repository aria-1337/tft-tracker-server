package main

import (
	"database/sql"
	"fmt"
    "log"
	_ "github.com/lib/pq"
)

func DBConnect() *sql.DB {
    connStr := "postgres://me:me@localhost/tft_tracker"
    db, connErr := sql.Open("postgres", connStr)
    if connErr != nil {
        fmt.Println("Error opening PG", connErr)
    }

    return db
}

func DBClose(conn *sql.DB) {
    defer conn.Close()
}

// Util functions - - - - Should move these later
// These are also only for testing will need to write a proper mirgrator
func LoadSchema(conn *sql.DB) {
    txn, err := conn.Begin(); 
    if err != nil {
        log.Fatal(err)
    }
    
    // IN ORDER OF PROPER CREATION ORDER
    // summoners
    // puuid - Source of truth for a given user immutable in the API 
    // game_name - Users game name (can change and needs to be updated)
    // tag_line - Users tagLine (a unique identifier to seperate gameName from others)
    // region - What region the user belongs too.
    // summoner_id Identifier
    // account_id Identifier
    // profile_icon_id INT represents what icon id they are using for their riot account.
    // summoner_level INT summoners current level

    // Notes:
    // We could store a history of game_name and tag_line if we want later on
    summoners := `
        CREATE TABLE summoners (
            puuid TEXT PRIMARY KEY,
            game_name TEXT,
            tag_line TEXT,
            region TEXT,
            summoner_id TEXT,
            account_id TEXT,
            profile_icon_id INT,
            summoner_level INT 
        );
    `
    if _, err := txn.Exec(summoners); err != nil {
        log.Fatal(err)
    }
    if err := txn.Commit(); err != nil {
        log.Fatal(err)
    }
}

func DestroySchema(conn *sql.DB) {
    destroy := `
        DROP TABLE summoner;
    `
    if _, err := conn.Exec(destroy); err != nil {
        log.Fatal(err)
    }
}
// - - - - - - - - - - - - - - - - - - - - - - -


/*
var (
    rowName string
    rowName2 int
)
rows, err := db.Query("SELECT rowName, rowName2 FROM table WHERE age = $1", age)
if err != nil {
    panic(err)
}
defer rows.Close()
for rows.Next() {
    err := rows.Scan(&rowName, &rowName2)
    if err != nil {
        panic(err)
    }
    fmt.Println("\n", rowName, rowName2)
}
err = rows.Err()
if err != nil {
    panic(err)
}
*/
