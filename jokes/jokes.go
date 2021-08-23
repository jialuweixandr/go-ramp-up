package jokes

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

func GetARandomJoke(ch chan Response) {
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


func GetRandomJokes(ch chan Response, n int) []Response{
	cache := Cache{dict: make(map[int]bool)}
	jokes := make([]Response, n)
	cnt := 0
	for {
		go GetARandomJoke(ch)

		if cnt == n {
			break
		}

		joke := <-ch

		if cache.checkVisited((joke.Id)) {
			fmt.Println("old joke... fetching you a new one")
		} else {
			jokes[cnt] = joke
			fmt.Println((cnt))
			fmt.Println("	Setup: ", jokes[cnt].Setup)
			fmt.Println("	Punchline: ", jokes[cnt].Punchline)
			cnt += 1
		}
	}
	return jokes
}
