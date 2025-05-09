package services

import "errors"

func SearchWord(word string, invIndex map[string]map[string]int) (map[string]int, error) {
	val, ok := invIndex[word]

	if ok {
		return val, nil
	} else {
		return make(map[string]int), errors.New("searched word not found")
	}
}
