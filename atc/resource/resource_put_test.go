package resource_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/atc/resource"
	"github.com/pf-qiu/concourse/v6/atc/runtime"
	"github.com/pf-qiu/concourse/v6/atc/runtime/runtimefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Resource Put", func() {
	var (
		ctx             context.Context
		someProcessSpec runtime.ProcessSpec
		fakeRunnable    runtimefakes.FakeRunner

		putVersionResult runtime.VersionResult

		source  atc.Source
		params  atc.Params
		version atc.Version

		resource resource.Resource

		putErr error
	)

	BeforeEach(func() {
		ctx = context.Background()

		source = atc.Source{"some": "source"}
		version = atc.Version{"some": "version"}
		params = atc.Params{"some": "params"}

		someProcessSpec.Path = "some/fake/path"
		someProcessSpec.Args = []string{"some/foo-dir"}
		someProcessSpec.StderrWriter = gbytes.NewBuffer()

		resource = resourceFactory.NewResource(source, params, version)

	})

	JustBeforeEach(func() {
		putVersionResult, putErr = resource.Put(ctx, someProcessSpec, &fakeRunnable)
	})

	Context("when Runnable -> RunScript succeeds and returns a Version", func() {
		BeforeEach(func() {
			fakeRunnable.RunScriptStub = func(i context.Context, s string, strings []string, bytes []byte, versionResult interface{}, writer io.Writer, b bool) error {
				err := json.Unmarshal([]byte(`{"version": {"ref":"v1"}}`), &versionResult)
				if err != nil {
					return err
				}

				return nil
			}
		})

		It("Invokes Runnable -> RunScript with the correct arguments", func() {
			actualCtx, actualSpecPath, actualArgs, actualInput,
				actualVersionResultRef, actualSpecStdErrWriter,
				actualRecoverableBool := fakeRunnable.RunScriptArgsForCall(0)

			signature, err := resource.Signature()
			Expect(err).ToNot(HaveOccurred())

			Expect(actualCtx).To(Equal(ctx))
			Expect(actualSpecPath).To(Equal(someProcessSpec.Path))
			Expect(actualArgs).To(Equal(someProcessSpec.Args))
			Expect(actualInput).To(Equal(signature))
			Expect(actualVersionResultRef).To(Equal(&putVersionResult))
			Expect(actualSpecStdErrWriter).To(Equal(someProcessSpec.StderrWriter))
			Expect(actualRecoverableBool).To(BeTrue())
		})

		It("doesnt return an error", func() {
			Expect(putErr).To(BeNil())
		})
	})
	Context("when Runnable -> RunScript succeeds and does NOT return a Version", func() {
		BeforeEach(func() {
			fakeRunnable.RunScriptReturns(nil)
		})
		It("returns a corresponding error", func() {
			Expect(putErr).To(MatchError("resource script (" + someProcessSpec.Path + " " + strings.Join(someProcessSpec.Args, " ") + ") output a null version"))
		})
	})

	Context("when Runnable -> RunScript returns an error", func() {
		var disasterErr = errors.New("there was an issue")
		BeforeEach(func() {
			fakeRunnable.RunScriptReturns(disasterErr)
		})
		It("returns the error", func() {
			Expect(putErr).To(Equal(disasterErr))
		})
	})

})
