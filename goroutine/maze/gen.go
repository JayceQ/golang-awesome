package main
import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func init(){
	//以时间作为初始化种子
	rand.Seed(time.Now().UnixNano())
}
func ts() {
	fileName := "goroutine/maze/maze.dat"
	//file, _ := os.OpenFile(fileName,os.O_CREATE|os.O_APPEND|os.O_RDWR,0777)
	file, _ := os.Create(fileName)
	defer file.Close()
	row, col := 50, 50
	fmt.Fprintf(file, "%d %d\n", row, col)
	//ch := make(chan int,1)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			//select {
			//case ch <- 0:
			//case ch <- 1:
			//}
			i2 := rand.Float32()
			var num int
			if i2 > 0.8{
				num = 1
			}else {
				num = 0
			}
			if  i == row -1 && j == col -1{
				fmt.Fprintf(file, "%d",0)
			}else if i == 0 && j == 0{
				fmt.Fprintf(file, "%d  ",0)
			}else if j == col -1{
				fmt.Fprintf(file, "%d",num)
			}else{
				fmt.Fprintf(file, "%d  ",num)
			}

		}
		fmt.Fprintf(file, "\n")
	}


}


