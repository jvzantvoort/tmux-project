package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Ask(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", question)
	text, _ := reader.ReadString('\n')
	return strings.TrimSuffix(text, "\n")
}
