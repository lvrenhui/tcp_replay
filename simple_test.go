package main

import (
	"log"
	"strings"
	"testing"
)

func TestInput(t *testing.T) {
	log.Println(strings.Split("127.0.0.1:4000", ":")[1])
}
