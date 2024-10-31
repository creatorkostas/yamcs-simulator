package tm_responcer

import (
	"acubesat/ops/yamcs-simulator/internal/tc_decoder"
	"acubesat/ops/yamcs-simulator/internal/tools"
	"encoding/hex"
	"fmt"
	"sync"
)

// var tm_17_2 string = "0801c000000a201102000000010f2bd13a"
var loaded_commands tools.Commands
var command_map = map[string]tools.Command{}

type TMResponder struct {
	Commands_yaml_path string
	Run                bool
	Debug              bool
	Wg                 *sync.WaitGroup
}

func hexToBytes(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return b
}

// func RespondTM(tc tc_decoder.TCData, TM_data chan []byte) {
// 	if tc.ServiceTypeID == 17 && tc.MessageSubtypeID == 1 {
// 		TM_data <- hexToBytes(tm_17_2)
// 	}
// }

func (tm_responder *TMResponder) Start(TM_data chan []byte, decoded_channel chan tc_decoder.TCData) error {
	loaded_commands = tools.Load_configs(tm_responder.Commands_yaml_path)
	for _, item := range loaded_commands.CommandList {
		command_map[item.Name] = item
	}

	loaded_commands.CommandList = nil

	go func(tm_responder *TMResponder) {
		fmt.Println("Start TM Responder")
		var run bool = tm_responder.Run
		var debug bool = tm_responder.Debug
		var data string
		var data_bytes []byte
		for run {
			if !tm_responder.Run {
				break
			}

			tc := <-decoded_channel

			data = command_map["TC("+fmt.Sprintf("%v", tc.ServiceTypeID)+","+fmt.Sprintf("%v", tc.MessageSubtypeID)+")"].TM
			data_bytes = hexToBytes(data)
			if debug {
				fmt.Println(data_bytes)
			}
			TM_data <- data_bytes

		}
		fmt.Println("Stop TM Responder")
		tm_responder.Wg.Done()
	}(tm_responder)
	return nil
}

func (tm_responder *TMResponder) Stop() {
	if tm_responder.Run {
		tm_responder.Run = false
	}
}
