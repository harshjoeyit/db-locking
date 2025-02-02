package main

import (
	"database/sql"
	"fmt"

	"github.com/bxcodec/faker/v3"
	_ "github.com/go-sql-driver/mysql"
)

func newConn() *sql.DB {
	user := "root"
	pass := "mysql"
	dbName := "airline"

	dsn := fmt.Sprintf("%s:%s@/%s", user, pass, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic("Failed to connect: " + err.Error())
	}

	return db
}

// genererate 120 users and insert them into the database
func populateUsers() {
	db := newConn()
	defer db.Close()

	for i := 0; i < 120; i++ {
		_, err := db.Exec("INSERT INTO users (name) VALUES (?)", faker.Name())
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Users populated successfully")
}

func populateSeats() {
	db := newConn()
	defer db.Close()

	// Total 120 seats
	// Seats in 1 row = 6
	// Total rows = 120/6 = 20
	// Seat naming convention  = <row>-A, <row>-B, <row>-C, <row>-D, <row>-E, <row>-F

	tripID := 1

	for row := 1; row <= 20; row++ {
		for j := 0; j < 6; j++ {
			seat := string(rune(65 + j))
			seatName := fmt.Sprintf("%d-%s", row, seat)

			_, err := db.Exec("INSERT INTO seats (name, trip_id) VALUES (?, ?)", seatName, tripID)
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("Trips populated successfully")
}

func main() {
	populateUsers()
	// populate trips
	// INSERT INTO trips (name) VALUES ('6E-6523');
	populateSeats()

	fmt.Println("Database populated successfully")
}
