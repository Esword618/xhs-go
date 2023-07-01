package main

import (
	"fmt"
	"github.com/Esword618/xhs-go/consts"
	"github.com/Esword618/xhs-go/xhs"
	"time"
)

func main() {
	const cookieStr = "输入你的cookies"

	XHS := xhs.NewXhs()
	err := XHS.Initialize(cookieStr, false)
	if err != nil {
		fmt.Println("Failed to initialize XhsClient:", err)
		return
	}
	fmt.Println("--------------------------------》GetNoteById")
	newjson, err := XHS.GetNoteById("649989da0000000027028451")
	if err != nil {
		return
	}
	fmt.Println(*newjson)
	fmt.Println("--------------------------------》GetUserInfo")
	newjson, err = XHS.GetUserInfo("5ff0e6410000000001008400")
	if err != nil {
		return
	}
	fmt.Println(*newjson)
	fmt.Println("--------------------------------》GetUserNotes")
	newjson, err = XHS.GetUserNotes("5ff0e6410000000001008400")
	if err != nil {
		return
	}
	fmt.Println(*newjson)
	fmt.Println("--------------------------------》GetNoteByKeyword")
	newjson, err = XHS.GetNoteByKeyword("音乐节")
	if err != nil {
		return
	}
	fmt.Println(*newjson)
	fmt.Println("--------------------------------》GetHomeFeed")
	newjson, err = XHS.GetHomeFeed(consts.Recommend)
	if err != nil {
		return
	}
	fmt.Println(*newjson)

	go func() {
		time.Sleep(10 * time.Second)
		XHS.SendCloseSignal()
		fmt.Println("-------------------》CloseCh")
	}()
	//阻塞程序，等待关闭命令
	XHS.WaitForCloseSignal()
	// 关闭客户端
	XHS.Close()
}
