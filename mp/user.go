// user
package mp

import ()

const (
	groupCreateUri       = "/groups/create"
	groupQueryUri        = "groups/get"
	UserGroupUri         = "/groups/getid"
	groupUpdateUri       = "/groups/update"
	groupMemberUpdateUri = "/groups/members/update"

	userInfoUri      = "/user/info"
	userFollowersUri = "/user/get"
)

type GroupContent struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type GroupCreateResponse struct {
	Group GroupContent `json:"group"`
	Error
}

type GroupQueryResponse struct {
	Groups []GroupContent `json:"groups"`
	Error
}

type UserGroupResponse struct {
	GroupId int `json:"groupid"`
	Error
}

type UserInfoResponse struct {
	Subscribe     bool   `json:"subscribe"`
	OpenId        string `json:"openid"`
	Nickname      string `json:"nickname"`
	Sex           bool   `json:"sex"`
	Language      string `json:"language"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	HeadImgUrl    string `json:"headimgurl"`
	SubscribeTime uint64 `json:"subscribe_time"`
	Error
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
