package simplessh

import (
	"bytes"
	"io"
	"net"
	"os"
	"time"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

type SimpleSsh_t struct {
	Client *ssh.Client
}

func (s *SimpleSsh_t) RemoteRun(cmd string, showRemoteLog bool) error {
	session, err := s.Client.NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	if showRemoteLog {
		stdout, _ := session.StdoutPipe()
		stderr, _ := session.StderrPipe()

		go io.Copy(os.Stdout, stdout)
		go io.Copy(os.Stdout, stderr)
	}

	err = session.Run(cmd)
	return err
}

func (s *SimpleSsh_t) RemoteRunGetResponse(cmd string) (string, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()
	output := &bytes.Buffer{}

	stdout, _ := session.StdoutPipe()
	stderr, _ := session.StderrPipe()

	go io.Copy(output, stdout)
	go io.Copy(output, stderr)

	err = session.Run(cmd)
	return output.String(), err
}

func (s *SimpleSsh_t) Copy(src string, dest string) error {
	session, err := s.Client.NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	err = scp.CopyPath(src, dest, session)
	if err != nil {
		return err
	}
	return nil
}

// timeout = 0 means no timeout
func Connect(address string, port string, user string, password string, timeout time.Duration) (SimpleSsh_t, error) {

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: timeout,
	}

	client, err := ssh.Dial("tcp", address+":"+port, config)
	if err != nil {
		return SimpleSsh_t{
			Client: nil,
		}, err
	}

	return SimpleSsh_t{
		Client: client,
	}, nil
}
