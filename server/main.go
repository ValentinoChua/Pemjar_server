package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var klien = make(map[net.Conn]string)
var mutex sync.Mutex

func main() {
	ln, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal("gagal membuka port 9001: ", err)
	}
	defer ln.Close()

	log.Println("server chat mendengarkan di port 9001")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("gagal menerima sambungan:", err)
			continue
		}
		go layani(conn)
	}
}

func layani(conn net.Conn) {
	defer func() {
		mutex.Lock()
		nama, terdaftar := klien[conn]
		delete(klien, conn)
		mutex.Unlock()

		conn.Close()

		if terdaftar {
			log.Println("klien terputus:", nama)
			siarkan(fmt.Sprintf("%s meninggalkan obrolan", nama))
		}
	}()

	// Menggunakan alamat koneksi sebagai identitas klien
	nama := conn.RemoteAddr().String()
	log.Println("klien tersambung:", nama)

	mutex.Lock()
	klien[conn] = nama
	mutex.Unlock()

	siarkan(fmt.Sprintf("%s bergabung ke obrolan", nama))

	scanner := bufio.NewScanner(conn)

	// Membaca pesan dari klien
	for scanner.Scan() {
		baris := strings.TrimSpace(scanner.Text())

		if baris == "" {
			continue
		}

		log.Printf("dari %s: %q", nama, baris)

		if strings.EqualFold(baris, "KELUAR") {
			return
		}

		// Mengirim pesan ke semua klien
		siarkan(fmt.Sprintf("%s: %s", nama, baris))
	}

	if err := scanner.Err(); err != nil {
		log.Println("kesalahan koneksi", nama, ":", err)
	}
}

func siarkan(pesan string) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn := range klien {
		if _, err := fmt.Fprintf(conn, "%s\n", pesan); err != nil {
			log.Println("gagal mengirim pesan ke", conn.RemoteAddr(), err)
		}
	}
}
