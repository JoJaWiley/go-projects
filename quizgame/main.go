package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

func greeting() string {
	return fmt.Sprintf("Welcome to Quiz Game!")
}

func end(score int) string {
	return fmt.Sprintf("You've reached the end of the game with a score of %d!", score)
}

type Question struct {
	Question string
	Answer   string
}

func parseQuestions(file *os.File) []Question {
	// initialize he CSV reader
	reader := csv.NewReader(file)

	var questions []Question

	// Loop through CSV records line by line
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			log.Fatalf("Error reading line: %s", err)
		}

		// map to struct
		question := Question{
			Question: record[0],
			Answer:   record[1],
		}

		// append mapped question to list of questions
		questions = append(questions, question)
	}

	// return the list of questions
	return questions
}

func removeQuestion(i int, questions []Question) []Question {
	// overwrite the removed element and append the rest
	questions = append(questions[:i], questions[i+1:]...)

	// Clear the last element
	questions[len(questions)-1] = Question{Question: "", Answer: ""}

	return questions[:len(questions)-1]
}

func randomizeQuestions(questions []Question) []Question {
	var newQuestions []Question
	for i := range questions {
		// take a random question out and place it at the end, but don't touch any that have already been removed and placed
		n := rand.IntN(len(questions) - i)
		m := questions[n]
		lessQuestions := removeQuestion(n, questions)
		newQuestions = append(lessQuestions, m)
	}
	return newQuestions
}

func main() {
	basePath := "testdata/"

	var fileName string
	flag.StringVar(&fileName, "file", "problems.csv", "The inputted CSV file of questions and answers from the testdata folder.")

	var timer int
	flag.IntVar(&timer, "timer", 30, "Enter the time limit in seconds.")

	var rand bool
	flag.BoolVar(&rand, "rand", false, "Whether the questions are randomized.")
	flag.Parse()

	// open the CSV file
	file, err := os.Open(basePath + fileName)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// return the questions as a list parsed from CSV file
	questions := parseQuestions(file)
	if rand {
		questions = randomizeQuestions(questions)
	}

	score := 0

	gameDuration := time.Duration(timer) * time.Second

	// This channel receives a signal precisely when the time expires
	gameOverTimer := time.After(gameDuration)

	// greet the user
	fmt.Println(greeting())
	for i, question := range questions {
		fmt.Printf("Question #%d: %v?", i+1, question.Question)

		// create a channel to receive the user's answer
		answerCh := make(chan string)

		// run the blocking input scan in a separate goroutine
		go func() {
			var answer string
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				answer = scanner.Text()
			}
			answerCh <- answer
		}()

		// Wait for either the timer to run out or the user to answer
		select {
		case <-gameOverTimer:
			fmt.Println("\nTime's up!")
			fmt.Println(end(score))
			return
		case answer := <-answerCh:
			trimmedAnswer := strings.Trim(answer, "!.,;:\"'! ")
			spaceTrimmedAnswer := strings.TrimSpace(trimmedAnswer)
			lowerAnswer := strings.ToLower(spaceTrimmedAnswer)

			if lowerAnswer == question.Answer {
				score++
			}
		}
	}

	// End the game
	fmt.Println("\nYou answered all questions!")
	fmt.Println(end(score))
}
