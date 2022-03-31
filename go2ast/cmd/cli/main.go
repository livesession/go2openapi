package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/livesession/restflix/go2ast"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	src := ""

	for scanner.Scan() {
		src += scanner.Text()
	}

	if err := go2ast.Generate(src, os.Stdout); err != nil {
		fmt.Println(err)
	}
}
