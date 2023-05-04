package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/gorilla/mux"
	"syreclabs.com/go/faker"
)

const (
	appServingPort string = "8000"
)

type crawlResult struct {
	url    string
	status int
}

type crawlErr struct {
	url string
	err error
}

func (c crawlResult) String() string {
	return fmt.Sprintf("url: %s, resp status: %d", c.url, c.status)
}

func (c crawlErr) String() string {
	return fmt.Sprintf("url: %s, err: %v", c.url, c.err)
}

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

func crawlPages() {

	var (
		urls = []string{
			"http://www.golang.org",
			"http://www.google.com",
			"http://www.somefakewebpageurl.xyz",
		}
		wg    sync.WaitGroup
		resCh = make(chan crawlResult)
		errCh = make(chan crawlErr)
	)

	wg.Add(len(urls))
	for _, url := range urls {
		go queryUrl(url, resCh, errCh, &wg)
	}

	for i := 0; i < len(urls); i++ {
		select {
		case result := <-resCh:
			log.Println(result)
		case errResult := <-errCh:
			log.Println(errResult)
		}
	}

	close(resCh)
	close(errCh)

	log.Print("done!")
}

func queryUrl(url string, ch chan<- crawlResult, errCh chan<- crawlErr, wg *sync.WaitGroup) {
	var (
		res    = crawlResult{url: url}
		client = http.Client{Timeout: time.Second * 2}
		req, _ = http.NewRequest("GET", url, nil)
		err    error
	)

	defer wg.Done()

	log.Printf("querying: %s", url)
	resp, err := client.Do(req)
	if err != nil {
		errCh <- crawlErr{url: url, err: err}
		return
	}

	if resp == nil {
		errCh <- crawlErr{url: url, err: fmt.Errorf("empty response")}
		return
	}

	defer resp.Body.Close()

	res.status = resp.StatusCode
	ch <- res
}

func capitalizeString(s string) string {
	const sep = " "

	if s == "" {
		return s
	}

	parts := strings.Split(strings.ToLower(s), sep)
	runes := []rune(parts[0])
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
		parts[0] = string(runes)
	}

	return strings.Join(parts, sep)
}
func testStrings() {
	var (
		ss = "hello"
	)

	log.Printf("-- ss --")
	for _, v := range ss {
		log.Printf("%T, %v", v, string(v))
	}

	runes := []rune(ss)
	runes[0] = 'H'
	ss = string(runes)

	log.Printf("ss: %s", ss)
}

type person struct {
	id   int
	name string
}

func testSlices() {
	var (
		p1 = person{id: 1, name: "Bob"}
		p2 = person{id: 2, name: "John"}
		p3 = person{id: 2, name: "Ann"}
		s  = []person{p2, p1, p3}
	)

	sort.SliceStable(s, func(i, j int) bool { return s[i].name < s[j].name })

	log.Printf("slice: %v", s)

	//alterSlice(s)
	//log.Print(s)
}

func appendToSlice(s []person, values ...person) {
	s = append(s, values...)
}

func alterSlice(s []person) {
	s[0].name = "Peter"
}

func insideDefer() {
	defer func() { log.Printf("defer inside 1") }()
	defer func() { log.Printf("defer inside 2") }()
}

func testDefer() {
	defer func() { log.Printf("defer 1") }()
	defer func() { log.Printf("defer 2") }()
	defer func() { log.Printf("defer 3") }()
	defer func() { log.Printf("defer 4") }()

	insideDefer()
}

func testContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t := time.After(time.Second)

	for {
		select {
		case <-ctx.Done():
			log.Print("done")
			return
		case <-t:
			log.Print("timed")
			return
		default:
			log.Print("working")
		}
	}
}

func testWg() {
	var (
		wg sync.WaitGroup
		s  = []string{"1", "2", "3", "4"}
	)

	wg.Add(len(s))

	for _, v := range s {
		vv := v
		go func(value string, wg *sync.WaitGroup) {
			defer wg.Done()
			log.Printf("processing %s", vv)
			time.Sleep(time.Second * 1)
		}(vv, &wg)
	}

	wg.Wait()

	log.Print("done!")
}

func testSorting() {
	var s = []int{1, 2, 3, 10, -7, 100}
	sort.Ints(s)

	log.Print(s)

}

func testArrays() {
	var (
		first  = [...]int{1, 2, 3}
		second = [...]int{4, 5, 6, 7, 8}
		third  = [len(first) + len(second)]int{}
	)

	for idx, v := range second {
		third[idx] = v
	}

	log.Printf("third: %v", third)

}

func main() {
	//serveHttp()
	//crawlPages()
	//testStrings()
	//testSlices()
	//testDefer()
	//testContext()
	//testWg()
	//testSorting()
	testArrays()
}
