package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "mydata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS notes(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	note := "Мій перший запис у базі"
	_, err = db.Exec("INSERT INTO notes (content) VALUES (?)", note)
	if err != nil {
		log.Fatal("Помилка додавання запису:", err)
	}
	fmt.Println("Запис додано!")

	// Читання всіх записів
	rows, err := db.Query("SELECT id, content FROM notes")
	if err != nil {
		log.Fatal("Помилка читання записів:", err)
	}
	defer rows.Close()

	fmt.Println("Усі записи:")
	for rows.Next() {
		var id int
		var content string
		rows.Scan(&id, &content)
		fmt.Printf("[%d] %s\n", id, content)
	}
}
