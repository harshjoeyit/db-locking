package main

import (
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
)

func newConn() *sql.DB {
	user := "root"
	pass := "mysql"
	dbName := "test"

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

var DB *sql.DB

func transferMoney(from_user int, to_user int, amount int) {
	txn, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	q1 := `UPDATE account SET balance = balance - ? WHERE user_id = ?`

	_, err = txn.Exec(q1, amount, from_user)
	if err != nil {
		txn.Rollback()
		log.Fatal(err)
	}

	q2 := `UPDATE account SET balance = balance + ? WHERE user_id = ?`
	_, err = txn.Exec(q2, amount, to_user)
	if err != nil {
		txn.Rollback()
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func displayTotalBalance() {
	row := DB.QueryRow("SELECT SUM(balance) FROM account")
	var totalBalance int
	err := row.Scan(&totalBalance)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total Balance: %d\n", totalBalance)
}

func insertData() {
	// Truncate the account table
	_, err := DB.Exec("TRUNCATE TABLE account")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < N; i++ {
		_, err := DB.Exec("INSERT INTO account (balance) VALUES (?)", 1000)
		if err != nil {
			log.Fatal(err)
		}
	}
}

var N int

func main() {
	DB = newConn()

	N = 1000

	insertData()

	type Transfer struct {
		from   int
		to     int
		amount int
	}

	var tt []Transfer

	// Total N transfers
	for i := 1; i <= N; i++ {
		tt = append(tt, Transfer{1, i, 1})
	}

	displayTotalBalance()

	var wg sync.WaitGroup
	sem := make(chan struct{}, 50) // Limit to 10 concurrent goroutines

	for _, t := range tt {
		wg.Add(1)
		sem <- struct{}{} // Acquire a token

		go func(f int, t int, a int) {
			defer wg.Done()
			defer func() { <-sem }() // Release the token

			transferMoney(f, t, a)
		}(t.from, t.to, t.amount)
	}

	wg.Wait()

	displayTotalBalance()
}
