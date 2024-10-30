package tc_decoder

import (
	"fmt"
	"os"
	"sync"
)

var m = sync.RWMutex{}

type TCData struct {
	TCHeader
}

type TCHeader struct {
	PacketVersionNumber  uint8  // 3
	PacketType           uint8  // 1
	SecondaryHeaderFlag  uint8  // 1
	ApplicationProcessID uint16 // 11 unsigned
	SequenceFlags        uint8  //2
	PacketSequenceCount  uint16 //14
	PacketDataLength     uint16 //16
	TCSecondaryHeader
}

type TCSecondaryHeader struct {
	PUSVersionNumber uint8  // 4 unsigned
	CompletionAck    uint8  // 1
	ProgressAck      uint8  // 1
	StartAck         uint8  // 1
	AcceptanceAck    uint8  // 1
	ServiceTypeID    uint8  // 8 unsigned
	MessageSubtypeID uint8  // 8 unsigned
	SourceID         uint16 // 16 unsigned
}

// 3 +1+1+11+2+14+16+4+1+1+1 = 55

type TCDecoder struct {
	Run           bool
	ExitInWarning bool
	Wg            *sync.WaitGroup
}

func DecodeTC(data []byte) (TCData, error) {
	var err error = nil
	tc_header := TCHeader{}
	//TODO: make bytes to bits and change DataConvert for bits
	tc_header.PacketVersionNumber = 0  //, err = tools.DataConvert[uint8](data[:3])
	tc_header.PacketType = 0           //, err = tools.DataConvert[uint8](data[3:4])
	tc_header.SecondaryHeaderFlag = 0  //, err = tools.DataConvert[uint8](data[4:5])
	tc_header.ApplicationProcessID = 0 //, err = tools.DataConvert[uint16](data[5:16])
	tc_header.SequenceFlags = 0        //, err = tools.DataConvert[uint8](data[16:18])
	tc_header.PacketSequenceCount = 0  //, err = tools.DataConvert[uint16](data[18:32])
	tc_header.PacketDataLength = 0     //, err = tools.DataConvert[uint16](data[32:48])

	tc_SecondaryHeader := TCSecondaryHeader{}
	tc_SecondaryHeader.PUSVersionNumber = 0       //, err = tools.DataConvert[uint8](data[48:52])
	tc_SecondaryHeader.CompletionAck = 0          //, err = tools.DataConvert[uint8](data[52:53])
	tc_SecondaryHeader.ProgressAck = 0            //, err = tools.DataConvert[uint8](data[53:54])
	tc_SecondaryHeader.StartAck = 0               //, err = tools.DataConvert[uint8](data[54:55])
	tc_SecondaryHeader.AcceptanceAck = 0          //, err = tools.DataConvert[uint8](data[55:56])
	tc_SecondaryHeader.ServiceTypeID = data[7]    //, err = tools.DataConvert[uint8](data[56:64])
	tc_SecondaryHeader.MessageSubtypeID = data[8] //, err = tools.DataConvert[uint8](data[64:72])
	tc_SecondaryHeader.SourceID = 0               //, err = tools.DataConvert[uint16](data[72:88])

	tc_header.TCSecondaryHeader = tc_SecondaryHeader

	tc_data := TCData{TCHeader: tc_header}
	return tc_data, err

}

func (tm_decoder *TCDecoder) Start(channel chan []byte) error {
	tm_decoder.Run = true
	tm_decoder.Wg.Add(1)
	go func(tm_decoder *TCDecoder) {
		fmt.Println("Start TC Decoder")
		var run bool = tm_decoder.Run

		for run {
			if !tm_decoder.Run {
				break
			}

			data := <-channel
			tm_data, err := DecodeTC(data)
			if err != nil {
				fmt.Println("Error decoding TC:", err.Error())
				if tm_decoder.ExitInWarning {
					os.Exit(1)
				}
			}

			if tm_decoder.Run {
				fmt.Println("TC:", tm_data.ServiceTypeID, " , ", tm_data.MessageSubtypeID)
			}

		}
		fmt.Println("Stop TC Decoder")
		tm_decoder.Wg.Done()
	}(tm_decoder)
	return nil
}

func (tm_decoder *TCDecoder) Stop() {
	if tm_decoder.Run {
		tm_decoder.Run = false
	}
}
