// message
package mp

import (
	"encoding/xml"
	"net/http"
	"time"
)

type MsgType string
type EventType string

const (
	MsgTypeText     MsgType = "text"
	MsgTypeImage            = "image"
	MsgTypeVoice            = "voice"
	MsgTypeVideo            = "video"
	MsgTypeMusic            = "music"
	MsgTypeNews             = "news"
	MsgTypeLocation         = "location"
	MsgTypeLink             = "link"
	MsgTypeEvent            = "event"

	MsgTypeEventSubscribe   = "subscribe"
	MsgTypeEventUnsubscribe = "unsubscribe"
	MsgTypeEventScan        = "SCAN"
	MsgTypeEventLocation    = "LOCATION"
	MsgTypeEventClick       = "CLICK"
)

type MsgHeader struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
}

type ServiceMsgHeader struct {
	ToUser  string `json:"touser"`
	MsgType string `json:"msgtype"`
}

type TitleDesc struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Music struct {
	MusicURL     string `json:"musicurl"`
	HQMusicUrl   string `json:"hqmusicurl"`
	ThumbMediaId string `json:"thumb_media_id"`
}

type Article struct {
	TitleDesc
	PicUrl string `json:"picurl"`
	Url    string `json:"url"`
}

type Message struct {
	MsgHeader
	TitleDesc
	Music
	MsgId        uint64
	Content      string
	Url          string
	PicUrl       string
	MediaId      string
	ArticleCount int
	Articles     []Article
	Format       string
	Recognition  string
	LocationX    float64 `xml:"Location_X"`
	LocationY    float64 `xml:"Location_Y"`
	Scale        int
	Label        string
	Event        string
	EventKey     string
	Ticket       string
	Latitude     float64
	Longitude    float64
	Precision    float64
}

type MessageSender interface {
	SendText(touser string, content string) error
	SendImage(touser string, mediaId string) error
	SendVoice(touser string, mediaId string) error
	SendVideo(touser string, mediaId string, info TitleDesc) error
	SendMusic(touser string, info TitleDesc, music Music) error
	SendImageText(touser string, articles []Article) error
}

type messageReply struct {
	fromUserName string
	w            http.ResponseWriter
}

func (r *messageReply) reply(v interface{}) error {
	data, err := xml.Marshal(v)
	if err != nil {
		return err
	}
	r.w.Header().Set("Content-Type", xmlContentType)
	_, err = r.w.Write(data)
	return err
}

func (r *messageReply) SendText(touser string, content string) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		Content string
	}

	data.MsgType = string(MsgTypeText)
	data.ToUserName = touser
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.Content = content

	return r.reply(&data)
}

func (r *messageReply) SendImage(touser string, mediaId string) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		Image struct {
			MediaId string
		}
	}

	data.MsgType = string(MsgTypeImage)
	data.ToUserName = touser
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.Image.MediaId = mediaId

	return r.reply(&data)
}

func (r *messageReply) SendVoice(touser string, mediaId string) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		Voice struct {
			MediaId string
		}
	}

	data.MsgType = string(MsgTypeVoice)
	data.ToUserName = touser
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.Voice.MediaId = mediaId

	return r.reply(&data)
}

func (r *messageReply) SendVideo(touser string, mediaId string, info TitleDesc) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		Video struct {
			MediaId string
			TitleDesc
		}
	}

	data.MsgType = string(MsgTypeVideo)
	data.ToUserName = touser
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.Video.MediaId = mediaId
	data.Video.TitleDesc = info

	return r.reply(&data)
}

func (r *messageReply) SendMusic(touser string, info TitleDesc, music Music) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		M struct {
			TitleDesc
			Music
		} `xml:"Music"`
	}

	data.MsgType = string(MsgTypeMusic)
	data.ToUserName = touser
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.M.TitleDesc = info
	data.M.Music = music

	return r.reply(&data)
}

func (r *messageReply) SendImageText(touser string, articles []Article) error {
	var data struct {
		MsgHeader
		ArticleCount int
		Articles     []Article `xml:"Articles>item"`
	}

	data.MsgType = string(MsgTypeNews)
	data.ToUserName = touser
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.Articles = articles

	return r.reply(&data)
}
