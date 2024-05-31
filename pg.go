package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
    "os"
    "strings"
    "encoding/json"
)

type PG struct {
    Queries map[string]string
    Conn *sql.DB
}

func newPG() PG {
    connStr := "postgres://me:me@localhost/tft_tracker"
    db, connErr := sql.Open("postgres", connStr)
    if connErr != nil {
        fmt.Println("Error opening PG", connErr)
    }

    return PG{
        Queries: make(map[string]string),
        Conn: db,
    }
}

func (pg *PG) disconnect() {
    pg.Conn.Close()
}

func (pg *PG) addQueriesToMap() {
    queries, err := os.ReadDir("./queries")
    if err != nil {
        fmt.Println("error reading dir", err)
    }

    for _, file := range queries {
        data, err := os.ReadFile("./queries/" + file.Name())
        if err != nil {
            fmt.Println("error reading file", err)
        }
        splitName := strings.Split(file.Name(), ".sql")
        pg.Queries[splitName[0]] = string(data)
    }
}

/*
This is kinda messy basically we are assuming our queries all build a json object and not multiple rows

I need to test the performance of this to see how much it affects performance.
from very simple inital tests it seems to be slower but not by a huge margin.
*/
func (pg *PG) query(queryName string, args ...interface{}) interface{} {
    rows, err := pg.Conn.Query(pg.Queries[queryName], args...)
    if err != nil {
        fmt.Println("error querying", pg.Queries[queryName], err)
    }
    defer rows.Close()
    var data interface{}
    for rows.Next() {
        var result []byte 
        err := rows.Scan(&result)
        if err != nil {
            fmt.Println("error scanning row into result", err)
        }
        json.Unmarshal(result, &data)
    }
    return data 
}
