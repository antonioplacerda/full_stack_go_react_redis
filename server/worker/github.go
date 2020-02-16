package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

// JobsGithubJSON ...
type JobsGithubJSON []struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	CreatedAt   string `json:"created_at"`
	Company     string `json:"company"`
	CompanyURL  string `json:"company_url"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Description string `json:"description"`
	HowToApply  string `json:"how_to_apply"`
	CompanyLogo string `json:"company_logo"`
}

// Job ...
type Job struct {
	source      string `json:"source"`
	URL         string `json:"url"`
	CreatedAt   string `json:"createdAt"`
	Company     string `json:"company"`
	CompanyURL  string `json:"companyURL"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Result ...
type Result struct {
	page     int
	results  int
	filtered int
	jobs     []Job
}

var jobs = make(chan int, 10)
var results = make(chan Result, 10)

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

// GitHubFetchPage ...
func GitHubFetchPage(page int) (int, JobsGithubJSON) {
	baseURL := "https://jobs.github.com/positions.json"
	res, err := http.Get(baseURL + "?page=" + strconv.Itoa(page))
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var target JobsGithubJSON
	json.Unmarshal(body, &target)

	return len(target), target
}

func filterGitHubJobs(allJobs JobsGithubJSON) (int, []Job) {
	var filteredJobs []Job
	for _, el := range allJobs {
		if !(strings.Contains(strings.ToLower(el.Description), "sr ") ||
			strings.Contains(strings.ToLower(el.Description), "sr. ") ||
			strings.Contains(strings.ToLower(el.Description), "senior ") ||
			strings.Contains(strings.ToLower(el.Description), "architect ") ||
			strings.Contains(strings.ToLower(el.Title), "sr ") ||
			strings.Contains(strings.ToLower(el.Title), "sr. ") ||
			strings.Contains(strings.ToLower(el.Title), "senior ") ||
			strings.Contains(strings.ToLower(el.Title), "architect ")) {

			filteredJobs = append(filteredJobs, Job{
				"github",
				el.URL,
				el.CreatedAt,
				el.Company,
				el.CompanyURL,
				el.Location,
				el.Title,
				el.Description})
		}
	}
	return len(filteredJobs), filteredJobs
}

func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		hits, resJSON := GitHubFetchPage(job)
		filteredHits, filteredJobs := filterGitHubJobs(resJSON)
		output := Result{job, hits, filteredHits, filteredJobs}
		results <- output
	}
	wg.Done()
}
func createWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()
	close(results)
}
func allocate(noOfJobs int) {
	for i := 1; i < noOfJobs; i++ {
		jobs <- i
	}
	close(jobs)
}
func result(done chan bool) {
	var allResults []Job
	for result := range results {
		fmt.Printf("Page %d -> num of results %d -> num of filtered %d\n", result.page, result.results, result.filtered)
		allResults = append(allResults, result.jobs...)
	}

	fmt.Println("Total results %d", len(allResults))

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(allResults)
	out := string(buf.String())

	setRedis("github", string(out))

	done <- true
}

// GithubJobs ...
func GithubJobs() {
	startTime := time.Now()
	noOfPages := 10
	go allocate(noOfPages)

	done := make(chan bool)
	go result(done)

	noOfWorkers := 10
	createWorkerPool(noOfWorkers)
	<-done

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds")
}

func main() {
	GithubJobs()
}
