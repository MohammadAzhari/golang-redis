package main

import (
	"fmt"
	"net"
	"strings"
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

	for {
		resp := NewResp(conn)

		aof, err := NewAof("database.aof")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer aof.Close()

		v, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}

		if v.typ != "array" || len(v.array) == 0 {
			fmt.Println("Excpecting array of length more than 0")
			continue
		}

		command := strings.ToUpper(v.array[0].bulk)
		handler, ok := handlers[command]

		w := NewWriter(conn)

		if !ok {
			fmt.Println("Invalid command: ", command)
			w.Write(Value{typ: "string", str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(v)
		}

		args := v.array[1:]
		result := handler(args)

		err = w.Write(result)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
