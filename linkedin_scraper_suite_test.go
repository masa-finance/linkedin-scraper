package linkedinscraper_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLinkedinScraper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LinkedinScraper Suite")
}
