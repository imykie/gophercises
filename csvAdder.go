package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("problems.csv")
	count := 0
	questions := 0

	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)

	for {
		record, err := r.Read()
		var v int

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("What is %v, sir? ", record[0])
		fmt.Scan(&v)
		ans, _ := strconv.Atoi(record[1])
		if v == ans {
			count++
		}
		questions++
	}

	fmt.Printf("You answered %d out of %d correctly", count, questions)
}
