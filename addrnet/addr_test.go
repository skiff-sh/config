package addrnet

import (
	"net"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AddrTestSuite struct {
	suite.Suite
}

func (p *AddrTestSuite) TestHost() {
	type test struct {
		Given      []byte
		Expected   string
		IsIP       bool
		ExpectedIP net.IP
	}

	tests := []test{
		{
			Given:      []byte("localhost"),
			Expected:   "localhost",
			ExpectedIP: net.IP{127, 0, 0, 1},
		},
		{
			Given:      []byte{127, 0, 0, 1},
			Expected:   "127.0.0.1",
			IsIP:       true,
			ExpectedIP: []byte{127, 0, 0, 1},
		},
		{
			Given:      []byte("::"),
			Expected:   "127.0.0.1",
			IsIP:       true,
			ExpectedIP: []byte{127, 0, 0, 1},
		},
		{
			Given:      []byte(":"),
			Expected:   ":",
			ExpectedIP: []byte{0, 0, 0, 0},
		},
	}

	for i, v := range tests {
		p.Equal(v.Expected, Host(v.Given).String(), "test %d given %s", i, v.Given)
		p.Equal(v.IsIP, Host(v.Given).IsIP(), "test %d given %s", i, v.Given)
		p.Equal(v.ExpectedIP, Host(v.Given).AsIP(), "test %d given %s", i, v.Given)
	}
}

func (p *AddrTestSuite) TestNewAddr() {
	type test struct {
		Given    Addr
		Expected string
	}

	tests := []test{
		{
			Given:    NewAddr(ProtoTCP, net.IP{127, 0, 0, 1}, 8080),
			Expected: "127.0.0.1:8080",
		},
		{
			Given:    NewAddr(ProtoTCP, Host("localhost"), 8080),
			Expected: "localhost:8080",
		},
		{
			Given:    NewAddr(ProtoUnknown, net.IP{127, 0, 0, 1}, 8080),
			Expected: "127.0.0.1:8080",
		},
		{
			Given:    Addr(":8080"),
			Expected: ":8080",
		},
		{
			Given:    NewAddr(ProtoTCP, nil, 8080),
			Expected: ":8080",
		},
		{
			Given:    NewAddr(ProtoHTTP, Host("localhost"), 8080),
			Expected: "http://localhost:8080",
		},
		{
			Given:    NewAddr(ProtoHTTPS, Host("localhost"), 0),
			Expected: "https://localhost",
		},
	}

	for i, v := range tests {
		p.Equal(v.Expected, v.Given.String(), "test %d given %s", i, v.Given)
	}
}

func (p *AddrTestSuite) TestSplit() {
	type test struct {
		Given         Addr
		ExpectedHost  Host
		ExpectedPort  uint16
		ExpectedProto Proto
	}

	tests := []test{
		{
			Given:        Addr("::8080"),
			ExpectedHost: Host("localhost"),
			ExpectedPort: 8080,
		},
		{
			Given:         Addr("tcp://localhost:8080"),
			ExpectedHost:  Host("localhost"),
			ExpectedPort:  8080,
			ExpectedProto: ProtoTCP,
		},
		{
			Given:        Addr("localhost:8080"),
			ExpectedHost: Host("localhost"),
			ExpectedPort: 8080,
		},
		{
			Given:         Addr("tcp://localhost"),
			ExpectedHost:  Host("localhost"),
			ExpectedProto: ProtoTCP,
		},
		{
			Given:        Addr("localhost"),
			ExpectedHost: Host("localhost"),
		},
		{
			Given:        Addr("127.0.0.1:8080"),
			ExpectedHost: Host("127.0.0.1"),
			ExpectedPort: 8080,
		},
		{
			Given:        Addr(":8080"),
			ExpectedHost: Host(":"),
			ExpectedPort: 8080,
		},
		{
			Given:         Addr("http://localhost:8080"),
			ExpectedHost:  Host("localhost"),
			ExpectedProto: ProtoHTTP,
			ExpectedPort:  8080,
		},
	}

	for i, v := range tests {
		proto, host, port := v.Given.Split()
		p.Equal(host, v.ExpectedHost, "test %d given %s", i, v.Given)
		p.Equal(port, v.ExpectedPort, "test %d given %s", i, v.Given)
		p.Equal(proto, v.ExpectedProto, "test %d given %s", i, v.Given)
	}
}

func TestAddrTestSuite(t *testing.T) {
	suite.Run(t, new(AddrTestSuite))
}
