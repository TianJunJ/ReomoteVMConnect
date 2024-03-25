package model

import (
	"bytes"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"sync"
)

//本文件定义了与虚拟机进行交互的相关对象

// 线程安全的缓冲区
type OutputDataBuffer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (w *OutputDataBuffer) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

func (w *OutputDataBuffer) Flush(ws *websocket.Conn, connectID string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.buffer.Len() != 0 {
		//err := ws.WriteJSON(map[string]string{"Code": "200", "connectID": connectID, "cmdResponseData": w.buffer.String(), "message": ""})
		err := ws.WriteJSON(map[string]string{"Code": "200", "connectID": connectID, "resData": w.buffer.String(), "message": ""})
		if err != nil {
			return err
		}
		w.buffer.Reset()
	}
	log.Printf("缓冲区数据已经刷新到WebSocket连接中")
	return nil
}

func (w *OutputDataBuffer) WriteToVM(stdinPipe io.WriteCloser) {
	w.mu.Lock()
	defer w.mu.Unlock()
	stdinPipe.Write(w.buffer.Bytes())
	w.buffer.Reset()
}
