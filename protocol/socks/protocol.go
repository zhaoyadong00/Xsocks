package socks

import (
	"encoding/binary"
	"errors"
	"net"
)

const (
	// Ver is socks protocol version
	Ver byte = 0x05

	// MethodNone is none method
	MethodNone byte = 0x00
	// MethodGSSAPI is gssapi method
	MethodGSSAPI byte = 0x01 // MUST support // todo
	// MethodUsernamePassword is username/assword auth method
	MethodUsernamePassword byte = 0x02 // SHOULD support
	// MethodUnsupportAll means unsupport all given methods
	MethodUnsupportAll byte = 0xFF

	// UserPassVer is username/password auth protocol version
	UserPassVer byte = 0x01
	// UserPassStatusSuccess is success status of username/password auth
	UserPassStatusSuccess byte = 0x00
	// UserPassStatusFailure is failure status of username/password auth
	UserPassStatusFailure byte = 0x01 // just other than 0x00

	// CmdConnect is connect command
	CmdConnect byte = 0x01
	// CmdBind is bind command
	CmdBind byte = 0x02
	// CmdUDP is UDP command
	CmdUDP byte = 0x03

	// ATYPIPv4 is ipv4 address type
	ATYPIPv4 byte = 0x01 // 4 octets
	// ATYPDomain is domain address type
	ATYPDomain byte = 0x03 // The first octet of the address field contains the number of octets of name that follow, there is no terminating NUL octet.
	// ATYPIPv6 is ipv6 address type
	ATYPIPv6 byte = 0x04 // 16 octets

	// RepSuccess means that success for repling
	RepSuccess byte = 0x00
	// RepServerFailure means the server failure
	RepServerFailure byte = 0x01
	// RepNotAllowed means the request not allowed
	RepNotAllowed byte = 0x02
	// RepNetworkUnreachable means the network unreachable
	RepNetworkUnreachable byte = 0x03
	// RepHostUnreachable means the host unreachable
	RepHostUnreachable byte = 0x04
	// RepConnectionRefused means the connection refused
	RepConnectionRefused byte = 0x05
	// RepTTLExpired means the TTL expired
	RepTTLExpired byte = 0x06
	// RepCommandNotSupported means the request command not supported
	RepCommandNotSupported byte = 0x07
	// RepAddressNotSupported means the request address not supported
	RepAddressNotSupported byte = 0x08
)

//协商请求
type ConsultRequest struct {
	VER      uint8
	NMETHODS uint8
	METHODS  []uint8
}


func NewConsult() *ConsultRequest {
	return &ConsultRequest{}
}
//生成 协商数据包
func (c *ConsultRequest)BuildConsultRequest(ver uint8 , m []uint8) ([]byte , error) {
	c.VER = ver
	c.NMETHODS = byte(len(m))
	c.METHODS = m
	buf := make([]byte , 0)
	buf = append(buf, c.VER , c.NMETHODS)
	buf = append(buf, c.METHODS...)
	return buf , nil
}

//协商加密方式
//func (p *ConsultRequest)CheckConsultRequest(b []byte) (*ConsultReply , error) {
//	n := len(b)
//	if n < 3 {
//		return nil, errors.New("协议错误, sNMETHODS不对")
//	}
//	p.VER = b[0] //ReadByte reads and returns a single byte，第一个参数为socks的版本号
//	if p.VER != Ver {
//		return nil, errors.New("协议错误, version版本不为5!")
//	}
//	p.NMETHODS = b[1]
//	if n != int(2 + p.NMETHODS) {
//		return nil, errors.New("协议错误, sNMETHODS不对")
//	}
//	p.METHODS = b[2 : 2+p.NMETHODS] //读取指定长度信息，读取正好len(buf)长度的字节。如果字节数不是指定长度，则返回错误信息和正确的字节数
//
//	useMethod := byte(0x00) //默认不需要密码
//	for _, v := range p.METHODS {
//		if v == MethodUsernamePassword {
//			useMethod = MethodUsernamePassword
//		}
//	}
//	if p.VER != Ver {
//		return nil, errors.New("该协议不是socks5协议")
//	}
//
//	if useMethod != MethodUsernamePassword {
//		return nil, errors.New("协议错误, 加密方法不对")
//	}
//	//协商返回数据结构
//	reply := &ConsultReply{
//		Ver: Ver,
//		Method: useMethod,
//	}
//	return reply, nil
//}
//协商返回结构
type ConsultReply struct {
	Ver uint8
	Method uint8
}

func NewConsultReply() *ConsultReply {
	return &ConsultReply{}
}

//结构转byte
func (c *ConsultReply)ToByte() []byte {
	buf := make([]byte , 0)
	buf = append(buf, c.Ver , c.Method)
	return buf
}

//用户密码认证 数据包
type Socks5AuthUPasswdRequest struct {
	VER    uint8
	ULEN   uint8
	UNAME  string
	PLEN   uint8
	PASSWD string
}

func NewAuthPassword() *Socks5AuthUPasswdRequest {
	return &Socks5AuthUPasswdRequest{}
}

func (s *Socks5AuthUPasswdRequest)BuildRequest(user , pass string) []byte {
	s.VER = Ver
	s.ULEN = uint8(len(user))
	s.UNAME = user
	s.PLEN = uint8(len(pass))
	s.PASSWD = pass
	buf := make([]byte , 0)
	buf = append(buf , s.VER , s.ULEN )
	buf = append(buf, []byte(s.UNAME)...)
	buf = append(buf, s.ULEN)
	buf = append(buf, []byte(s.PASSWD)...)
	return buf
}

type Socks5AuthUPasswdReply struct {
	Ver    byte
	Status byte
}

func NewSocks5AuthUPasswdReply(d []byte) (*Socks5AuthUPasswdReply , error) {
	if len(d) != 2 {
		return nil , errors.New("params error")
	}
	if d[0] != UserPassVer {
		return nil , errors.New("Invalid Version of Username Password Auth")
	}
	return &Socks5AuthUPasswdReply{
		Ver: d[0],
		Status: d[1],
	}, nil
}

func (s *Socks5AuthUPasswdRequest) HandleAuth(b []byte) ([]byte, error) {
	n := len(b)

	s.VER = b[0]
	if s.VER != 5 {
		return nil, errors.New("该协议不是socks5协议")
	}

	s.ULEN = b[1]
	s.UNAME = string(b[2 : 2+s.ULEN])
	s.PLEN = b[2+s.ULEN+1]
	s.PASSWD = string(b[n-int(s.PLEN) : n])

	resp := []byte{Ver, 0x00}

	return resp, nil
}

type Socks5Resolution struct {
	VER       uint8
	CMD       uint8
	RSV       uint8
	ATYP      uint8
	DSTADDR   []byte
	DSTPORT   uint16
	DSTDOMAIN string
	RAWADDR   *net.TCPAddr
}

func (s *Socks5Resolution) LSTRequest(b []byte) ([]byte, error) {
	n := len(b)
	if n < 7 {
		return nil, errors.New("请求协议错误")
	}
	s.VER = b[0]
	if s.VER != Ver {
		return nil, errors.New("该协议不是socks5协议")
	}

	s.CMD = b[1]
	if s.CMD != 1 {
		return nil, errors.New("客户端请求类型不为代理连接, 其他功能暂时不支持.")
	}
	s.RSV = b[2] //RSV保留字端，值长度为1个字节

	s.ATYP = b[3]

	switch s.ATYP {
	case 1:
		// IP V4 address: X'01'
		s.DSTADDR = b[4 : 4+net.IPv4len]
	case 3:
		// DOMAINNAME: X'03'
		s.DSTDOMAIN = string(b[5 : n-2])
		ipAddr, err := net.ResolveIPAddr("ip", s.DSTDOMAIN)
		if err != nil {
			return nil, err
		}
		s.DSTADDR = ipAddr.IP
	case 4:
		// IP V6 address: X'04'
		s.DSTADDR = b[4 : 4+net.IPv6len]
	default:
		return nil, errors.New("IP地址错误")
	}

	s.DSTPORT = binary.BigEndian.Uint16(b[n-2 : n])
	// DSTADDR全部换成IP地址，可以防止DNS污染和封杀
	s.RAWADDR = &net.TCPAddr{
		IP:   s.DSTADDR,
		Port: int(s.DSTPORT),
	}
	//net.DialTCP("tcp" ,  , s.RAWADDR)
	/**
	  回应客户端,响应客户端连接成功
	      +----+-----+-------+------+----------+----------+
	      |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	      +----+-----+-------+------+----------+----------+
	      | 1  |  1  | X'00' |  1   | Variable |    2     |
	      +----+-----+-------+------+----------+----------+
	*/
	resp := []byte{Ver, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	// conn.Write(resp)

	return resp, nil
}

//go func() {
//	defer wg.Done()
//	defer dstServer.Close()
//	io.Copy(dstServer, client)
//}()
//
//go func() {
//	defer wg.Done()
//	defer client.Close()
//	io.Copy(client, dstServer)
//}()