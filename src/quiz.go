package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("problems.csv")
	if err != nil {
		log.Fatal("Something went wrong when opening the file", err)
		os.Exit(1)
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()

	if err != nil {
		log.Fatal("Something went wrong when reading csv", err)
		os.Exit(1)
	}

	total := len(records)
	var correctCount int

	for _, r := range records {
		if len(r) != 2 {
			log.Fatal("Expected csv rows to have 2 values but got ", r)
			os.Exit(1)
		}
		exercise := r[0]
		correctAnswer := r[1]
		if strings.Contains(exercise, "?") {
			fmt.Println(exercise)
		} else {
			fmt.Printf("what %v, sir?\n", exercise)
		}
		var answer string

		// Taking input from user
		fmt.Scanln(&answer)
		if answer == correctAnswer {
			correctCount++
		}
	}
	fmt.Printf("Correct answers %v out of  %v\n", correctCount, total)
}
