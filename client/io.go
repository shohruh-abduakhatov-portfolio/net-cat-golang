package client

import (
	"bufio"
	"io"
	"net"
	"strings"
)

type IO struct {
	Conn   net.Conn
	input  chan string
	output chan string
}

func StartClient(ip, port string) (*IO, error) {

	// tcpAddr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	// if err != nil {
	// 	return nil, err
	// }
	// conn, err := net.DialTCP("tcp", nil, tcpAddr)

	// connected to server
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return nil, err
	}
	o, err := NewIO(conn)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func NewIO(conn net.Conn) (*IO, error) {
	o := &IO{
		Conn:   conn,
		input:  make(chan string),
		output: make(chan string),
	}
	go o.handleInput()
	go o.listenServer()
	return o, nil
}

func (o *IO) handleInput() {
	for {
		select {
		// case <-stop:
		// 	return
		case msg := <-o.input:
			if msg == "" {
				continue
			}
			r := strings.NewReader(msg + "\n")
			_, err := io.Copy(oi.Conn, r)
			if err != nil {
				continue
			}
			continue
		}
	}
}

func (o *IO) listenServer() {
	for {
		// scanner := bufio.NewScanner(o.Conn)
		// message := ""
		// for scanner.Scan() {
		// 	message += scanner.Text()
		// }
		// if err := scanner.Err(); err != nil {
		// 	fmt.Fprintln(os.Stderr, "reading standard input:", err)
		// }
		// V2
		// buf := new(bytes.Buffer)
		// buf.ReadFrom(o.Conn)
		// message := buf.String()
		//
		// v3
		line, _, err := bufio.NewReader(o.Conn).ReadLine()
		if err != nil {
			continue
		}
		message := string(line)
		message = strings.Trim(strings.TrimSpace(message), "\r\n")

		// v4
		// Calling Pipe method
		// pipeReader, pipeWriter := io.Pipe()

		// // Using Fprint in go function to write
		// // data to the file
		// go func() {
		// 	fmt.Fprint(pipeWriter, oi.Conn)
		// 	// Using Close method to close write
		// 	pipeWriter.Close()
		// }()

		// // Creating a buffer
		// buffer := new(bytes.Buffer)

		// // Calling ReadFrom method and writing
		// // data into buffer
		// buffer.ReadFrom(pipeReader)

		// // Prints the data in buffer
		// message := buffer.String()

		o.output <- message
	}
}
