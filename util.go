package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}


func fetchURL(URL string){
	resp, err := http.Get(URL)

	//if there was an error, report it and exit
	if err != nil {
		//.Fatalf() prints the error and exits the process
		log.Fatalf("error fetching URL: %v\n", err)
	}

	//make sure the response body gets closed
	defer resp.Body.Close()
	//check response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("response status code was %d\n", resp.StatusCode)
	}

	//check response content type
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		log.Fatalf("response content type was %s not text/html\n", ctype)
	}
}