package context
import (
    "encoding/json"
    "fmt"
    "github.com/zccdy/wxapi/util"
)


const (
    commitCodeURL = "https://api.weixin.qq.com/wxa/commit?access_token=%s"
    submitReviewURL = "https://api.weixin.qq.com/wxa/submit_audit?access_token=%s"
    getLastAuditStatusURL = "https://api.weixin.qq.com/wxa/get_latest_auditstatus?access_token=%s"
    miniReleaseURL = "https://api.weixin.qq.com/wxa/release?access_token=%s"
)


/*
{
    "extEnable": true,
    "extAppid": "wxf9c4501a76931b33",
    "directCommit": false,
    "ext": {
        "name": "wechat",
        "attr": {
                "host": "open.weixin.qq.com",
                "users": [
                    "user_1",
                    "user_2"
                ]
        }
    },
    "extPages": {
        "pages/logs/logs": {
            "navigationBarTitleText": "logs"
        }
    },
    "window":{
        "backgroundTextStyle":"light",
        "navigationBarBackgroundColor": "#fff",
        "navigationBarTitleText": "Demo",
        "navigationBarTextStyle":"black"
    },
    "tabBar": {
        "list": [{
            "pagePath": "pages/index/index",
            "text": "首页"
            }, {
            "pagePath": "pages/logs/logs",
            "text": "日志"
            }]
        },
        "networkTimeout": {
            "request": 10000,
            "downloadFile": 10000
        }
}
*/
type CommitMiniProgramCodeParam struct {
    AccessToken string  `json:"access_token"`    //小程序接口调用令牌
    TemplateId  string  `json:"template_id"`     //代码库中的代码模板 ID
    ExtJson     string  `json:"ext_json"`        //第三方自定义的配置
    Version     string  `json:"user_version"`    //代码版本号，开发者可自定义（长度不要超过 64 个字符）
    Desc        string  `json:"user_desc"`       //代码描述，开发者可自定义
}

// 上传代码 https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Mini_Programs/code/commit.html
func (ctx *Context) CommitMiniProgramCode(param *CommitMiniProgramCodeParam,miniAppId string) error {
    if param.ExtJson=="" {
        param.ExtJson=fmt.Sprintf(`{"extEnable": false, "extAppid": "%s",}`,miniAppId)
    }
    url := fmt.Sprintf(commitCodeURL, param.AccessToken)
    body, err := util.PostJSON(url, param)
    if err != nil {
        return err
    }

    var ret struct {
        Code  int		   `json:"errcode"`
        ErrMsg string      `json:"errmsg"`
    }

    if err := json.Unmarshal(body, &ret); err != nil {
        return err
    }

    if ret.Code!=0&&ret.ErrMsg!="" {
        return fmt.Errorf("ErrCode=%d ErrMsg=%s",ret.Code,ret.ErrMsg)
    }

    return nil
}





// 提交审核 https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Mini_Programs/code/submit_audit.html
func (ctx *Context) CommitMiniProgram2Review(accessToken string) (string, error) {

    var req struct{
        AccessToken  string      `json:"access_token"`
    }
    req.AccessToken=accessToken
    url := fmt.Sprintf(submitReviewURL,accessToken )
    body, err := util.PostJSON(url, &req)
    if err != nil {
        return "", err
    }

    var ret struct {
        Code  int		   `json:"errcode"`
        ErrMsg string      `json:"errmsg"`
        AuditId string     `json:"auditid"`
    }

    if err := json.Unmarshal(body, &ret); err != nil {
        return "",err
    }

    if ret.Code!=0&&ret.ErrMsg!="" {
        return "", fmt.Errorf("ErrCode=%d ErrMsg=%s",ret.Code,ret.ErrMsg)
    }
    return ret.AuditId,nil
}

type LastAuditStatus struct {
    Code        int		   `json:"errcode"`
    ErrMsg      string     `json:"errmsg"`
    AuditId     string     `json:"auditid"` //最新的审核 ID
    Status      int        `json:"status"`  //0--审核成功,1--审核被拒绝,2--审核中,3--已撤回
    Reason      string     `json:"reason"`  //当审核被拒绝时，返回的拒绝原因
    ScreenShot  string     `json:"ScreenShot"`  //当审核被拒绝时，会返回审核失败的小程序截图示例。用 | 分隔的 media_id 的列表，可通过获取永久素材接口拉取截图内容
}

//获取小程序提交代码审核状态  https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Mini_Programs/code/get_latest_auditstatus.html
func (ctx *Context) getMiniProgramLastAuditStatus(accessToken string) (*LastAuditStatus, error) {
    url := fmt.Sprintf(getLastAuditStatusURL,accessToken )
    body, err := util.HTTPGet(url)
    if err != nil {
        return nil, err
    }

    var ret LastAuditStatus

    if err := json.Unmarshal(body, &ret); err != nil {
        return nil,err
    }

    if ret.Code!=0&&ret.ErrMsg!="" {
        return nil, fmt.Errorf("ErrCode=%d ErrMsg=%s",ret.Code,ret.ErrMsg)
    }
    return &ret,nil

}




//小程序发布 https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Mini_Programs/code/release.html
func (ctx *Context) MiniProgramRelease(accessToken string) error{

    var req struct{
        AccessToken  string      `json:"access_token"`
    }
    req.AccessToken=accessToken
    url := fmt.Sprintf(miniReleaseURL,accessToken )
    body, err := util.PostJSON(url, &req)
    if err != nil {
        return  err
    }

    var ret struct {
        Code  int		   `json:"errcode"`
        ErrMsg string      `json:"errmsg"`
    }

    if err := json.Unmarshal(body, &ret); err != nil {
        return err
    }

    if ret.Code!=0&&ret.ErrMsg!="" {
        return  fmt.Errorf("ErrCode=%d ErrMsg=%s",ret.Code,ret.ErrMsg)
    }
    return nil
}