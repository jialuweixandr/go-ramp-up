package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

// A Response struct to map the Entire Response
type Response struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

type Cache struct {
	dict map[int]bool
	mux  sync.Mutex
}

func (c *Cache) checkVisited(id int) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.dict[id]
	if ok == false {
		c.dict[id] = true
		return false
	}
	return true
}

func getARandomJoke(ch chan Response) {
	response, err := http.Get("https://official-joke-api.appspot.com/random_joke")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	var responseObject Response
	json.Unmarshal(responseData, &responseObject)
	//fmt.Println("Here? responseObj: ", responseObject)
	ch <- responseObject
}

// input: take number n from command line
// output: a list of n jokes: Setup and punchline (format tbd)
func main() {

	var n int
	fmt.Print("Enter a number: ")
	_, err := fmt.Scanf("%d", &n)
	fmt.Println("Generating ", n, " jokes...")
	fmt.Println("------")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	jokes := make([]Response, n)

	ch := make(chan Response, n)
	cache := Cache{dict: make(map[int]bool)}

	cnt := 0
	for {
		go getARandomJoke(ch)

		if cnt == n {
			break
		}

		res := <-ch

		if cache.checkVisited((res.Id)) {
			fmt.Println("old joke... fetching you a new one")
		} else {
			jokes[cnt] = res
			fmt.Println((cnt))
			fmt.Println("	Setup: ", jokes[cnt].Setup)
			fmt.Println("	Punchline: ", jokes[cnt].Punchline)
			cnt += 1
		}
	}

	// for i := 0; i < n; i++ {
	// 	go getARandomJoke(ch)

	// 	res := <-ch
	// 	if cache.checkVisited((res.Id)) {
	// 		fmt.Printf("old joke....")
	// 		continue
	// 	}
	// 	jokes[i] = res

	// 	fmt.Println(i + 1)
	// 	fmt.Println("	Setup: ", jokes[i].Setup)
	// 	fmt.Println("	Punchline: ", jokes[i].Punchline)
	// }
}
