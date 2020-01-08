package main

import (
	"fmt"
)

type Token struct {
	data      string
	recipient int
	ttl       int
}

var createdNum int
var tokenReturning bool

func inputHandler() int {
	for {
		var N int
		_, err := fmt.Scanf("%d", &N)
		if err == nil {
			fmt.Println("Input accepted")
			return N
		}
		fmt.Println("There was an error with handling your number. Please, try again.")
	}
}

func messageInputHandler() Token {
	token := Token{}
	return token
}

func tokenRing(upperChan chan Token, upperID int, maxN int) {
	var lowerChan chan Token
	childCreated := false

	innerID := upperID + 1
	fmt.Printf("Inner ID:%d\n", innerID)

	token := <-upperChan

	if token == (Token{}) && maxN != createdNum {
		lowerChan = make(chan Token)
		childCreated = true
		fmt.Printf("Created new goroutine #%d\n", innerID)
		createdNum++
		go tokenRing(lowerChan, innerID, maxN)
		lowerChan <- (Token{})
		fmt.Printf("If this works, golang is crap or I'm stupid %d\n", innerID)
		token = <-lowerChan
	} else if createdNum == maxN {
		fmt.Printf("Returning control from %d to %d\n", innerID, upperID)
		upperChan <- (Token{ttl: -2})
	}

	if token.ttl == -2 {
		fmt.Printf("Returning control from ttl statement %d to %d\n", innerID, upperID)
		upperChan <- (Token{ttl: -2})
	}

	for {
		fmt.Printf("Entered loop at: %d\n", innerID)

		if tokenReturning {
			fmt.Printf("Got token from lower lvls at %d\n", innerID)
			token = <-lowerChan
		} else {
			fmt.Printf("Got token from upper lvls at %d\n", innerID)
			token = <-upperChan
		}

		fmt.Printf("Taking token from somewhere %d\n", innerID)

		if token.recipient == innerID {
			fmt.Printf("Delivered at: %d\n", innerID)
			token.data = "The message has been delivered."
			token.ttl = -1
			tokenReturning = true
			upperChan <- token
		} else if token.ttl == 0 {
			fmt.Printf("Ttl ran out on: %d\n", innerID)
			token.data = "The message hasn't been delivered due to ending of ttl."
			tokenReturning = true
			upperChan <- token
		} else if token.ttl == -1 {
			fmt.Printf("Returning token on upper level at: %d\n", innerID)
			tokenReturning = true
			upperChan <- token
		} else if childCreated {
			token.ttl--
			// fmt.Println("Sending token lower")
			fmt.Printf("Froze the routine #%d\n", innerID)
			lowerChan <- token
			fmt.Printf("Got unfrozen %d, ttl: %d\n", innerID, token.ttl)
		} else {
			fmt.Printf("Didn't find a recipient at: %d\n", innerID)
			token.data = "The message hasn't been delivered because there was no such recipient."
			tokenReturning = true
			upperChan <- token
		}
	}
}

func main() {
	fmt.Println("Started")

	fmt.Print("Enter number of go routines: ")
	N := inputHandler()

	fmt.Println("Pre goroutine message")
	messages := make(chan Token)

	sampleToken := Token{}

	createdNum = 1
	go tokenRing(messages, 0, N)
	tokenReturning = false
	messages <- sampleToken

	fmt.Printf("Created %d routine(s)\n", N)

	fmt.Println("Retrieving sample token")

	sampleToken = <-messages
	tokenReturning = false

	if sampleToken.ttl == -2 {
		fmt.Println("Tis good")
	}

	fmt.Println("Sending token")

	token := Token{data: "kappa123", recipient: 1, ttl: 4}
	messages <- token
	msg := <-messages

	fmt.Println(msg.data)
	close(messages)
}
