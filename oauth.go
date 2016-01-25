package rdapi

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func New(rid int, consumerKey string, consumerSecret string, secondSecret string) *RDConfig {
	return &RDConfig{
		ConsumerKey:     consumerKey,
		ConsumerSecret:  consumerSecret,
		SecondSecret:    secondSecret,
		RestaurantID:    rid,
		Endpoint:        "http://uk.rdbranch.com/OAuth/V10a",
		ServiceEndpoint: "http://uk.rdbranch.com/WebServices/Epos/v1",
	}
}

func (conf *RDConfig) doOAuth(params []string, sig string, pos int) (*authKeys, error) {
	var body []string
	body = append(body, params[:pos]...)
	body = append(body, "oauth_signature="+url.QueryEscape(sig))
	body = append(body, params[pos:]...)

	req, err := http.NewRequest("POST", conf.Endpoint, strings.NewReader(strings.Join(body, "&")))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf8")
	client := &http.Client{
		Timeout: time.Duration(3 * time.Second),
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, _ := ioutil.ReadAll(res.Body)

	query, err := url.ParseQuery(string(resBody))
	if err != nil {
		return nil, err
	}

	keys := authKeys{}
	keys.Parse(query)
	if !keys.Valid() {
		return nil, errors.New("OAuth Tokens Invalid")
	}

	return &keys, nil
}

func (conf *RDConfig) doFirstAuth() error {
	params := []string{
		"oauth_consumer_key=" + url.QueryEscape(conf.ConsumerKey),
		"oauth_nonce=" + genNonce(),
		"oauth_signature_method=HMAC-SHA1",
		"oauth_timestamp=" + strconv.Itoa(int(time.Now().Unix())),
		"oauth_version=1.0",
		"scope=" + url.QueryEscape("http://app.restaurantdiary.com/WebServices/Epos/v1"),
		"second_secret=" + url.QueryEscape(conf.SecondSecret),
	}

	signature := generateOAuthKey("POST", conf.Endpoint, params, conf.ConsumerSecret, "")
	keys, err := conf.doOAuth(params, signature, 3)
	if keys != nil {
		conf.firstAuth = *keys
	}

	return err
}

func (conf *RDConfig) doSecondAuth() error {
	params := []string{
		"oauth_consumer_key=" + url.QueryEscape(conf.ConsumerKey),
		"oauth_nonce=" + genNonce(),
		"oauth_signature_method=HMAC-SHA1",
		"oauth_timestamp=" + strconv.Itoa(int(time.Now().Unix())),
		"oauth_token=" + url.QueryEscape(conf.firstAuth.Token),
		"oauth_version=1.0",
	}

	signature := generateOAuthKey("POST", conf.Endpoint, params, conf.ConsumerSecret, conf.firstAuth.Secret)
	keys, err := conf.doOAuth(params, signature, 4)
	if keys != nil {
		conf.secondAuth = *keys
	}

	return err
}

func (conf *RDConfig) Authenticate() error {
	if err := conf.doFirstAuth(); err != nil {
		return err
	}

	return conf.doSecondAuth()
}

func (conf *RDConfig) SetKeys(token string, secret string) {
	conf.secondAuth.Token = token
	conf.secondAuth.Secret = secret
}

func (conf *RDConfig) RestaurantRequest(method, urlStr string, body io.Reader) (*http.Client, *http.Request, error) {
	return conf.NewRequest(method, "/Restaurant/"+strconv.Itoa(conf.RestaurantID)+urlStr, body)
}

func (conf *RDConfig) NewRequest(method, urlStr string, body io.Reader) (*http.Client, *http.Request, error) {
	req, err := http.NewRequest(method, conf.ServiceEndpoint+urlStr, nil)

	req.Header.Set("Authorization", conf.GetAuthorization(method, urlStr))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	return client, req, err
}
