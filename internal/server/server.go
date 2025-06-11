package server

import (
	"devtrail/internal/storage"
	"log"
	"net/http"
)

func StartServer() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	RegisterRoutes(mux)

	// Запускаємо сервер
	log.Println("Server started at http://127.0.0.1:8080")
	err := http.ListenAndServe("127.0.0.1:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
	storage.CloseDB()
}
