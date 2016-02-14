// package main accepts a sudoku puzzle on standard in and outputs a solved
// puzzle on standard out.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var nums = [10]int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}

var v = flag.Bool("v", false, "verbose output")

func main() {
	flag.Parse()

	var board [81]int

	buildBoard(bufio.NewScanner(os.Stdin), &board)

	if *v == true {
		fmt.Println("Input:")
		printBoard(board)
	}

	solve(&board, 0)

	if *v == true {
		fmt.Println("Output:")
	}

	printBoard(board)
}

// buildBoard builds a "board" (an array) by scanning input, splitting comma-
// separated integers and inserting them into an array.
func buildBoard(input *bufio.Scanner, board *[81]int) {
	l := 0

	for input.Scan() {
		for i, n := range strings.Split(input.Text(), ",") {
			var val int

			// If i is a dash, val is 0
			if n == "-" {
				val = 0
			} else {
				// Convert i to an int
				val2, err := strconv.Atoi(n)
				if err != nil {
					fmt.Println(os.Stderr, err)
					os.Exit(2)
				}
				val = val2
			}

			board[i+9*l] = val
		}

		l++
	}
}

// Print the board in an attractive manner.
func printBoard(board [81]int) {
	for i, n := range board {
		fmt.Print(n)

		if i%9 == 8 {
			fmt.Printf("\n")
		} else {
			fmt.Printf(",")
		}
	}
}

// Solves the board in-place by:
// 	* Returns true if the board is complete
//  * If the current space is non-zero and valid,
//    calling itself recursively on the next space.
//  * Otherwise, going through the numbers 9 - 0 (DESC) and seeing if they're
//    valid.
//  * If they are valid, try solving the next space, and return true if it
// 		works.
//  * Otherwise, keep trying numbers.
func solve(board *[81]int, ptr int) bool {
	if ptr >= 81 {
		return true
	}

	if board[ptr] != 0 && checkSpace(*board, ptr) == true {
		return solve(board, ptr+1)
	}

	for _, n := range nums {
		board[ptr] = n

		if board[ptr] == 0 {
			return false
		}

		if checkSpace(*board, ptr) == false {
			continue
		}

		if solve(board, ptr+1) == true {
			return true
		}
	}

	return true
}

// Calculate the x,y coordinate of a position in the board array.
func x_y(ptr int) (int, int) {
	x := ptr % 9
	y := ptr / 9

	return x, y
}

// isSameSector checks whether two pairs of x,y coordinates are in the
// same sector of the game board.
func inSameSector(x, y, ptr_x, ptr_y int) bool {
	return x/3 == ptr_x/3 && y/3 == ptr_y/3
}

// checkSpace checks a space on the board to see whether the current column,
// row and sector are valid. It ignores zeros, and so relies on the algorithm
// to never leave a space zero.
func checkSpace(board [81]int, ptr int) bool {
	ptr_x, ptr_y := x_y(ptr)

	row := make(map[int]bool)
	col := make(map[int]bool)
	sec := make(map[int]bool)

	for i := 0; i < len(board); i++ {
		// If the space is a zero, ignore it.
		if board[i] == 0 {
			continue
		}

		val := board[i]
		x, y := x_y(i)

		if y == ptr_y && axisIsInvalid(val, row) {
			return false
		}

		if x == ptr_x && axisIsInvalid(val, col) {
			return false
		}

		if inSameSector(x, y, ptr_x, ptr_y) && axisIsInvalid(val, sec) {
			return false
		}
	}

	return true
}

// axisIsInvalid checks to see whether an axis (x, y or sector) is invalid.
// It is invalid if it contains more than one number.
func axisIsInvalid(val int, seen map[int]bool) bool {
	if seen[val] == true {
		return true
	}

	seen[val] = true

	return false
}
