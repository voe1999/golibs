// 提供企业微信的io.Writer接口
package extensions

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type QyWeiXinWriter interface {
	io.Writer
}

type qyWeiXinWriter struct {
	webhook string
}

func NewQyWeiXinWriter(webhook string) QyWeiXinWriter {
	return &qyWeiXinWriter{webhook: webhook}
}

type MessageOptions struct {
	MsgType QyWeiXinMessageType `json:"msgtype"`
	Text    TextMessageOptions  `json:"text"`
}

type QyWeiXinMessageType string

const (
	Text     QyWeiXinMessageType = "text"
	Markdown QyWeiXinMessageType = "markdown"
	Image    QyWeiXinMessageType = "image"
	News     QyWeiXinMessageType = "news"
	File     QyWeiXinMessageType = "file"
)

type TextMessageOptions struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list"`
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

func (q qyWeiXinWriter) Write(p []byte) (int, error) {
	params := MessageOptions{
		MsgType: "text",
		Text: TextMessageOptions{
			Content:             string(p),
			MentionedList:       []string{"@all"},
			MentionedMobileList: []string{"@all"},
		},
	}
	b, err := json.Marshal(params)
	if err != nil {
		return 0, err
	}
	res, err := http.Post(q.webhook, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return 0, err
	}
	return int(res.ContentLength), nil
}
