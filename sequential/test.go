package sequential

import "fmt"

var m = [9][9]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 3, 0, 8, 5},
	{0, 0, 1, 0, 2, 0, 0, 0, 0},
	{0, 0, 0, 5, 0, 7, 0, 0, 0},
	{0, 0, 4, 0, 0, 0, 1, 0, 0},
	{0, 9, 0, 0, 0, 0, 0, 0, 0},
	{5, 0, 0, 0, 0, 0, 0, 7, 3},
	{0, 0, 2, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 0, 4, 0, 0, 0, 9},
}

func ok(i, j int) bool {
	val := m[i][j]

	// Check row
	for q := 0; q < 9; q++ {
		if q != j && m[i][q] == val && val != 0 {
			return false
		}
	}

	// Check column
	for q := 0; q < 9; q++ {
		if q != i && m[q][j] == val && val != 0 {
			return false
		}
	}

	// Check 3x3 subgrid
	rowIdx := i / 3 * 3
	clmIdx := j / 3 * 3
	for q := 0; q < 3; q++ {
		for h := 0; h < 3; h++ {
			if (rowIdx+q != i || clmIdx+h != j) && m[rowIdx+q][clmIdx+h] == val && val != 0 {
				return false
			}
		}
	}
	return true
}

func display() {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%d ", m[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func solve(i, j int) bool {
	if j == 9 {
		j = 0
		i++
	}
	if i == 9 && j == 0 {
		display()
		return true
	}

	if m[i][j] > 0 {
		return solve(i, j+1)
	}

	for v := 1; v <= 9; v++ {
		m[i][j] = v
		if ok(i, j) {
			if solve(i, j+1) {
				return true
			}
		}
		m[i][j] = 0
	}
	return false
}

func main() {
	solve(0, 0)
}
