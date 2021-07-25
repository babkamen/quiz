package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/hako/durafmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type quizResult struct {
	totalQuestions int
	correctAnswers int
	acceptAnswers  bool
}

func (qr *quizResult) stopAcceptingAnswers() {
	qr.acceptAnswers = false
}

func (qr *quizResult) correctAnswer() {
	if qr.acceptAnswers {
		qr.correctAnswers++
	}
}

func (qr *quizResult) printResults() {
	fmt.Printf("Correct answers %v out of  %v\n", qr.correctAnswers, qr.totalQuestions)
}

type problem struct {
	question string
	answer   string
}

type problems []problem

func (p problems) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })
}

var (
	problemsFilePath = flag.String("file", "problems.csv", "filepath to csv with questions")
	quizTime         = flag.Duration("quiz-time-limit", 30*time.Second, "quiz total time limit")
	shuffle          = flag.Bool("shuffle", true, "flag that shows where to shuffle questions")
)

func main() {

	problems := readProblemsFile(*problemsFilePath)
	if *shuffle {
		problems.shuffle()
	}

	qr := quizResult{totalQuestions: len(problems), acceptAnswers: true}

	startQuiz(qr.totalQuestions)

	timer := time.NewTimer(*quizTime)
	go func(quizRes *quizResult) {
		<-timer.C
		quizRes.stopAcceptingAnswers()
		quizRes.printResults()
		fmt.Println("Out of time")
		os.Exit(0)
	}(&qr)

	processQuiz(problems, &qr)

	qr.printResults()
}

func startQuiz(totalQuestions int) {
	duration, err := durafmt.ParseString(quizTime.String())
	logFatal("Something went wrong when parsing duration", err)
	fmt.Printf("question contains %v questions. You have %v to answer them\nPress enter to start quiz\n",
		totalQuestions, duration)
	_, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Fatal("Something went wrong when reading input ", err)
	}
}

func logFatal(message string, err error) {
	if err != nil {
		log.Fatal(message, " ", err)
	}
}

func readProblemsFile(filepath string) problems {
	file, err := os.Open(filepath)
	logFatal("Something went wrong when opening the file ", err)

	r := csv.NewReader(file)
	defer file.Close()
	records, err := r.ReadAll()
	logFatal("Something went wrong when reading csv ", err)
	//validate file
	result := make([]problem, len(records))
	for i, r := range records {
		if len(r) != 2 {
			log.Fatalf("Expected csv rows to have 2 values but got %v on line %v ", r, i+1)
		}
		result[i] = problem{r[0], r[1]}
	}

	return result
}

func processQuiz(problems []problem, quizResult *quizResult) {
	for i, p := range problems {
		if strings.Contains(p.question, "?") {
			fmt.Println(p.question)
		} else {
			fmt.Printf("%v. what %v, sir?\n", i+1, p.question)
		}
		answer := getInputFromUser()
		if strings.EqualFold(answer, p.answer) {
			quizResult.correctAnswer()
		}
	}
}

func getInputFromUser() string {
	var result string
	reader := bufio.NewReader(os.Stdin)
	result, err := reader.ReadString('\n')
	logFatal("Something went wrong when reading input ", err)
	result = strings.TrimSpace(result)
	return result
}
