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

type authKeys struct {
	Token  string
	Secret string
}

func (keys *authKeys) Parse(values url.Values) {
	keys.Token = values.Get("oauth_token")
	keys.Secret = values.Get("oauth_token_secret")
}

func (keys *authKeys) Valid() bool {
	return len(keys.Token) > 0 && len(keys.Secret) > 0
}

func (conf *RDConfig) GetAuthorization(method string, endpoint string) string {
	if !conf.secondAuth.Valid() {
		return "NOT AUTHENTICATED"
	}

	nonce := genNonce()
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	u, _ := url.Parse(conf.ServiceEndpoint + endpoint)

	params := []string{
		"oauth_consumer_key=" + conf.ConsumerKey,
		"oauth_nonce=" + nonce,
		"oauth_signature_method=HMAC-SHA1",
		"oauth_timestamp=" + timestamp,
		"oauth_token=" + url.QueryEscape(conf.secondAuth.Token),
		"oauth_version=1.0",
	}

	// prepend querystring params to signature parameters
	query := u.Query()
	for k, v := range query {
		p := []string{k + "=" + url.QueryEscape(v[0])}
		params = append(p, params...)
	}

	// Remove querystring data from url
	u.RawQuery = ""
	u.Fragment = ""

	signature := generateOAuthKey(method, u.String(), params, conf.ConsumerSecret, conf.secondAuth.Secret)

	return "OAuth " + strings.Join([]string{
		"oauth_token=\"" + url.QueryEscape(conf.secondAuth.Token) + "\"",
		"oauth_consumer_key=\"" + conf.ConsumerKey + "\"",
		"oauth_nonce=\"" + nonce + "\"",
		"oauth_signature_method=\"HMAC-SHA1\"",
		"oauth_signature=\"" + url.QueryEscape(signature) + "\"",
		"oauth_version=\"1.0\"",
		"oauth_timestamp=\"" + timestamp + "\"",
	}, ",")
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
