/*
 * @Author: liziwei01
 * @Date: 2024-05-01 20:29:10
 * @LastEditors: liziwei01
 * @LastEditTime: 2024-05-01 20:31:28
 * @Description: file content
 */
package main

import (
	"container/heap"
	"fmt"
)

const maxDepth = 100 // 设置搜索的最大深度

type State struct {
	board  [][]int // 当前的棋盘状态
	moves  int     // 到达当前状态所需要的步数
	parent *State  // 当前状态的父状态，即上一步的状态
	score  int     // A*算法中的评估分数
}

type PriorityQueue []*State // 优先队列，用于存储待搜索的状态

// 以下几个函数是实现heap.Interface接口所必需的
func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].score < pq[j].score } // 优先队列以评估分数作为优先级
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*State)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func main() {
	board := [][]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{14, 15, 13, 0},
	}
	if isSolvable(board) { // 判断棋盘是否有解
		fmt.Println("有解")
		solution := solvePuzzle(board) // 求解
		fmt.Println("最少步骤解:")
		for _, state := range solution { // 打印解的每一步
			fmt.Println("步数:", state.moves)
			printBoard(state.board)
			fmt.Println()
		}
	} else {
		fmt.Println("无解")
	}
}

func solvePuzzle(board [][]int) []*State {
	initialState := &State{board, 0, nil, 0} // 初始状态
	priorityQueue := make(PriorityQueue, 0)  // 创建优先队列
	heap.Push(&priorityQueue, initialState)  // 将初始状态加入队列
	visited := make(map[string]bool)         // 记录已经访问过的状态
	visited[serialize(board)] = true

	for len(priorityQueue) > 0 { // 当队列不为空时，继续搜索
		currentState := heap.Pop(&priorityQueue).(*State) // 取出当前状态

		if currentState.moves > maxDepth { // 如果超过最大深度，停止搜索
			fmt.Println("超过深度限制")
			return nil
		}

		if isGoalState(currentState.board) { // 如果当前状态是目标状态，返回解
			solution := make([]*State, currentState.moves+1)
			for i := currentState.moves; i >= 0; i-- {
				solution[i] = currentState
				currentState = currentState.parent
			}
			return solution
		}

		emptyRow, emptyCol := findEmptySpace(currentState.board)                        // 找到空格的位置
		for _, move := range []struct{ dr, dc int }{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} { // 对每个可能的移动方向
			newRow, newCol := emptyRow+move.dr, emptyCol+move.dc
			if newRow >= 0 && newRow < len(currentState.board) &&
				newCol >= 0 && newCol < len(currentState.board[0]) { // 如果移动后的位置合法
				newBoard := make([][]int, len(currentState.board)) // 创建新的棋盘状态
				for i := range newBoard {
					newBoard[i] = make([]int, len(currentState.board[i]))
					copy(newBoard[i], currentState.board[i])
				}
				newBoard[emptyRow][emptyCol], newBoard[newRow][newCol] = newBoard[newRow][newCol], newBoard[emptyRow][emptyCol] // 移动空格
				newState := &State{newBoard, currentState.moves + 1, currentState, heuristic(newBoard)}                         // 创建新的状态并计算评估分数
				if !visited[serialize(newBoard)] {                                                                              // 如果新的状态没有被访问过
					heap.Push(&priorityQueue, newState) // 将新的状态加入队列
					visited[serialize(newBoard)] = true // 标记新的状态为已访问
				}
			}
		}
	}

	return nil // 如果搜索完所有状态都没有找到解，返回nil
}

// 将棋盘状态转化为字符串，用于记录已访问的状态
func serialize(board [][]int) string {
	result := ""
	for _, row := range board {
		for _, cell := range row {
			result += fmt.Sprintf("%d,", cell)
		}
	}
	return result
}

// 判断当前状态是否是目标状态
func isGoalState(board [][]int) bool {
	n := len(board)
	expected := 1
	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			if row == n-1 && col == n-1 {
				expected = 0
			}
			if board[row][col] != expected {
				return false
			}
			expected++
		}
	}
	return true
}

// 找到空格的位置
func findEmptySpace(board [][]int) (int, int) {
	for row := 0; row < len(board); row++ {
		for col := 0; col < len(board[row]); col++ {
			if board[row][col] == 0 {
				return row, col
			}
		}
	}
	return -1, -1
}

// 打印棋盘
func printBoard(board [][]int) {
	for _, row := range board {
		fmt.Println(row)
	}
}

// 判断棋盘是否有解
func isSolvable(board [][]int) bool {
	// 将棋盘展平为一维数组
	var flatBoard []int
	for _, row := range board {
		flatBoard = append(flatBoard, row...)
	}

	// 计算逆序数
	var inversions int
	for i := 0; i < len(flatBoard)-1; i++ {
		for j := i + 1; j < len(flatBoard); j++ {
			if flatBoard[i] != 0 && flatBoard[j] != 0 && flatBoard[i] > flatBoard[j] {
				inversions++
			}
		}
	}

	// 如果棋盘大小是奇数，逆序数必须是偶数
	// 如果棋盘大小是偶数，考虑从底部计算的空格的奇偶性
	n := len(board)
	if n%2 == 1 {
		return inversions%2 == 0
	} else {
		emptyRow, _ := findEmptySpace(board)
		return (n-emptyRow)%2 != inversions%2
	}
}

// 计算当前状态到目标状态的估计代价（曼哈顿距离）
func heuristic(board [][]int) int {
	n := len(board)
	distance := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			val := board[i][j]
			if val != 0 {
				row := (val - 1) / n
				col := (val - 1) % n
				distance += abs(row-i) + abs(col-j)
			}
		}
	}
	return distance
}

// 返回整数的绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
