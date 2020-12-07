package flaghelpers_test

import (
	"github.com/pf-qiu/concourse/v6/atc"
	. "github.com/pf-qiu/concourse/v6/fly/commands/internal/flaghelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobFlag", func() {
	var (
		flag *JobFlag
	)

	BeforeEach(func() {
		flag = &JobFlag{}
	})

	Describe("UnmarshalFlag", func() {
		for _, tt := range []struct {
			desc        string
			flag        string
			pipelineRef atc.PipelineRef
			jobName     string
			err         string
		}{
			{
				desc:        "basic",
				flag:        "some-pipeline/some-job",
				pipelineRef: atc.PipelineRef{Name: "some-pipeline"},
				jobName:     "some-job",
			},
			{
				desc: "instance vars",
				flag: "some-pipeline/branch:master,foo.bar:baz/some-job",
				pipelineRef: atc.PipelineRef{
					Name:         "some-pipeline",
					InstanceVars: atc.InstanceVars{"branch": "master", "foo": map[string]interface{}{"bar": "baz"}},
				},
				jobName: "some-job",
			},
			{
				desc: "instance var with a '/'",
				flag: "some-pipeline/branch:feature/do_things,foo:bar/some-job",
				pipelineRef: atc.PipelineRef{
					Name:         "some-pipeline",
					InstanceVars: atc.InstanceVars{"branch": "feature/do_things", "foo": "bar"},
				},
				jobName: "some-job",
			},
			{
				desc: "instance var with special chars",
				flag: `some-pipeline/foo."bar.baz":'abc,def:ghi'/some-job`,
				pipelineRef: atc.PipelineRef{
					Name: "some-pipeline",
					InstanceVars: atc.InstanceVars{
						"foo": map[string]interface{}{
							"bar.baz": "abc,def:ghi",
						},
					},
				},
				jobName: "some-job",
			},
			{
				desc: "only pipeline specified",
				flag: "some-pipeline",
				err:  "argument format should be <pipeline>/<job>",
			},
			{
				desc: "job name not specified",
				flag: "some-pipeline/",
				err:  "argument format should be <pipeline>/<job>",
			},
			{
				desc: "malformed instance var",
				flag: "some-pipeline/branch=master/some-job",
				err:  "argument format should be <pipeline>/<key:value>/<job>",
			},
		} {
			tt := tt
			It(tt.desc, func() {
				err := flag.UnmarshalFlag(tt.flag)
				if tt.err == "" {
					Expect(err).ToNot(HaveOccurred())
					Expect(flag.PipelineRef).To(Equal(tt.pipelineRef))
					Expect(flag.JobName).To(Equal(tt.jobName))
				} else {
					Expect(err).To(MatchError(tt.err))
				}
			})
		}
	})
})
