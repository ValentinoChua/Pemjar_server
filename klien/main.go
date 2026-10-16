package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var mutex sync.Mutex

func main() {
	alamat := "localhost:9001"
	if len(os.Args) > 1 {
		alamat = os.Args[1]
	}

	conn, err := net.Dial("tcp", alamat)
	if err != nil {
		fmt.Println("gagal tersambung ke", alamat, ":", err)
		return
	}
	defer conn.Close()

	fmt.Println("tersambung ke", alamat)
	fmt.Println("ketik pesan lalu Enter. Ketik KELUAR untuk berhenti.")

	// Membaca pesan dari server
	go func() {
		dariServer := bufio.NewScanner(conn)

		for dariServer.Scan() {
			mutex.Lock()
			fmt.Printf("\n%s\n", dariServer.Text())
			fmt.Print("> ")
			mutex.Unlock()
		}

		if err := dariServer.Err(); err != nil {
			fmt.Println("\nkesalahan membaca dari server:", err)
		}
	}()

	// Membaca input pengguna
	keyboard := bufio.NewScanner(os.Stdin)

	for {
		if !keyboard.Scan() {
			return
		}

		baris := strings.TrimSpace(keyboard.Text())

		if strings.EqualFold(baris, "KELUAR") {
			return
		}

		if _, err := fmt.Fprintf(conn, "%s\n", baris); err != nil {
			fmt.Println("gagal mengirim:", err)
			return
		}
	}
}
