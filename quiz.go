package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
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
	quiz, err := loadProblems()

	if err == nil {
		askQuiz(quiz)
	}
}

func loadProblems() (*Quiz, error) {
	file, err := os.Open("problems.csv")
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

func askQuiz(quiz *Quiz) {
	fmt.Println("Answer the following questions")

	for _, problem := range quiz.problems {
		var ans string
		fmt.Print(problem.question, " = ")
		_, _ = fmt.Scanln(&ans)

		if problem.answer == ans {
			quiz.totalScore++
		}
	}

	fmt.Println("Out of: ", len(quiz.problems), ", your total is: ", quiz.totalScore)
}
