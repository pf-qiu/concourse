package algorithm_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"

	"github.com/pf-qiu/concourse/v6/atc/db"
	"github.com/pf-qiu/concourse/v6/atc/db/lock"
	"github.com/pf-qiu/concourse/v6/atc/metric"
	"github.com/pf-qiu/concourse/v6/atc/postgresrunner"
	"github.com/pf-qiu/concourse/v6/tracing"
)

var (
	postgresRunner postgresrunner.Runner
	dbProcess      ifrit.Process

	lockFactory lock.LockFactory
	teamFactory db.TeamFactory

	dbConn db.Conn

	exporter *jaeger.Exporter
)

func TestAlgorithm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Algorithm Suite")
}

var _ = BeforeSuite(func() {
	jaegerURL := os.Getenv("JAEGER_URL")

	if jaegerURL != "" {
		c := tracing.Config{
			Jaeger: tracing.Jaeger{
				Endpoint: jaegerURL + "/api/traces",
				Service:  "algorithm_test",
			},
		}

		err := c.Prepare()
		Expect(err).ToNot(HaveOccurred())
	}

	postgresRunner = postgresrunner.Runner{
		Port: 5433 + GinkgoParallelNode(),
	}

	dbProcess = ifrit.Invoke(postgresRunner)

	postgresRunner.CreateTestDB()
})

var _ = BeforeEach(func() {
	postgresRunner.Truncate()

	dbConn = postgresRunner.OpenConn()

	lockFactory = lock.NewLockFactory(postgresRunner.OpenSingleton(), metric.LogLockAcquired, metric.LogLockReleased)
	teamFactory = db.NewTeamFactory(dbConn, lockFactory)
})

var _ = AfterEach(func() {
	err := dbConn.Close()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	dbProcess.Signal(os.Interrupt)
	<-dbProcess.Wait()

	if exporter != nil {
		exporter.Flush()
	}
})
