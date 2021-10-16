package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var wg sync.WaitGroup

	sqlite_db, err := sql.Open("sqlite3", "../fakedata.db")

	if err != nil {
		log.Fatal(err)
	}

	defer sqlite_db.Close()

	psqlconn := fmt.Sprintf("host=localhost port=5432 user=postgres password=postgres dbname=postgres")
	pg_db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		log.Fatal(err)
	}

	defer pg_db.Close()

	start := time.Now()

	_, err = pg_db.Exec("DROP TABLE IF EXISTS users")

	if err != nil {
		log.Fatal(err)
	}

	_, err = pg_db.Exec(`
		CREATE TABLE users (
			user_id INTEGER PRIMARY KEY,
			user_name TEXT,
			user_age INTEGER,
			user_address TEXT
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	var rowcount int
	rows, _ := sqlite_db.Query("SELECT COUNT(*) FROM users")
	for rows.Next() {
		err := rows.Scan(&rowcount)
		if err != nil {
			log.Fatal(err)
		}
	}

	chunksize := 25_000

	for i := 0; i*chunksize <= rowcount; i++ {
		wg.Add(1)
		go copyChunk(&wg, i*chunksize, chunksize, sqlite_db, pg_db)
	}

	wg.Wait()
	fmt.Printf("Finished in %.2f seconds\n", time.Since(start).Seconds())
}

func copyChunk(wg *sync.WaitGroup, startidx int, chunksize int, from_db *sql.DB, to_db *sql.DB) {
	defer wg.Done()

	var user_id int
	var user_name string
	var user_age int
	var user_address string

	rows, err := from_db.Query("SELECT * FROM users LIMIT $1 OFFSET $2;", chunksize, startidx)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	tx, err := to_db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO users (user_id, user_name, user_age, user_address)
		VALUES            (     $1,        $2,       $3,           $4)
	`)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		rows.Scan(&user_id, &user_name, &user_age, &user_address)
		stmt.Exec(user_id, user_name, user_age, user_address)
	}

	err = tx.Commit()

	if err != nil {
		log.Fatal(err)
	}
}
