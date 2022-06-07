package main

import (
	"bytes"
	"io"
	"net/http"
)

func getView(index int, id string, channel chan Result, client *http.Client) (err error) {

	resp := requestView(id, "eng", "6", client)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	resp.Body.Close()

	reader := bytes.NewReader(body)

	data, err := ParseView(reader, id, "6")
	if err != nil {
		panic(err)
	}

	data.Alpha = index

	channel <- data

	return err
}
