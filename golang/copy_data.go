package main

import (
    "database/sql"
    "fmt"
    "time"

    _ "github.com/mattn/go-sqlite3"
    _ "github.com/lib/pq"
)


func main() {
    sqlite, err := sql.Open("sqlite3", "../fakedata.db")
    CheckError(err)
    defer sqlite.Close()

    rows, err := sqlite.Query("SELECT * FROM users LIMIT 5000")
    CheckError(err)

    psqlconn := fmt.Sprintf("host=localhost port=5432 user=postgres password=postgres dbname=postgres")
    postgres, err := sql.Open("postgres", psqlconn)
    CheckError(err)
    defer postgres.Close()

    drop_stmt, err := postgres.Prepare("DROP TABLE IF EXISTS users")
    CheckError(err)
    drop_stmt.Exec()

    create_str := `
        CREATE TABLE users (
            user_id INTEGER PRIMARY KEY,
            user_name TEXT,
            user_age INTEGER,
            user_address TEXT
        )
    `
    create_stmt, err := postgres.Prepare(create_str)
    create_stmt.Exec()

    insert_str := `
        INSERT INTO users (user_id, user_name, user_age, user_address)
        VALUES            (     $1,        $2,       $3,           $4)
    `
    stmt, err := postgres.Prepare(insert_str)
    CheckError(err)

    var user_id int
    var user_name string
    var user_age int
    var user_address string

    start := time.Now()
    for rows.Next() {
        rows.Scan(&user_id, &user_name, &user_age, &user_address)
        stmt.Exec(user_id, user_name, user_age, user_address)
    }
    elapsed := time.Since(start).Seconds()

    fmt.Printf("Finished in %.2f seconds\n", elapsed)
}


func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}
