package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestBackChar(t *testing.T) {

	str := strings.Repeat("\b", 10)

	fmt.Printf("123456%s7890abc\n", str)
}
