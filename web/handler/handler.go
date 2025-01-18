package handler

import (
	"encoding/json"
	"fmt"
	logs "goauth/pkg/logger"
	"goauth/pkg/oauth"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func LoginPageHandler(c *gin.Context) {
	ctx := c.Request.Context()
	l := logs.GetLoggerctx(ctx)
	tmpl, err := template.ParseFiles(filepath.Join(viper.GetString("app.uiTemplates"), "loginpage.html"))
	if err != nil {
		l.Sugar().Error("error parsing the file", err.Error())
		return
	}
	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		l.Sugar().Error("error in executing ther template", err)
		return
	}

}

// google auth handler
func GoogleauthHandler(c *gin.Context) {
	ctx := c.Request.Context()
	l := logs.GetLoggerctx(ctx)
	state, err := oauth.GenerateState()
	if err != nil {
		l.Sugar().Error("state generation failed", err)
		return
	}

	clientID := []byte(os.Getenv("CLIENT-ID"))
	params := url.Values{}
	params.Add("client_id", string(clientID))
	params.Add("redirect_url", viper.GetString("oauth.redirectURL"))
	params.Add("response_type", "code")         // it gives us back the authorization code using which we can trade to get our access token
	params.Add("scope", "openid email profile") // i might access openid, email,profile
	params.Add("state", state)
	authURL := fmt.Sprintf("%s?%s", viper.GetString("oauth.googleAuthURL"), params.Encode())
	http.Redirect(c.Writer, c.Request, authURL, http.StatusTemporaryRedirect) // we are redirecting it to our auth url
}

func CallbackHandler(c *gin.Context) {
	ctx := c.Request.Context()
	l := logs.GetLoggerctx(ctx)
	//state := c.Request.URL.Query().Get("state")
	// compare our state in the db with the state received
	// if they dont match then network has been tampered
	// so dont allow them to login
	code := c.Query("code") //code is the authorisation code with which we will trade to get access token
	if code == "" {
		l.Sugar().Error("authorization code is empty")
		return
	}

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", os.Getenv("CLIENT-ID"))
	data.Set("client_secret", os.Getenv("OAUTH-CLIENT-SECRET"))
	data.Set("redirect_uri", viper.GetString("oauth.redirectURL"))
	data.Set("grant_type", "authorization_code")

	//back channel post
	resp, err := http.PostForm(viper.GetString("oauth.googleTokenURL"), data)
	if err != nil {
		l.Sugar().Error("form post request failed", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Sugar().Error("reading response body failed", err)
		return
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"` // this idtoken is a jwt token
	}

	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		l.Sugar().Error("unmarshall json failed", err)
		return
	}

	// main puropse of oauth that is to authorize
	// to get user info
	req, _ := http.NewRequest("GET", viper.GetString("oauth.googleUserInfoURL"), nil)

	req.Header.Add("Authorization", "Bearer"+tokenResp.AccessToken)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		l.Sugar().Error("get user request failed", err)
		return
	}
	defer resp.Body.Close()

	http.Redirect(c.Writer, c.Request, "/web/home", http.StatusTemporaryRedirect)
}

func HomePageHandler(c *gin.Context) {
	ctx := c.Request.Context()
	l := logs.GetLoggerctx(ctx)
	tmpl, err := template.ParseFiles(filepath.Join(viper.GetString("app.uiTemplates"), "home.html"))
	if err != nil {
		l.Sugar().Error("error parsing the file", err.Error())
		return
	}
	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		l.Sugar().Error("error in executing ther template", err)
		return
	}

}
