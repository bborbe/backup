package pkg_test

import (
	"github.com/bborbe/backup/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path", func() {
	var directory pkg.Path
	BeforeEach(func() {
		directory = "/backup"
	})
	Context("Join", func() {
		var elems []string
		var result pkg.Path
		JustBeforeEach(func() {
			result = directory.Join(elems...)
		})
		Context("simple", func() {
			BeforeEach(func() {
				elems = []string{"foo"}
			})
			It("returns correct path", func() {
				Expect(result).To(Equal(pkg.Path("/backup/foo")))
			})
		})
	})
})
