package main

import "fmt"

//前后元素依次比较
func bubbleSort(a []int){
	for i := 0; i<len(a);i++{
		for j := 1;j < len(a) -i;j++{
			if a[j] < a[j-1]{
				a[j],a[j-1] = a[j-1],a[j]
			}
		}
	}
}

//每次选一个最小或最大的出来，
func selectSort(a []int){
	for i :=0;i<len(a);i++{
		var min int =i
		for j := i+1;j<len(a);j++{
			if a[min] > a[j]{
				min = j
			}
		}
		a[i],a[min] = a[min],a[i]
	}
}

//插入排序，将元素一个个插入到有序序列
func insertSort(a []int){
	for i:=1;i<len(a);i++{
		for j:= i;j>0 ;j--{
			if a[j] > a[j-1]{
				break
			}
			a[j],a[j-1] = a[j-1],a[j]
		}
	}
}

//每次确定一个元素的位置，小的放左边，大的放右边，递归
func quickSort(a []int,left,right int){
	if left >= right{
		return
	}

	val := a[left]
	k := left
	for i := left +1; i<len(a);i++{
		if a[i] < val {
			a[k] = a[i]
			a[i] = a[k+1]
			k++
		}

	}
	a[k] = val
	quickSort(a,left,k-1)
	quickSort(a,k+1,right)

}

func main() {
	a := [...]int{8,7,5,4,3,10,15}
	//bubbleSort(a[:])
	//selectSort(a[:])
	//insertSort(a[:])
	quickSort(a[:],0,len(a)-1)
	fmt.Println(a)
}
