package rdapi

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (conf *RDConfig) GetAuthorization(method string, endpoint string) string {
	params := getParamArray([][]string{
		{"oauth_consumer_key", conf.ConsumerKey},
		{"oauth_nonce", genNonce()},
		{"oauth_signature_method", "HMAC-SHA1"},
		{"oauth_timestamp", strconv.Itoa(int(time.Now().Unix()))},
		{"oauth_token", conf.secondAuth.Token},
		{"oauth_version", "1.0"},
	})

	signature := generateOAuthKey(method, conf.ServiceEndpoint+endpoint, params, conf.ConsumerSecret, conf.secondAuth.Secret)

	paramsBody := []string{params[0], params[1], "oauth_signature=" + url.QueryEscape(signature)}
	paramsBody = append(paramsBody, params[2:]...)

	return "OAuth " + strings.Join(paramsBody, ",")
}

func genNonce() string {
	b := make([]rune, 16)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateOAuthKey(method string, endpoint string, params []string, secret string, tokenSecret string) string {
	base := []byte(method + "&" + url.QueryEscape(endpoint) + "&" + url.QueryEscape(strings.Join(params, "&")))
	mac := hmac.New(sha1.New, []byte(url.QueryEscape(secret)+"&"+url.QueryEscape(tokenSecret)))
	mac.Write(base)

	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func getParamArray(src [][]string) []string {
	var dst = make([]string, len(src))
	for i, v := range src {
		dst[i] = v[0] + "=" + url.QueryEscape(v[1])
	}

	return dst
}
