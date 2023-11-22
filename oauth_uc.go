package ucgo

import "github.com/jiebutech/uc/model"

type OauthClient struct {
	config OauthConfig
}

func NewOauthClient(cfg OauthConfig) *OauthClient {
	return &OauthClient{config: cfg}
}

func (o *OauthClient) Login(code string) (string, model.UserEntity, error) {
	return "", nil, nil
}
