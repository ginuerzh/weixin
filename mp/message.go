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
	MsgText             MsgType = "text"
	MsgImage                    = "image"
	MsgVoice                    = "voice"
	MsgVideo                    = "video"
	MsgMusic                    = "music"
	MsgNews                     = "news"
	MsgLocation                 = "location"
	MsgLink                     = "link"
	MsgEvent                    = "event"
	MsgSubscribeEvent           = "event.subscribe"
	MsgUnsubscribeEvent         = "event.unsubscribe"
	MsgScanEvent                = "event.SCAN"
	MsgLocationEvent            = "event.LOCATION"
	MsgClickEvent               = "event.CLICK"

	EventSubscribe   EventType = "subscribe"
	EventUnsubscribe           = "unsubscribe"
	EventScan                  = "SCAN"
	EventLocation              = "LOCATION"
	EventClick                 = "CLICK"
)

type MsgHeader struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	Type         string
}

type ServiceMsgHeader struct {
	ToUser string `json:"touser"`
	Type   string `json:"msgtype"`
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

type Replyer interface {
	ReplyText(content string) error
	ReplyImage(mediaId string) error
	ReplyVoice(mediaId string) error
	ReplyVideo(mediaId string, info TitleDesc) error
	ReplyMusic(info TitleDesc, music Music) error
	ReplyImageText(articles []Article) error
}

type messageReply struct {
	fromUserName string
	toUserName   string
	w            http.ResponseWriter
	replied      bool
}

func (r *messageReply) reply(v interface{}) error {
	r.replied = true

	data, err := xml.Marshal(v)
	if err != nil {
		return err
	}
	r.w.Header().Set("Content-Type", xmlContentType)
	_, err = r.w.Write(data)
	return err
}

func (r *messageReply) ReplyText(content string) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		Content string
	}

	data.Type = string(MsgText)
	data.ToUserName = r.toUserName
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.Content = content

	return r.reply(&data)
}

func (r *messageReply) ReplyImage(mediaId string) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		Image struct {
			MediaId string
		}
	}

	data.Type = string(MsgImage)
	data.ToUserName = r.toUserName
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.Image.MediaId = mediaId

	return r.reply(&data)
}

func (r *messageReply) ReplyVoice(mediaId string) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		Voice struct {
			MediaId string
		}
	}

	data.Type = string(MsgVoice)
	data.ToUserName = r.toUserName
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.Voice.MediaId = mediaId

	return r.reply(&data)
}

func (r *messageReply) ReplyVideo(mediaId string, info TitleDesc) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		Video struct {
			MediaId string
			TitleDesc
		}
	}

	data.Type = string(MsgVideo)
	data.ToUserName = r.toUserName
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.Video.MediaId = mediaId
	data.Video.TitleDesc = info

	return r.reply(&data)
}

func (r *messageReply) ReplyMusic(info TitleDesc, music Music) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		M struct {
			TitleDesc
			Music
		} `xml:"Music"`
	}

	data.Type = string(MsgMusic)
	data.ToUserName = r.toUserName
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.M.TitleDesc = info
	data.M.Music = music

	return r.reply(&data)
}

func (r *messageReply) ReplyImageText(articles []Article) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		ArticleCount int
		Articles     []Article `xml:"Articles>item"`
	}

	data.Type = string(MsgNews)
	data.ToUserName = r.toUserName
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.ArticleCount = len(articles)
	data.Articles = articles

	return r.reply(&data)
}
