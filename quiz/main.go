package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Error reading file as CSV: ", err)
	}

	return records
}

// func countDown(seconds int, stopGame *bool) {
// 	for i := seconds; i > 0; i-- {
// 		time.Sleep(time.Second)
// 	}
// 	*stopGame = true
// }
// This countdown is very naive, accuracy is low, we need to use channels to communicate between goroutines to send the stop signal instantly setting a flag will only work in the next iteration at the moment
// and use time.Ticker to get accurate time intervals

func countDown(timeLimit int, done chan bool) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	select {
	case <-timer.C:
		fmt.Println("\nTime's up!")
		done <- true
		return
	}
}

func main() {
	records := readCsvFile("problems.csv")
	score := 0
	timeLimit := 5
	done := make(chan bool)

	fmt.Println("Welcome to the Quiz Game!")
	fmt.Println("You will be asked a series of questions and you need to answer them correctly to score points.")
	fmt.Printf("Press enter to start the game, A timer will start and you will have %d seconds to answer all the questions.\n", timeLimit)
	fmt.Scanln()

	go countDown(timeLimit, done)
	questionDone := make(chan bool)

	go func() {
		for _, row := range records {
			fmt.Println("Question: ", row[0])
			answer := ""
			fmt.Print("Your answer: ")
			fmt.Scanln(&answer)

			if answer == row[1] {
				score++
			}
		}
		questionDone <- true
	}()

	select {
	case <-done:
		fmt.Println("Game Over!")
		fmt.Println("Your score is: ", score)

	case <-questionDone:
		fmt.Println("You have answered all the questions!")
		fmt.Println("Your score is: ", score)
	}
}
