package rdapi

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type RDConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	SecondSecret   string

	Endpoint        string
	ServiceEndpoint string

	firstAuth  authKeys
	secondAuth authKeys
}

type authKeys struct {
	Token  string
	Secret string
}

func New(consumerKey string, consumerSecret string, secondSecret string) *RDConfig {
	return &RDConfig{
		ConsumerKey:     consumerKey,
		ConsumerSecret:  consumerSecret,
		SecondSecret:    secondSecret,
		Endpoint:        "http://uk.rdbranch.com/OAuth/V10a",
		ServiceEndpoint: "http://uk.rdbranch.com/WebServices/Epos/V1",
	}
}

func (keys *authKeys) Parse(values url.Values) {
	keys.Token = values.Get("oauth_token")
	keys.Secret = values.Get("oauth_token_secret")
}

func (keys *authKeys) Valid() bool {
	return len(keys.Token) > 0 && len(keys.Secret) > 0
}

func (conf *RDConfig) doOAuth(params []string) (url.Values, error) {
	req, _ := http.NewRequest("POST", conf.Endpoint, strings.NewReader(strings.Join(params, "&")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf8")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	client.Transport = &http.Transport{
		DisableCompression: true,
	}

	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)

	return url.ParseQuery(string(body))
}

func (conf *RDConfig) doFirstAuth() error {
	params := getParamArray([][]string{
		{"oauth_consumer_key", conf.ConsumerKey},
		{"oauth_nonce", genNonce()},
		{"oauth_signature_method", "HMAC-SHA1"},
		{"oauth_timestamp", strconv.Itoa(int(time.Now().Unix()))},
		{"oauth_version", "1.0"},
		{"scope", "http://app.restaurantdiary.com/WebServices/Epos/v1"},
		{"second_secret", conf.SecondSecret},
	})

	signature := generateOAuthKey("POST", conf.Endpoint, params, conf.ConsumerSecret, "")

	paramsBody := []string{params[0], params[1], "oauth_signature=" + url.QueryEscape(signature)}
	paramsBody = append(paramsBody, params[2:]...)

	values, err := conf.doOAuth(paramsBody)
	if err != nil {
		fmt.Println("Could not fetch first set of OAuth keys")
		return err
	}

	conf.firstAuth.Parse(values)
	if !conf.firstAuth.Valid() {
		// abort
		return errors.New("Oauth Tokens invalid")
	}

	return nil
}

func (conf *RDConfig) doSecondAuth() error {
	params := getParamArray([][]string{
		{"oauth_consumer_key", conf.ConsumerKey},
		{"oauth_nonce", genNonce()},
		{"oauth_signature_method", "HMAC-SHA1"},
		{"oauth_timestamp", strconv.Itoa(int(time.Now().Unix()))},
		{"oauth_token", conf.firstAuth.Token},
		{"oauth_version", "1.0"},
	})

	signature := generateOAuthKey("POST", conf.Endpoint, params, conf.ConsumerSecret, conf.firstAuth.Secret)

	paramsBody := []string{params[0], params[1], "oauth_signature=" + url.QueryEscape(signature)}
	paramsBody = append(paramsBody, params[2:]...)

	values, err := conf.doOAuth(paramsBody)
	if err != nil {
		fmt.Println("Could not fetch second set of OAuth keys")
		return err
	}

	conf.secondAuth.Parse(values)
	if !conf.secondAuth.Valid() {
		// abort
		return errors.New("Oauth Tokens invalid")
	}

	return nil
}

func (conf *RDConfig) Authenticate() error {
	if err := conf.doFirstAuth(); err != nil {
		return err
	}
	if err := conf.doSecondAuth(); err != nil {
		return err
	}

	fmt.Println("Save these:", conf.secondAuth.Token, conf.secondAuth.Secret)

	return nil
}

func (conf *RDConfig) SetKeys(token string, secret string) {
	conf.secondAuth.Token = token
	conf.secondAuth.Secret = secret
}

func (conf *RDConfig) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, conf.ServiceEndpoint+urlStr, nil)
	req.Header.Set("Authorization", conf.GetAuthorization(method, urlStr))
	req.Header.Set("Accept", "application/json")

	return req, err
}
