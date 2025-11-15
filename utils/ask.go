package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ask prompts the user with a question and returns their input as a string
func Ask(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", question)
	text, _ := reader.ReadString('\n')
	return strings.TrimSuffix(text, "\n")
}
