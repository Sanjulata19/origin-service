package static_host

import (
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"testing"
)

func Test_configUnmarshalJSON(t *testing.T) {
	g := goblin.Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("config.UnmarshalJSON", func() {
	})
}

func Test_configValidate(t *testing.T) {
	g := goblin.Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
}
