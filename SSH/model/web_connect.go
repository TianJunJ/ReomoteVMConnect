package model

//本文件封装了用于网络的所有对象（网络请求对象、网络响应对象）

/*
VMConnectRequest 请求对象。该对象封装了移动端（前端）在进行远程连接时要向后端传递的参数
Host: 进行虚拟机连接时要连接的虚拟机地址 （如 192.168.163.129）
Port: 进行连接时要连接的虚拟机端口（默认22）
Username: 进行连接时 登录虚拟机的哪个用户
Password: 该用户的密码
todo 增设 ConnectType 用户可以选择连接方式
*/
type VMConnectRequest struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	//todo ConnectType string `json:"connectType"` 现在默认为ssh
}

type VMConnectResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

/*
ShellRequest  该对象封装了用户在与 远程虚拟机 交互时的Linux shell指令。该对象会在WebSocket中传输
Command 封装了交互时的Linux shell指令
*/
type ShellRequest struct {
	Command string `json:"data"`
}

/*
ShellResponse  该对象封装了 用户的Linux shell指令被虚拟机运行后的结果
ConnectID 每个连接都有一个唯一的ID
CmdResponse 虚拟机在执行完指令后的返回结果
Code 状态响应码 200-成功  400-不成功
Message 响应的信息 防止与虚拟机的连接突然中断或者出现其他异常
*/
type ShellResponse struct {
	ConnectID       string `json:"connectID"`
	CmdResponseData string `json:"cmdResponseData"`
	Code            int    `json:"code"`
	Message         string `json:"message"`
}
