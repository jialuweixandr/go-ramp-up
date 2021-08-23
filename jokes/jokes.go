package jokes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)
 

type JokeResult struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
	Error 	  error
}


// A thread-safe map
type Cache struct {
	dict map[int]bool
	mux  sync.Mutex
}


// check if id is in map
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


// call the joke api and return
func GetARandomJoke(ch chan JokeResult, joke_type string) {
	// return value
	var jkres JokeResult = JokeResult{}

	// build url with query
	//query := "general"
	// TODO: put an error check for "type"
	path := fmt.Sprintf("https://official-joke-api.appspot.com/jokes/%s/random", joke_type)
	//fmt.Println("path: ", path)

	// call API
	//response, err := http.Get("https://official-joke-api.appspot.com/random_joke")
	//response, err := http.Get("https://official-joke-api.appspot.com/jokes/programming/random")
	response, err := http.Get(path)
	if err != nil {
		jkres.Error = err
		ch <- jkres
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// parse the response
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		jkres.Error = err
		ch <- jkres
		fmt.Print(err.Error())
		os.Exit(1)
	}


	// json parsing
	var jkres_arr []JokeResult
	json.Unmarshal(responseData, &jkres_arr)

	// put result struct into the channel - the asynchronous way of returning(?)
	ch <- jkres_arr[0]
}

// get n non-repeating jokes, concurrently
func GetRandomJokes(ch chan JokeResult, n int, joke_type string) []JokeResult{
	// TODO: figure out what to do with joke_type.. maybe make it into an array?
	
	// thread-safe dict
	cache := Cache{dict: make(map[int]bool)}

	// return value
	jkres_list := make([]JokeResult, n)

	// breaks only when gathered n non-repeaing jokes. A better way to write this?
	cnt := 0
	for {
		if cnt == n {
			break
		}

		// call go routine and error check
		go GetARandomJoke(ch, joke_type)
		jkres := <-ch
		if jkres.Error != nil {
			fmt.Print(jkres.Error.Error())
			os.Exit(1)
		}

		// check duplicate. increment only when it's a new joke
		if cache.checkVisited((jkres.Id)) {
			fmt.Println("old joke... fetching you a new one")
			continue
		} else {
			jkres_list[cnt] = jkres
			// fmt.Println((cnt))
			// fmt.Println("	Setup: ", jkres_list[cnt].Setup)
			// fmt.Println("	Punchline: ", jkres_list[cnt].Punchline)
			cnt += 1
		}
	}
	return jkres_list
}
