package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	SERVER_HOST = ""
	SERVER_PORT = "2215"
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 1024
	GROUP_NAME  = "SIUUU"
	HTML_RES_ID = `<html><body><h1>Halo, kami dari SIUUU</h1></body></html>`
	HTML_RES_EN = `<html><body><h1>Hello, we are from SIUUU</h1></body></html>`
)

type HttpRequest struct {
	Method          string
	Uri             string
	Version         string
	Host            string
	Accept          string
	AcceptLanguange string
}

type HttpResponse struct {
	Version         string
	StatusCode      string
	ContentType     string
	ContentLanguage string
	Data            string
}

type Student struct {
	Name      string
	StudentId string
}

func main() {
	tcpAddr, server := initServer()

	defer server.Close()

	fmt.Println("Listening on", tcpAddr.String())
	fmt.Println("Waiting for client...")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error:", err.Error())
		} else {
			go handleConnection(conn)
		}
	}
}

func initServer() (tcpAddr *net.TCPAddr, server *net.TCPListener) {
	tcpAddr, err := net.ResolveTCPAddr(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}

	server, err = net.ListenTCP(SERVER_TYPE, tcpAddr)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	return
}

func handleConnection(connection net.Conn) {
	buffer := make([]byte, BUFFER_SIZE)
	bufLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	request := requestDecoder(buffer[:bufLen])
	responseBody := handleRequest(request)
	_, err = connection.Write(ResponseEncoder(responseBody))

	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}

	defer connection.Close()
}

func handleRequest(request HttpRequest) (responseBody HttpResponse) {
	responseBody.Version = "HTTP/1.1"

	if request.Method == "GET" {
		if (request.Uri == "/" || request.Uri == "/greeting" || request.Uri == "/?name=SIUUU") && (request.Accept == "text/html") {
			responseBody.StatusCode = "200"
			responseBody.ContentType = "text/html"
			responseBody.ContentLanguage = request.AcceptLanguange

			if strings.Contains(request.AcceptLanguange, "en-US") {
				responseBody.Data = HTML_RES_EN
			} else if strings.Contains(request.AcceptLanguange, "id-ID") {
				responseBody.Data = HTML_RES_ID
			} else {
				responseBody.ContentLanguage = "en-US"
				responseBody.Data = HTML_RES_EN
			}

			return
		}

		if (request.Uri == "/data") && (request.Accept == "application/xml" || request.Accept == "application/json") {
			responseBody.StatusCode = "200"
			responseBody.ContentType = request.Accept
			responseBody.ContentLanguage = request.AcceptLanguange

			students := []Student{
				{
					Name:      "Dhafin Raditya Juliawan",
					StudentId: "2106650304",
				},
				{
					Name:      "Rifqi Farel Muhammad",
					StudentId: "2106650310",
				},
				{
					Name:      "Fadhlan Hasyim",
					StudentId: "2106652215",
				},
			}

			var responseData []byte
			if strings.Contains(request.Accept, "application/xml") {
				responseData, _ = xml.Marshal(students)
			} else {
				responseData, _ = json.Marshal(students)
			}
			responseBody.Data = string(responseData)

			return
		}
	}

	responseBody.StatusCode = "404"
	responseBody.ContentType = ""
	responseBody.ContentLanguage = request.AcceptLanguange

	return
}

func requestDecoder(bytestream []byte) (request HttpRequest) {
	requestData := string(bytestream)
	splittedRequestData := strings.Split(requestData, "\r\n")
	informationAboutRequest := strings.Split(splittedRequestData[0], " ")

	request = HttpRequest{
		Method:          informationAboutRequest[0],
		Uri:             informationAboutRequest[1],
		Version:         informationAboutRequest[2],
		Host:            strings.Split(splittedRequestData[1], " ")[1],
		Accept:          strings.Split(splittedRequestData[2], " ")[1],
		AcceptLanguange: strings.Split(splittedRequestData[3], " ")[1],
	}

	return
}

func ResponseEncoder(response HttpResponse) []byte {
	data := fmt.Sprintf(
		"%s %s\r\nContent-Type: %s\r\nContent-Language: %s\r\nData: %s",
		response.Version, response.StatusCode, response.ContentType, response.ContentLanguage, response.Data,
	)

	return []byte(data)
}
