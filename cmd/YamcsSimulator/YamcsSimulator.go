package main

import (
	"acubesat/ops/yamcs-simulator/internal/connection"
	"acubesat/ops/yamcs-simulator/internal/tc_decoder"
	"acubesat/ops/yamcs-simulator/internal/tm_responcer"
	"fmt"
	"strings"
	"sync"
)

var tcp_wg sync.WaitGroup
var TC_data = make(chan []byte, 1000)
var TC_decoded_data = make(chan tc_decoder.TCData, 1000)
var TM_data = make(chan []byte, 256)
var Events = make(chan string, 5)
var connections_list = [9]*connection.TCPClient{}

func main() {
	var tm_client connection.TCPClient = connection.TCPClient{Port: 10013, Host: "localhost", Description: "TM", Debug: true, ExitInWarning: false, Wg: &tcp_wg}
	var tm_client1 connection.TCPClient = connection.TCPClient{Port: 10014, Host: "localhost", Description: "TM", Debug: true, ExitInWarning: false, Wg: &tcp_wg}
	var tm_client2 connection.TCPClient = connection.TCPClient{Port: 10015, Host: "localhost", Description: "TM", Debug: true, ExitInWarning: false, Wg: &tcp_wg}
	var tm_client3 connection.TCPClient = connection.TCPClient{Port: 10016, Host: "localhost", Description: "TM", Debug: true, ExitInWarning: false, Wg: &tcp_wg}
	var tm_client4 connection.TCPClient = connection.TCPClient{Port: 10012, Host: "localhost", Description: "TM", Debug: true, ExitInWarning: false, Wg: &tcp_wg}
	var tc_client connection.TCPClient = connection.TCPClient{Port: 10025, Host: "localhost", Description: "TC", Debug: true, ExitInWarning: false, Wg: &tcp_wg}
	var tc_client1 connection.TCPClient = connection.TCPClient{Port: 10023, Host: "localhost", Description: "TC", Debug: true, ExitInWarning: false, Wg: &tcp_wg}
	var tc_client2 connection.TCPClient = connection.TCPClient{Port: 10026, Host: "localhost", Description: "TC", Debug: true, ExitInWarning: false, Wg: &tcp_wg}
	var tc_client3 connection.TCPClient = connection.TCPClient{Port: 10024, Host: "localhost", Description: "TC", Debug: true, ExitInWarning: false, Wg: &tcp_wg}

	connections_list[0] = &tc_client
	connections_list[1] = &tc_client1
	connections_list[2] = &tc_client2
	connections_list[3] = &tc_client3
	connections_list[4] = &tm_client
	connections_list[5] = &tm_client1
	connections_list[6] = &tm_client2
	connections_list[7] = &tm_client3
	connections_list[8] = &tm_client4

	for _, con := range connections_list {
		con.Connect()
	}

	for i, con := range connections_list {
		fmt.Println("Starting client", i)
		if i == 0 || i == 1 || i == 2 || i == 3 {
			con.StartRead(TC_data)
		} else {
			con.StartWrite(TM_data)
		}
	}

	var TM_decoder tc_decoder.TCDecoder = tc_decoder.TCDecoder{Run: true, ExitInWarning: false, Wg: &tcp_wg}
	TM_decoder.Start(TC_data, TC_decoded_data)
	var TM_responder tm_responcer.TMResponder = tm_responcer.TMResponder{Run: true, Wg: &tcp_wg}
	TM_responder.Start(TM_data, TC_decoded_data)

	cli()
	fmt.Println("Disconnecting...")
	TM_decoder.Stop()
	TM_responder.Stop()
	tcp_wg.Wait()
	fmt.Println("Done")
}

func cli() {
	var command string
	for {
		command = print_and_get(">> ")
		fmt.Println(command)
		switch command {
		case "exit":
			for _, con := range connections_list {
				con.Disconnect()
			}
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("exit - exit the program")
			fmt.Println("help - print this help")
		default:
			fmt.Println("Unknown command")
		}

		if command == "exit" {
			fmt.Println("Exiting...")
			break
		}
	}
}

func print_and_get(print string) string {
	var str string = ""
	fmt.Print(print)
	fmt.Scanln(&str)
	return strings.Replace(str, print, "", 1)
}
