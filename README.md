naive sudoku solver that try all possible solutions

each cell is filled with a valid number, if a cell has no valid number available, then need to change the previous cell value e resume

parallel solution:
  - naive approach: add random number to the puzzle and launch a go routine per subproblem
  - whenever a goroutine is free, a new instance of the problem is computed and the goroutine restart
  - go routines run up to a solution
  - other go routines has to be notified about the digit to don't try for some cell

  - expect an improvement up to 8 times (num of core in my machine)
  - not recursive algorithm should improve a little bit the seq. solver

concurrent solution:
  - launch a number of go-routines with a reduced problem (hint on the solution)
  - each go-routine could have a hint on a different cell
  - if the go routine terminate successfully, other go routines should terminate
  - upon termination, no solution was found it notifies other routines about the number on the cell not to use

  - each go routine is both an observer and an observable


other possible concurrent solution:
  - each empty cell is a goroutine which receives limitation for such a cell. Keep iterating until just one number is available and after having notified the others return
  - each go routine should use 22 channels for:
    - sending its value, whenever found
    - receiving neighborhood restrictions

  Channels owners:
    - the goroutine for channels used to notify neighbors (closed when its terminates)






GO RUNTIME doesn't clean go routine, just upon termination. So context switching still be executed and it wastes resources
Go routine should be terminated 

CONVENTION: If a goroutine is responsible for creating a goroutine, it is also responsible for ensuring it can stop the goroutine
