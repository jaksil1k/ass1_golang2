package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number!"},
		{"not prime", 8, false, "8 is not a prime number because it is divisible by 2!"},
		{"zero", 0, false, "0 is not prime, by definition!"},
		{"one", 1, false, "1 is not prime, by definition!"},
		{"negative number", -11, false, "Negative numbers are not prime, by definition!"},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)
		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}

func Test_checkNumbers(t *testing.T) {

	primeTests := []struct {
		name       string
		testString string
		expected   bool
		msg        string
	}{
		{"quit", "q", true, ""},
		{"not num", "some", false, "Please enter a whole number!"},
		{"prime", "7", false, "7 is a prime number!"},
		{"not prime", "8", false, "8 is not a prime number because it is divisible by 2!"},
	}

	for _, e := range primeTests {
		s := bufio.NewScanner(strings.NewReader(e.testString))
		msg, result := checkNumbers(s)
		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}

func Test_prompt(t *testing.T) {
	expectedString := "-> "

	realStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	prompt()
	w.Close()

	result, _ := io.ReadAll(r)

	resultString := string(result)

	if resultString != expectedString {
		t.Errorf("expected: %s but got %s", expectedString, resultString)
	}
	os.Stdout = realStdout
}

func Test_intro(t *testing.T) {
	expectedString := "Is it Prime?\n------------\nEnter a whole number, and we'll tell you if it is a prime number or not. Enter q to quit.\n-> "

	realStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	intro()
	w.Close()

	result, _ := io.ReadAll(r)

	resultString := string(result)

	if resultString != expectedString {
		t.Errorf("expected: %s but got %s", expectedString, resultString)
	}
	os.Stdout = realStdout
}

func Test_readUserInput(t *testing.T) {

	testString := "7\nq\n"

	expectedString := "7 is a prime number!\n-> "

	reader := strings.NewReader(testString)
	doneChan := make(chan bool)

	realStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// start a goroutine to read user input and run program
	go readUserInput(reader, doneChan)

	// block until the doneChan gets a value
	<-doneChan

	w.Close()

	result, _ := io.ReadAll(r)

	resultString := string(result)

	if resultString != expectedString {
		t.Errorf("expected: %s but got %s", expectedString, resultString)
	}
	os.Stdout = realStdout

	// close the channel
	close(doneChan)
}

func Test_main(t *testing.T) {

	expectedResult := "7 is a prime number!\n-> Goodbye.\n"

	input := []byte("7\nq\n")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}

	defer func(v *os.File) { os.Stdin = v; os.Stdout = v }(os.Stdin)
	os.Stdin = r
	os.Stdout = w

	main()

	w.Close()

	result, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	resultString := string(result)

	if expectedResult != resultString {
		t.Errorf("expected: %s but got:%s", expectedResult, resultString)
	}
}
