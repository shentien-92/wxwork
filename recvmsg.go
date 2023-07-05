package wxwork

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// RecvMessage 消息接收基础结构
type RecvMessage struct {
	ToUsername   string `xml:"ToUserName"`   // 企业微信CorpID
	FromUsername string `xml:"FromUserName"` // 成员UserID
	CreateTime   int64  `xml:"CreateTime"`   // 消息创建时间（整型）
	MsgType      string `xml:"MsgType"`      // 消息类型
	AgentID      int    `xml:"AgentId"`      // 企业应用的id，整型
	MsgID        int64  `xml:"MsgId"`        // 消息id，64位整型

	Content string `xml:"Content"` // 文本消息

	// 多媒体消息
	MediaID      string `xml:"MediaId"`      // 媒体文件id
	PicURL       string `xml:"PicUrl"`       // 图片链接
	Format       string `xml:"Format"`       // 语音格式
	ThumbMediaID string `xml:"ThumbMediaId"` // 视频缩略图的媒体id

	// 位置消息
	LocationX float64 `xml:"Location_X"` // 地理位置纬度
	LocationY float64 `xml:"Location_Y"` // 地理位置经度
	Scale     int     `xml:"Scale"`      // 地图缩放大小
	Label     string  `xml:"Label"`      // 地理位置信息

	// 链接消息
	Title       string `xml:"Title"`       // 标题
	Description string `xml:"Description"` // 链接描述
	Url         string `xml:"Url"`         // 链接跳转的url

	RecvEvent
}

type CallbackData struct {
	AgentId string `json:"agentId"`
	Content string `json:"content"`
	Ldap    string `json:"ldap"`
}

// ParseRecvMessage 解析接收到的消息
func (a *Agent) ParseRecvMessage(signature, timestamp, nonce string, data []byte) (recv RecvMessage, err error) {
	msg, cryptErr := a.crypt.DecryptMsg(signature, timestamp, nonce, data)
	if nil != cryptErr {
		fmt.Println(string(data))
		var callbackData CallbackData
		err = json.Unmarshal(json.RawMessage(data), &callbackData)
		if err != nil {
			return recv, fmt.Errorf("DecryptMsg fail: %v, \n json dump error %v", cryptErr, err)
		}
		fmt.Println(callbackData.Content)
		recv.AgentID, err = strconv.Atoi(callbackData.AgentId)
		recv.Content = callbackData.Content
		recv.FromUsername = callbackData.Ldap
		recv.ToUsername = "diyQABot"
		recv.CreateTime = time.Now().Unix()
		return recv, nil
	}
	err = xml.Unmarshal(msg, &recv)
	return recv, err
}

// CallbackVerify 回调配置验证URL有效性
func (a *Agent) CallbackVerify(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	verifySignature := query.Get("msg_signature")
	verifyTimestamp := query.Get("timestamp")
	verifyNonce := query.Get("nonce")
	verifyEchoStr := query.Get("echostr")

	echoStr, cryptErr := a.crypt.VerifyURL(verifySignature, verifyTimestamp, verifyNonce, verifyEchoStr)

	if nil != cryptErr {
		log.Println("verifyUrl fail", cryptErr)
	}

	w.Write(echoStr)
}
