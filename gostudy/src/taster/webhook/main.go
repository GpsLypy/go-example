package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Message struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func SendMessage(msg string) {
	var m Message
	m.MsgType = "text"
	m.Text.Content = msg
	jsons, err := json.Marshal(m)
	if err != nil {
		fmt.Println("1 erro")
		return
	}
	resp := string(jsons)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=53ee27a2-9866-4bf1-892e-f1784167efba", strings.NewReader(resp))
	if err != nil {
		fmt.Println("1 erro")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		fmt.Println("1 erro")
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("1 erro")
		return
	}
	fmt.Println("1 yes", string(body))
}

func main() {
	SendMessage("tasks lost")
}
