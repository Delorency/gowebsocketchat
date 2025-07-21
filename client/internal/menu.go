package internal

import "fmt"

func DrawMenu() uint8 {
	fmt.Println("--------Menu-------")
	fmt.Println("(1) New chat")
	fmt.Println("(2) Connect to chat")
	fmt.Println("(0) Exit")

	var dec uint8
	fmt.Print(": ")
	fmt.Scan(&dec)

	return dec

}

func DrawInvite() uint8 {
	fmt.Println("-------Invite------")
	fmt.Println("(0) Back")
	fmt.Println("-------------------")
	var dec uint8
	fmt.Print("chat id: ")
	fmt.Scan(&dec)

	return dec
}
