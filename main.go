package main

import (
	"fmt"
	"os"
)

// Проверяем, является ли число допустимым для данной позиции
func isValid(board [][]int, row, col int, num int) bool {
	// Проверка строки
	for x := 0; x < 9; x++ {
		if board[row][x] == num {
			return false
		}
	}

	// Проверка столбца
	for x := 0; x < 9; x++ {
		if board[x][col] == num {
			return false
		}
	}

	// Проверка 3x3 квадрата
	startRow := row - row%3
	startCol := col - col%3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i+startRow][j+startCol] == num {
				return false
			}
		}
	}

	return true
}

// Решаем судоку с помощью бэктрекинга
func solveSudoku(board [][]int) bool {
	row, col := -1, -1
	isEmpty := false

	// Находим пустую клетку
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] == 0 {
				row = i
				col = j
				isEmpty = true
				break
			}
		}
		if isEmpty {
			break
		}
	}

	// Если пустых клеток нет, судоку решено
	if !isEmpty {
		return true
	}

	// Пробуем числа от 1 до 9
	for num := 1; num <= 9; num++ {
		if isValid(board, row, col, num) {
			board[row][col] = num
			if solveSudoku(board) {
				return true
			}
			board[row][col] = 0
		}
	}
	return false
}

// Проверяем валидность входных данных
func validateInput(args []string) bool {
	if len(args) != 9 {
		return false
	}

	for _, row := range args {
		if len(row) != 9 {
			return false
		}
		for _, ch := range row {
			if ch != '.' && (ch < '1' || ch > '9') {
				return false
			}
		}
	}
	return true
}

// Преобразуем входные данные в матрицу
func parseBoard(args []string) [][]int {
	board := make([][]int, 9)
	for i := range board {
		board[i] = make([]int, 9)
		for j, ch := range args[i] {
			if ch == '.' {
				board[i][j] = 0
			} else {
				board[i][j] = int(ch - '0')
			}
		}
	}
	return board
}

// Проверяем, есть ли конфликты в начальном состоянии судоку
func hasInitialConflicts(board [][]int) bool {
	// Проверяем каждую заполненную ячейку
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] != 0 {
				temp := board[i][j]
				board[i][j] = 0
				if !isValid(board, i, j, temp) {
					return true
				}
				board[i][j] = temp
			}
		}
	}
	return false
}

// Подсчитываем количество возможных решений
func countSolutions(board [][]int, solutions *int) {
	row, col := -1, -1
	isEmpty := false

	// Находим пустую клетку
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] == 0 {
				row = i
				col = j
				isEmpty = true
				break
			}
		}
		if isEmpty {
			break
		}
	}

	// Если пустых клеток нет, значит найдено решение
	if !isEmpty {
		*solutions++
		return
	}

	// Пробуем числа от 1 до 9
	for num := 1; num <= 9; num++ {
		if isValid(board, row, col, num) {
			board[row][col] = num
			countSolutions(board, solutions)
			board[row][col] = 0
		}
		// Если найдено больше одного решения, выходим
		if *solutions > 1 {
			return
		}
	}
}

// Проверяем, есть ли больше одного решения
func hasMultipleSolutions(board [][]int) bool {
	solutions := 0
	countSolutions(board, &solutions)
	return solutions > 1
}

func main() {
	// Проверяем аргументы командной строки
	if !validateInput(os.Args[1:]) {
		fmt.Println("Error")
		return
	}

	// Преобразуем входные данные в матрицу
	board := parseBoard(os.Args[1:])

	// Проверяем начальное состояние на конфликты
	if hasInitialConflicts(board) {
		fmt.Println("Error")
		return
	}

	// Проверяем наличие более одного решения
	if hasMultipleSolutions(board) {
		fmt.Println("Error: multiple solutions")
		return
	}

	// Решаем судоку
	if !solveSudoku(board) {
		fmt.Println("Error: no solution found")
		return
	}

	// Выводим решение
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%d", board[i][j])
			if j < 8 {
				fmt.Print(" ")
			}
		}
		fmt.Print("")
		fmt.Println()
	}
	fmt.Println("")
}
