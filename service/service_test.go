package service

import (
	"context"
	"github.com/franela/goblin"
	"github.com/nullserve/origin-service/config"
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
			host := "1234.sites.example.com"
			mServer := server{
				config: &config.OriginService{
					HostSuffix: suffix,
					AppPrefix: "site",
					RefPrefix: "sites",
				},
			}

			// When
			result, err := mServer.hostToDeploymentId(context.Background(), host, suffix)

			// Then
			Expect(err).Should(BeNil())
			Expect(*result).Should(Equal("1234"))
		})

		g.It("should error if the suffix doesn't exist", func() {
			// Given
			suffix := "example.com"
			host := "1234.sites.duh.com"
			mServer := server{
				config: &config.OriginService{
					HostSuffix: suffix,
					AppPrefix: "site",
					RefPrefix: "sites",
				},
			}

			// When
			result, err := mServer.hostToDeploymentId(context.Background(), host, suffix)

			// Then
			Expect(result).Should(BeNil())
			Expect(err).Should(Equal(errNoSuchSuffix))
		})
	})
}
