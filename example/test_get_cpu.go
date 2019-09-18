package main

import (
	"runtime"
	"fmt"
)

func main() {
	TestGetCPU()
}

func TestGetCPU() {
	cpuNums:=runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("逻辑cpu个数：",cpuNums)
}
