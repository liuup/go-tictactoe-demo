package main

import (
	"main/router"
)

// 一个基于websocket的井字棋对战后端
func main() {
	r := router.Router()

	r.Run(":8080")
}
