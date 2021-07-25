package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/hako/durafmt"
	"log"
	"os"
	"strings"
	"time"
)

type QuizResult struct {
	totalQuestions int
	correctAnswers int
}

var quizTime time.Duration

func init() {
	flag.DurationVar(&quizTime, "quiz-time", 30*time.Second, "quiz total time")
	flag.Parse()
}

func main() {

	records := readProblemsFile()

	quizResult := QuizResult{totalQuestions: len(records)}

	startQuiz(quizResult.totalQuestions)

	timer := time.NewTimer(quizTime)
	go func(quizResult *QuizResult) {
		<-timer.C
		quizResult.printResults()
		os.Exit(0)
	}(&quizResult)

	processQuiz(records, &quizResult)

	quizResult.printResults()
}

func startQuiz(totalQuestions int) {
	duration, err := durafmt.ParseString(quizTime.String())
	if err != nil {
		log.Fatal("Something went wrong when parsing duration ", err)
	}
	fmt.Printf("Quiz contains %v questions. You have %v to answer them\nPress enter to start quiz\n",
		totalQuestions, duration)
	_, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Fatal("Something went wrong when reading input ", err)
	}
}

func readProblemsFile() [][]string {
	file, err := os.Open("problems.csv")
	if err != nil {
		log.Fatal("Something went wrong when opening the file ", err)
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()

	if err != nil {
		log.Fatal("Something went wrong when reading csv ", err)
	}
	//validate file
	for i, r := range records {
		if len(r) != 2 {
			log.Fatalf("Expected csv rows to have 2 values but got %v on line %v ", r, i+1)
		}
	}

	return records
}

func processQuiz(records [][]string, quizResult *QuizResult) {
	for i, r := range records {
		exercise := r[0]
		correctAnswer := r[1]
		if strings.Contains(exercise, "?") {
			fmt.Println(exercise)
		} else {
			fmt.Printf("%v. what %v, sir?\n", i+1, exercise)
		}
		answer := getInputFromUser()
		if strings.EqualFold(answer, correctAnswer) {
			quizResult.correctAnswers++
		}
	}
}

func (quizResult QuizResult) printResults() {
	fmt.Printf("Correct answers %v out of  %v\n", quizResult.correctAnswers, quizResult.totalQuestions)
}

func getInputFromUser() string {
	var result string
	reader := bufio.NewReader(os.Stdin)
	result, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Something went wrong when reading input ", err)
	}
	result = strings.TrimSpace(result)
	fmt.Printf("User entered %q\n", result)
	return result
}
