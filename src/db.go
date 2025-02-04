package main

import (	
	"fmt"
	"log"
	"database/sql"
	"os"
	
	_ "github.com/lib/pq"
)


func ConfigureDB() *sql.DB {
	var dbtest *sql.DB
	return dbtest
	//connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
			log.Fatal(err)
	}

	//create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id NUMERIC PRIMARY KEY, password TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS persons (id NUMERIC PRIMARY KEY, name TEXT, email TEXT, dni TEXT, phone TEXT, user_id NUMERIC REFERENCES users (id))")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS transactions (id NUMERIC PRIMARY KEY, description TEXT, buyMoney boolean, money NUMERIC, owner_id NUMERIC REFERENCES persons (id), match_id NUMERIC REFERENCES persons (id))")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func SeedDB(db *sql.DB){
	fmt.Println("Populate DB")
	//populate
	_, err := db.Exec("INSERT INTO users (id, password, email) VALUES (1, 'pass1','prueba1@hola.cl') ON CONFLICT DO NOTHING")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO users (id, password, email) VALUES (2, 'pass2','prueba2@hola.cl') ON CONFLICT DO NOTHING")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO persons (id, name, email, dni, phone, user_id) VALUES (1,'Usuario Prueba 1','prueba1@hola.cl','XXXXXXXX-X','999999999',1) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO persons (id, name, email, dni, phone, user_id) VALUES (2,'Usuario Prueba 2','prueba2@hola.cl','AXXXXXXX-X','888888888',2) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO transactions (id, description, buyMoney, money, owner_id, match_id) VALUES (1,'Transacción 1',TRUE,10000,1,NULL) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO transactions (id, description, buyMoney, money, owner_id, match_id) VALUES (2,'Transacción 2',FALSE,15000,1,2) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO transactions (id, description, buyMoney, money, owner_id, match_id) VALUES (3,'Transacción 3',TRUE,20000,2,NULL) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Fatal(err)
	}


}


func getUserByEmail(email string, db *sql.DB) UserLoginData {
	var u UserLoginData
	err := db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", email).Scan(&u.Id, &u.Email, &u.Password)
	if err != nil {
		log.Fatal(err)
	}
	
	return u
}

func getTransactions(db *sql.DB) []Transaction {

	rows, err := db.Query("SELECT id, description, buyMoney, money, owner_id, match_id FROM transactions")
	if err != nil {
			log.Fatal(err)
	}
	defer rows.Close()

	var Owner_Id sql.NullFloat64
	var Match_Id sql.NullFloat64
	transactions := []Transaction{}
	for rows.Next() {
			var t Transaction
			if err := rows.Scan(&t.Id, &t.Description, &t.BuyMoney, &t.Money, &Owner_Id, &Match_Id); err != nil {
					log.Fatal(err)
			}
			if Owner_Id.Valid {
        t.Owner_Id = Owner_Id.Float64
			}
			if Match_Id.Valid {
					t.Match_Id = Match_Id.Float64
			}
			transactions = append(transactions, t)
	}

	return transactions
}