package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/MohammadAzhari/golang-redis/resp"
)

func main() {
	fmt.Println("Listening on port :6379")

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer aof.Close()

	err = aof.Read((func(value resp.Value) {
		handleValue(value)
	}))

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		respReader := resp.NewRespReader(conn)

		input, err := respReader.Read()

		if err != nil {
			fmt.Println(err)
			return
		}

		output := handleValue(input)
		err = aof.Write(input)
		if err != nil {
			fmt.Println(err)
			return
		}

		respWriter := resp.NewRespWriter(conn)

		err = respWriter.Write(output)

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func handleValue(value resp.Value) resp.Value {
	if value.Type != "array" || len(value.Array) == 0 {
		fmt.Println("Excpecting array of length more than 0")
		return resp.Value{Type: "error", String: "Excpecting array of length more than 0"}
	}

	command := strings.ToUpper(value.Array[0].String)
	handler, ok := handlers[command]

	if !ok {
		return resp.Value{Type: "error", String: "Invalid command: " + command}
	}

	args := value.Array[1:]
	return handler(args)
}
