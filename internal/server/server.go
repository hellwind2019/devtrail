package server

import (
	"log"
	"net/http"
)

func StartServer() {
	mux := http.NewServeMux()

	RegisterRoutes(mux)

	// Запускаємо сервер
	log.Println("Сервер стартує на http://127.0.0.1:8080")
	err := http.ListenAndServe("127.0.0.1:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
