package vaulthcplib

import (
	"golang.org/x/oauth2"
	"time"
)

type TestTokenSource struct{}

func (*TestTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: "Test.Access.Token",
		Expiry:      time.Now().Add(time.Hour),
	}, nil
}
