package fakedns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/uuid"
)

const (
	TestDomain  = "fake-dns-test.v2fly.org"
	TestDomain2 = "fake-dns-test.v2fly.org"
)

func TestNewFakeDnsHolder(_ *testing.T) {
	_, err := NewFakeDNSHolder()
	common.Must(err)
}

func TestFakeDnsHolderCreateMapping(t *testing.T) {
	fkdns, err := NewFakeDNSHolder()
	common.Must(err)

	addr := fkdns.GetFakeIPForDomain(TestDomain)
	assert.Equal(t, "240.0.0.0", addr[0].IP().String())
}

func TestFakeDnsHolderCreateMappingMany(t *testing.T) {
	fkdns, err := NewFakeDNSHolder()
	common.Must(err)

	addr := fkdns.GetFakeIPForDomain(TestDomain)
	assert.Equal(t, "240.0.0.0", addr[0].IP().String())

	addr2 := fkdns.GetFakeIPForDomain(TestDomain2)
	assert.Equal(t, "240.0.0.1", addr2[0].IP().String())
}

func TestFakeDnsHolderCreateMappingManyAndResolve(t *testing.T) {
	fkdns, err := NewFakeDNSHolder()
	common.Must(err)

	{
		addr := fkdns.GetFakeIPForDomain(TestDomain)
		assert.Equal(t, "240.0.0.0", addr[0].IP().String())
	}

	{
		addr2 := fkdns.GetFakeIPForDomain(TestDomain2)
		assert.Equal(t, "240.0.0.1", addr2[0].IP().String())
	}

	{
		result := fkdns.GetDomainFromFakeDNS(net.ParseAddress("240.0.0.0"))
		assert.Equal(t, "fakednstest.v2fly.org", result)
	}

	{
		result := fkdns.GetDomainFromFakeDNS(net.ParseAddress("240.0.0.1"))
		assert.Equal(t, "fakednstest2.v2fly.org", result)
	}
}

func TestFakeDnsHolderCreateMappingManySingleDomain(t *testing.T) {
	fkdns, err := NewFakeDNSHolder()
	common.Must(err)

	addr := fkdns.GetFakeIPForDomain(TestDomain)
	assert.Equal(t, "240.0.0.0", addr[0].IP().String())

	addr2 := fkdns.GetFakeIPForDomain(TestDomain2)
	assert.Equal(t, "240.0.0.0", addr2[0].IP().String())
}

func TestFakeDnsHolderCreateMappingAndRollOver(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping DNS Holder RollOver test in short mode. ~190s")
	}

	fkdns, err := NewFakeDNSHolder()
	common.Must(err)

	{
		addr := fkdns.GetFakeIPForDomain(TestDomain)
		assert.Equal(t, "240.0.0.0", addr[0].IP().String())
	}

	{
		addr2 := fkdns.GetFakeIPForDomain(TestDomain2)
		assert.Equal(t, "240.0.0.1", addr2[0].IP().String())
	}

	for i := 0; i <= 33554432; i++ {
		{
			result := fkdns.GetDomainFromFakeDNS(net.ParseAddress("240.0.0.0"))
			assert.Equal(t, TestDomain, result)
		}

		{
			result := fkdns.GetDomainFromFakeDNS(net.ParseAddress("240.0.0.1"))
			assert.Equal(t, TestDomain2, result)
		}

		{
			u := uuid.New()
			domain := u.String() + ".fake-dns-test.v2fly.org"
			addr := fkdns.GetFakeIPForDomain(domain)
			resultAddr := addr[0].IP().String()

			result := fkdns.GetDomainFromFakeDNS(net.ParseAddress(resultAddr))
			assert.Equal(t, domain, result)
		}
	}
}
