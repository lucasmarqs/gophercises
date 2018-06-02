package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// Question holds the question text and its answer
type Question struct {
	Text   string
	Answer string
}

// Control holds game informations
type Control struct {
	Corrects       int
	Timeout        bool
	TotalQuestions int
}

// Missed results on how many questions the gamer missed
func (c Control) Missed() int {
	return c.TotalQuestions - c.Corrects
}

func main() {
	quiz := flag.String("quiz", "problems.csv", "Optional. Define alternative problems to the quiz.")
	duration := flag.Int("timer", 30, "In seconds. The maxium time to complete the quiz.")
	flag.Parse()

	questions, err := parseCSV(*quiz)
	if err != nil {
		log.Fatalln(err)
	}

	timer := time.NewTimer(time.Duration(*duration) * time.Second)
	control := Control{TotalQuestions: len(*questions)}

	fmt.Println("[[[ Welcome to Quiz Game ]]]")
	reader := bufio.NewReader(os.Stdin)

GameLoop:
	for idx, q := range *questions {
		fmt.Printf("(Question %d) %s ? ", idx+1, q.Text)
		answerCh := make(chan string)

		go func() {
			answer, _ := reader.ReadString('\n')
			answerCh <- answer[:len(answer)-1]
		}()

		select {
		case <-timer.C:
			control.Timeout = true
			break GameLoop
		case answer := <-answerCh:
			if strings.Compare(answer, q.Answer) == 0 {
				control.Corrects++
			}
		}
	}

	gameOver(control)
}

func gameOver(c Control) {
	if c.Timeout {
		fmt.Printf("\nYour time is gone!\n")
	}
	fmt.Printf("You got %d of questions and missed %d\n", c.Corrects, c.Missed())
}

func parseCSV(filename string) (*[]Question, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)
	questions := []Question{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		q := Question{Text: row[0], Answer: row[1]}
		questions = append(questions, q)
	}
	return &questions, nil
}
