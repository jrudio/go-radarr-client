package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jrudio/go-radarr-client"
)

func main() {
	url := flag.String("url", "", "REQUIRED url pointing to radarr app")
	apikey := flag.String("apikey", "", "REQUIRED apikey for radarr")

	flag.Parse()

	if *url == "" || *apikey == "" {
		// fmt.Println("a url ")
		flag.Usage()
		os.Exit(1)
	}

	client := radarr.New(*url, *apikey)

	results, err := client.Search("lord of the rings")

	if err != nil {
		fmt.Printf("search failed: %v\n", err)
		os.Exit(1)
	}

	for _, result := range results {
		fmt.Println(result.Title)
	}
}
