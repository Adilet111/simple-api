package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "322206"
	TB_NAME     = "movies"
)

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, TB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

type Movie struct {
	MovieID   string `json:"movieid"`
	MovieName string `json:"moviename"`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Movie `json:"data"`
	Message string  `json:"message"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/movies", GetMovies).Methods("GET")
	router.HandleFunc("/hello", func(resp http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(resp, "Hello there!")
	})

	fmt.Println("Server at 8080")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Get all movies

// response and request handlers
func GetMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting movies...")

	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM movies")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var movies []Movie

	// Foreach movie
	for rows.Next() {
		var id int
		var movieID string
		var movieName string

		err = rows.Scan(&id, &movieID, &movieName)

		// check errors
		checkErr(err)

		movies = append(movies, Movie{MovieID: movieID, MovieName: movieName})
	}

	var response = JsonResponse{Type: "success", Data: movies}

	json.NewEncoder(w).Encode(response)
}
