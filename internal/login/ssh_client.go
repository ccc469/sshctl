package login

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
		err      error
		password string
	)

	if internal.Pass {
		password = askPass("Enter SSH Password: ")
	}

	client, err := NewSSHClient(internal.Pass, internal.User, password, internal.Ip, internal.Port, internal.PrivateKey)
	if err != nil {
		fmt.Println("Connection failed, please check your connection setting")
		return
	}

	if internal.AliasName == "" {
		internal.AliasName = internal.Ip
	}

	if internal.Save {
		SaveToLocal(internal.Pass, internal.Ip, internal.Port, internal.User, password, internal.PrivateKey, internal.AliasName)
	}

	defer client.Close()
	RunTerminal(client, os.Stdout, os.Stdin)
}

func NewSSHClient(pass bool, user string, password string, ip string, port int, privateKey string) (*goph.Client, error) {
	var (
		auth goph.Auth
		err  error
	)

	if pass {
		auth = goph.Password(password)
	} else {
		auth, err = goph.Key(privateKey, "")
		if err != nil {
			return nil, err
		}
	}

	client, err := goph.NewConn(&goph.Config{
		User:     user,
		Addr:     ip,
		Port:     uint(port),
		Auth:     auth,
		Callback: VerifyHost,
	})

	return client, err
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

func SaveToLocal(hasPass bool, ip string, port int, user string, password string, privateKey string, aliasName string) {

	var (
		authType internal.AuthType
		key      string
	)
	if _, err := os.Stat(internal.HomePath); err != nil {
		fmt.Printf("%s path not exists，create now\n", internal.HomePath)
		err := os.MkdirAll(internal.HomePath, 0755)
		if err != nil {
			log.Println(err)
			return
		}
	}

	file, err := os.OpenFile(internal.ConfigPath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil && os.IsNotExist(err) {
		fmt.Printf("%s file not exists，create now\n", internal.ConfigPath)
		os.Create(internal.ConfigPath)
	}
	defer file.Close()

	results, err := internal.ReadFile(internal.ConfigPath)
	if err != nil {
		panic(err)
	}

	for _, item := range results {
		server := strings.Split(item, internal.Symbol)
		if aliasName == server[0] {
			fmt.Printf("Waring: server alias name [%s] already exits\n\n", aliasName)
			return
		}
	}

	if hasPass {
		authType = internal.Username
		key = password
	} else {
		authType = internal.SSHkey
		key = privateKey
	}

	content := strings.Join([]string{aliasName, ip, user, fmt.Sprintf("%d", port), fmt.Sprintf("%v", authType), key}, internal.Symbol) + "\n"
	if len(results) != 0 {
		writer := bufio.NewWriter(file)
		writer.WriteString(content)
		writer.Flush()
	} else {
		err := ioutil.WriteFile(internal.ConfigPath, []byte(content), 0666)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Your Connect Saved to %s\n\n", internal.ConfigPath)
}
