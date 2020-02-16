package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func exampleClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

// func getJson(url string, target interface{}) error {
// 	r, err := myClient.Get(url)
// 	if err != nil {
// 		return err
// 	}
// 	defer r.Body.Close()

// 	return json.NewDecoder(r.Body).Decode(target)
// }

func gitHubFetchAll() JobsGithub {
	var allJobs JobsGithub
	for i := 1; i < 10; i++ {
		// jobs := gitHubFetchPage(i)
		go GitHubFetchPage(i)
		// if len(jobs) > 0 {
		// 	allJobs = append(allJobs, jobs...)
		// } else {
		// 	return allJobs
		// }
	}
	time.Sleep(5 * time.Second)
	return allJobs
}

func setRedis(key string, value string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := rdb.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func delRedis(key string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := rdb.Del(key).Err()
	if err != nil {
		panic(err)
	}
}

func fetchJSON() {
	//Simple Employee JSON object which we will parse
	empJSON := `{
		"id" : 11,
		"name" : "Irshad",
		"department" : "IT",
		"designation" : "Product Manager"
	}`

	// Declared an empty interface
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(empJSON), &result)

	//Reading each value by its key
	fmt.Println("Id :", result["id"],
		"\nName :", result["name"],
		"\nDepartment :", result["department"],
		"\nDesignation :", result["designation"])
}

func main() {
	// pageBody := gitHubFetchPage(1)
	pageBody := gitHubFetchAll()
	fmt.Println(pageBody)
	GithubJobs()
}
