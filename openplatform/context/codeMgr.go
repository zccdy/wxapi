package context
import (
    "fmt"
    "github.com/zccdy/wxapi/util"
)


const (
    commitCodeURL = "https://api.weixin.qq.com/wxa/commit?access_token=%s"

)


type CommitMiniProgramCodeParam struct {
    AccessToken string  `json:"access_token"`    //小程序接口调用令牌
    TemplateId  string  `json:"template_id"`     //代码库中的代码模板 ID
    ExtJson     string  `json:"ext_json"`        //第三方自定义的配置
    Version     string  `json:"user_version"`    //代码版本号，开发者可自定义（长度不要超过 64 个字符）
    Desc        string  `json:"user_desc"`       //代码描述，开发者可自定义
}

func (ctx *Context) CommitMiniProgramCode(param *CommitMiniProgramCodeParam) error {

    url := fmt.Sprintf(commitCodeURL, param.AccessToken)
    data, err := util.PostJSON(url, param)
    if err != nil {
        return err
    }
    return util.DecodeWithCommonError(data, "component/fastregisterweapp?action=create")
}