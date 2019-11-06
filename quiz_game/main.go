package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var scvFile string
var quizTimeout int

var results []answer

type answer struct {
	question      string
	userAnswer    string
	correctAnswer string
}

func init() {
	flag.StringVar(&scvFile, "scv", "./problems.csv", "File with problems")
	flag.IntVar(&quizTimeout, "limit", 30, "Quiz time")
	flag.Parse()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func checkCSVRecord(record []string) {
	if len(record) != 2 {
		panic(errors.New("Invalid scv"))
	}
}
func main() {
	f, err := ioutil.ReadFile(scvFile)
	checkErr(err)

	reader := readAnswers(string(f))
	ticker := time.NewTicker(time.Second * time.Duration(quizTimeout))
loop:
	for {
		select {
		case ans, ok := <-reader:
			if !ok {
				break loop
			}
			results = append(results, ans)
		case <-ticker.C:
			ticker.Stop()
			break loop
		}
	}

	var score int
	for _, s := range results {
		if s.userAnswer == s.correctAnswer {
			score++
		}
	}
	fmt.Printf("You scored %d out of %d\n", score, len(results))
}

func readAnswers(questions string) <-chan answer {
	resChan := make(chan answer)

	go func() {
		defer close(resChan)
		r := csv.NewReader(strings.NewReader(string(questions)))
		scaner := bufio.NewScanner(os.Stdin)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			checkErr(err)
			checkCSVRecord(record)

			fmt.Print(record[0] + " = ")
			scaner.Scan()

			resChan <- answer{
				question:      record[0],
				userAnswer:    scaner.Text(),
				correctAnswer: record[1],
			}

		}
	}()

	return resChan
}
