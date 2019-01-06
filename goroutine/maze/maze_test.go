package main

import "testing"

func TestMaze(t *testing.T){

	maze := readMaze( "maze.dat")
	walk(maze,point{0,0},point{len(maze) - 1,len(maze[0]) - 1})

}
