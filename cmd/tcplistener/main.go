package main

import (
	"fmt"
	"http/internal/request"
	"log"
	"net"
	"os"
)

// delete
// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	out := make(chan string)

// 	go func() {
// 		defer close(out)
// 		defer f.Close()

// 		str := ""

// 		for {
// 			readBytes := make([]byte, 8)

// 			n, err := f.Read(readBytes)
// 			if err != nil {
// 				if err == io.EOF {
// 					break
// 				}
// 				log.Println("error reading file: ", err)
// 				return
// 			}

// 			if i := bytes.IndexByte(readBytes, '\n'); i != -1 {
// 				str += string(readBytes[:i])
// 				readBytes = readBytes[i+1:]
// 				out <- str
// 				str = string(readBytes)
// 			} else {
// 				str += string(readBytes[:n])
// 			}

// 		}

// 		if len(str) != 0 {
// 			out <- str
// 		}
// 	}()

// 	return out
// }

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("failed to establish a TCP connection: ", err)
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			fmt.Println("error while closing TCP connection: ", err)
		}
	}()

	if err := os.Truncate("./out/output.txt", 0); err != nil {
		log.Println("failed to truncate:", err)
	}

	file, err := os.OpenFile("./out/output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error openning file: %s", err)
	}
	defer file.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("error establishing connection: %s", err)
		}
		fmt.Println("a connection has been accepted")

		request, err := request.RequestFromReader(connection)
		if err != nil {
			log.Fatal("error", err)
		}

		fmt.Println("Method: ", request.RequestLine.Method)
		fmt.Println("Target: ", request.RequestLine.RequestTarget)
		fmt.Println("Version: ", request.RequestLine.HttpVersion)

	}

}
