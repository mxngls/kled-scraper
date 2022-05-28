package main

import (
	"fmt"
	"net/http"
)

func requestView(id string, lang string, langCode string, client *http.Client) (resp *http.Response) {

	url := fmt.Sprintf("https://krdict.korean.go.kr/%s/dicSearch/SearchView?ParaWordNo=%s&nation=%s&nationCode=%s&viewType=A&blockCount=10&viewTypes=on",
		lang,
		id,
		lang,
		langCode,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}
