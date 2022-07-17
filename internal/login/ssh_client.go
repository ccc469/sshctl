package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/ccc469/sshctl/internal"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func Run() {
	var (
		auth goph.Auth
		err  error
	)

	if internal.Pass {
		auth = goph.Password(askPass("Enter SSH Password: "))
	} else {
		auth, err = goph.Key(internal.Key, getPassphrase(internal.Passphrase))
		if err != nil {
			panic(err)
		}
	}

	client, err := goph.NewConn(&goph.Config{
		User:     internal.User,
		Addr:     internal.Ip,
		Port:     uint(internal.Port),
		Auth:     auth,
		Callback: VerifyHost,
	})

	if err != nil {
		panic(err)
	}

	defer func() {
		client.Close()
		// 保存到本地
		
	}()
	RunTerminal(client, os.Stdout, os.Stdin)
}

func askPass(msg string) string {
	fmt.Print(msg)
	pass, err := term.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	fmt.Println("")
	return strings.TrimSpace(string(pass))
}
func getPassphrase(ask bool) string {
	if ask {
		return askPass("Enter Private Key Passphrase: ")
	}
	return ""
}

func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {

	//
	// If you want to connect to new hosts.
	// here your should check new connections public keys
	// if the key not trusted you shuld return an error
	//

	// hostFound: is host in known hosts file.
	// err: error if key not in known hosts file OR host in known hosts file but key changed!
	hostFound, err := goph.CheckKnownHost(host, remote, key, "")

	// Host in known hosts but key mismatch!
	// Maybe because of MAN IN THE MIDDLE ATTACK!
	if hostFound && err != nil {
		return err
	}

	// handshake because public key already exists.
	if hostFound && err == nil {
		return nil
	}

	// Ask user to check if he trust the host public key.
	if !askIsHostTrusted(host, key) {

		// Make sure to return error on non trusted keys.
		return errors.New("you typed no, aborted")
	}

	// Add the new host to known hosts file.
	return goph.AddKnownHost(host, remote, key, "")
}

func askIsHostTrusted(host string, key ssh.PublicKey) bool {

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Unknown Host: %s \nFingerprint: %s \n", host, ssh.FingerprintSHA256(key))
	fmt.Print("Would you like to add it? type yes or no: ")

	a, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	return strings.ToLower(strings.TrimSpace(a)) == "yes"
}

func RunTerminal(client *goph.Client, stdout, stderr io.Writer) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)

	session.Stdout = stdout
	session.Stderr = stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := term.GetSize(fd)
	if err != nil {
		panic(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}
	session.Shell()
	session.Wait()
	return nil
}
