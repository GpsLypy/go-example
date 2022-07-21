package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	_ "log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	keepaliveClient *http.Client
)

// type Message struct {
// 	MsgType string `json:"msgtype"`
// 	Text    struct {
// 		Content string `json:"content"`

// 	} `json:"text"`
// }

// type Message struct {
// 	MsgType string `json:"msgtype"`
// 	Text    struct {
// 		Content               string   `json:"content"`
// 		Mentioned_mobile_list []string `json:"mentioned_mobile_list"`
// 	} `json:"text"`
// }

type Message struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

func init() {
	keepaliveClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,
			MaxConnsPerHost: 5,
			// DialContext: (&net.Dialer{
			// 	Timeout:   30 * time.Second,
			// 	KeepAlive: 30 * time.Second,
			// }).DialContext,
		},
		Timeout: time.Second * 2,
	}
}

func main() {
	SendMessage("tasks lost")
}

func doRequestKeepalive(method, reqUrl string, args map[string]string, headers map[string]string, body []byte) (int, []byte, error) {
	var requestBody io.Reader
	if nil != body {
		requestBody = bytes.NewBuffer(body)
	}
	// get need rewrite url
	if nil != args &&
		0 != len(args) {
		u, _ := url.Parse(strings.Trim(reqUrl, "/"))
		q := u.Query()
		if nil != args {
			for arg, val := range args {
				q.Add(arg, val)
			}
		}

		u.RawQuery = q.Encode()
		reqUrl = u.String()
		//log.Debugf("Request url after process: %s", reqUrl)
	}
	// new request
	req, err := http.NewRequest(method, reqUrl, requestBody)
	if nil != err {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// Set custom headers
	if nil != headers {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	rsp, err := keepaliveClient.Do(req)
	if nil != err {
		return 0, nil, err
	}
	defer rsp.Body.Close()

	result, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return 0, nil, err
	}

	return rsp.StatusCode, result, nil
}

func SendMessage(msg string) {

	webhook := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=53ee27a2-9866-4bf1-892e-f1784167efba"
	var m Message
	m.MsgType = "markdown"
	m.Markdown.Content = fmt.Sprintf("DTS丢失任务量<font color=\"warning\">%d</font>,请相关同事注意。\n", 10)

	//m.Text.Content = fmt.Sprintf("%d tasks have been lost!", 5)
	//m.Text.Content = fmt.Sprintf("DTS任务丢失提醒\n丢失任务量%d\n请相关同事注意!", 10)
	//m.Text.Mentioned_mobile_list = append(m.Text.Mentioned_mobile_list, "@all")

	jsons, err := json.Marshal(m)
	if err != nil {
		//log.Errorf("json.Marshal error: %v", err)
		return
	}
	resp := string(jsons)
	statusCode, _, err := doRequestKeepalive(http.MethodPost,
		webhook,
		nil,
		nil,
		[]byte(resp))
	if nil != err {
		// log.Errorf("Send webhook request (%s) failed, error: %v",
		// 	webhook, err)
	}
	if statusCode != http.StatusOK {
		// log.Errorf("Request for webhook (%s), response with status code %d",
		// 	webhook, statusCode)
	}

	// var m Message
	// m.MsgType = "text"
	// m.Text.Content = msg
	// jsons, err := json.Marshal(m)
	// if err != nil {
	// 	fmt.Println("1 erro")
	// 	return
	// }
	// resp := string(jsons)

	// webhook := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=53ee27a2-9866-4bf1-892e-f1784167efba"

	// statusCode, _, err := doRequestKeepalive(http.MethodPost,
	// 	webhook,
	// 	nil,
	// 	nil,
	// 	[]byte(resp))
	// if nil != err {
	// 	// log.Errorf("Send webhook request (%s) failed, error: %v",
	// 	// 	receiver.Webhook, err)
	// }
	// if statusCode != http.StatusOK {
	// 	// log.Errorf("Request for webhook (%s), response with status code %d",
	// 	// 	receiver.Webhook, statusCode)
	// }
}
