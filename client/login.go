package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net/url"
	"regexp"

	"github.com/chomosuke/cf-tool/cookiejar"
	"github.com/chomosuke/cf-tool/util"
	"github.com/fatih/color"
)

// genFtaa generate a random one
func genFtaa() string {
	return util.RandString(18)
}

// genBfaa generate a bfaa
func genBfaa() string {
	return "f1b3f18c715565b589b7823cda7448ce"
}

// ErrorNotLogged not logged in
var ErrorNotLogged = "Not logged in"

// findHandle if logged return (handle, nil), else return ("", ErrorNotLogged)
func findHandle(body []byte) (string, error) {
	reg := regexp.MustCompile(`handle = "([\s\S]+?)"`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New(ErrorNotLogged)
	}
	return string(tmp[1]), nil
}

func findCsrf(body []byte) (string, error) {
	reg := regexp.MustCompile(`csrf='(.+?)'`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New("Cannot find csrf")
	}
	return string(tmp[1]), nil
}

// Login codeforces with handler and password
func (c *Client) Login() (err error) {
	color.Cyan("Login %v...\n", c.HandleOrEmail)

	password, err := c.DecryptPassword()
	if err != nil {
		return
	}

	jar, _ := cookiejar.New(nil)
	c.client.Jar = jar
	body, err := util.GetBody(c.client, c.host+"/enter")
	if err != nil {
		return
	}

	csrf, err := findCsrf(body)
	if err != nil {
		return
	}

	ftaa := genFtaa()
	bfaa := genBfaa()

	body, err = util.PostBody(c.client, c.host+"/enter", url.Values{
		"csrf_token":    {csrf},
		"action":        {"enter"},
		"ftaa":          {ftaa},
		"bfaa":          {bfaa},
		"handleOrEmail": {c.HandleOrEmail},
		"password":      {password},
		"_tta":          {"176"},
		"remember":      {"on"},
	})
	if err != nil {
		return
	}

	handle, err := findHandle(body)
	if err != nil {
		return
	}

	c.Ftaa = ftaa
	c.Bfaa = bfaa
	c.Handle = handle
	c.Jar = jar
	color.Green("Welcome %v~", handle)
	return c.save()
}

func createHash(key string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hasher.Sum(nil)
}

func encrypt(handle, password string) (ret string, err error) {
	block, err := aes.NewCipher(createHash("glhf" + handle + "233"))
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	text := gcm.Seal(nonce, nonce, []byte(password), nil)
	ret = hex.EncodeToString(text)
	return
}

func decrypt(handle, password string) (ret string, err error) {
	data, err := hex.DecodeString(password)
	if err != nil {
		err = errors.New("Cannot decode the password")
		return
	}
	block, err := aes.NewCipher(createHash("glhf" + handle + "233"))
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonceSize := gcm.NonceSize()
	nonce, text := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, text, nil)
	if err != nil {
		return
	}
	ret = string(plain)
	return
}

// DecryptPassword get real password
func (c *Client) DecryptPassword() (string, error) {
	if len(c.Password) == 0 {
		return "", errors.New("empty password")
	}
	if len(c.HandleOrEmail) == 0 {
		return "", errors.New("empty handle/email")
	}
	return decrypt(c.HandleOrEmail, c.Password)
}

func (c *Client) Setup(handle string, rawPassword string) (err error) {
	c.Handle = handle
	c.HandleOrEmail = handle

	c.Password, err = encrypt(c.HandleOrEmail, rawPassword)
	if err != nil {
		return
	}

	return c.Login()
}
