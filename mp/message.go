// message
package mp

import (
	"encoding/xml"
	"net/http"
	"time"
)

const (
	customMessageUri = "/message/custom/send"
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

type TitleDesc struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Music struct {
	MusicURL     string
	HQMusicUrl   string
	ThumbMediaId string
}

type Article struct {
	TitleDesc
	PicUrl string
	Url    string
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

type MessageReplyer interface {
	ReplyText(content string) error
	ReplyImage(mediaId string) error
	ReplyVoice(mediaId string) error
	ReplyVideo(mediaId string, info TitleDesc) error
	ReplyMusic(info TitleDesc, music Music) error
	ReplyImageText(articles []Article) error
}

type messageReply struct {
	toUserName   string
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

func (r *messageReply) ReplyText(content string) error {
	var data struct {
		XMLName xml.Name `xml:"xml"`
		MsgHeader
		Content string
	}

	data.MsgType = string(MsgTypeText)
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

	data.MsgType = string(MsgTypeImage)
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

	data.MsgType = string(MsgTypeVoice)
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

	data.MsgType = string(MsgTypeVideo)
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

	data.MsgType = string(MsgTypeMusic)
	data.ToUserName = r.toUserName
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.M.TitleDesc = info
	data.M.Music = music

	return r.reply(&data)
}

func (r *messageReply) ReplyImageText(articles []Article) error {
	var data struct {
		MsgHeader
		ArticleCount int
		Articles     []Article `xml:"Articles>item"`
	}

	data.MsgType = string(MsgTypeNews)
	data.ToUserName = r.toUserName
	data.FromUserName = r.fromUserName
	data.CreateTime = time.Now().Unix()
	data.Articles = articles

	return r.reply(&data)
}

// custom service message response

type ServiceMsgHeader struct {
	ToUser  string `json:"touser"`
	MsgType string `json:"msgtype"`
}

type TextContent struct {
	Content string `json:"content"`
}
type ServiceTextMsgResponse struct {
	ServiceMsgHeader
	Text TextContent `json:"text"`
}

type MediaIdContent struct {
	MediaId string `json:"media_id"`
}
type ServiceImageMsgResponse struct {
	ServiceMsgHeader
	Image MediaIdContent `json:"image"`
}

type ServiceVoiceMsgResponse struct {
	ServiceMsgHeader
	Voice MediaIdContent `json:"voice"`
}

type VideoContent struct {
	MediaIdContent
	TitleDesc
}
type ServiceVideoMsgResponse struct {
	ServiceMsgHeader
	Video VideoContent `json:"video"`
}

type MusicContent struct {
	TitleDesc
	MusicUrl     string `json:"musicurl"`
	HQMusicUrl   string `json:"hqmusicurl"`
	ThumbMediaId string `json:"thumb_media_id"`
}
type ServiceMusicMsgResponse struct {
	ServiceMsgHeader
	Music MusicContent `json:"music"`
}

type NewsContent struct {
	Articles []Article `json:"articles"`
}
type ServiceImageTextMsgResponse struct {
	ServiceMsgHeader
	News NewsContent `json:"news"`
}
