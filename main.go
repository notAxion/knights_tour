package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	verbose := flag.Bool("v", false, "[v]erbose prints the board after tour completes")
	sizeStr := flag.String("size", "1x1", "[size] of the tour \n\teg. --size 5x5")
	flag.Parse()

	if *sizeStr == "1x1" {
		fmt.Println("add a board size \neg. --size 5x5")
		return
	}
	boardSize, err := getBoardSize(*sizeStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// gg := heuristicSolve(boardSize)
	gg := bruteForce(boardSize)
	if *verbose {
		fmt.Println(gg)
	}
}

func bruteForce(boardSize pos) *game {
	gg := newGame(boardSize)

	backtracking(gg, NewPos(0, 0))
	fmt.Println("game over you visited", gg.visits)

	return gg
}

func backtracking(gg *game, to pos) bool {
	gg.moveKnight(to)
	// fmt.Printf("%s\n", "\n"+gg.String())
	if gg.visits >= gg.boardArea {
		return true
	}

	moves := gg.findMoves()
	for _, move := range moves {
		if backtracking(gg, move) {
			return true
		}
		gg.undoMove()
		// fmt.Printf("%s\n", "\n"+gg.String())
	}

	return false
}

func heuristicSolve(boardSize pos) *game {
	gg := newGame(boardSize)
	gg.moveKnight(NewPos(0, 0))

	moves := gg.findMoves()
	for i := 0; i < gg.boardArea*2; i++ {
		if gg.visits >= gg.boardArea || len(moves) == 0 {
			break
		}

		nextMove := pos{}
		min := 9
		for _, move := range moves {
			futureMoves := gg.findMovesFrom(move)
			if min > len(futureMoves) {
				min = len(futureMoves)
				moves = futureMoves
				nextMove = move
			}
		}

		gg.moveKnight(nextMove)
	}
	fmt.Println("game over you visited:", gg.visits)

	return gg
}

type game struct {
	board [][]int
	// boardSize int
	boardSize pos
	boardArea int
	knight    pos

	// how many quares the knight has visited
	visits int

	undoStack []pos
}

func newGame(boardSize pos) *game {
	g := &game{
		board:     newBoard(boardSize),
		boardSize: boardSize,
		boardArea: boardSize.x * boardSize.y,
		knight:    NewPos(-1, -1),
		undoStack: make([]pos, 0),
	}

	return g
}

func (g *game) moveKnight(move pos) {
	// TODO in gui this will hinder
	// undo before the 2nd move
	if g.knight != NewPos(-1, -1) {
		g.undoStack = append(g.undoStack, g.knight)
	}

	g.visits++
	g.board[move.x][move.y] = g.visits
	g.knight = move
}

func (g *game) undoMove() {
	if len(g.undoStack) == 0 {
		return
	}

	move := g.undoStack[len(g.undoStack)-1]
	g.undoStack = g.undoStack[:len(g.undoStack)-1]

	g.visits--
	g.board[g.knight.x][g.knight.y] = 0
	g.knight = move
}

func (g *game) findMovesFrom(from pos) []pos {
	moves := []pos{}

	for _, movement := range movements {
		move := from.Add(movement)
		if g.inside(move) && g.board[move.x][move.y] == 0 {
			moves = append(moves, move)
		}
	}

	return moves
}

func (g *game) findMoves() []pos {
	return g.findMovesFrom(g.knight)
}

func (g *game) inside(move pos) bool {
	if move.x < 0 || move.y < 0 {
		return false
	} else if move.x >= g.boardSize.x || move.y >= g.boardSize.y {
		return false
	}

	return true
}

func (g *game) String() string {
	var str strings.Builder
	str.Grow(g.visits)
	d := fmt.Sprint(int(math.Log10(float64(g.visits)) + 2))
	for _, row := range g.board {
		str.WriteByte('|')
		for _, col := range row {
			str.WriteString(fmt.Sprintf("%"+d+"d|", col))
		}
		str.WriteByte('\n')
	}

	return str.String()
}

func newBoard(boardSize pos) [][]int {
	board := make([][]int, boardSize.x)

	for i := range board {
		board[i] = make([]int, boardSize.y)
	}

	return board
}

func getBoardSize(sizeStr string) (pos, error) {
	numsStr := strings.Split(sizeStr, "x")
	num1, err := strconv.Atoi(numsStr[0])
	if err != nil {
		return pos{}, err
	}
	num2, err := strconv.Atoi(numsStr[1])
	if err != nil {
		return pos{}, err
	}

	return NewPos(num1, num2), nil
}

var movements = []pos{
	NewPos(1, 2), NewPos(2, 1),
	NewPos(1, -2), NewPos(2, -1),
	NewPos(-1, -2), NewPos(-2, -1),
	NewPos(-1, 2), NewPos(-2, 1),
}

type pos struct {
	x, y int
}

func NewPos(x, y int) pos {
	return pos{
		x: x,
		y: y,
	}
}

func (p pos) Add(to pos) pos {
	return p.AddXY(to.x, to.y)
}

func (p pos) AddXY(x, y int) pos {
	return NewPos(p.x+x, p.y+y)
}
