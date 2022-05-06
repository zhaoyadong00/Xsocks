package socks

import (
	"Xsocks-core/util/logs"
	"bufio"
	"errors"
	"go.uber.org/zap"
	"io"
	"net"
	"time"
)

type Client struct {
	Server   string
	UserName string
	Password string
	TCPConn       *net.TCPConn
	UDPConn       *net.UDPConn
	RemoteAddress net.Addr
	TCPTimeout    int
	UDPTimeout    int
}

func NewClient(addr , user , pass string) (*Client, error) {
	return &Client{Server: addr , UserName: user , Password: pass}, nil
}

func (c *Client)Dail(network , dst string) (net.Conn , error) {
	return c.DailWithLocalAddr(network , "" , dst)
}

func (c *Client)DailWithLocalAddr(network , src , dst string) (net.Conn , error) {
	c = &Client{
		Server:              c.Server,
		UserName:            c.UserName,
		Password:            c.Password,
		TCPTimeout:          c.TCPTimeout,
		UDPTimeout:          c.UDPTimeout,
	}
	var err error
	if network == "tcp" {
		if c.RemoteAddress == nil {
			c.RemoteAddress, err = net.ResolveTCPAddr("tcp", dst)
			if err != nil {
				return nil, err
			}
		}
		var la *net.TCPAddr
		if src != "" {
			la, err = net.ResolveTCPAddr("tcp", src)
			if err != nil {
				return nil, err
			}
		}
		err = c.handshake(la)
		if err != nil {
			logs.Logger.Info("handshake error" , zap.Error(err))
			return nil, err
		}
		return c , nil
	}else{
		return nil, errors.New("Temporary does not support")
	}
}

func (c *Client)handshake(laddr *net.TCPAddr) error {
	raddr, err := net.ResolveTCPAddr("tcp", c.Server)
	if err != nil {
		logs.Logger.Info("handshake ResolveTCPAddr" , zap.Error(err))
		return err
	}
	c.TCPConn, err = net.DialTCP("tcp", laddr, raddr)
	if err != nil {
		logs.Logger.Info("handshake DialTCP" , zap.Error(err))
		return err
	}
	if c.TCPTimeout != 0 {
		if err := c.TCPConn.SetDeadline(time.Now().Add(time.Duration(c.TCPTimeout) * time.Second)); err != nil {
			return err
		}
	}
	//认证协商
	conReq := NewConsult()
	m := MethodNone
	if c.UserName != "" && c.Password != "" {
		m = MethodUsernamePassword
	}
	req , _ := conReq.BuildConsultRequest(Ver , []byte{m})
	if _ , err = c.TCPConn.Write(req); err != nil {//发送协商数据
		logs.Logger.Info("TCPConn.Write err:" , zap.Error(err))
		return err
	}
	//接收协商数据
	buf := make([]byte , 0)
	_ , err = io.ReadFull(c.TCPConn , buf)
	if err != nil {
		logs.Logger.Info("handshake ReadFull err:" , zap.Error(err))
		return err
	}
	if buf[0] != Ver {
		logs.Logger.Info("handshake consult Ver err:" , zap.Error(err))
		return errors.New("handshake version error")
	}
	if buf[1] != m {
		return errors.New("Unsupport method")
	}

	if m == MethodUsernamePassword {//验证用户
		auth := NewAuthPassword()
		_ , err = c.TCPConn.Write(auth.BuildRequest(c.UserName , c.Password))
		if err != nil {
			logs.Logger.Info("MethodUsernamePassword Write err" , zap.Error(err))
			return err
		}
		authReply := make([]byte , 2)
		_ , err = io.ReadFull(bufio.NewReader(c.TCPConn) , authReply)
		if err != nil {
			logs.Logger.Info("auth reply reader err" , zap.Error(err))
			return err
		}
		reply , err := NewSocks5AuthUPasswdReply(authReply)
		if err != nil {
			logs.Logger.Info("AuthPasswordReply err" , zap.Error(err))
			return err
		}
		if reply.Status != UserPassStatusSuccess {
			return errors.New("Invalid Username or Password for Auth")
		}
	}
	return nil
}



func (c *Client) Read(b []byte) (int, error) {
	if c.UDPConn == nil {
		return c.TCPConn.Read(b)
	}
	n, err := c.UDPConn.Read(b)
	if err != nil {
		return 0, err
	}
	//d, err := NewDatagramFromBytes(b[0:n])
	//if err != nil {
	//	return 0, err
	//}
	//n = copy(b, d.Data)
	return n, nil
}

func (c *Client) Write(b []byte) (int, error) {
	if c.UDPConn == nil {
		return c.TCPConn.Write(b)
	}
	//a, h, p, err := ParseAddress(c.RemoteAddress.String())
	//if err != nil {
	//	return 0, err
	//}
	//if a == ATYPDomain {
	//	h = h[1:]
	//}
	////d := NewDatagram(a, h, p, b)
	//b1 := d.Bytes()
	//n, err := c.UDPConn.Write(b1)
	//if err != nil {
	//	return 0, err
	//}
	//if len(b1) != n {
	//	return 0, errors.New("not write full")
	//}
	return len(b), nil
}

func (c *Client) Close() error {
	if c.UDPConn == nil {
		return c.TCPConn.Close()
	}
	if c.TCPConn != nil {
		c.TCPConn.Close()
	}
	return c.UDPConn.Close()
}

func (c *Client) LocalAddr() net.Addr {
	if c.UDPConn == nil {
		return c.TCPConn.LocalAddr()
	}
	return c.UDPConn.LocalAddr()
}

func (c *Client) RemoteAddr() net.Addr {
	return c.RemoteAddress
}

func (c *Client) SetDeadline(t time.Time) error {
	if c.UDPConn == nil {
		return c.TCPConn.SetDeadline(t)
	}
	return c.UDPConn.SetDeadline(t)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	if c.UDPConn == nil {
		return c.TCPConn.SetReadDeadline(t)
	}
	return c.UDPConn.SetReadDeadline(t)
}

func (c *Client) SetWriteDeadline(t time.Time) error {
	if c.UDPConn == nil {
		return c.TCPConn.SetWriteDeadline(t)
	}
	return c.UDPConn.SetWriteDeadline(t)
}

func (c *Client) Process() error {
	//conn , err := net.DialTCP("tcp" , "" , "")
	//if err != nil {
	//	logs.Logger.Info("DialTCP err:" , zap.Error(err))
	//	return err
	//}
	return nil
}
