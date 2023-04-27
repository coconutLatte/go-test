package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"wh-test/log"
)

func main() {
	t()
}

func t() {
	engine := gin.Default()

	engine.GET("/ws", func(c *gin.Context) {
		upgrade := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				//检查origin 限定websocket 被其他的域名访问
				return true
			},
		}

		protocol := c.Request.Header.Get("Sec-Websocket-Protocol")

		ws, err := upgrade.Upgrade(c.Writer, c.Request, http.Header{
			"Sec-Websocket-Protocol": {protocol},
		})
		if err != nil {
			log.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		fullData, _ := ioutil.ReadFile("./video_record")

		parseNalsAndSend(fullData, ws)
	})

	engine.Run(":1234")
}

type queue struct {
	mu  *sync.Mutex
	buf []byte
}

func NewQueue() *queue {
	return &queue{
		mu:  &sync.Mutex{},
		buf: make([]byte, 0),
	}
}

func (q *queue) append(p []byte) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.buf = append(q.buf, p...)
}

func (q *queue) pop() []byte {
	q.mu.Lock()
	defer q.mu.Unlock()

	sc, exist, err := findStartCode(q.buf)
	if err != nil {
		log.Error(err)
		return nil
	}

	if !exist {
		//log.Info("not exist")
		return nil
	}

	log.Info(sc)

	if sc.startIndex == 0 {
		q.buf = q.buf[sc.endIndex+1:]
		return nil
	}
	b := make([]byte, sc.startIndex)

	// pop front
	copy(b, q.buf[:sc.startIndex])

	// delete separator and front in q.buf
	q.buf = q.buf[sc.endIndex+1:]

	return b
}

type startCode struct {
	startIndex int
	endIndex   int
	length     int
}

func findStartCode(raw []byte) (startCode, bool, error) {
	index := 0
	found := false

	// start code is 00 00 00 01 or 00 00 01
	rank := 0

	reader := bytes.NewReader(raw)

loop:
	for ; index < len(raw); index++ {
		b, err := reader.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return startCode{}, false, fmt.Errorf("read byte failed, %w", err)
		}

		switch b {
		case 0x00:
			if rank <= 2 {
				rank++
			} else {
				return startCode{}, false, fmt.Errorf("invalid stream format contains more than 3 consecutive zeros")
			}

		case 0x01:
			if rank == 2 || rank == 3 {
				found = true
				break loop
			} else {
				rank = 0
			}

		default:
			rank = 0
		}

	}

	if !found {
		return startCode{}, false, nil
	}

	return startCode{
		startIndex: index - rank,
		endIndex:   index,
		length:     rank + 1,
	}, true, nil
}
