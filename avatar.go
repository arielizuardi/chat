package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrNoAvatarURL is
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar is
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar auth
type AuthAvatar struct{}

// UseAuthAvatar use
var UseAuthAvatar AuthAvatar

// GetAvatarURL get
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}

	return "", ErrNoAvatarURL
}

// GravatarAvatar gravatar
type GravatarAvatar struct{}

// UseGravatar use
var UseGravatar GravatarAvatar

// GetAvatarURL get
func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("https://www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}

	return "", ErrNoAvatarURL
}
