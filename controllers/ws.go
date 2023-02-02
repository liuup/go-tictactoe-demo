package controllers

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"net/http"

	"main/game"
	"main/utils/token"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

/*
1. 建立多房间
2. 检测连接数量

*/

type WsServe struct {
	MsgChan  chan Msg
	WsConnMp map[Client]bool
	mutex    sync.Mutex // 哈希表互斥锁
}

type Client struct {
	Id           int
	Address      string
	Ws           *websocket.Conn
	ClientPlayer *game.Player
}

type Msg struct {
	Type     int     // 消息类型
	LastStep int     // 消息具体内容
	Content  string  // 传递的额外信息
	Sender   string  // 消息的发送方
	Board    [][]int // 棋盘状态
}

// 初始化处理ws连接服务
var (
	wsserve = wsServeInit()
)

// WsServe初始化
func wsServeInit() *WsServe {
	return &WsServe{
		MsgChan:  make(chan Msg, 100),
		WsConnMp: make(map[Client]bool),
		mutex:    sync.Mutex{},
	}
}

func WsStart(c *gin.Context) {
	// 权限校验
	_, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 允许跨域
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 升级协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}

	// 首先获取输入，选择黑白方
	_, p, err := ws.ReadMessage()
	if err != nil {
		log.Println(err)
	}
	pr, _ := strconv.Atoi(string(p))

	// 创建客户端
	client := Client{
		Id:           -1,
		Address:      c.Request.RemoteAddr,
		Ws:           ws,
		ClientPlayer: &game.Player{pr},
	}
	log.Println(client, client.ClientPlayer.Identity)

	// log.Println(client)

	// 从客户端读线程
	go read(client)

	// 往客户端写线程
	go write()

	// for {

	// }

	// wsserve.mutex.Unlock()

	// wsconn[ws] = true

	// 需要检测ws连接数量，两个的时候才能参加对战

	// 消息处理线程
	// go readMsg(ws)

	// 广播消息线程
	// go boardcastMsg()

	// for {
	// 	time.Sleep(1 * time.Second)
	// 	fmt.Println(len(wsserve.MsgChan), wsserve.WsConnMp)
	// }

	// go getSystemStatus()

	// 防止主线程结束
	select {}
}

// 从客户端读消息
func read(client Client) {
	defer disconnect(client)

	// 限制在线人数
	err := handleOnlinePeoples(client)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		msgType, coord, err := client.Ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// 转为坐标
		laststep, err := strconv.Atoi(string(coord))
		if err != nil {
			log.Println(err)
			break
		}

		// 改变棋盘状态
		err = game.Game.ChangeStatus(*client.ClientPlayer, laststep)
		if err != nil {
			log.Println(err)
			// 再次尝试获取玩家输入
			continue
		}

		if game.Game.CheckWin() {
			// 向双方广播胜者
			content := ""
			if game.Game.Finalwin == game.BLACK_PLAYER {
				content = "Black win"
			} else if game.Game.Finalwin == game.WHITE_PLAYER {
				content = "White win"
			}

			msg := Msg{
				Type:     msgType,
				LastStep: laststep,
				Content:  content,
				Sender:   client.Address,
				Board:    game.Game.Board,
			}
			writeToAll(msg)
			disconnectAll()
		}
		// 是否超过棋局限制
		game.Game.Stepcount++
		if game.Game.Stepcount == 9 {
			disconnectAll() // 断掉所有连接
			return
		}

		// 正常情况写入消息
		wsserve.MsgChan <- Msg{
			Type:     msgType,
			LastStep: laststep,
			Content:  "",
			Sender:   client.Address,
			Board:    game.Game.Board,
		}
	}
}

// 往客户端广播消息
func write() {
	for {
		msg := <-wsserve.MsgChan
		// log.Println(msg)
		for client := range wsserve.WsConnMp {
			// 不给自己广播消息
			if client.Address != msg.Sender {
				// 直接把结构体转为json
				client.Ws.WriteJSON(msg)
			}
		}
	}
}

// 限制在线人数
func handleOnlinePeoples(client Client) error {
	// 当前在线人数限制为2人
	if len(wsserve.WsConnMp) >= 2 {
		disconnect(client)
		return errors.New("在线人数已满")
	}

	wsserve.WsConnMp[client] = true
	return nil
}

// 一次性向所有连接写消息
func writeToAll(msg Msg) {
	for client := range wsserve.WsConnMp {
		client.Ws.WriteJSON(msg)
	}
}

// 关闭特定客户端连接
func disconnect(client Client) {
	client.Ws.Close() // 关闭连接
	wsserve.mutex.Lock()
	delete(wsserve.WsConnMp, client) // 删除
	wsserve.mutex.Unlock()
}

// 关闭所有客户端连接
func disconnectAll() {
	for client := range wsserve.WsConnMp {
		client.Ws.Close()
		wsserve.mutex.Lock()
		delete(wsserve.WsConnMp, client)
		wsserve.mutex.Unlock()
	}
}

// 获取系统状态信息
func getSystemStatus() {
	fmt.Println(len(wsserve.MsgChan), wsserve.WsConnMp)
}

// 接收来自客户端的消息
// func readMsg(ws *websocket.Conn) {
// 	defer func() {
// 		ws.Close()
// 		delete(wsconn, ws)
// 		// wsconn[ws] = false
// 	}()

// 	for {
// 		msg_type, sb, err := ws.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			break
// 		}
// 		fmt.Println(msg_type, sb)

// 		msg_chan <- Msg{msg_type, sb, ws.RemoteAddr().String()}
// 	}
// }

// 广播消息处理
// func boardcastMsg() {
// 	for {
// 		msg := <-msg_chan
// 		for ws, v := range wsconn {
// 			if v {
// 				// 不给自己广播消息
// 				if ws.RemoteAddr().String() == msg.Sender {
// 					continue
// 				}
// 				ws.WriteMessage(msg.Type, msg.Content)
// 			}
// 		}
// 	}
// }
