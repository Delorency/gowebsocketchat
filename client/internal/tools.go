package internal

import "fmt"

func ClearLastNRows(n int) {
	fmt.Printf("\033[%dA\033[K", n)
	for range n {
		fmt.Println("                      ")
	}
	fmt.Printf("\033[%dA\033[K", n)
}

func СlearConsole() {
	fmt.Print("\033[H\033[2J")
}
