package main

import (
	"SSH/model"
	"SSH/utils"
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 不推荐在生产中使用
	},
}

func sshConnect(w http.ResponseWriter, r *http.Request) {

	host := r.URL.Query().Get("host")
	port := r.URL.Query().Get("port")
	password := r.URL.Query().Get("password")
	username := r.URL.Query().Get("username")

	vmInfo := model.VMConnectRequest{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}

	config := &ssh.ClientConfig{
		User: vmInfo.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(vmInfo.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//测试连接
	flag, connctID := utils.VerifyConnect(&vmInfo)
	if flag != true {
		//todo 返回错误码
		log.Fatalln("001测试连接失败")
		return
	}
	//连接成功

	//创建客户端与会话
	conn, err := ssh.Dial("tcp", vmInfo.Host+":"+vmInfo.Port, config)
	if err != nil {
		log.Fatalf("002unable to connect: %v", err)
		//todo 返回错误码
		return
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("003unable to create session: %v", err)
		//todo 返回错误码
		return
	}
	defer session.Close()

	//设置流
	stdinPipe, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
		//todo 返回错误码
	}
	defer stdinPipe.Close()

	// Prepare pipes for capturing output
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		//todo 返回错误码
	}
	stderrPipe, err := session.StderrPipe()
	if err != nil {
		log.Fatal(err)
		//todo 返回错误码
	}

	//升级为websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	// Request a pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, ssh.TerminalModes{}); err != nil {
		log.Fatalf("request for pseudo terminal failed: %v", err)
		//todo 返回错误码
	}

	// Start a shell
	if err := session.Shell(); err != nil {
		log.Fatal(err)
		//todo 返回错误码
	}

	commands := []string{" "}
	for _, cmd := range commands {
		_, err = stdinPipe.Write([]byte(cmd))
		if err != nil {
			log.Fatal(err)
		}
	}

	/// Custom buffers
	var stdoutBuf model.OutputDataBuffer
	// Capture stdout and stderr
	go io.Copy(&stdoutBuf, stdoutPipe)
	go io.Copy(&stdoutBuf, stderrPipe)

	//将虚拟机的初始化信息输出到ws上
	//延时
	time.Sleep(50 * time.Millisecond)
	stdoutBuf.Flush(ws, connctID)

	// 接受前端的数据

	var inputBuffer model.OutputDataBuffer

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var cmdData model.ShellRequest
		err = json.Unmarshal(message, &cmdData)
		if err != nil {
			log.Println("前端请求对象解析错误:", err)
			continue // 解析失败时继续监听下一条消息
		}
		// 将接收到的消息放入缓冲区
		inputBuffer.Write([]byte(cmdData.Command))

		//将缓冲区的数据写入虚拟机
		inputBuffer.WriteToVM(stdinPipe)

		//延时
		time.Sleep(50 * time.Millisecond)

		// 虚拟机的数据写入websocket
		stdoutBuf.Flush(ws, connctID)

	}

}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	var inputBuffer bytes.Buffer

	conn.WriteMessage(websocket.TextMessage, []byte("开始websocket"))

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		// 将接收到的消息放入缓冲区
		inputBuffer.Write(message)

		var operationBuffer bytes.Buffer

		// inputBuffer数据被读入operationBuffer
		if _, err := operationBuffer.Write(inputBuffer.Bytes()); err != nil {
			log.Println("error copying input buffer to operation buffer:", err)
			break
		}
		inputBuffer.Reset()

		// 假设我们立即处理消息并准备响应
		// 在实际应用中，你可能需要根据消息内容进行一些处理
		time.Sleep(1 * time.Second) // 模拟处理延时

		// operationBuffer的数据写入websocket
		if err := conn.WriteMessage(websocket.TextMessage, operationBuffer.Bytes()); err != nil {
			log.Println("write:", err)
			break
		}

		// 清空缓冲区，为下一次消息接收准备
		inputBuffer.Reset()
	}
}

func main() {
	http.HandleFunc("/api/terminal/sshConnect", sshConnect)
	log.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
