package creds_test

import (
	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/atc/creds"
	"github.com/pf-qiu/concourse/v6/vars"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Evaluate", func() {
	var source creds.Source

	BeforeEach(func() {
		variables := vars.StaticVariables{
			"some-param": "lol",
		}
		source = creds.NewSource(variables, atc.Source{
			"some": map[string]interface{}{
				"source-key": "((some-param))",
			},
		})
	})

	Describe("Evaluate", func() {
		It("parses variables", func() {
			result, err := source.Evaluate()
			Expect(err).NotTo(HaveOccurred())

			Expect(result).To(Equal(atc.Source{
				"some": map[string]interface{}{
					"source-key": "lol",
				},
			}))
		})
	})
})
