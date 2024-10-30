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
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the amount of time, in seconds, given to take the quiz.")

	flag.Parse()

	file, err := os.Open(*csvFilename)

	if err != nil {
		fmt.Printf("failed to open the CSV file:  %s\n", *csvFilename)
		os.Exit(1)
	}

	// The csv.NewReader() function is called in
	// which the object os.File passed as its parameter
	// and this creates a new csv.Reader that reads
	// from the file
	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Failed to parse CSV file.")
	}

	problems := pareLines(lines)

	fmt.Println("Staring quiz...")

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	defer timer.Stop()

	correct := 0

myLoop: // allows breakpoint for nested components
	for i, p := range problems {
		fmt.Printf("Problem %d: %s = ", i+1, p.q)

		answerCh := make(chan string) //make answer channel
		go func() {
			var ans string
			fmt.Scanln(&ans)
			answerCh <- ans // send answer to channel
		}()

		// this utilizes 'channels' in golang to send the signal, not requiring memory allocation.
		// the answer and the timer end signal are sent through separate channels, the switch statement
		// either ends the quiz or checks for correct answer depending on which channel is received first.

		select {
		case <-timer.C: // received timer end from channel
			fmt.Println()
			break myLoop
		case ans := <-answerCh: // received answer from channel
			if ans == p.a {
				correct++
			}
		}

	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func pareLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}
