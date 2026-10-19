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

	// problem number
	problem_num := 1
	// correct answer counter
	correct_answers := 0

	// channel for processCSV function state
	processState := make(chan bool, 1)

	// wait for the user to press Enter to start
	fmt.Print("Press 'Enter' to start...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// goroutine for giving problems to the user in the background
	go func() {
		/* if the user answers all the question before the time runs out,
		'true' will be send to the processState channel*/
		processCSV(reader, &problem_num, &correct_answers)
		processState <- true
	}()

	// Wait until one of the conditions is met
	select {
	// finish monitoring is processState changes to 'true'
	case <-processState:
		break
	// finish monitoring if the user runs out of time
	case <-time.After(time.Duration(*limitFlag) * time.Second):
		break
	}

	// write the user's score at the end of the program
	fmt.Printf("\nYou scored %d out of %d\n", correct_answers, problem_num)
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
func processCSV(reader *csv.Reader, problem_num, correct_answers *int) {
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
		fmt.Printf("Problem №%d: %s = ", *problem_num, record[0])

		var user_input int
		// get user input
		fmt.Scan(&user_input)

		// compare the user input to the answer to the problem stored in the csv file
		// if the answer is correct, increase the correct answer count
		if strconv.Itoa(user_input) == record[1] {
			*correct_answers++
		}

		*problem_num++
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
