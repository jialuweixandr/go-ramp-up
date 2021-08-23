package jokes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"errors"
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


// TODO: make errors consistent
// call the joke api and return
func GetARandomJoke(ch chan JokeResult, joke_type string) {
	if joke_type != "programming" && joke_type != "general" {
		//fmt.Println("Bad joke type! Shold be programming or general")
		ch <- JokeResult{Error: errors.New("Bad joke type! Shold be programming or general")}
	}

	// return value
	var jkres JokeResult = JokeResult{}

	// build url with query
	//query := "general"
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
	if len(jkres_arr) == 0 {
		ch <- JokeResult{Error: errors.New("bad array length 0")}
		os.Exit(1)
	}

	// put result struct into the channel - the asynchronous way of returning(?)
	ch <- jkres_arr[0]
}

// TODO: figure out what to do with joke_type.. maybe make it into an array?
// get n non-repeating jokes, concurrently
func GetRandomJokes(num_jokes int, joke_type string) ([]JokeResult, error){
	// parameters error checking
	if joke_type != "programming" && joke_type != "general" {
		return nil, errors.New("Bad joke type! Shold be programming or general")
	}
	if num_jokes < 0 {
		return nil, errors.New("Invalid num_jokes value!")
	}

	// create a channel
	ch := make(chan JokeResult, num_jokes)
	// return value
	jkres_list := make([]JokeResult, num_jokes)
	// thread-safe dict
	cache := Cache{dict: make(map[int]bool)}

	// asynchrounously call GetARandomJoke. n non-repeaing jokes. A better way to write this?
	cnt := 0
	for {
		if cnt == num_jokes {
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
	return jkres_list, nil
}
