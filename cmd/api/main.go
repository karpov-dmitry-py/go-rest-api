package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"syreclabs.com/go/faker"
)

const (
	appServingPort string = "8000"
)

type Article struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Author   Author `json:"author"`
	Views    int    `json:"views"`
	Comments int    `json:"comments"`
}

type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

var (
	authors = []Author{
		{ID: faker.RandomInt64(1, 1000), Name: faker.RandomString(30)},
		{ID: faker.RandomInt64(1, 1000), Name: faker.RandomString(30)},
		{ID: faker.RandomInt64(1, 1000), Name: faker.RandomString(30)},
		{ID: faker.RandomInt64(1, 1000), Name: faker.RandomString(30)},
		{ID: faker.RandomInt64(1, 1000), Name: faker.RandomString(30)},
	}

	articles = []Article{
		{
			ID:       faker.RandomInt64(1, 1000),
			Name:     faker.RandomString(100),
			Author:   authors[faker.RandomInt(0, len(authors)-1)],
			Views:    faker.RandomInt(50, 1000),
			Comments: faker.RandomInt(50, 1000),
		},
		{
			ID:       faker.RandomInt64(1, 1000),
			Name:     faker.RandomString(100),
			Author:   authors[faker.RandomInt(0, len(authors)-1)],
			Views:    faker.RandomInt(50, 1000),
			Comments: faker.RandomInt(50, 1000),
		},
		{
			ID:       faker.RandomInt64(1, 1000),
			Name:     faker.RandomString(100),
			Author:   authors[faker.RandomInt(0, len(authors)-1)],
			Views:    faker.RandomInt(50, 1000),
			Comments: faker.RandomInt(50, 1000),
		},
		{
			ID:       faker.RandomInt64(1, 1000),
			Name:     faker.RandomString(100),
			Author:   authors[faker.RandomInt(0, len(authors)-1)],
			Views:    faker.RandomInt(50, 1000),
			Comments: faker.RandomInt(50, 1000),
		},
	}

	articlesDict = map[string]any{
		"items":       articles,
		"total_count": len(articles),
	}
)

func serveHttp() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(respContentTypeMiddleware)

	router.HandleFunc("/", getHealthCheck).Methods("GET")
	router.HandleFunc("/articles", listArticles).Methods("GET")

	log.Printf("serving app on port %s", appServingPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appServingPort), router))
}

func respContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getHealthCheck(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func listArticles(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(articlesDict)
}

func main() {
	serveHttp()
}
