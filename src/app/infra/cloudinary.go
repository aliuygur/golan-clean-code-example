package infra

import (
	"app/interfaces/errs"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func NewCloudinary(curl string) *Cloudinary {
	u, _ := url.Parse(curl)
	c := Cloudinary{URL: u}
	return &c
}

type Cloudinary struct {
	URL *url.URL
}

func (c *Cloudinary) ApiKey() string {
	return c.URL.User.Username()
}

func (c *Cloudinary) ApiSecret() string {
	p, _ := c.URL.User.Password()

	return p
}

func (c *Cloudinary) uploadURL() string {
	return fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", c.URL.Host)
}

func (c *Cloudinary) Upload(file string) (*http.Response, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	data := url.Values{
		"api_key":   []string{c.ApiKey()},
		"timestamp": []string{timestamp},
		"file":      []string{file},
	}

	// Signature
	hash := sha1.New()
	part := fmt.Sprintf("timestamp=%s%s", timestamp, c.ApiSecret())
	hash.Write([]byte(part))
	data.Set("signature", fmt.Sprintf("%x", hash.Sum(nil)))

	resp, err := http.PostForm(c.uploadURL(), data)

	if err != nil {
		return nil, errs.WrapMsg(err, "form can't posted")
	}

	if resp.StatusCode != http.StatusOK {
		var errMsg struct {
			Error struct {
				Message string
			}
		}
		json.NewDecoder(resp.Body).Decode(&errMsg)
		return nil, errs.BadRequest(errMsg.Error.Message)
	}

	return resp, nil
}
