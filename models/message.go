package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	FromId   int64  //发送者
	TargetId int64  //接收者
	Type     int    //发送类型 群聊 私聊 广播
	Media    int    //消息类型 文字 图片 音视频
	Content  string //消息内容
	Pic      string //图片
	Url      string //Url地址
	Desc     string //描述信息
	Amount   int    //其他数字统计
}

func (table *Message) TableName() string {
	return "message"

}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

//  需要: 发送者 接收者 消息类型 消息内容 发送类型

func Chat(write http.ResponseWriter, request *http.Request) {
	//1. 获取参数并校验token 等合法性
	//token := query.Get("token")
	query := request.URL.Query()
	userId, _ := strconv.ParseInt(query.Get("userId"), 10, 64)
	//msgType := query.Get("type")
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	isvalida := true //check Token 待完成

	conn, err := (&websocket.Upgrader{
		//Token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(write, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//2. 获取链接conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	//3.用户关系

	//4.userid 跟 node绑定 并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	//5.完成发送的逻辑
	go sendProc(node)
	//6.完成接收的逻辑
	go recvProc(node)

	sendMsg(userId, []byte("欢迎进入聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}

		}
	}
}
func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data) //广播消息
		fmt.Println("[ws] <<< ", data)
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成udp数据发送协程
func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(10, 211, 0, 1),
		Port: 3000,
		Zone: "",
	})
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			_, err1 := conn.Write(data)
			if err1 != nil {
				fmt.Println(err1)
				return
			}

		}
	}
}

// 完成udp数据接收协程
func udpRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	for {
		var buf [512]byte
		n, err2 := conn.Read(buf[0:])
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		sendMsg(msg.TargetId, data)
		//case 2:
		//	sendGroupMsg()
		//case 3:
		//	sendAllMsg()
		//case 4:
		//默认

	}

}

func sendMsg(targetId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[targetId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
