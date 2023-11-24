# Socket Programming
a simple API with server socket written from scratch in Golang.

## Technology Stack
[![Tech Stacks](https://skillicons.dev/icons?i=golang)](https://skillicons.dev)

## How to Run the Program
1. Run ```go run main.go``` on the server package to run the server
2. Run ```go run main.go``` on the client package to run the client
3. Input the URL, accepted content type, and accepted content language
4. The response will be sent by the server based on the input

### List of URL
- SERVER_HOST:SERVER_PORT/ (note: accepted content type must be text/html)
- SERVER_HOST:SERVER_PORT/greeting (note: accepted content type must be text/html)
- SERVER_HOST:SERVER_PORT/?name=SIUUU (note: accepted content type must be text/html)
- SERVER_HOST:SERVER_PORT/data (note: accepted content type must be application/xml or application/json)

### List of Accepted Content Type
- text/html
- application/xml
- application/json

### List of Accepted Content Language
- en-US
- id-ID

## Authors
- Dhafin Raditya Juliawan
- Fadhlan Hasyim
- Rifqi Farel Muhammad
