package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "A CSV file in the format of 'questions and answers'")
	timeLimit := flag.Int("limit", 10, "To set your time limit for the quiz in Seconds")
	flag.Parse()
	// _ = csvFileName

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFileName))
	}

	// CSV Reader
	r := csv.NewReader(file)

	// Parse the CSV
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	// fmt.Println(lines) // Only for test parsing

	problems := parseLines(lines)
	// fmt.Println(problems) // Print the parsed problems

	// Timer is put under the parsedLines so that users won't lose time because of the programs need to parse the problems
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0 // counter for correct answers
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, problem.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("You have answered %d correct answers of %d questions\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == problem.answer {
				correct++
				fmt.Println("Correct!")
			}
		}
	}
}

func parseLines(lines [][]string) []problem {
	result := make([]problem, len(lines))
	for i, line := range lines {
		result[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]), // To trim the space before the first character, hence more answerable
		}
	}
	return result
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
