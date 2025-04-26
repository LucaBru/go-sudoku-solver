package concurrent

import "fmt"

type Address struct {
	row int
	clm int
}

func (a *Address) rowNeighbors() map[Address]struct{} {
	rowNeighbors := map[Address]struct{}{}
	// leftMostClm := a.clm / 3 * 3
	/* 	for i := range leftMostClm {
	   		address := Address{row: a.row, clm: i}
	   		if network[address] == nil {
	   			continue
	   		}
	   		rowNeighbors[address] = struct{}{}
	   	}
	   	for i := range 9 - (leftMostClm + 3) {
	   		address := Address{row: a.row, clm: leftMostClm + i + 3}
	   		if network[address] == nil {
	   			continue
	   		}
	   		rowNeighbors[address] = struct{}{}
	   	} */
	for i := range 9 {
		address := Address{row: a.row, clm: i}
		if network[address] == nil {
			continue
		}
		rowNeighbors[address] = struct{}{}
	}

	delete(rowNeighbors, *a)
	return rowNeighbors
}

func (a *Address) clmNeighbors() map[Address]struct{} {
	clmNeighbors := map[Address]struct{}{}
	// upperRow := a.row / 3 * 3
	/* 	for i := range upperRow {
	   		address := Address{row: i, clm: a.clm}
	   		if network[address] == nil {
	   			continue
	   		}
	   		clmNeighbors[address] = struct{}{}
	   	}
	   	for i := range 9 - (upperRow + 3) {
	   		address := Address{row: upperRow + i + 3, clm: a.clm}
	   		if network[address] == nil {
	   			continue
	   		}
	   		clmNeighbors[address] = struct{}{}
	   	} */
	for i := range 9 {
		address := Address{row: i, clm: a.clm}
		if network[address] == nil {
			continue
		}
		clmNeighbors[address] = struct{}{}
	}
	delete(clmNeighbors, *a)
	return clmNeighbors
}

func (a *Address) boxNeighbors() map[Address]struct{} {
	boxNeighbors := map[Address]struct{}{}
	upperRow := a.row / 3 * 3
	leftMostClm := a.clm / 3 * 3
	for i := range 3 {
		for j := range 3 {
			address := Address{row: upperRow + i, clm: leftMostClm + j}
			if network[address] == nil {
				continue
			}
			boxNeighbors[address] = struct{}{}
		}
	}
	delete(boxNeighbors, Address{row: a.row, clm: a.clm})
	return boxNeighbors
}

func (a Address) String() string {
	return fmt.Sprintf("[%d,%d]", a.row, a.clm)
}
