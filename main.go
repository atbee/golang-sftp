package main

import (
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	rawurl := os.Getenv("SFTP_URL")

	parsedUrl, err := url.Parse(rawurl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse SFTP To Go URL: %s\n", err)
		os.Exit(1)
	}

	user := parsedUrl.User.Username()
	pass, _ := parsedUrl.User.Password()
	host, port, _ := net.SplitHostPort(parsedUrl.Host)

	fmt.Fprintf(os.Stdout, "Connecting to %s ...\n", host)

	var auths []ssh.AuthMethod

	if pass != "" {
		auths = append(auths, ssh.Password(pass))
	}

	config := ssh.ClientConfig{
		User:            user,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connecto to [%s]: %v\n", addr, err)
		os.Exit(1)
	}

	defer conn.Close()

	sc, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start SFTP subsystem: %v\n", err)
		os.Exit(1)
	}

	defer sc.Close()
}
