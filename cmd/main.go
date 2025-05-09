package main

import (
	"fmt"
	"sync"
	"webScraper/internal/services"
)

func main() {
	url := "https://usf-cs272-s25.github.io/top10/"
	fmt.Println("Starting scraper...")

	htmlChan := make(chan services.HTMLPage)
	invIndex := make(map[string]map[string]int)

	crawler := services.NewCrawlService(url, htmlChan)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		crawler.Start()
	}()

	go func() {
		defer wg.Done()
		for page := range htmlChan {
			services.ExtractService(page, &invIndex)
		}
	}()

	wg.Wait()

	for {
		fmt.Println("1. Do you want to search a word\n2. Do you want to exit?\n(Enter 1/2): ")
		var choice int
		fmt.Scanln(&choice)

		if choice == 2 {
			fmt.Println("Exiting the program...")
			break
		} else if choice == 1 {
			fmt.Println("Enter the word you want to search: ")
			var word string
			fmt.Scanln(&word)

			cleanWord := services.CleanWord(word)
			result, err := services.SearchWord(cleanWord, invIndex)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Searched for the word:", cleanWord, "Result is as follows:")
				for k, v := range result {
					fmt.Printf("Link: %s, Count: %d\n", k, v)
				}
			}
		} else {
			fmt.Println("Invalid choice. Please enter either choice 1 or 2.")
		}
	}
}
