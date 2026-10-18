package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	// set a -csv flag so that a user can specify their own file
	csvFlag := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer' (default \"problems.csv\")")
	// limitFlag := flag.Int("limit", 30, "the time limit for the quiz in seconds (default 30)")

	// parse command-line flags
	flag.Parse()

	data, err := readCSVFile(*csvFlag)
	// check for errors
	if err != nil {
		CSVError(err, "file_read")
		os.Exit(1)
	}

	reader, err := parseCSV(data)
	// check for errors
	if err != nil {
		CSVError(err, "parse")
		os.Exit(1)
	}

	// problem number
	problem_num := 1
	// correct answer counter
	correct_answers := 0

	processCSV(reader, &problem_num, &correct_answers)

	fmt.Printf("\nYou scored %d out of %d\n", correct_answers, problem_num-1)
}

/* The function opens a csv file, reads all the data in bytes and returns it*/
func readCSVFile(filename string) ([]byte, error) {
	// open the csv file
	f, err := os.Open(filename)
	// check for errors
	if err != nil {
		return CSVError(err, "default")
	}

	// close the file in the end of the function
	defer f.Close()

	data, err := io.ReadAll(f)
	// check for errors
	if err != nil {
		return CSVError(err, "default")
	}

	return data, nil
}

// The function converts bytes into a readable stream
func parseCSV(data []byte) (*csv.Reader, error) {
	reader := csv.NewReader(bytes.NewReader(data))
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
		*problem_num++

		var user_input int
		// get user input
		fmt.Scan(&user_input)

		// compare the user input to the answer to the problem stored in the csv file
		// if the answer is correct, increase the correct answer count
		if strconv.Itoa(user_input) == record[1] {
			*correct_answers++
		}
	}
}

// A function returns the errors related to csv files
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
