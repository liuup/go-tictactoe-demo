package game

import (
	"errors"
	// "log"
	// "strconv"
	// "main/controllers"
)

// 双方 棋子
const BLACK_PLAYER = 1
const WHITE_PLAYER = -1

type GameStatus struct {
	Board     [][]int // 棋盘状态
	Finalwin  int
	Stepcount int
}

type Player struct {
	Identity int // 黑方或者白方
}

func GameInit() *GameStatus {
	return &GameStatus{
		Board: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		Finalwin:  0,
		Stepcount: 0,
	}
}

func PlayerInit(identity int) *Player {
	return &Player{
		Identity: identity,
	}
}

// 初始化棋局
var (
	Game = GameInit()
)

/*
棋盘坐标如下：

	{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
*/
func (g GameStatus) ChangeStatus(player Player, coord int) error {
	x, y := 0, 0
	if coord == 1 {
		x, y = 0, 0
	} else if coord == 2 {
		x, y = 0, 1
	} else if coord == 3 {
		x, y = 0, 2
	} else if coord == 4 {
		x, y = 1, 0
	} else if coord == 5 {
		x, y = 1, 1
	} else if coord == 6 {
		x, y = 1, 2
	} else if coord == 7 {
		x, y = 2, 0
	} else if coord == 8 {
		x, y = 2, 1
	} else if coord == 9 {
		x, y = 2, 2
	} else {
		return errors.New("坐标选择错误！")
	}

	if g.Board[x][y] != 0 {
		return errors.New("当前位置已有棋子！")
	}

	if player.Identity == BLACK_PLAYER {
		g.Board[x][y] = BLACK_PLAYER
	} else if player.Identity == WHITE_PLAYER {
		g.Board[x][y] = WHITE_PLAYER
	}

	return nil
}

// 检查棋盘当前是否存在赢家
func (g *GameStatus) CheckWin() bool {
	// 检查3个横排
	for i := 0; i < len(g.Board); i++ {
		if g.Board[i][0] == g.Board[i][1] && g.Board[i][1] == g.Board[i][2] && g.Board[i][0] != 0 {
			g.Finalwin = g.Board[i][0]
			return true
		}
	}

	// 检查3个竖排
	for i := 0; i < len(g.Board); i++ {
		if g.Board[0][i] == g.Board[1][i] && g.Board[1][i] == g.Board[2][i] && g.Board[0][i] != 0 {
			g.Finalwin = g.Board[0][i]
			return true
		}
	}

	// 检查主对角线
	if g.Board[0][0] == g.Board[1][1] && g.Board[1][1] == g.Board[2][2] && g.Board[0][0] != 0 {
		g.Finalwin = g.Board[0][0]
		return true
	}

	// 检查副对角线
	if g.Board[0][2] == g.Board[1][1] && g.Board[1][1] == g.Board[2][0] && g.Board[0][2] != 0 {
		g.Finalwin = g.Board[0][2]
		return true
	}

	return false
}
