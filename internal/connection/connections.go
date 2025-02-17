package connection

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup
var m = sync.RWMutex{}

type TCPClient struct {
	Port          int
	Host          string
	Description   string
	Debug         bool
	Ok            bool
	Conn          *net.TCPListener
	ExitInWarning bool
	Run           bool
	Wg            *sync.WaitGroup
}

func (tcp_client *TCPClient) GetRun() bool {
	return tcp_client.Run
}

func (tcp_client *TCPClient) Init(port int, host string, description string, debug bool, exitInWarning bool, wg *sync.WaitGroup) error {
	tcp_client.Port = port
	tcp_client.Host = host
	tcp_client.Description = description
	tcp_client.Debug = debug
	tcp_client.ExitInWarning = exitInWarning
	tcp_client.Conn = nil
	tcp_client.Ok = false
	tcp_client.Run = false
	tcp_client.Wg = wg
	return nil
}

func (tcp_client *TCPClient) Connect() bool {
	fmt.Println("Listening: ", tcp_client.Description, " in ", tcp_client.Host+":", fmt.Sprint(tcp_client.Port))
	tcpAddr, err := net.ResolveTCPAddr("tcp", tcp_client.Host+":"+fmt.Sprint(tcp_client.Port))
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		if tcp_client.ExitInWarning {
			os.Exit(1)
		}
	}

	conn, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		println("ListenTCP failed:", err.Error())
		if tcp_client.ExitInWarning {
			os.Exit(1)
		}
	} else {
		tcp_client.Conn = conn
		tcp_client.Ok = true
	}

	return tcp_client.Ok
}

func (tcp_client *TCPClient) Disconnect() {
	if tcp_client.Ok || tcp_client.Conn != nil || tcp_client.Run {
		fmt.Println("Disconnecting ", tcp_client.Description, " in ", tcp_client.Host, ":", fmt.Sprint(tcp_client.Port))
		tcp_client.Ok = false
		tcp_client.Run = false
		tcp_client.Conn.Close()
		tcp_client.Conn = nil
	}
}

func (tcp_client *TCPClient) StartRead(channel chan []byte) error {
	tcp_client.Run = true
	tcp_client.Wg.Add(1)
	go func(tcp_client *TCPClient) {
		fmt.Println("TCP_Read: Started")
		m.Lock()
		var run bool = tcp_client.Ok
		m.Unlock()
		var tmp []byte = make([]byte, 1024)
		// var buff *bytes.Buffer

		c, err := tcp_client.Conn.Accept()
		if err != nil {
			println("Accept failed:", err.Error())
			if tcp_client.ExitInWarning {
				os.Exit(1)
			}
		}
		for run {
			if !tcp_client.Ok {
				break
			}

			i, err := c.Read(tmp)
			if err != nil {
				println("Read failed:", err.Error())
				if tcp_client.ExitInWarning {
					os.Exit(1)
				}
			} else {
				fmt.Println("Read:", i, " bytes")
			}

			if channel != nil {
				m.Lock()
				channel <- tmp
				m.Unlock()
			}

			if tcp_client.Debug {
				fmt.Println("TC:", hex.EncodeToString(tmp)[:i*2])
			}

		}
		fmt.Println("TCP_Read: Disconnected")
		tcp_client.Wg.Done()
	}(tcp_client)
	return nil
}

func (tcp_client *TCPClient) StartWrite(channel chan []byte) {
	tcp_client.Run = true
	tcp_client.Wg.Add(1)
	go func(tcp_client *TCPClient) {
		fmt.Println("TCP_Write: Started")
		m.Lock()
		var run bool = tcp_client.Ok
		m.Unlock()

		c, err := tcp_client.Conn.Accept()
		if err != nil {
			println("Accept failed:", err.Error())
			if tcp_client.ExitInWarning {
				os.Exit(1)
			}
		}

		for run {
			if !tcp_client.Ok {
				break
			}

			data := <-channel

			if tcp_client.Ok && err == nil {
				_, err := c.Write(data)
				if err != nil {
					println("Write to server failed:", err.Error())
				}
			}

		}
		fmt.Println("TCP_Write: Disconnected")
		tcp_client.Wg.Done()
	}(tcp_client)

}
