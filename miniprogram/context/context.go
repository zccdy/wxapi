package context

import (
	"github.com/zccdy/wxapi/credential"
	"github.com/zccdy/wxapi/miniprogram/config"
)

// Context struct
type Context struct {
	*config.Config
	credential.AccessTokenHandle
}
