package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"a21hc3NpZ25tZW50/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Initialize the services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}

// var store = sessions.NewCookieStore([]byte("my-key"))

// func getSession(r *http.Request) *sessions.Session {
// 	session, _ := store.Get(r, "chat-session")
// 	return session
// }

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the Hugging Face token from the environment variables
	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
	}

	// Set up the router
	router := mux.NewRouter()

	// File upload endpoint
	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		if !strings.HasSuffix(handler.Filename, ".csv") {
			http.Error(w, "Only .csv files are allowed", http.StatusInternalServerError)
			return
		}

		// Membaca file content
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file); err != nil {
			http.Error(w, "Failed to read file content", http.StatusInternalServerError)
			return
		}
		fileContent := buf.String()

		// process file
		parsedData, err := fileService.ProcessFile(fileContent)
		if err != nil {
			http.Error(w, "Error processing file", http.StatusInternalServerError)
			return
		}

		queries := []string{
			"Find the least electricity usage appliance.",
			"Find the most electricity usage appliance.",
		}

		// analisi data
		answer, err := aiService.AnalyzeFile(parsedData, queries, token)
		if err != nil {
			http.Error(w, "Failed to analyze data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"status": "success",
			"answer": answer,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}).Methods("POST")

	// Chat endpoint
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		// TODO: answer here
	}).Methods("POST")

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow your React app's origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
