package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	// 记录重定向 URI，用于调试
	zaplog.Logger.Info("OAuth callback with redirectURI: " + redirectURI)

	// Exchange code for token
	tokenData, err := s.exchangeCodeForToken(*providerConfig, code, redirectURI)
	if err != nil {
		zaplog.Logger.Error("Failed to exchange code for token: " + err.Error())
		return nil, err
	}

	// 记录 token 信息（注意不要记录完整的 token，仅记录是否存在）
	zaplog.Logger.Info("Token exchange successful, access_token exists: " + strconv.FormatBool(tokenData["access_token"] != ""))

	// Get user info using the token
	userInfo, err := s.fetchUserInfo(*providerConfig, tokenData["access_token"])
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
func (s *UserService) exchangeCodeForToken(config config.OAuthProviderConfig, code, redirectURI string) (map[string]string, error) {
	data := url.Values{}
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	// 记录请求参数（不包含敏感信息）
	zaplog.Logger.Info("Exchanging code for token with provider: " + config.Name)
	zaplog.Logger.Info("Token URL: " + config.TokenURL)
	zaplog.Logger.Info("Redirect URI: " + redirectURI)

	req, err := http.NewRequest("POST", config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		zaplog.Logger.Error("Failed to create token request: " + err.Error())
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zaplog.Logger.Error("Failed to send token request: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zaplog.Logger.Error("Failed to read token response: " + err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		zaplog.Logger.Error("Token exchange failed with status code: " + strconv.Itoa(resp.StatusCode))
		zaplog.Logger.Error("Response body: " + string(body))
		return nil, errors.New("failed to exchange code for token: " + string(body))
	}

	var tokenData map[string]string
	if err := json.Unmarshal(body, &tokenData); err != nil {
		zaplog.Logger.Error("Failed to parse token response as JSON: " + err.Error())
		
		// Some providers might return non-JSON format
		values, err := url.ParseQuery(string(body))
		if err != nil {
			zaplog.Logger.Error("Failed to parse token response as URL query: " + err.Error())
			return nil, err
		}

		tokenData = make(map[string]string)
		for k, v := range values {
			if len(v) > 0 {
				tokenData[k] = v[0]
			}
		}
	}

	zaplog.Logger.Info("Token exchange successful")
	return tokenData, nil
}

// fetchUserInfo fetches user information from the OAuth provider
func (s *UserService) fetchUserInfo(config config.OAuthProviderConfig, accessToken string) (map[string]interface{}, error) {
	// 记录请求参数（不包含敏感信息）
	zaplog.Logger.Info("Fetching user info from provider: " + config.Name)
	zaplog.Logger.Info("User info URL: " + config.UserInfoURL)

	req, err := http.NewRequest("GET", config.UserInfoURL, nil)
	if err != nil {
		zaplog.Logger.Error("Failed to create user info request: " + err.Error())
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zaplog.Logger.Error("Failed to send user info request: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zaplog.Logger.Error("Failed to read user info response: " + err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		zaplog.Logger.Error("User info request failed with status code: " + strconv.Itoa(resp.StatusCode))
		zaplog.Logger.Error("Response body: " + string(body))
		return nil, errors.New("failed to fetch user info: " + string(body))
	}

	var userInfo map[string]any
	if err := json.Unmarshal(body, &userInfo); err != nil {
		zaplog.Logger.Error("Failed to parse user info response as JSON: " + err.Error())
		return nil, err
	}

	// 记录成功获取用户信息（不记录具体内容，保护隐私）
	zaplog.Logger.Info("Successfully fetched user info")
	
	// 记录字段映射信息，帮助调试
	zaplog.Logger.Info("Field mapping - identifier: " + config.FieldMappingIdentifier)
	zaplog.Logger.Info("Field mapping - display name: " + config.FieldMappingDisplayName)
	zaplog.Logger.Info("Field mapping - email: " + config.FieldMappingEmail)
	
	// 检查必要字段是否存在
	_, hasIdentifier := userInfo[config.FieldMappingIdentifier]
	_, hasDisplayName := userInfo[config.FieldMappingDisplayName]
	_, hasEmail := userInfo[config.FieldMappingEmail]
	
	zaplog.Logger.Info("Field exists - identifier: " + strconv.FormatBool(hasIdentifier))
	zaplog.Logger.Info("Field exists - display name: " + strconv.FormatBool(hasDisplayName))
	zaplog.Logger.Info("Field exists - email: " + strconv.FormatBool(hasEmail))

	return userInfo, nil
}
