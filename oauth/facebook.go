package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jiebutech/uc/model"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const fbGetAccessTokenURL = "https://graph.facebook.com/v18.0/oauth/access_token"
const fbGetUserInfoURL = "https://graph.facebook.com/me"

type fbCli struct {
	clientId     string
	clientSecret string
	redirectURI  string
	httpCli      http.Client
}
type fbAccessTokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type fbUserInfoResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (r fbUserInfoResp) ToOauthUser() *model.OauthUserInfo {
	return &model.OauthUserInfo{
		Openid:   r.ID,
		Nickname: r.Name,
	}
}

func (cli *fbCli) GetAccessToken(code string) (string, error) {
	param := url.Values{}
	param.Add("client_id", cli.clientId)
	param.Add("client_secret", cli.clientSecret)
	param.Add("redirect_uri", cli.redirectURI)
	param.Add("code", code)
	req, err := http.NewRequest(http.MethodGet, fbGetAccessTokenURL, strings.NewReader(param.Encode()))
	if err != nil {
		return "", err
	}
	resp, err := cli.httpCli.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	data := new(fbAccessTokenResp)
	content, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(content, data)
	if err != nil {
		return "", err
	}
	if data.AccessToken == "" {
		return "", fmt.Errorf("access token is empty")
	}
	return data.AccessToken, nil

}

func (cli *fbCli) GetOauthUserInfo(token string) (*model.OauthUserInfo, error) {
	param := url.Values{}
	param.Add("access_token", token)
	req, err := http.NewRequest(http.MethodGet, fbGetUserInfoURL, strings.NewReader(param.Encode()))
	if err != nil {
		return nil, err
	}
	resp, err := cli.httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	data := new(fbUserInfoResp)
	content, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, err
	}
	if data.ID == "" {
		return nil, fmt.Errorf("facebook user openid is empty")
	}
	return data.ToOauthUser(), nil
}

func NewFacebookOauth(cfg Config) Oauth {
	return &fbCli{
		clientId:     cfg.ClientId(),
		clientSecret: cfg.ClientSecret(),
		redirectURI:  cfg.RedirectURI(),
		httpCli:      http.Client{},
	}
}
