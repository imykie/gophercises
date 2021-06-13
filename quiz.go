package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Countdown struct {
	t int
	d int
	h int
	m int
	s int
	n int
}

type Problem struct {
	q string
	a string
}

func getRemainingTime(t time.Time) Countdown {
	currentTime := time.Now()
	difference := t.Sub(currentTime)

	total := int(difference.Seconds())
	days := total / (60 * 60 * 24)
	hours := total / (60 * 60) % 24
	minutes := (total / 60) % 60
	seconds := total % 60

	return Countdown{
		t: total,
		d: days,
		h: hours,
		m: minutes,
		s: seconds,
	}
}

func countdown(t chan Countdown, timer time.Time) {
	for range time.Tick(1 * time.Second) {
		tc := getRemainingTime(timer)
		t <- tc

		if tc.t <= 0 {
			break
		}
	}
	close(t)
}

func main() {
	file := flag.String("file", "problems.csv", "A CSV file in format of (question, answer)")
	duration := flag.Int("duration", 30, "Duration for the quiz in seconds")
	flag.Parse()
	f, err := os.Open(*file)
	count := 0
	timer := time.NewTimer(time.Second * time.Duration(*duration))
	defer timer.Stop()

	if err != nil {
		log.Fatal(err)
		exit(fmt.Sprintf("Can't open csv file: %s", *file))
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()

	if err != nil {
		log.Fatal(err)
		exit("Error occurred while reading csv file")
	}

	problems := parseLines(records)
	tt := make(chan Countdown)
	tm := time.Now().Add(time.Second * 5)

	fmt.Println("You have", *duration, "seconds to complete the tasks")
	go countdown(tt, tm)

	for t := range tt {
		fmt.Printf("You quiz starts in %d seconds \n", t.s)
	}
	fmt.Println("Start!")

problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem %d:  %v = ", i+1, p.q)
		ansCh := make(chan string)
		go func() {
			var ans string
			fmt.Scanf("%s", &ans)
			ansCh <- ans
		}()

		select {
		case <-timer.C:
			fmt.Println("\nYour time is up!")
			break problemLoop
		case ans := <-ansCh:
			if ans == p.a {
				count++
			}
		}
	}

	fmt.Printf("You answered %d out of %d correctly", count, len(records))
}

func parseLines(lines [][]string) []Problem {
	data := make([]Problem, len(lines))
	for i, line := range lines {
		data[i] = Problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return data
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
