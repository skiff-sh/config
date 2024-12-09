package addrnet

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

var _ fmt.Stringer = Proto(0)

// Proto encompasses both network protocols and schemes.
type Proto byte

func ParseProto(str string) Proto {
	switch strings.TrimSpace(strings.ToLower(str)) {
	case "udp":
		return ProtoUDP
	case "tcp":
		return ProtoTCP
	case "http":
		return ProtoHTTP
	case "https":
		return ProtoHTTPS
	case "k8spf":
		return ProtoK8SPF
	}
	return ProtoUnknown
}

func (p Proto) IsScheme() bool {
	return p != ProtoTCP && p != ProtoUDP
}

// IANANumber returns the official number according to the IANA ->
// https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml. 255 is returned if proto is unknown.
func (p Proto) IANANumber() uint8 {
	switch p {
	case ProtoTCP:
		return 6
	case ProtoUDP:
		return 17
	default:
		return 255
	}

}

func (p Proto) String() string {
	switch p {
	case ProtoUDP:
		return "udp"
	case ProtoTCP:
		return "tcp"
	case ProtoHTTP:
		return "http"
	case ProtoHTTPS:
		return "https"
	case ProtoK8SPF:
		return "k8spf"
	default:
		return ""
	}
}

const (
	ProtoUnknown Proto = iota
	ProtoTCP
	ProtoUDP
	ProtoHTTP
	ProtoHTTPS
	ProtoK8SPF
)

type Host []byte

func (h Host) String() string {
	i, conv := h.asIP()
	if i == nil || conv && string(h) != "::" {
		return string(h)
	}
	return i.String()
}

func (h Host) asIP() (net.IP, bool) {
	str := string(h)
	if str == "::" || str == "localhost" {
		return net.IP{127, 0, 0, 1}, true
	}
	if str == ":" {
		return net.IP{0, 0, 0, 0}, true
	}
	i := net.IP(h)
	v4 := i.To4()
	if v4 != nil {
		return v4, false
	}
	return i.To16(), false
}

func (h Host) AsIP() net.IP {
	i, _ := h.asIP()
	if i != nil {
		return i
	}
	addrs, _ := net.DefaultResolver.LookupIPAddr(context.Background(), string(h))
	if len(addrs) == 0 {
		return nil
	}
	return addrs[len(addrs)-1].IP.To4()
}

func (h Host) IsIP() bool {
	i, conv := h.asIP()
	return string(h) == "::" || (i != nil && !conv)
}

var _ net.Addr = Addr("")

type Addr string

func NewTCPAddr(host string, port uint16) Addr {
	return NewAddr(ProtoTCP, []byte(host), port)
}

func NewAddr(proto Proto, host []byte, port uint16) Addr {
	builder := strings.Builder{}
	if proto != ProtoUnknown {
		builder.WriteString(proto.String())
		builder.WriteString("://")
	}

	builder.WriteString(Host(host).String())

	if !proto.IsScheme() || (proto.IsScheme() && port != 0) {
		builder.WriteString(":")
		builder.WriteString(strconv.Itoa(int(port)))
	}

	return Addr(builder.String())
}

func (a Addr) Network() string {
	proto := a.Proto()
	switch proto {
	case ProtoHTTP, ProtoHTTPS, ProtoK8SPF:
		return ProtoTCP.String()
	default:
		return a.Proto().String()
	}
}

func (a Addr) IsSocket() bool {
	return filepath.IsAbs(string(a))
}

// Key returns the Addr as a string to be used as a key.
func (a Addr) Key() string {
	return string(a)
}

// String returns an addressable string to be used in things like net.Dial.
func (a Addr) String() string {
	if a.IsSocket() {
		return string(a)
	}

	proto, host, port := a.Split()
	if host.String() == ":" {
		host = Host("")
	}

	if proto.IsScheme() {
		return string(a)
	}

	builder := strings.Builder{}

	builder.WriteString(host.String())
	builder.WriteString(":" + strconv.Itoa(int(port)))

	return builder.String()
}

func (a Addr) Split() (proto Proto, host Host, port uint16) {
	schemeIdx := strings.Index(string(a), "://")
	parse := a
	if schemeIdx >= 0 {
		proto = ParseProto(string(parse[:schemeIdx]))
		parse = parse[schemeIdx+3:]
	}

	if strings.HasPrefix(string(parse), "::") {
		host = Host("localhost")
		if len(parse) > 2 {
			num, _ := strconv.Atoi(string(parse[2:]))
			port = uint16(num)
		}
	} else {
		h, p, err := net.SplitHostPort(string(parse))
		if err != nil {
			return proto, Host(parse), 0
		}
		if h == "" {
			h = ":"
		}

		host = Host(h)
		num, _ := strconv.Atoi(p)
		port = uint16(num)
	}

	return
}

func (a Addr) Host() Host {
	_, host, _ := a.Split()
	return host
}

func (a Addr) Port() uint16 {
	_, _, port := a.Split()
	return port
}

func (a Addr) Proto() Proto {
	p, _, _ := a.Split()
	return p
}
