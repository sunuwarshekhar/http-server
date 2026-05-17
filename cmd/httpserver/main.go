package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"main.go/internal/request"
	"main.go/internal/response"
	"main.go/internal/server"
)

const port = 42069

func response400() []byte {
	return []byte(`<html>
	<head>
	<title>400 Bad Request</title>
	</head>
	<body>
		<h1>Bad Request</h1>
		<p>Your request is bad.</p>
	</body>
	</html>`)
}
func response500() []byte {
	return []byte(`<html>
	<head>
	<title>500 Server Error</title>
	</head>
	<body>
		<h1>Server Error</h1>
		<p>Something went wrong.</p>
	</body>
	</html>`)
}
func response200() []byte {
	return []byte(`<html>
	<head>
	<title>400 Bad Request</title>
	</head>
	<body>
		<h1>Success</h1>
		<p>Good Request</p>
	</body>
	</html>`)
}
func main() {
	s, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		h := response.GetDefaultHeaders(0)
		body := response200()
		status := response.StatusOk
		if req.RequestLine.RequestTarget == "/yourproblem" {

			body = response400()
			status = response.StatusBadRequest

		} else if req.RequestLine.RequestTarget == "/myproblem" {
			body = response500()
			status = response.StatusInternalServerError

		}
		h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		h.Replace("Content-Type", "text/html")
		w.WriteStatusLine(status)
		w.WriteHeaders(*h)
		w.WriteBody(body)
	})
	if err != nil {
		log.Fatalf("Error starting server:%v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
