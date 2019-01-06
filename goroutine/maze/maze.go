package main

import (
	"fmt"
	"os"
)

func readMaze(fileName string)[][]int{
	file, err := os.Open(fileName)
	if err != nil{
		panic(err)
	}

	var row, col int
	fmt.Fscanf(file,"%d %d", &row, &col)

	maze := make([][]int,row)
	for i := range maze{
		maze[i] = make([]int, col)
		for j := range maze[i]{
			fmt.Fscanf(file,"%d",&maze[i][j])
		}
	}
	return maze
}

type point struct {
	i, j int
}

//directions, up -> left -> down -> right
var dirs = [4]point{
	{-1,0},{0,-1},{1,0},{0,1}}


func (p point) add(r point) point{
	return point{p.i +r.i,p.j+r.j}
}

func (p point) at(gird [][]int) (int, bool){
	if p.i < 0 || p.i >= len(gird){
		return 0,false
 	}

	if p.j < 0 || p.j >= len(gird[p.i]){
		return 0,false
	}
	return gird[p.i][p.j],true
}

func walk(maze [][]int,start,end point) [][]int{
	steps := make([][]int,len(maze))
	for i:= range steps{
		steps[i] = make([]int,len(maze[i]))
	}

	Q := []point{start}
	for len(Q) > 0{
		cur := Q[0]
		Q = Q[1:]

		if cur == end{
			break
		}
		for _, dir := range dirs{
			next := cur.add(dir)
			// maze at next is zero
			val, ok := next.at(maze)
			if !ok || val == 1 {
				continue
			}
			// and step at next is zero
			val, ok = next.at(steps)
			if !ok || val !=0 {
				continue
			}
			// and next != start
			if next == start {
				continue
			}
			curSteps,_ := cur.at(steps)
			steps[next.i][next.j] = curSteps + 1
			Q = append(Q,next)
		}
	}
	return steps
}



func main(){
	maze := readMaze( "goroutine/maze/maze.dat")
	steps := walk(maze,point{0,0},point{len(maze) - 1,len(maze[0]) - 1})

	for _,row := range steps{
		for _, val := range row {
			fmt.Printf("%3d",val)
		}
		fmt.Println()
	}
}
