package main

import (
	"arknights_bot/cmd"
	"log"
)

func main() {
	cmd.Execute()
}

// 设置日志格式
func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}
