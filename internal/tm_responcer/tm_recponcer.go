package tm_responcer

import (
	"acubesat/ops/yamcs-simulator/internal/tc_decoder"
	"encoding/hex"
	"fmt"
	"sync"
)

var tm_17_2 string = "0801c000000a201102000000010f2bd13a"

type TMResponder struct {
	Run bool
	Wg  *sync.WaitGroup
}

func hexToBytes(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return b
}

func RespondTM(tc tc_decoder.TCData, TM_data chan []byte) {
	if tc.ServiceTypeID == 17 && tc.MessageSubtypeID == 1 {
		TM_data <- hexToBytes(tm_17_2)
	}
}

func (tm_responder *TMResponder) Start(TM_data chan []byte, decoded_channel chan tc_decoder.TCData) error {
	go func(tm_responder *TMResponder) {
		fmt.Println("Start TM Responder")
		var run bool = tm_responder.Run

		for run {
			if !tm_responder.Run {
				break
			}

			tc := <-decoded_channel

			if tc.ServiceTypeID == 17 && tc.MessageSubtypeID == 1 {
				TM_data <- hexToBytes(tm_17_2)
			}

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
