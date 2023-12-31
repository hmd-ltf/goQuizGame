package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Problem struct {
	question string
	answer   string
}

type Quiz struct {
	problems   []*Problem
	totalScore int
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "A csv file in format of question,answer (default = problem.csv)")
	timeLimit := flag.Int("limit", 30, "Time limit of the quiz in seconds (default = 30)")
	flag.Parse()

	quiz, err := loadQuiz(*csvFileName)

	if err == nil {
		startQuiz(timeLimit, quiz)
		fmt.Printf("\nOut of %d questions, your got %d correct\n", len(quiz.problems), quiz.totalScore)
	}
}

func loadQuiz(fileName string) (*Quiz, error) {
	file, err := os.Open(fileName)
	var quiz *Quiz

	if err != nil {
		log.Fatal("Error while reading the file", err)
		return quiz, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Error while closing the file", err)
		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		log.Fatal("Error while reading records", err)
		return quiz, err
	}

	var problems []*Problem

	for _, eachRecord := range records {
		problems = append(problems, &Problem{eachRecord[0], eachRecord[1]})
	}

	quiz = &Quiz{problems, 0}

	return quiz, nil
}

func startQuiz(timeLimit *int, quiz *Quiz) {
	fmt.Printf("You will have %d seconds to answer all questions", *timeLimit)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	scoreChan := make(chan int)
	go askQuestions(quiz, scoreChan)

	select {
	case <-timer.C:
		fmt.Println("Oops it seams the time ended")
	case <-scoreChan:
		fmt.Println("Well done you answer all questions")
	}
}

func askQuestions(quiz *Quiz, scoreChan chan int) {
	fmt.Println("Answer the following questions")

	for _, problem := range quiz.problems {
		var ans string
		fmt.Printf("%s = ", problem.question)
		_, _ = fmt.Scanln(&ans)

		if strings.EqualFold(strings.TrimSpace(problem.answer), strings.TrimSpace(ans)) {
			quiz.totalScore++
		}
	}

	scoreChan <- quiz.totalScore
}
