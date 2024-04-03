package util

import (
	"strings"
)

const (
	admin      = "ADMIN"
	moder      = "MODERATOR"
	superadmin = "SUPERADMIN"
	user       = "USER"
)

func IsSupportedRole(role string) bool {
	switch role {
	case admin, moder, user:
		return true
	}
	return false
}
func IsSupportedInstagram(url string) bool {
	switch {
	case strings.Contains(url, "https://www.instagram.com"), strings.Contains(url, "https://instagram.com"), strings.Contains(url, "www.instagram.com"), strings.Contains(url, "instagram.com"):
		return true
	}
	return false
}
func IsSupportedFacebook(url string) bool {
	switch {
	case strings.Contains(url, "https://www.facebook.com"), strings.Contains(url, "https://facebook.com"), strings.Contains(url, "www.facebook.com"), strings.Contains(url, "facebook.com"):
		return true
	}
	return false
}
