package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	Init()
	file, err := os.OpenFile("server-log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("can't open log file %e", err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("can't close log file %e", err)
		}
	}()

	log.SetOutput(file)
	log.Print("start application\n")

	host := "0.0.0.0"
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "9999"
	}
	err = startServer(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatal(err)
	}

}

func startServer(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("can't listen 0.0.0.0:9999 %v", err)
		return err
	}
	defer func() {
		err = listener.Close()
		if err != nil {
			log.Printf("Can't close Listener: %v", err)
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't connect client")
			continue
		}
		fmt.Println("connected")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	readString, _ := reader.ReadString('\n')
	split := strings.Split(strings.TrimSpace(readString), " ")
	if len(split) != 3 {
		log.Print("incorrect request")
		return
	}
	meth, request, protocol := split[0], split[1], split[2]
	if meth == "GET" && protocol == "HTTP/1.1" {
		ResponseToHttp(request, conn)
	}

}
var ContentType = make(map[string]string)

func Init() {
	ContentType[""] = "Content-Type: text/html\r\n"
	ContentType["png"] = "Content-Type: image/png\r\n"
	ContentType["png?download"] = "Content-Disposition: attachment; filename=down.png\r\n"
	ContentType["jpg"] = "Content-Type: image/jpg\r\n"
	ContentType["jpg?download"] = "Content-Disposition: attachment; filename=down.jpg\r\n"
	ContentType["html"] = "Content-Type: text/html\r\n"
	ContentType["html?download"] = "Content-Disposition: attachment; filename=down.html\r\n"
	ContentType["txt"] = "Content-Type: text/html\r\n"
	ContentType["txt?download"] = "Content-Disposition: attachment; filename=down.txt\r\n"
	ContentType["pdf"] = "Content-Type: application/pdf\r\n"
	ContentType["jpg?download"] = "Content-Disposition: attachment; filename=down.jpg\r\n"
}
func ResponseToHttp(http string, conn net.Conn) {
	var index string
	download := ""
	var file string
	if strings.HasSuffix(http, "?download") {
		http = http[1 : len(http)-9]
		file = http

		download = "?download"
	} else {
		file = http[1:]
	}
	if strings.HasSuffix(http, ".html") {
		index = "html"
	} else
	if strings.HasSuffix(http, ".pdf") {
		index = "pdf"
	} else
	if strings.HasSuffix(http, ".png") {
		index = "png"
	} else
	if strings.HasSuffix(http, ".jpg") {
		index = "jpg"
	} else
	if strings.HasSuffix(http, ".txt") {
		index = "txt"
	} else {
		index = ""
	}
	index += download
	if file == "" {
		file = "commands.html"
	}

	all, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("can't read file %v", err)
		file = "404.html"
		all, err = ioutil.ReadFile(file)
		if err != nil {
			log.Printf("can't read file %v", err)
			return
		}
	}
	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString("HTTP/1.1 200 OK\r\n")
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	_, err = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(all)))
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	_, err = writer.WriteString(ContentType[index])
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	_, err = writer.WriteString("Connection: Close\r\n")
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	_, err = writer.WriteString("\r\n")
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	_, err = writer.Write(all)
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		log.Printf("can't send string %v", err)
		return
	}
}