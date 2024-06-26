package utils

import (
	"SSH/model"
	"github.com/rs/xid"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
)

//本文件定义了在通过ssh的方式连接虚拟机时用到的工具函数

func VerifyConnect(vmInfo *model.VMConnectRequest) (bool, string) {
	sshClient, _ := CreateSSHClient(vmInfo, ssh.Password(vmInfo.Password))
	defer sshClient.Close()
	if tempSession, err := sshClient.NewSession(); err == nil {
		defer tempSession.Close()
		//判断是否能够连接成功
		if err := tempSession.Run("whoami"); err == nil {
			//连接成功
			//生成唯一的SSHConnectID
			// 生成一个新的 xid
			myID := xid.New()
			// 将 xid 转换为字符串
			connectID := myID.String()
			return true, connectID
		} else {
			//连接不成功，返回错误
			return false, ""
		}
	} else {
		tempSession.Close()
		return false, ""
	}

}

// CreateSSHClient  创建SSH连接的客户端
func CreateSSHClient(vmInfo *model.VMConnectRequest, auth ssh.AuthMethod) (*ssh.Client, error) {

	//var hostKey ssh.PublicKey

	//实现认证，目前只有密码认证，TODO 后续可以添加其他认证方式（如密钥认证）
	config := &ssh.ClientConfig{
		User: vmInfo.Username,
		Auth: []ssh.AuthMethod{
			auth,
		},
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			//十分危险的方法，仅用于当前测试，无论远程服务器的主机密钥是什么，都会被接受，不会引发错误
			return nil
		},
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(vmInfo.Host, vmInfo.Port), config)
	if err != nil {
		log.Fatalf("001.创建SSHClient失败！ %s", err)
		return nil, err
	}
	log.Println("成功创建ssh.client")
	return client, nil
}
