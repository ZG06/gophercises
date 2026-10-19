package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func main() {
	// set the -csv flag so that a user can specify their own file
	csvFlag := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer' (default \"problems.csv\")")
	// add a -limit flag to set the timer duration
	limitFlag := flag.Int("limit", 30, "the time limit for the quiz in seconds (default 30)")

	// parse command-line flags
	flag.Parse()

	f, err := readCSVFile(*csvFlag)
	// check for errors
	if err != nil {
		CSVError(err, "file_read")
		os.Exit(1)
	}

	defer f.Close()

	reader, err := parseCSV(f)
	// check for errors
	if err != nil {
		CSVError(err, "parse")
		os.Exit(1)
	}
	// channel for processCSV function state
	processChan := make(chan bool, 1)
	answerChan := make(chan bool, 1)

	// wait for the user to press Enter to start
	fmt.Print("Press 'Enter' to start...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// goroutine for giving problems to the user in the background
	go func() {
		/* if the user answers all the question before the time runs out,
		'true' will be send to the processChan channel*/
		processCSV(reader, answerChan)
		processChan <- true
	}()

	problem_num := 0
	correct_answers := 0

	// Wait until one of the conditions is met
	for {
		select {
		case isCorrect := <-answerChan:
			// increate the problem_num counter when answerChan receives data
			problem_num++
			// is the answer was correct, increase the correct_answers counter
			if isCorrect {
				correct_answers++
			}

		case <-processChan:
			fmt.Printf("\nYou scored %d out of %d\n", correct_answers, problem_num)
			return // exit main when all questions are answered

		case <-time.After(time.Duration(*limitFlag) * time.Second):
			fmt.Printf("\nYou scored %d out of %d\n", correct_answers, problem_num)
			return // exit main when timer expires
		}
	}
}

/* The function opens a csv file, reads all the data in bytes and returns it*/
func readCSVFile(filename string) (*os.File, error) {
	// open the csv file
	f, err := os.Open(filename)
	// check for errors
	if err != nil {
		CSVError(err, "default")
		os.Exit(1)
	}

	return f, nil
}

// The function converts bytes into a readable stream
func parseCSV(data io.Reader) (*csv.Reader, error) {
	reader := csv.NewReader(data)
	return reader, nil
}

// The function processes the data from a csv file and evaluating the user's answers
func processCSV(reader *csv.Reader, answerChan chan<- bool) {
	problemNum := 1

	// read every line in the given csv file
	for {
		record, err := reader.Read()

		// check for errors
		if err == io.EOF {
			break
		} else if err != nil {
			CSVError(err, "csv_read")
			break
		}

		// print a problem number and a question
		fmt.Printf("Problem №%d: %s = ", problemNum, record[0])

		var user_input int
		// get user input
		fmt.Scan(&user_input)

		// compare the user input to the answer to the problem stored in the csv file
		// if the answer is correct, increase the correct answer count
		if strconv.Itoa(user_input) == record[1] {
			answerChan <- true
		} else {
			answerChan <- false
		}

		problemNum++
	}
}

// The function returns the errors related to csv files
func CSVError(err error, err_type string) ([]byte, error) {
	switch err_type {
	case "file_read":
		fmt.Println("Error reading file:", err)
		return nil, err

	case "parse":
		fmt.Println("Error creating CSV reader:", err)
		return nil, err

	case "csv_read":
		fmt.Println("Error reading CSV data:", err)
		return nil, err

	default:
		return nil, err
	}
}
