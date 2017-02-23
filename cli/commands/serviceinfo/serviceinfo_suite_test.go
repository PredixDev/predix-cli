package serviceinfo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServiceinfo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Serviceinfo Suite")
}
