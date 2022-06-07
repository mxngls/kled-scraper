package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Initialize used variables
var Words []Result
var idArr []string

// Custom func that creates JSON from every data type and every rune encoding
func createEncoder(w io.Writer) (enocer *json.Encoder) {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder
}

type ByIndex []Result

func main() {

	// Initiate the http client
	client, err := createClient()
	if err != nil {
		panic(err)
	}

	// Read the word ids
	data, _ := ioutil.ReadFile("dict/wordIDs.json")
	err = json.Unmarshal(data, &idArr)
	if err != nil {
		panic(err)
	}

	// Open the file where to save our dictionary entries
	// If no such file exits yet create one
	f, err := os.OpenFile("dict/dict_FULL.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Create a custom JSON encoder
	encoder := createEncoder(f)

	// Set the number of goroutines that should be run
	start := 0
	end := start + 100

	// Initiate Channel
	c := make(chan Result)

	// f.WriteString("[")

	currIndex := start

	// Loop over all the ids and write them to file
	for {

		// Set the end point of the last loop the the number of word ids
		if end > len(idArr) {
			end = len(idArr)
		}

		// Start looping
		// In order to provide enough time to fetch a response wait for 50ms after
		// a new goroutine is run
		for i := start; i < end; i++ {
			id := string(idArr[i])
			time.Sleep(time.Millisecond * 30)
			go getView(id, c, client)

			// fmt.Printf("\rFetched and parsed %d from %d entries; running...", i+1, len(idArr))
		}

		arr := &[]Result{}

		// Store the fetched and parsed response bodies to file
		for i := start; i < end; i++ {
			*arr = append(*arr, <-c)
			currIndex++
		}

		encoder.Encode(*arr)
		encoder.Encode(",")

		fmt.Printf("\rCurrent Write: %d (Id: %d) (of %d writes in total); running...", currIndex, (*arr)[len(*arr)-1].Id, len(idArr))

		// Break out of the loop when the end point is equal to the number of word ids
		if end == 200 {

			f.WriteString("]")

			f.Close()

			fmt.Printf("\rWrote %d from %d entries to file: 'dict_FULL.JSON'; finished", currIndex, len(idArr))

			break
		}

		// Increment the start point
		start += 100
		end += 100
	}
}
