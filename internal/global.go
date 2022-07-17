package internal

import (
	"github.com/melbahja/goph"
	"github.com/pkg/sftp"
)

var (
	Auth       goph.Auth
	Client     *goph.Client
	Ip         string
	User       string
	Port       int
	Key        string
	Pass       bool
	Passphrase bool
	Sftpc      *sftp.Client
	Save       bool
)
