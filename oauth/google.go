package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const googleGetUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
const googleLoginURL = "https://accounts.google.com/o/oauth2/v2/auth?scope=https://www.googleapis.com/auth/userinfo.email&include_granted_scopes=true&response_type=token&redirect_uri=%s&client_id=%s"

type googleCli struct {
	clientId     string
	clientSecret string
	redirectURI  string
	httpCli      http.Client
}

type googleUserInfoResp struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func (r googleUserInfoResp) ToOauthUser() *UserInfo {
	return &UserInfo{
		Openid: r.ID,
		Avatar: r.Picture,
		Email:  r.Email,
	}
}

// GetAccessToken access token
func (g *googleCli) GetAccessToken(code string) (string, error) {
	// google登录的code本身就是access token
	return code, nil
}

// GetOauthUserInfo 获取用户信息
func (g *googleCli) GetOauthUserInfo(token string) (*UserInfo, error) {
	param := url.Values{}
	param.Add("access_token", token)
	req, err := http.NewRequest(http.MethodGet, googleGetUserInfoURL+"?"+param.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := g.httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := new(googleUserInfoResp)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, err
	}
	return data.ToOauthUser(), nil
}

// GenAuthLoginURL 生成授权地址
func (g *googleCli) GenAuthLoginURL() string {
	return fmt.Sprintf(googleLoginURL, g.redirectURI, g.clientId)
}

func NewGoogleOauth(cfg Config) Oauth {
	return &googleCli{
		clientId:     cfg.ClientId(),
		clientSecret: cfg.ClientSecret(),
		redirectURI:  cfg.RedirectURI(),
		httpCli:      http.Client{},
	}
}
