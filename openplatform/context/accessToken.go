package context

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zccdy/wxapi/util"
)

const (
	componentAccessTokenURL = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
	getPreCodeURL           = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=%s"
	queryAuthURL            = "https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=%s"
	refreshTokenURL         = "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=%s"
	getComponentInfoURL     = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?component_access_token=%s"
	//TODO 获取授权方选项信息
	getComponentConfigURL = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_option?component_access_token=%s"
	//TODO 获取已授权的账号信息
	getuthorizerListURL = "POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_list?component_access_token=%s"
)

// ComponentAccessToken 第三方平台
type ComponentAccessToken struct {
	AccessToken string `json:"component_access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	StoreAt     int64  `json:"store_at"`
}

// GetComponentAccessToken 获取 ComponentAccessToken
func (ctx *Context) GetComponentAccessToken() (string, error) {
	s,_:=ctx.GetComponentAccessTokenStruct()

	return s.AccessToken, nil
}


// GetComponentAccessTokenStruct 获取 struct
func (ctx *Context) GetComponentAccessTokenStruct() (ComponentAccessToken, error) {
	accessTokenCacheKey := fmt.Sprintf("component_access_token_%s", ctx.AppID)
	cat:=ComponentAccessToken{}
	if err:= ctx.Cache.GetAndUnmarshal(accessTokenCacheKey,&cat);err!=nil{
		return cat, err
	}

	return cat, nil
}

// SetComponentAccessToken 通过component_verify_ticket 获取 ComponentAccessToken
func (ctx *Context) SetComponentAccessToken(verifyTicket string) (*ComponentAccessToken, error) {
	body := map[string]string{
		"component_appid":         ctx.AppID,
		"component_appsecret":     ctx.AppSecret,
		"component_verify_ticket": verifyTicket,
	}
	respBody, err := util.PostJSON(componentAccessTokenURL, body)
	if err != nil {
		return nil, err
	}

	at := &ComponentAccessToken{}
	//fmt.Println("------->respBody=",string(respBody))
	if err := json.Unmarshal(respBody, at); err != nil {
		return nil, err
	}
	at.StoreAt=time.Now().Unix()
	accessTokenCacheKey := fmt.Sprintf("component_access_token_%s", ctx.AppID)
	expires := at.ExpiresIn
	if err := ctx.Cache.Set(accessTokenCacheKey, at, time.Duration(expires)*time.Second); err != nil {
		return nil, nil
	}
	return at, nil
}

// GetPreCode 获取预授权码
func (ctx *Context) GetPreCode() (string, error) {
	cat, err := ctx.GetComponentAccessToken()
	if err != nil {
		return "", err
	}
	req := map[string]string{
		"component_appid": ctx.AppID,
	}
	uri := fmt.Sprintf(getPreCodeURL, cat)
	body, err := util.PostJSON(uri, req)
	if err != nil {
		return "", err
	}

	var ret struct {
		PreCode string `json:"pre_auth_code"`
	}
	//fmt.Println("------->GetPreCode respBody=",string(body))
	if err := json.Unmarshal(body, &ret); err != nil {
		return "", err
	}

	return ret.PreCode, nil
}

// ID 微信返回接口中各种类型字段
type ID struct {
	ID int `json:"id"`
}

// AuthBaseInfo 授权的基本信息
type AuthBaseInfo struct {
	AuthrAccessToken
	FuncInfo []AuthFuncInfo `json:"func_info"`
}

// AuthFuncInfo 授权的接口内容
type AuthFuncInfo struct {
	FuncscopeCategory ID `json:"funcscope_category"`
}

// AuthrAccessToken 授权方AccessToken
type AuthrAccessToken struct {
	Appid        string `json:"authorizer_appid"`
	AccessToken  string `json:"authorizer_access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"authorizer_refresh_token"`
}

// QueryAuthCode 使用授权码换取公众号或小程序的接口调用凭据和授权信息
func (ctx *Context) QueryAuthCode(authCode string) (*AuthBaseInfo, error) {
	cat, err := ctx.GetComponentAccessToken()
	if err != nil {
		return nil, err
	}

	req := map[string]string{
		"component_appid":    ctx.AppID,
		"authorization_code": authCode,
	}
	uri := fmt.Sprintf(queryAuthURL, cat)
	body, err := util.PostJSON(uri, req)
	if err != nil {
		return nil, err
	}
	//{"errcode":61009,"errmsg":"code is invalid rid: 5f34ffda-15f5225d-519bc4eb"}
	var ret struct {
		Code  int		   `json:"errcode"`
		ErrMsg string      `json:"errmsg"`
		Info *AuthBaseInfo `json:"authorization_info"`
	}

	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}

	if ret.Info==nil&&ret.ErrMsg!="" {
		return nil,fmt.Errorf("ErrCode=%d ErrMsg=%s",ret.Code,ret.ErrMsg)
	}

	return ret.Info, nil
}

// RefreshAuthrToken 获取（刷新）授权公众号或小程序的接口调用凭据（令牌）
func (ctx *Context) RefreshAuthrToken(appid, refreshToken string) (*AuthrAccessToken, error) {
	cat, err := ctx.GetComponentAccessToken()
	if err != nil {
		return nil, err
	}

	req := map[string]string{
		"component_appid":          ctx.AppID,
		"authorizer_appid":         appid,
		"authorizer_refresh_token": refreshToken,
	}
	uri := fmt.Sprintf(refreshTokenURL, cat)
	body, err := util.PostJSON(uri, req)
	if err != nil {
		return nil, err
	}

	ret := &AuthrAccessToken{}
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}

	authrTokenKey := "authorizer_access_token_" + appid
	if err := ctx.Cache.Set(authrTokenKey, ret.AccessToken, time.Minute*80); err != nil {
		return nil, err
	}
	return ret, nil
}

// GetAuthrAccessToken 获取授权方AccessToken
func (ctx *Context) GetAuthrAccessToken(appid string) (string, error) {
	authrTokenKey := "authorizer_access_token_" + appid
	val := ctx.Cache.Get(authrTokenKey)
	if val == nil {
		return "", fmt.Errorf("cannot get authorizer %s access token", appid)
	}
	return val.(string), nil
}

type MiniNetWork struct {
	RequestDomain       []string `json:"RequestDomain"`
	WsRequestDomain     []string `json:"WsRequestDomain"`
	UploadDomain 		[]string `json:"UploadDomain"`
	DownloadDomain     	[]string `json:"DownloadDomain"`
}

type MiniCategory struct {
	First        	string `json:"first"`
	Second         	string `json:"second"`
}
type MiniProgramInfo struct {
	NetWork        	MiniNetWork `json:"network"`
	Categories     	[]MiniCategory `json:"categories"`
	VisitStatus 	int     	`json:"visit_status"`
}

// AuthorizerInfo 授权方详细信息
type AuthorizerInfo struct {
	NickName        string `json:"nick_name"`
	HeadImg         string `json:"head_img"`
	ServiceTypeInfo ID     `json:"service_type_info"`
	VerifyTypeInfo  ID     `json:"verify_type_info"`
	UserName        string `json:"user_name"`
	PrincipalName   string `json:"principal_name"`
	BusinessInfo    struct {
		OpenStore int `json:"open_store"`
		OpenScan  int `json:"open_scan"`
		OpenPay   int `json:"open_pay"`
		OpenCard  int `json:"open_card"`
		OpenShake int `json:"open_shake"`
	}  `json:"business_info"`
	Alias     string `json:"alias"`
	QrcodeURL string `json:"qrcode_url"`
	MiniProgramInfo		*MiniProgramInfo `json:"MiniProgramInfo"`
}

// GetAuthrInfo 获取授权方的帐号基本信息
func (ctx *Context) GetAuthrInfo(appid string) (*AuthorizerInfo, *AuthBaseInfo, error) {
	cat, err := ctx.GetComponentAccessToken()
	if err != nil {
		return nil, nil, err
	}

	req := map[string]string{
		"component_appid":  ctx.AppID,
		"authorizer_appid": appid,
	}

	uri := fmt.Sprintf(getComponentInfoURL, cat)
	body, err := util.PostJSON(uri, req)
	if err != nil {
		return nil, nil, err
	}

	var ret struct {
		AuthorizerInfo    *AuthorizerInfo `json:"authorizer_info"`
		AuthorizationInfo *AuthBaseInfo   `json:"authorization_info"`
	}
	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, nil, err
	}

	return ret.AuthorizerInfo, ret.AuthorizationInfo, nil
}
