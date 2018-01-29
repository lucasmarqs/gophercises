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
)

type Question struct {
	Text   string
	Answer string
}

func main() {
	quiz := flag.String("quiz", "problems.csv", "Optional. Define alternative problems to the quiz.")
	flag.Parse()
	questions, err := parseCSV(*quiz)
	if err != nil {
		log.Fatalln(err)
	}

	control := struct {
		Corrects int
		Wrongs   int
	}{Corrects: 0, Wrongs: 0}

	fmt.Println("[[[ Welcome to Quiz Game ]]]")
	reader := bufio.NewReader(os.Stdin)
	for idx, q := range *questions {
		fmt.Printf("(Question %d) %s ? ", idx+1, q.Text)
		answer, _ := reader.ReadString('\n')
		answer = answer[:len(answer)-1]
		if strings.Compare(answer, q.Answer) == 0 {
			control.Corrects += 1
		} else {
			control.Wrongs += 1
		}
	}

	fmt.Printf("You got %d questions and missed %d\n", control.Corrects, control.Wrongs)
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
