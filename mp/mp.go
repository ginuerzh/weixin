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
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	baseUrl  = "https://api.weixin.qq.com/cgi-bin"
	tokenUri = "/token"

	menuCreateUri = "/menu/create"
	menuQueryUri  = "/menu/get"
	menuDelUri    = "/menu/delete"
)

const (
	xmlContentType  = "application/xml; charset=utf-8"
	jsonContentType = "application/json; charset=utf-8"
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
	menu      *Menu
	routes    map[string]HandlerFunc
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
}

func (mp *MP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	signature := r.FormValue("signature")
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")
	log.Println(signature, timestamp, nonce)

	if !mp.checkSignature(signature, timestamp, nonce) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		fmt.Fprint(w, r.FormValue("echostr"))
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := Message{}
	if err := xml.Unmarshal(data, &msg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if handle, ok := mp.routes[msg.MsgType]; ok {
		reply := &messageReply{toUserName: msg.ToUserName,
			fromUserName: msg.FromUserName, w: w}
		handle(reply, &msg)
	}
}

func (mp *MP) Run() error {
	http.Handle(mp.url, mp)
	return http.ListenAndServe(":8080", nil)
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
		mp.token.expire = 3 * time.Second
		return err
	}
	if err = checkCode(response.Error); err != nil {
		mp.token.expire = 3 * time.Second
		return err
	}

	mp.token.token = response.AccessToken
	mp.token.expire = time.Duration(response.Expire * int64(time.Second))
	log.Println("get token success!", mp.token.token)

	return nil
}

func (mp *MP) SetMenu(menu *Menu) (err error) {
	b, err := json.Marshal(&menu.buttons)
	if err != nil {
		return
	}

	var response Error
	url := baseUrl + menuCreateUri + fmt.Sprintf("?access_token=%s", mp.token.token)
	if err = post(url, jsonContentType, bytes.NewBuffer(b), &response); err != nil {
		return
	}
	if err = checkCode(response); err != nil {
		return
	}

	mp.menu = menu
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
