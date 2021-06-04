// Forked From: https://github.com/cunnie/sslip.io/
// (Golang-based DNS server which maps DNS records with embedded IP addresses to those addresses)
// by Brian Cunnie (https://github.com/cunnie/)

package xip_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestXip(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Xip Suite")
}
