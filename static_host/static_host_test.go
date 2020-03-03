package static_host

import (
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"testing"
)

func Test_hostToDeploymentId(t *testing.T) {
	g := goblin.Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("hostToDeploymentId", func() {
		g.It("should extract the prefix successfully", func() {
			// Given
			suffix := "example.com"
			host := "1234.example.com"

			// When
			result, err := hostToDeploymentId(host, suffix)

			// Then
			Expect(*result).Should(Equal("1234"))
			Expect(err).Should(BeNil())
		})

		g.It("should error if the suffix doesn't exist", func() {
			// Given
			suffix := "example.com"
			host := "1234.duh.com"

			// When
			result, err := hostToDeploymentId(host, suffix)

			// Then
			Expect(result).Should(BeNil())
			Expect(err).Should(Equal(errNoSuchSuffix))
		})
	})
}
