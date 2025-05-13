package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func StartServer() {
	godotenv.Load()
	var sessionKey = []byte(os.Getenv("SESSION_KEY"))
	fmt.Println(sessionKey)
	mux := http.NewServeMux()

	RegisterRoutes(mux)

	// Запускаємо сервер
	log.Println("Сервер стартує на http://127.0.0.1:8080")
	err := http.ListenAndServe("127.0.0.1:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
