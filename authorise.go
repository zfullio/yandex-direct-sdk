package yandex_direct_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type AuthByCode struct {
	DeviceCode      string `json:"device_code"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	UserCode        string `json:"user_code"`
	VerificationUrl string `json:"verification_url"`
}

type errorResponse struct {
	Err string `json:"error"`
	Msg string `json:"error_description"`
}

func (e errorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Msg)
}

var ErrAuthorisationPending = errorResponse{
	Err: "authorization_pending",
	Msg: "User has not yet authorized your application",
}

type Token struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (c *client) Authorise() (err error) {
	authData, err := c.GetAccessCode()
	if err != nil {
		return err
	}

	fmt.Printf("user code: %s | %s \n", authData.UserCode, authData.VerificationUrl)

	token := Token{}
	done := make(chan bool, 1)
	duration := time.Duration(authData.Interval) * time.Second
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

outerLoop:
	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			token, err = c.GetTokenByCode(authData)
			if errors.Is(err, ErrAuthorisationPending) {
				fmt.Println(err)
				continue
			} else if err != nil {
				return err
			}
			if token.AccessToken != "" {
				done <- true
				break outerLoop
			}
		}
	}

	file, err := os.Create(fmt.Sprintf("%s.txt", c.Login))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer file.Close()

	_, err = file.WriteString(token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println(token.AccessToken)

	return nil
}

func (c *client) GetAccessCode() (auth AuthByCode, err error) {
	host := "oauth.yandex.ru"

	reqAccessURL := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/device/code",
	}
	param := url.Values{}
	param.Add("client_id", c.App.ID)

	resp, err := c.Tr.PostForm(reqAccessURL.String(), param)
	if err != nil {
		return auth, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return auth, fmt.Errorf("StatusCode: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return auth, err
	}

	err = json.Unmarshal(body, &auth)
	if err != nil {
		return auth, err
	}

	return auth, err
}

func (c *client) GetTokenByCode(code AuthByCode) (token Token, err error) {
	host := "oauth.yandex.ru"

	reqAccessURL := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "token",
	}
	param := url.Values{
		"client_id":     {c.App.ID},
		"grant_type":    {"device_code"},
		"code":          {code.DeviceCode},
		"client_secret": {c.App.Secret},
	}

	resp, err := c.Tr.PostForm(reqAccessURL.String(), param)
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()

	//if resp.StatusCode != http.StatusOK {
	//	return token, fmt.Errorf("StatusCode: %v", resp.StatusCode)
	//}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		return token, err
	}

	respErr := errorResponse{}
	err = json.Unmarshal(body, &respErr)
	if err != nil {
		return token, err
	}

	if respErr.Err != "" {
		return token, respErr
	}

	return token, err
}
