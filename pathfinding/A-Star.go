package pathfinding

import "github.com/rilcal/Wildlife-Simulator/structs"

// Astar follows the a* search algorithm to quickly generate the fastest path to a given target point from a given start point in a matrix of given movement matrix
func Astar(start structs.Point, target structs.Point, matrix [][]int) (path []structs.Point) {
	openList := make([]structs.Point, 0, 50)
	openList = append(openList, start)
	closedList := make([]structs.Point, 0, 50)

	for i := 1; i > 0; {
		neighbors := make([]node, 0, 8)
		for i := 0; i < 9; i++ {

		}
	}
	return
}

type node struct {
	pos       structs.Point
	parentPos structs.Point
	neighbors []structs.Point
}
