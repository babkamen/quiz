package quiz

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	problemsFilePath = flag.String("file", "problems.csv", "filepath to csv with questions")
	quizTime         = flag.Duration("quiz-time-limit", 30*time.Second, "quiz total time limit")
	shuffle          = flag.Bool("shuffle", true, "flag that shows where to shuffle questions")
)

func main() {
	if contains(os.Args, "--help") {
		printHelp()
	}
	flag.Parse()

	problems := readProblemsFile(*problemsFilePath)
	if *shuffle {
		problems.shuffle()
	}

	qr := quizResult{totalQuestions: len(problems)}

	startQuiz(qr.totalQuestions)

	processQuiz(problems, &qr)

	qr.printResults()
}

func printHelp() {
	fmt.Fprintf(os.Stderr, "usage: %s [flags] <paths...>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}
