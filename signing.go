package rdapi

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (conf *RDConfig) GetAuthorization(method string, endpoint string) string {
	nonce := genNonce()
	timestamp := strconv.Itoa(int(time.Now().Unix()))

	params := []string{
		"oauth_consumer_key=" + conf.ConsumerKey,
		"oauth_nonce=" + nonce,
		"oauth_signature_method=HMAC-SHA1",
		"oauth_timestamp=" + timestamp,
		"oauth_token=" + url.QueryEscape(conf.secondAuth.Token),
		"oauth_version=1.0",
	}

	signature := generateOAuthKey(method, conf.ServiceEndpoint+endpoint, params, conf.ConsumerSecret, conf.secondAuth.Secret)

	return "OAuth oauth_token=\"" + url.QueryEscape(conf.secondAuth.Token) + "\",oauth_consumer_key=\"" + conf.ConsumerKey + "\",oauth_nonce=\"" + nonce + "\",oauth_signature_method=\"HMAC-SHA1\",oauth_signature=\"" + url.QueryEscape(signature) + "\",oauth_version=\"1.0\",oauth_timestamp=\"" + timestamp + "\""
}

func genNonce() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}

func generateOAuthKey(method string, endpoint string, params []string, secret string, tokenSecret string) string {
	base := []byte(method + "&" + url.QueryEscape(endpoint) + "&" + url.QueryEscape(strings.Join(params, "&")))

	mac := hmac.New(sha1.New, []byte(url.QueryEscape(secret)+"&"+url.QueryEscape(tokenSecret)))
	mac.Write(base)

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
