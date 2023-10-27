package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type HttpRequest struct {
	Method         string
	Uri            string
	Version        string
	Host           string
	Accept         string
	AcceptLanguage string
}

type HttpResponse struct {
	Version         string
	StatusCode      string
	ContentType     string
	ContentLanguage string
	Data            string
}

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "2215"
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 1024
	GROUP_NAME  = "SIUUU"
)

func main() {
	conn := handshake()

	defer conn.Close()

	uri, accept, acceptLanguage := handleRequestInput()

	request := HttpRequest{
		Method:         "GET",
		Uri:            uri,
		Version:        "HTTP/1.1",
		Host:           SERVER_HOST,
		Accept:         accept,
		AcceptLanguage: acceptLanguage,
	}

	response := fetch(request, conn)
	fmt.Println("Status Code: " + response.StatusCode)
	fmt.Println("Body: " + response.Data)
}

func handshake() (conn *net.TCPConn) {
	tcpAddr, err := net.ResolveTCPAddr(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	conn, err = net.DialTCP(SERVER_TYPE, nil, tcpAddr)
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	return
}

func handleRequestInput() (uri, accept, acceptLanguage string) {
	fmt.Print("Input the url: ")
	url, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	uriWithNewLine := strings.Split(url, ":"+SERVER_PORT)[1]
	uri = strings.Split(uriWithNewLine, "\r\n")[0]

	fmt.Print("Input the data type: ")
	accept, err = bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	fmt.Print("Input the language: ")
	acceptLanguage, err = bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	return
}

func fetch(request HttpRequest, connection net.Conn) (response HttpResponse) {
	_, err := connection.Write([]byte(requestEncoder(request)))
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, BUFFER_SIZE)
	bufLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	response = responseDecoder(buffer[:bufLen])
	return
}

func responseDecoder(bytestream []byte) (response HttpResponse) {
	receivedMessage := string(bytestream)
	splittedReceivedMessage := strings.Split(receivedMessage, "\r\n")
	informationAboutResponse := strings.Split(splittedReceivedMessage[0], " ")

	response = HttpResponse{
		Version:         informationAboutResponse[0],
		StatusCode:      informationAboutResponse[1],
		ContentType:     splittedReceivedMessage[1],
		ContentLanguage: splittedReceivedMessage[2],
		Data:            strings.Split(splittedReceivedMessage[3], ": ")[1],
	}

	return
}

func requestEncoder(request HttpRequest) []byte {
	data := fmt.Sprintf(
		"%s %s %s\r\nHost: %s\r\nAccept: %sAccept-Language: %s",
		request.Method, request.Uri, request.Version, request.Host, request.Accept, request.AcceptLanguage,
	)

	return []byte(data)
}
