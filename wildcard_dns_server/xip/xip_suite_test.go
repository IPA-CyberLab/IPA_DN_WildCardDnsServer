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
