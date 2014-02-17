// error
package mp

import (
	"errors"
	"strconv"
)

const (
	SystemBusy = -1   // 系统繁忙
	Success    = iota // 请求成功
)

const (
	_                           = 40000 + iota
	AppSecret                           // 1 获取access_token时AppSecret错误，或者access_token无效
	GrantTypeInvalid                    // 2 不合法的凭证类型
	OpenIDInvalid                       // 3 不合法的OpenID
	MediaTypeInvalid                    // 4 不合法的媒体文件类型
	FileTypeInvalid                     // 5 不合法的文件类型
	FileSizeInvalid                     // 6 不合法的文件大小
	MediaIdInvalid                      // 7 不合法的媒体文件id
	MsgTypeInvalid                      // 8 不合法的消息类型
	ImageSizeInvalid                    // 9 不合法的图片文件大小
	AudioSizeInvalid                    // 10 不合法的语音文件大小
	VideoSizeInvalid                    // 11 不合法的视频文件大小
	ThumbnailSizeInvalid                // 1	2 不合法的缩略图文件大小
	AppIdInvalid                        // 13 不合法的APPID
	AccessTokenInvalid                  // 14 不合法的access_token
	MenuTypeInvalid                     // 15 不合法的菜单类型
	ButtonNumInvalid                    // 16 不合法的按钮个数
	ButtonTypeInvalid                   // 17 不合法的按钮个数
	ButtonNameLenInvalid                // 18 不合法的按钮名字长度
	ButtonKeyLenInvalid                 // 19 不合法的按钮KEY长度
	ButtonUrlLenInvalid                 // 20 不合法的按钮URL长度
	MenuVerInvalid                      // 21 不合法的菜单版本号
	SubMenuDegreeInvalid                // 22 不合法的子菜单级数
	SubMenuButtonNumInvalid             // 23 不合法的子菜单按钮个数
	SubMenuButtonTypeInvalid            // 24 不合法的子菜单按钮类型
	SubMenuButtonNameLenInvalid         // 25 不合法的子菜单按钮名字长度
	SubMenuButtonKeyLenInvalid          // 26 不合法的子菜单按钮KEY长度
	SubMenuButtonUrlLenInvalid          // 27 不合法的子菜单按钮URL长度
	CustomMenuUserInvalid               // 28 不合法的自定义菜单使用用户
	OAuthCodeInvalid                    // 29 不合法的oauth_code
	RefreshTokenInvalid                 // 30 不合法的refresh_token
	OpenIdListInvalid                   // 31 不合法的openid列表
	OpenIdListLenInvalid                // 32 不合法的openid列表长度
	CharacterInvalid                    // 33 不合法的请求字符，不能包含\uxxxx格式的字符
	_                                   // 34
	ParamInvalid                        // 35 不合法的参数
	_                                   // 36
	_                                   // 37
	FormatInvalid                       // 38 不合法的请求格式
	UrlLenInvalid                       // 39 不合法的URL长度
	GroupIdInvalid              = 40050 // 50 不合法的分组id
	GroupNameInvalid            = 40051 // 51 分组名字不合法
)

const (
	_                   = 41000 + iota
	AccessTokenMissing  // 1 缺少access_token参数
	AppIdMissing        // 2 缺少appid参数
	RefreshTokenMissing // 3 缺少refresh_token参数
	SecretMissing       // 4 缺少secret参数
	MediaMissing        // 5 缺少多媒体文件数据
	MediaIdMissing      // 6 缺少media_id参数
	SubMenuMissing      // 7 缺少子菜单数据
	OAuthCodeMissing    // 8 缺少oauth code
	OpenIdMissing       // 9 缺少openid
)

const (
	_                   = 42000 + iota
	AccessTokenTimeout  // 1 access_token超时
	RefreshTokenTimeout // 2 refresh_token超时
	OAuthCodeTimeout    // 3 oauth_code超时
)

const (
	_              = 43000 + iota
	GetNeeded      // 1 需要GET请求
	PostNeeded     // 2 需要POST请求
	HttpsNeeded    // 3 需要HTTPS请求
	ReceiverFollow // 4 需要接收者关注
	FriendNeeded   // 5 需要好友关系
)

const (
	_             = 44000 + iota
	MediaEmpty    // 1 多媒体文件为空
	PostEmpty     // 2 POST的数据包为空
	ImageMsgEmpty // 3 图文消息内容为空
	TextMsgEmpty  // 4 文本消息内容为空
)

const (
	_                   = 45000 + iota
	MediaSizeExceeded   // 1 多媒体文件大小超过限制
	MsgExceeded         // 2 消息内容超过限制
	TitleExceeded       // 3 标题字段超过限制
	DescriptionExceeded // 4 描述字段超过限制
	UrlExceeded         // 5 链接字段超过限制
	ImageUrlExceeded    // 6 图片链接字段超过限制
	AudioTimeExceeded   // 7 语音播放时间超过限制
	ImageMsgExceeded    // 8 图文消息超过限制
	ApiCallingExceeded  // 9 接口调用超过限制
	MenuNumExceeded     // 10 创建菜单个数超过限制
	_
	_
	_
	_
	ReplyTimeExceeded     // 15 回复时间超过限制
	SysGroupNotPermitted  // 16 系统分组，不允许修改
	GroupNameTooLong      // 17 分组名字过长
	GroupNumLimitExceeded // 18 分组数量超过上限
)

const (
	_               = 46000 + iota
	MediaNotExist   // 1 不存在媒体数据
	MenuVerNotExist // 2 不存在的菜单版本
	MenuDataExist   // 3 不存在的菜单数据
	UserNotExist    // 4 不存在的用户
)

const (
	JsonXmlParser       = 47001 // 解析JSON/XML内容错误
	ApiUnauthorized     = 48001 // api功能未授权
	ApiUserUnauthorized = 50001 // 用户未授权该api
)

type Error struct {
	Code int    `json:"errcode,omitempty"`
	Msg  string `json:"errmsg,omitempty"`
}

func (err *Error) String() string {
	return strconv.Itoa(err.Code) + ": " + err.Msg
}

func checkCode(err Error) error {
	if err.Code != Success {
		return errors.New(err.String())
	}
	return nil
}
