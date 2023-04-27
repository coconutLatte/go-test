package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/gorilla/websocket"

	"wh-test/log"
)

func parseNalsAndSend(fullData []byte, wsConn *websocket.Conn) {
	nalus, err := splitNalus(fullData)
	if err != nil {
		log.Error(err)
		return
	}

	//f, _ := os.Create("mini_video_record")
	//defer f.Close()
	count := 1
	for _, nalu := range nalus {
		log.Info(nalu[:5])
		log.Info(count)

		if checkIfBad(nalu) {
			log.Info("catch it")
		}
		log.Info("")

		//if count <= 4 {
		//	f.Write(nalu)
		//}

		wsConn.WriteMessage(websocket.BinaryMessage, nalu)
		count++
	}
}

func checkIfBad(nalu []byte) bool {
	return bytes.Contains(nalu, []byte{0x00, 0x00, 0x1c, 0x02, 0xa5, 0x87, 0x00, 0x01, 0x00, 0x00, 0x03, 0x00, 0xe0})
}

func splitNalus(full []byte) ([][]byte, error) {
	reader := bytes.NewReader(full)

	res := make([][]byte, 0)

	rank := 0
	nalu := make([]byte, 0)
	isFirst := true
	for i := 0; i < len(full); i++ {
		b, err := reader.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				res = append(res, nalu)
				break
			}

			log.Error(err)
			return nil, err
		}

		nalu = append(nalu, b)
		switch b {
		case 0x00:
			if rank <= 2 {
				rank++
			} else {
				return nil, fmt.Errorf("invalid stream format contains more than 3 consecutive zeros")
			}
		case 0x01:
			if rank == 2 {
				if isFirst {
					isFirst = false
				} else {
					tmp := make([]byte, len(nalu)-3)
					copy(tmp, nalu[:len(nalu)-3])
					res = append(res, nalu[:len(nalu)-3])
					nalu = nalu[len(nalu)-3:]
				}
			}

			if rank == 3 {
				if isFirst {
					isFirst = false
				} else {
					//res := make([]byte, 0)
					tmp := make([]byte, len(nalu)-4)
					copy(tmp, nalu[:len(nalu)-4])
					res = append(res, nalu[:len(nalu)-4])
					nalu = nalu[len(nalu)-4:]
				}
			}

			rank = 0
		default:
			rank = 0
		}
	}

	return res, nil
}
