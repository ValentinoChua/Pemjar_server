package main

import "fmt"

func main() {
	fmt.Println("Hello")

	//cara 1:
	// var bilangan int
	// bilangan = 1

	//cara 2: lebih praktis bisa otomatis merubah tipe variable data tidak perlu declare
	bilangan := 10

	fmt.Println(bilangan)

	if bilangan > 5 {
		fmt.Println(bilangan)
	}

}
