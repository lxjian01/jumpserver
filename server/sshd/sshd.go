package sshd

import (
	"encoding/binary"
	"fmt"
	"github.com/creack/pty"
	"golang.org/x/crypto/ssh"
	"jumpserver/config"
	"jumpserver/log"
	"jumpserver/pools"
	"jumpserver/utils/terminals"
	"net"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"unsafe"
)
const (
	mysqlPrompt = "Enter password: "

	mysqlShellFilename = "mysql"
)
func StartSshdServer(c *config.SshdConfig)  {

	linuxServer := terminals.LinuxServer{
		Host: c.Host,
		Port: c.Port,
	}
	config := linuxServer.GetSshServerConfig()
	// Once a ServerConfig has been configured, connections can be accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:2200")
	if err != nil {
		log.Error("Failed to listen on 2200", err)
	}

	// Accept all connections
	log.Info("Listening on 2200...")
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Error("Failed to accept incoming connection", err)
			continue
		}
		// Before use, a handshake must be performed on the incoming net.Conn.
		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		if err != nil {
			log.Error("Failed to handshake", err)
			continue
		}

		log.Infof("New SSH connection from %s (%s) \n", sshConn.RemoteAddr(), sshConn.ClientVersion())
		// Discard all global out-of-band Requests
		pools.Pool.Submit(func() {
			ssh.DiscardRequests(reqs)
		})
		// Accept all channels
		pools.Pool.Submit(func() {
			for newChannel := range chans {
				handleChannel(newChannel)
			}
		})
	}
}

func handleChannel(newChannel ssh.NewChannel) {
	// Since we're handling a shell, we expect a
	// channel type of "session". The also describes
	// "x11", "direct-tcpip" and "forwarded-tcpip"
	// channel types.
	if t := newChannel.ChannelType(); t != "session" {
		newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
		return
	}

	// At this point, we have the opportunity to reject the client's
	// request for another logical connection
	channel, requests, err := newChannel.Accept()
	if err != nil {
		log.Error("Could not accept channel", err)
		return
	}

	// Fire up bash for this session
	execCmd := exec.Command("unshare")


	// Allocate a terminal for this channel
	log.Info("Creating pty...")
	pty, err := pty.Start(execCmd)
	if err != nil {
		log.Error("Could not start pty", err)
		close(execCmd,pty,channel)
		return
	}
	pty.WriteString("whelcaaaaxxxxxxxxxxxxxxxxxxxxxxxxxx")
	buf := make([]byte, 1024)
	var nr int
	nr, err = pty.Read(buf)
	if err != nil {
		close(execCmd,pty,channel)
	}

	log.Infof("buf is %d %s \n",nr,string(buf))

	pools.Pool.Submit(func() {
		reaPtyWriteChannel(execCmd,pty,channel)
	})
	pools.Pool.Submit(func() {
		readChannelWritePty(execCmd,pty,channel)
	})

	// Sessions have out-of-band requests such as "shell", "pty-req" and "env"
	pools.Pool.Submit(func() {
		sessionRequest(requests,pty)
	})
}

func sessionRequest(requests <- chan *ssh.Request,pty *os.File){
	for req := range requests {
		switch req.Type {
		case "shell":
			// We only accept the default shell
			// (i.e. no command in the Payload)
			if len(req.Payload) == 0 {
				req.Reply(true, nil)
			}
			log.Info(1)
		case "pty-req":
			termLen := req.Payload[3]
			w, h := parseDims(req.Payload[termLen+4:])
			SetWinsize(pty.Fd(), w, h)
			// Responding true (OK) here will let the client
			// know we have a pty ready for input
			req.Reply(true, nil)
			log.Info(2)
		case "window-change":
			w, h := parseDims(req.Payload)
			SetWinsize(pty.Fd(), w, h)
			log.Info(3)
		}
	}
}

func close(execCmd *exec.Cmd,pty *os.File,connection ssh.Channel){
	var once sync.Once
	once.Do(func() {
		connection.Close()
		if pty != nil{
			pty.Close()
		}
		_, err := execCmd.Process.Wait()
		if err != nil {
			log.Error("Failed to exit bash", err)
		}
		log.Info("Session closed")
	})
}
func reaPtyWriteChannel(execCmd *exec.Cmd,pty *os.File,channel ssh.Channel){
	for{
		bufpty := make([]byte, 1024)
		_,err := pty.Read(bufpty)
		if err != nil{
			log.Error("Read pty error by",err)
			close(execCmd,pty,channel)
		}
		cmd := string(bufpty)
		log.Info("Pty buf is",cmd)
		_,err = channel.Write(bufpty)
		if err != nil{
			log.Error("Write pty error by",err)
			close(execCmd,pty,channel)
		}
	}
}
func readChannelWritePty(execCmd *exec.Cmd,pty *os.File,channel ssh.Channel){
	for{
		bufcon := make([]byte,1024)
		_,err := channel.Read(bufcon)
		if err != nil{
			log.Error("Read connect error by",err)
			close(execCmd,pty,channel)
		}
		cmd := string(bufcon)
		log.Info("Connect buf is",cmd)
		_,err = pty.Write(bufcon)
		if err != nil{
			log.Error("Write pty error by",err)
			close(execCmd,pty,channel)
		}
	}
}

// =======================

// parseDims extracts terminal dimensions (width x height) from the provided buffer.
func parseDims(b []byte) (uint32, uint32) {
	w := binary.BigEndian.Uint32(b)
	h := binary.BigEndian.Uint32(b[4:])
	return w, h
}

// ======================

// Winsize stores the Height and Width of a terminal.
type Winsize struct {
	Height uint16
	Width  uint16
	x      uint16 // unused
	y      uint16 // unused
}

// SetWinsize sets the size of the given pty.
func SetWinsize(fd uintptr, w, h uint32) {
	ws := &Winsize{Width: uint16(w), Height: uint16(h)}
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(ws)))
}
