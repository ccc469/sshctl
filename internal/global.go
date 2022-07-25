package internal

import (
	"os"
	"runtime"
	"strings"

	"github.com/melbahja/goph"
	"github.com/pkg/sftp"
)

var (
	Delimiter string = setDelimiter()

	Auth       goph.Auth
	Client     *goph.Client
	Ip         string
	User       string
	Port       int
	PrivateKey string
	Pass       bool
	Sftpc      *sftp.Client
	Save       bool
	AliasName  string

	ClearServer     bool
	RemoveAliasName string
	ShowServer      bool

	PrefixDir  string = ".sshctl"
	ConfigFile string = "config"
	Symbol     string = " "
	HomePath   string = strings.Join([]string{os.Getenv("HOME"), PrefixDir}, Delimiter)
	ConfigPath string = strings.Join([]string{os.Getenv("HOME"), PrefixDir, ConfigFile}, Delimiter)

	AuthTypeMap map[AuthType]string = map[AuthType]string{
		Username: "Username",
		SSHkey:   "SSHkey",
	}
)

type AuthType int

const (
	Username AuthType = iota // 0
	SSHkey                   // 1
)

const (
	Windows string = "\\"
	AnyUnix string = "/"
)

func setDelimiter() string {
	var delimiter string
	switch runtime.GOOS {
	case "windows":
		delimiter = string(Windows)
	default: // Unix
		delimiter = string(AnyUnix)
	}
	return delimiter
}
