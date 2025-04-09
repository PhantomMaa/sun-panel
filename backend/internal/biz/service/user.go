package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"sun-panel/internal/biz/repository"
	"sun-panel/internal/infra/config"
	"sun-panel/internal/infra/zaplog"
	"sun-panel/internal/util"
)

type UserService struct {
	itemGroupRepo repository.IItemIconGroupRepo
	userRepo      repository.IUserRepo
}

type IUserService interface {
	CreateUser(user *repository.User) error
	GetOAuthLoginURL(provider string, redirectURI string) (string, error)
	HandleOAuthCallback(provider, code, redirectURI string) (*repository.User, error)
}

func NewUserService(userRepo repository.IUserRepo, itemGroupRepo repository.IItemIconGroupRepo) *UserService {
	return &UserService{userRepo: userRepo, itemGroupRepo: itemGroupRepo}
}

func (s *UserService) CreateUser(user *repository.User) error {
	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	defaultGroup := repository.ItemIconGroup{
		Title:  "APP",
		UserId: user.ID,
		Icon:   "material-symbols:ad-group-outline",
	}

	if err := s.itemGroupRepo.Save(&defaultGroup); err != nil {
		return err
	}

	return nil
}

// GetOAuthLoginURL generates the OAuth login URL for the specified provider
func (s *UserService) GetOAuthLoginURL(provider string, redirectURI string) (string, error) {
	// Find the provider config
	var providerConfig *config.OAuthProviderConfig
	for _, p := range config.AppConfig.OAuth.Providers {
		if strings.EqualFold(p.Name, provider) {
			providerConfig = &p
			break
		}
	}

	if providerConfig == nil {
		return "", errors.New("unsupported OAuth provider")
	}

	// Build OAuth URL
	authURL, err := url.Parse(providerConfig.AuthURL)
	if err != nil {
		return "", err
	}

	q := authURL.Query()
	q.Set("client_id", providerConfig.ClientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", providerConfig.Scopes)
	q.Set("response_type", "code")
	q.Set("state", util.GenerateRandomString(16)) // State should be stored in session and verified on callback

	authURL.RawQuery = q.Encode()
	return authURL.String(), nil
}

// HandleOAuthCallback processes the OAuth callback and returns or creates a user
func (s *UserService) HandleOAuthCallback(provider, code, redirectURI string) (*repository.User, error) {
	// Find the provider config
	var providerConfig *config.OAuthProviderConfig
	for _, p := range config.AppConfig.OAuth.Providers {
		if strings.EqualFold(p.Name, provider) {
			providerConfig = &p
			break
		}
	}

	if providerConfig == nil {
		return nil, errors.New("unsupported OAuth provider")
	}

	// Exchange code for token
	accessToken, err := s.exchangeCodeForToken(*providerConfig, code, redirectURI)
	if err != nil {
		zaplog.Logger.Error("Failed to exchange code for token: " + err.Error())
		return nil, err
	}

	// Get user info using the token
	userInfo, err := s.fetchUserInfo(*providerConfig, accessToken)
	if err != nil {
		zaplog.Logger.Error("Failed to fetch user info: " + err.Error())
		return nil, err
	}

	// Extract user identifier from the provider's response
	identifier, ok := userInfo[providerConfig.FieldMappingIdentifier].(string)
	if !ok || identifier == "" {
		return nil, errors.New("failed to get user identifier from OAuth provider")
	}

	// Check if user already exists
	user, err := s.userRepo.GetByOAuthID(provider, identifier)
	if err == nil {
		// 用户已存在，直接返回
		return &user, nil
	}

	// User doesn't exist, create a new one
	displayName, _ := userInfo[providerConfig.FieldMappingDisplayName].(string)
	email, _ := userInfo[providerConfig.FieldMappingEmail].(string)

	// OAuth users don't need a password as they authenticate through the provider
	newUser := &repository.User{
		Username:      identifier,
		Password:      "", // No password needed for OAuth users
		Name:          displayName,
		Mail:          email,
		Status:        1, // Active
		Role:          2, // Regular user
		OauthProvider: provider,
		OauthID:       identifier,
	}

	if err := s.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// exchangeCodeForToken exchanges the authorization code for an access token
func (s *UserService) exchangeCodeForToken(config config.OAuthProviderConfig, code, redirectURI string) (string, error) {
	// 创建 OAuth2 配置
	conf := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  redirectURI,
		Endpoint: oauth2.Endpoint{
			TokenURL:  config.TokenURL,
			AuthURL:   config.AuthURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	// 使用超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 交换 code 获取 token
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		zaplog.Logger.Error("Failed to exchange token: " + err.Error())
		return "", errors.New("failed to exchange code for token: " + err.Error())
	}

	// 获取 access_token
	var accessToken string

	// 首先检查 token.AccessToken
	if token.AccessToken != "" {
		accessToken = token.AccessToken
	} else {
		// 尝试从 Extra 中获取
		if tokenStr, ok := token.Extra("access_token").(string); ok && tokenStr != "" {
			accessToken = tokenStr
		} else {
			zaplog.Logger.Error("No access_token found in token response")
			return "", errors.New("no access_token in response")
		}
	}

	zaplog.Logger.Info("Token exchange successful")
	return accessToken, nil
}

// fetchUserInfo fetches user information from the OAuth provider
func (s *UserService) fetchUserInfo(config config.OAuthProviderConfig, accessToken string) (map[string]interface{}, error) {
	// 创建 OAuth2 token
	token := &oauth2.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}

	// 创建 OAuth2 配置
	conf := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: config.TokenURL,
			AuthURL:  config.AuthURL,
		},
	}

	// 使用超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建带有 token 的客户端
	client := conf.Client(ctx, token)

	// 发送请求获取用户信息
	resp, err := client.Get(config.UserInfoURL)
	if err != nil {
		zaplog.Logger.Error("Failed to send user info request: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zaplog.Logger.Error("Failed to read user info response: " + err.Error())
		return nil, err
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		zaplog.Logger.Error("User info request failed with status code: " + strconv.Itoa(resp.StatusCode))
		zaplog.Logger.Error("Response body: " + string(body))
		return nil, errors.New("failed to fetch user info: " + string(body))
	}

	// 解析 JSON 响应
	var userInfo map[string]any
	if err := json.Unmarshal(body, &userInfo); err != nil {
		zaplog.Logger.Error("Failed to parse user info response as JSON: " + err.Error())
		return nil, err
	}

	return userInfo, nil
}
