package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// assignSeat assigns a seat to a user
func assignSeat(db *sql.DB, user User, tripID int) {
	// Assign seat to user
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Find a seat which is not assigned to anyone
	// and assign it to the user

	// Approach 1
	query := `SELECT id, name, trip_id, assigned_to FROM seats 
	WHERE trip_id = ? AND assigned_to IS NULL
	ORDER BY id LIMIT 1`

	// Approach 2
	/*
		query := `SELECT id, name, trip_id, assigned_to FROM seats
		WHERE trip_id = ? AND assigned_to IS NULL
		ORDER BY id LIMIT 1 FOR UPDATE`
	*/

	// Approach 3
	/*
		query := `SELECT id, name, trip_id, assigned_to FROM seats
		WHERE trip_id = ? AND assigned_to IS NULL
		ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED`
	*/

	var seat Seat
	err = txn.QueryRow(query, tripID).Scan(&seat.ID, &seat.Name, &seat.TripID, &seat.AssignedTo)
	if err != nil {
		log.Fatal(err)

	}

	_, err = txn.Exec("UPDATE seats SET assigned_to = ? WHERE id = ?", user.ID, seat.ID)
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s was assigned seat %s\n", user.Name, seat.Name)
}

func main() {
	start := time.Now()

	db := newConn()
	defer db.Close()

	users, err := fetchUsers(db)
	if err != nil {
		panic(err)
	}

	tripID := 1

	var wg sync.WaitGroup

	for _, u := range users {
		wg.Add(1)

		go func(user User) {
			assignSeat(db, user, tripID)

			wg.Done()
		}(u)
	}

	wg.Wait()

	printSeatStatus(db, tripID)

	fmt.Printf("\nTime taken: %v\n", time.Since(start))
}

type User struct {
	ID   int
	Name string
}

type Trip struct {
	ID   int
	Name string
}

type Seat struct {
	ID         int
	Name       string
	TripID     int
	AssignedTo sql.NullInt64
}

func printSeatStatus(db *sql.DB, tripID int) {
	// 120 seats, 20 rows with 6 seats each
	// displaying seat layout horizontally

	fmt.Println("\nSeat layout:")
	for col := 0; col < 6; col++ {
		for row := 0; row < 20; row++ {
			seatCol := string(rune(65 + col))
			seatName := fmt.Sprintf("%d-%s", row+1, seatCol)

			// query to get all for this seat
			var assignedTo sql.NullInt64
			query := "SELECT assigned_to FROM seats WHERE trip_id = ? AND name = ?"
			err := db.QueryRow(query, tripID, seatName).Scan(&assignedTo)
			if err != nil {
				if err == sql.ErrNoRows {
					fmt.Println("No rows found")
				} else {
					log.Fatal(err)
				}
			}

			// Check if the column is NULL
			if assignedTo.Valid {
				fmt.Printf("X ")
			} else {
				fmt.Printf("- ")
			}
		}

		fmt.Println()
	}
}

func fetchUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

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
