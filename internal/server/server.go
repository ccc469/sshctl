package server

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ccc469/sshctl/internal"
	"github.com/ccc469/sshctl/internal/login"
	"github.com/melbahja/goph"

	"github.com/manifoldco/promptui"
)

func Run() bool {
	if internal.ShowServer {
		ShowServer()
		return true
	}

	if internal.ClearServer {
		return true
	}

	if internal.RemoveAliasName != "" {
		return true
	}

	return false
}

type Server struct {
	AliasName   string
	Ip          string
	Port        string
	Username    string
	Password    string
	PrivateKey  string
	AuthType    internal.AuthType
	AuthTypeMsg string
}

func ShowServer() {
	var (
		server Server
		pass   bool
		port   int
		client *goph.Client
	)

	results, err := internal.ReadFile(internal.ConfigPath)
	if err != nil {
		panic(err)
	}

	servers := make([]Server, 0)
	for _, item := range results {
		row := strings.Split(item, internal.Symbol)
		authType, err := strconv.Atoi(row[4])
		if err != nil {
			panic(err)
		}

		servers = append(servers, Server{AliasName: row[0],
			Ip:          row[1],
			Port:        row[3],
			Username:    row[2],
			Password:    row[5],
			PrivateKey:  row[5],
			AuthTypeMsg: internal.AuthTypeMap[internal.AuthType(authType)],
			AuthType:    internal.AuthType(authType),
		})
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .AliasName | cyan }}",
		Inactive: "  {{ .AliasName | cyan }}  ",
		Selected: "\U0001F336 {{ .AliasName | red | cyan }}",
		Details: `
--------- Details ----------
{{ "AliasName:" | faint }}	{{ .AliasName }}
{{ "Ip       :" | faint }}	{{ .Ip }}
{{ "Port     :" | faint }}	{{ .Port }}
{{ "Username :" | faint }}	{{ .Username }}
{{ "AuthType :" | faint }}	{{ .AuthType }}`,
	}

	searcher := func(input string, index int) bool {
		pepper := servers[index]
		name := strings.Replace(strings.ToLower(pepper.AliasName), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Please choose the server",
		Items:     servers,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("You choose server %s\n\n", servers[i].AliasName)

	server = servers[i]
	if server.AuthType == internal.Username {
		pass = true
	} else {
		pass = false
	}

	port, err = strconv.Atoi(server.Port)
	if err != nil {
		panic(err)
	}
	client, err = login.NewSSHClient(pass, server.Username, server.Password, server.Ip, port, server.PrivateKey)
	if err != nil {
		panic(err)
	}

	defer client.Close()
	login.RunTerminal(client, os.Stdout, os.Stderr)
}
