package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// // Jobs ...
// type Jobs struct {

// }

// Job1 ...
type Job1 struct {
	source      string `json:"source"`
	URL         string `json:"url"`
	CreatedAt   string `json:"createdAt"`
	Company     string `json:"company"`
	CompanyURL  string `json:"companyURL"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func fetchFromDB() []Job1 {
	var jobs []Job1
	redisStr := getRedis("github")
	if err := json.Unmarshal([]byte(redisStr), &jobs); err != nil {
		panic(err)
	}
	return jobs
}

func getAllJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(fetchFromDB())
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hey Mary!")
}

func getRedis(key string) string {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	val, err := rdb.Get(key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/jobs", getAllJobs).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(router)))
}
