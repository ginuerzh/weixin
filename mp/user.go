// user
package mp

import ()

const ()

type Group struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count,omitempty"`
}

type User struct {
	Subscribe     int    `json:"subscribe"`
	OpenId        string `json:"openid"`
	Nickname      string `json:"nickname"`
	Sex           int    `json:"sex"`
	Language      string `json:"language"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	HeadImgUrl    string `json:"headimgurl"`
	SubscribeTime int64  `json:"subscribe_time"`
}

type OpenIdList struct {
	OpenIds []string `json:"openid"`
}

type UserFollowersResponse struct {
	Total      int        `json:"total"`
	Count      int        `json:"count"`
	Data       OpenIdList `json:"data"`
	NextOpenId string     `json:"next_openid"`
	Error
}
