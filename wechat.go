package wechat

import (

	"github.com/zccdy/wxapi/cache"
	"github.com/zccdy/wxapi/miniprogram"
	miniConfig "github.com/zccdy/wxapi/miniprogram/config"
	"github.com/zccdy/wxapi/officialaccount"
	offConfig "github.com/zccdy/wxapi/officialaccount/config"
	"github.com/zccdy/wxapi/openplatform"
	openConfig "github.com/zccdy/wxapi/openplatform/config"
	"github.com/zccdy/wxapi/pay"
	payConfig "github.com/zccdy/wxapi/pay/config"
)



// Wechat struct
type Wechat struct {
	cache cache.Cache
}

// NewWechat init
func NewWechat() *Wechat {
	return &Wechat{}
}

//SetCache 设置cache
func (wc *Wechat) SetCache(cahce cache.Cache) {
	wc.cache = cahce
}

//GetOfficialAccount 获取微信公众号实例
func (wc *Wechat) GetOfficialAccount(cfg *offConfig.Config) *officialaccount.OfficialAccount {
	if cfg.Cache == nil {
		cfg.Cache = wc.cache
	}
	return officialaccount.NewOfficialAccount(cfg)
}

// GetMiniProgram 获取小程序的实例
func (wc *Wechat) GetMiniProgram(cfg *miniConfig.Config) *miniprogram.MiniProgram {
	if cfg.Cache == nil {
		cfg.Cache = wc.cache
	}
	return miniprogram.NewMiniProgram(cfg)
}

// GetPay 获取微信支付的实例
func (wc *Wechat) GetPay(cfg *payConfig.Config) *pay.Pay {
	return pay.NewPay(cfg)
}

// GetOpenPlatform 获取微信开放平台的实例
func (wc *Wechat) GetOpenPlatform(cfg *openConfig.Config) *openplatform.OpenPlatform {
	return openplatform.NewOpenPlatform(cfg)
}
