package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	LoginPath        = "auth/login"
	LogoutPath       = "auth/logout"
	RefreshTokenPath = "auth/refresh"
)

// AuthService : Provides Auth services provided by the API.
type AuthService service

// AuthResponse : Typical AuthService response.
type AuthResponse struct {
	Result  string  `json:"result"`
	Token   token   `json:"token"`
	Message *string `json:"message,omitempty"`
}

func (ar AuthResponse) GetResult() string {
	return ar.Result
}

// token : MangaDex token. Includes session and refresh token.
type token struct {
	Session string `json:"session"`
	Refresh string `json:"refresh"`
}

// Login : Login to MangaDex.
// https://api.mangadex.org/docs.html#operation/post-auth-login
func (s *AuthService) Login(user, pwd string) error {
	return s.LoginContext(context.Background(), user, pwd)
}

// LoginContext : Login with custom context.
func (s *AuthService) LoginContext(ctx context.Context, user, pwd string) error {
	u, _ := url.Parse(BaseAPI)
	u.Path = LoginPath

	// Create required request body.
	req := map[string]string{
		"username": user,
		"password": pwd,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	var ar AuthResponse
	if err = s.client.RequestAndDecode(ctx, http.MethodPost, u.String(), bytes.NewBuffer(rBytes), &ar); err != nil {
		return err
	}

	// Set client token and header for authorization.
	s.SetRefreshToken(ar.Token.Refresh)
	s.client.header.Set("Authorization", fmt.Sprintf("Bearer %s", ar.Token.Session))
	return nil
}

// Logout : Logout of MangaDex and invalidates all tokens.
// https://api.mangadex.org/docs.html#operation/post-auth-logout
func (s *AuthService) Logout() error {
	return s.LogoutContext(context.Background())
}

// LogoutContext : Logout with custom context.
func (s *AuthService) LogoutContext(ctx context.Context) error {
	u, _ := url.Parse(BaseAPI)
	u.Path = LogoutPath

	var r Response
	if err := s.client.RequestAndDecode(ctx, http.MethodPost, u.String(), nil, &r); err != nil {
		return err
	}

	// Remove the stored client token and also authorization header.
	s.SetRefreshToken("")
	s.client.header.Del("Authorization")
	return nil
}

// RefreshSessionToken : Refresh session token using refresh token.
// https://api.mangadex.org/docs.html#operation/post-auth-refresh
func (s *AuthService) RefreshSessionToken() error {
	return s.RefreshSessionTokenContext(context.Background())
}

// RefreshSessionTokenContext : refreshToken with custom context.
func (s *AuthService) RefreshSessionTokenContext(ctx context.Context) error {
	u, _ := url.Parse(BaseAPI)
	u.Path = RefreshTokenPath

	// Create required request body.
	req := map[string]string{
		"token": s.client.refreshToken,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	var ar AuthResponse
	if err = s.client.RequestAndDecode(ctx, http.MethodPost, u.String(), bytes.NewBuffer(rBytes), &ar); err != nil {
		return err
	}

	// Update tokens
	s.SetRefreshToken(ar.Token.Refresh)
	s.client.header.Set("Authorization", fmt.Sprintf("Bearer %s", ar.Token.Session))
	return nil
}

// IsLoggedIn : Return true when client logged in and false otherwise.
func (s *AuthService) IsLoggedIn() bool {
	return s.client.header.Get("Authorization") != ""
}

// GetRefreshToken : Get the current refresh token of the client.
func (s *AuthService) GetRefreshToken() string {
	return s.client.refreshToken
}

// SetRefreshToken : Set the refresh token for the client.
func (s *AuthService) SetRefreshToken(refreshToken string) {
	s.client.refreshToken = refreshToken
}
