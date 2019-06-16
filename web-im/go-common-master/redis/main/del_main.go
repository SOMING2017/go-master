package main

import (
	"fmt"

	myRedis ".."
)

var pattern = "login:*"

func main() {
	reply, err := myRedis.Keys(pattern)
	if err != nil {
		fmt.Println("err is ", err)
		return
	}
	for _, item := range reply.([]interface{}) {
		itemString := string(item.([]uint8))
		err := myRedis.DelKey(itemString)
		if err != nil {
			fmt.Println("delete item fail,item is ", itemString)
			fmt.Println("delete,err is ", err)
		}
	}
	reply, err = myRedis.Keys(pattern)
	if err != nil {
		fmt.Println("keys err is ", err)
		return
	}
	replyNum := len(reply.([]interface{}))
	if replyNum == 0 {
		fmt.Println("删除成功...")
	} else {
		fmt.Println("仍有残留...")
	}
}
