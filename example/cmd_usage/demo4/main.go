package main

import (
	"bufio"
	"log"
	"os/exec"
)

func main() {
	//cmdName := "ping -c 4 localhost"
	cmdName := "ansible all -m ping"
	cmd := exec.Command("/bin/sh", "-c", cmdName)
	stdout, _ := cmd.StdoutPipe()

	cmd.Start()

	scan := bufio.NewScanner(stdout)
	for scan.Scan() {
		text := scan.Text()
		//fmt.Println(text) // 没有时间
		log.Printf("%s", text) // 有时间
	}

	cmd.Wait()
}
