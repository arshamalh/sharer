package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	srv := http.NewServeMux()

	// Routes
	file := []byte{}
	srv.HandleFunc("/upload", func(w http.ResponseWriter, req *http.Request) {
		fileName := req.URL.Query().Get("fileName")
		new_chunk, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprint(w, "There was something wrong")
		}
		defer req.Body.Close()
		file = append(file, new_chunk...)
		os.WriteFile(fileName, file, 0666)
		fmt.Fprint(w, "Chunk received!")
	})

	// File server 1
	uiServer := http.FileServer(http.Dir("./ui"))
	srv.Handle("/", uiServer)

	// File server 2
	staticFilesServer := http.FileServer(http.Dir("."))
	srv.Handle("/files/", staticFilesServer)

	port := "60"
	fmt.Printf("to receive: %s:%s", GetOutboundIP(), port)
	fmt.Printf("to serve: %s:%s/files", GetOutboundIP(), port)

	fmt.Println()
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatal(err)
	}
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	return strings.Split(conn.LocalAddr().String(), ":")[0]
}
