// mp
package mp

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	baseUrl              = "https://api.weixin.qq.com/cgi-bin"
	tokenUri             = "/token"
	customSendUri        = "/message/custom/send"
	menuCreateUri        = "/menu/create"
	menuQueryUri         = "/menu/get"
	menuDelUri           = "/menu/delete"
	groupCreateUri       = "/groups/create"
	groupQueryUri        = "/groups/get"
	GroupIdUri           = "/groups/getid"
	groupUpdateUri       = "/groups/update"
	groupMemberUpdateUri = "/groups/members/update"
	userInfoUri          = "/user/info"
	followersUri         = "/user/get"
	qrCodeCreateUri      = "/qrcode/create"
	mediaUploadUri       = "/media/upload"
	mediaDownloadUri     = "/media/get"
)

const (
	xmlContentType  = "application/xml; charset=utf-8"
	jsonContentType = "application/json; charset=utf-8"
)

type MediaType string

const (
	MediaImage MediaType = "image"
	MediaVoice           = "voice"
	MediaVideo           = "video"
	MediaThumb           = "thumb"
)

type LangType string

const (
	LangCN LangType = "Zh_CN"
	LangTW          = "zh_TW"
	LangEN          = "en"
)

type HandlerFunc func(reply MessageReplyer, m *Message)

type accessToken struct {
	token  string
	expire time.Duration
}

type MP struct {
	appId     string
	appSecret string
	url       string
	appToken  string
	token     accessToken
	routes    map[string]HandlerFunc
	menu      *Menu
	groups    []Group
}

func New(appId, appSecret, appToken string) *MP {
	mp := &MP{appId: appId, appSecret: appSecret, appToken: appToken}
	if err := mp.requestToken(); err != nil {
		log.Println(err)
	}
	go mp.refreshToken()

	return mp
}

func (mp *MP) Init(url string) {
	mp.url = url
}

func (mp *MP) HandleFunc(msgType MsgType, handler HandlerFunc) {
	if mp.routes == nil {
		mp.routes = make(map[string]HandlerFunc)
	}
	mp.routes[string(msgType)] = handler
	log.Println(mp.routes)
}

func (mp *MP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	signature := r.FormValue("signature")
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")
	log.Println(signature, timestamp, nonce)

	if !mp.checkSignature(signature, timestamp, nonce) {
		log.Println("checkSignature failed!")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		fmt.Fprint(w, r.FormValue("echostr"))
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := Message{}
	if err := xml.Unmarshal(data, &msg); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(msg.MsgType)
	log.Println(msg)

	if handle, ok := mp.routes[msg.MsgType]; ok {
		reply := &messageReply{fromUserName: msg.ToUserName,
			toUserName: msg.FromUserName, w: w}
		handle(reply, &msg)
	}
}

func (mp *MP) Run(port int) error {
	http.Handle(mp.url, mp)
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (mp *MP) checkSignature(signature, timestamp, nonce string) bool {
	list := []string{mp.appToken, timestamp, nonce}
	sort.Strings(list)

	h := sha1.New()
	io.WriteString(h, strings.Join(list, ""))

	return signature == fmt.Sprintf("%x", h.Sum(nil))
}

func (mp *MP) refreshToken() (err error) {
	for {
		select {
		case <-time.Tick(mp.token.expire):
			if err := mp.requestToken(); err != nil {
				log.Println(err)
				break
			}
		}
	}

	return
}

func (mp *MP) requestToken() (err error) {
	var response struct {
		AccessToken string `json:"access_token"`
		Expire      int64  `json:"expires_in"`
		Error
	}

	url := baseUrl + tokenUri +
		fmt.Sprintf("?grant_type=client_credential&appid=%s&secret=%s",
			mp.appId, mp.appSecret)
	if err = get(url, &response); err != nil {
		//log.Println(err)
		mp.token.expire = 3 * time.Second
		return err
	}
	if err = checkCode(response.Error); err != nil {
		//log.Println(err)
		mp.token.expire = 3 * time.Second
		return err
	}

	mp.token.token = response.AccessToken
	mp.token.expire = time.Duration(response.Expire * int64(time.Second))
	log.Println("get token success!", mp.token.token)

	return nil
}

func (mp *MP) SetMenu(menu *Menu) (err error) {
	if err = mp.sendJson(menuCreateUri, &menu.buttons); err == nil {
		mp.menu = menu
	}
	return
}

func (mp *MP) GetMenu() (menu *Menu, err error) {
	if mp.menu != nil {
		return mp.menu, nil
	}

	url := baseUrl + menuQueryUri + fmt.Sprintf("?access_token=%s", mp.token.token)
	if err = get(url, menu); err != nil {
		return
	}
	if err = checkCode(menu.Error); err != nil {
		return
	}

	mp.menu = menu
	return
}

func (mp *MP) DelMenu() (err error) {
	url := baseUrl + menuDelUri + fmt.Sprintf("?access_token=%s", mp.token.token)
	var response Error

	if err = get(url, &response); err != nil {
		return
	}

	if err = checkCode(response); err != nil {
		return
	}

	mp.menu = nil
	return nil
}

func (mp *MP) CreateGroup(name string) error {
	var req, resp struct {
		Grp Group `json:"group"`
		Error
	}

	req.Grp.Name = name
	b, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	url := baseUrl + groupCreateUri + fmt.Sprintf("?access_token=%s", mp.token.token)
	if err := post(url, jsonContentType, bytes.NewBuffer(b), &resp); err != nil {
		return err
	}
	if err := checkCode(resp.Error); err != nil {
		return err
	}

	mp.groups = append(mp.groups, resp.Grp)

	return nil
}

func (mp *MP) Groups() ([]Group, error) {
	var resp struct {
		Groups []Group `json:"groups"`
		Error
	}

	url := baseUrl + groupQueryUri + fmt.Sprintf("?access_token=%s", mp.token.token)
	if err := get(url, &resp); err != nil {
		return nil, err
	}
	if err := checkCode(resp.Error); err != nil {
		return nil, err
	}

	mp.groups = resp.Groups

	return mp.groups, nil
}

func (mp *MP) GroupId(uid string) (gid int, err error) {
	var req struct {
		Uid string `json:"openid"`
	}

	var resp struct {
		GroupId int `json:"groupid"`
		Error
	}

	req.Uid = uid
	b, err := json.Marshal(&req)
	if err != nil {
		return
	}
	url := baseUrl + GroupIdUri + fmt.Sprintf("?access_token=%s", mp.token.token)
	if err = post(url, jsonContentType, bytes.NewBuffer(b), &resp); err != nil {
		return
	}
	if err = checkCode(resp.Error); err != nil {
		return
	}

	return resp.GroupId, nil
}

func (mp *MP) UpdateGroup(group Group) error {
	var req struct {
		Grp Group `json:"group"`
	}

	return mp.sendJson(groupUpdateUri, &req)
}

func (mp *MP) MoveMember2Group(uid string, gid int) error {
	var req struct {
		Uid string `json:"openid"`
		Gid int    `json:"to_groupid"`
	}

	return mp.sendJson(groupMemberUpdateUri, &req)
}

func (mp *MP) UserInfo(uid string, lang LangType) (User, error) {
	var resp struct {
		User
		Error
	}

	url := baseUrl + userInfoUri +
		fmt.Sprintf("?access_token=%s&openid=%s&lang=%s",
			mp.token.token, uid, lang)
	if err := get(url, &resp); err != nil {
		return resp.User, err
	}
	if err := checkCode(resp.Error); err != nil {
		return resp.User, err
	}

	return resp.User, nil
}

func (mp *MP) Followers(start string) (int, []string, string, error) {
	var resp struct {
		Total int `json:"total"`
		Count int `json:"count"`
		Data  struct {
			OpenId []string `json:"openid"`
		} `json:"data"`
		Next string `json:"next_openid"`
		Error
	}

	url := baseUrl + followersUri + fmt.Sprintf("?access_token=%s", mp.token.token)
	if len(start) != 0 {
		url += fmt.Sprintf("&next_openid=%s", start)
	}
	if err := get(url, &resp); err != nil {
		return 0, nil, "", err
	}
	if err := checkCode(resp.Error); err != nil {
		return 0, nil, "", err
	}

	return resp.Total, resp.Data.OpenId, resp.Next, nil
}

// if expire != 0, return temp qrcode
func (mp *MP) QRCode(expire, sceneId int) (string, error) {
	var req struct {
		Expire int    `json:"expire_seconds,omitempty"`
		Action string `json:"action_name"`
		Info   struct {
			Scene struct {
				Id int `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}

	var resp struct {
		Ticket string `json:"ticket"`
		Expire int    `json:"expire_seconds"`
		Error
	}

	req.Expire = expire
	if expire == 0 {
		req.Action = "QR_LIMIT_SCENE"
	} else {
		req.Action = "QR_SCENE"
	}
	req.Info.Scene.Id = sceneId
	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	url := baseUrl + qrCodeCreateUri + fmt.Sprintf("?access_token=%s", mp.token.token)
	if err := post(url, jsonContentType, bytes.NewBuffer(b), &resp); err != nil {
		return "", err
	}
	if err := checkCode(resp.Error); err != nil {
		return "", err
	}

	return resp.Ticket, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createFormFile(writer *multipart.Writer, fieldname, filename, mime string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	if len(mime) == 0 {
		mime = "application/octet-stream"
	}
	h.Set("Content-Type", mime)
	return writer.CreatePart(h)
}

func makeFormData(filename, mimeType string, content io.Reader) (formData io.Reader, contentType string, err error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := createFormFile(writer, "media", filename, mimeType)
	//log.Println(filename, mimeType)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = io.Copy(part, content)
	if err != nil {
		log.Println(err)
		return
	}

	formData = buf
	contentType = writer.FormDataContentType()
	//log.Println(contentType)
	writer.Close()

	return
}

func (mp *MP) UploadMedia(mediaType MediaType, filename string, reader io.Reader) (mediaId string, err error) {
	var resp struct {
		Type      string `json:"type"`
		MediaId   string `json:"media_id"`
		CreatedAt int64  `json:"created_at"`
		Error
	}
	/*
		b := &bytes.Buffer{}
		writer := multipart.NewWriter(b)
		defer writer.Close()

		formFile, err := writer.CreateFormFile("media", filename)
		if err != nil {
			return
		}
		if _, err := io.Copy(formFile, reader); err != nil {
			return "", err
		}
	*/
	data, contentType, err := makeFormData(filename, "image/jpeg", reader)
	if err != nil {
		return "", err
	}
	url := baseUrl + mediaUploadUri +
		fmt.Sprintf("?access_token=%s&type=%s", mp.token.token, mediaType)
	if err := post(url, contentType, data, &resp); err != nil {
		return "", err
	}
	if err := checkCode(resp.Error); err != nil {
		return "", err
	}

	return resp.MediaId, nil
}

func (mp *MP) DownloadMedia(mediaId string) (io.Reader, error) {
	url := baseUrl + mediaDownloadUri +
		fmt.Sprintf("?access_token=%s&media_id=%s", mp.token.token, mediaId)

	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	b := &bytes.Buffer{}
	if _, err := io.Copy(b, r.Body); err != nil {
		return nil, err
	}

	var resp Error

	if err := json.Unmarshal(b.Bytes(), &resp); err == nil {
		return nil, checkCode(resp)
	}

	return b, nil
}

func (mp *MP) SendText(touser string, content string) error {
	var data struct {
		ServiceMsgHeader
		Text struct {
			Content string `json:"content"`
		} `json:"text"`
	}

	data.ToUser = touser
	data.MsgType = string(MsgTypeText)
	data.Text.Content = content

	return mp.sendJson(customSendUri, &data)
}

func (mp *MP) SendImage(touser string, mediaId string) error {
	var data struct {
		ServiceMsgHeader
		Image struct {
			MediaId string `json:"media_id"`
		} `json:"image"`
	}

	data.ToUser = touser
	data.MsgType = string(MsgTypeImage)
	data.Image.MediaId = mediaId

	return mp.sendJson(customSendUri, &data)
}

func (mp *MP) SendVoice(touser string, mediaId string) error {
	var data struct {
		ServiceMsgHeader
		Voice struct {
			MediaId string `json:"media_id"`
		} `json:"voice"`
	}

	data.ToUser = touser
	data.MsgType = string(MsgTypeVoice)
	data.Voice.MediaId = mediaId

	return mp.sendJson(customSendUri, &data)
}

func (mp *MP) SendVideo(touser string, mediaId string, info TitleDesc) error {
	var data struct {
		ServiceMsgHeader
		Video struct {
			MediaId string `json:"media_id"`
			TitleDesc
		} `json:"video"`
	}

	data.ToUser = touser
	data.MsgType = string(MsgTypeVideo)
	data.Video.MediaId = mediaId
	data.Video.TitleDesc = info

	return mp.sendJson(customSendUri, &data)
}

func (mp *MP) SendMusic(touser string, info TitleDesc, music Music) error {
	var data struct {
		ServiceMsgHeader
		M struct {
			TitleDesc
			Music
		} `json:"music"`
	}

	data.ToUser = touser
	data.MsgType = string(MsgTypeMusic)
	data.M.TitleDesc = info
	data.M.Music = music

	return mp.sendJson(customSendUri, &data)
}

func (mp *MP) SendImageText(touser string, articles []Article) error {
	var data struct {
		ServiceMsgHeader
		News struct {
			Articles []Article `json:"articles"`
		} `json:"news"`
	}

	data.ToUser = touser
	data.MsgType = string(MsgTypeNews)
	data.News.Articles = articles

	return mp.sendJson(customSendUri, &data)
}

func (mp *MP) sendJson(uri string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var result Error
	url := uri + fmt.Sprintf("?access_token=%s", mp.token.token)
	if err = post(url, jsonContentType, bytes.NewBuffer(data), &result); err != nil {
		return err
	}
	if err = checkCode(result); err != nil {
		return err
	}

	return nil
}

func get(url string, respStruct interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	return parse(resp, respStruct)
}

func post(url string, bodyType string, body io.Reader, respStruct interface{}) error {
	resp, err := http.Post(url, bodyType, body)
	if err != nil {
		return err
	}
	return parse(resp, respStruct)
}

func parse(resp *http.Response, respStruct interface{}) error {
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, respStruct); err != nil {
		return err
	}

	return nil
}
