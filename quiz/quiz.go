package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(fileName string) string {
	dat, err := ioutil.ReadFile(fileName)
	check(err)
	return string(dat)
}

func getQuestions(lines []string) (questions, answers []string) {
	questions = make([]string, length)
	answers = make([]string, length)
	for i := 0; i < length; i++ {
		line := lines[i]
		quizSplit := strings.Split(line, ",")
		questions[i] = quizSplit[0]
		answers[i] = quizSplit[1]
	}
	return questions, answers
}

func printQuestion(question, answer string) int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question)
	text, _ := reader.ReadString('\n')
	if strings.Trim(text, "\r\n") == answer {
		fmt.Println("Correct!")
		return 1
	}
	return 0
}

func startTimer(seconds int) *time.Timer {
	timer := time.NewTimer(10 * time.Second)
	return timer
}

func Generator() chan int {
	number := make(chan int)
	go func() {
		n := 0
		for {
			select {
			case number <- n:
				n++
			case <- number:
				return
			}
		}
	}()
	return number
}

var length = 0

func main() {
	dat := readFile("problems.csv")
	lines := strings.Split(dat, "\r\n")
	length = len(lines)
	questions, answers := getQuestions(lines)
	number := Generator()
	correctCounter := 0
	go func () {
		for {
			ind := <-number
			if ind >= length {
				fmt.Println("Quiz finished!")
				fmt.Println("Number Correct:", correctCounter)
				os.Exit(2)
			}
			correct := printQuestion(questions[ind], answers[ind])
			correctCounter += correct
		}
	}()
	timer := startTimer(10)
	<-timer.C
	number<-1
	fmt.Println("Times up")
	fmt.Println("Number Correct: ", correctCounter)
}
