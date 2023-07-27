package invariants

import (
	"context"
	"time"

	"k8s.io/client-go/rest"

	"github.com/openshift/origin/pkg/monitor/monitorapi"
	"github.com/openshift/origin/pkg/test/ginkgo/junitapi"
)

type InvariantTest interface {
	// StartCollection is responsible for setting up all resources required for collection of data on the cluster.
	// An error will not stop execution, but will cause a junit failure that will cause the job run to fail.
	// This allows us to know when setups fail.
	StartCollection(ctx context.Context, adminRESTConfig *rest.Config, recorder monitorapi.RecorderWriter) error

	// CollectData will only be called once near the end of execution, before all Intervals are inspected.
	// Errors reported will be indicated as junit test failure and will cause job runs to fail.
	// storageDir is for gathering data only, not for writing in this stage.  To store data, use WriteContentToStorage
	CollectData(ctx context.Context, storageDir string, beginning, end time.Time) (monitorapi.Intervals, []*junitapi.JUnitTestCase, error)

	// ConstructComputedIntervals is called after all InvariantTests have produced raw Intervals.
	// Order of ConstructComputedIntervals across different InvariantTests is not guaranteed.
	// Return *only* the constructed intervals.
	// Errors reported will be indicated as junit test failure and will cause job runs to fail.
	ConstructComputedIntervals(ctx context.Context, startingIntervals monitorapi.Intervals, recordedResources monitorapi.ResourcesMap, beginning, end time.Time) (constructedIntervals monitorapi.Intervals, err error)

	// EvaluateTestsFromConstructedIntervals is called after all Intervals are known and can produce
	// junit tests for reporting purposes.
	// Errors reported will be indicated as junit test failure and will cause job runs to fail.
	EvaluateTestsFromConstructedIntervals(ctx context.Context, finalIntervals monitorapi.Intervals) ([]*junitapi.JUnitTestCase, error)

	// WriteContentToStorage writes content to the storage directory that is collected by openshift CI.
	// Do not write.
	// 1. junits.  Those should be returned from EvaluateTestsFromConstructedIntervals
	// 2. intervals.  Those should be returned from CollectData and ConstructComputedIntervals
	// 3. tracked resources.  Those are written by some default invariantTests.
	// You *may* choose to store state in CollectData that you later persist via this method. An example might be
	// code that scans audit logs and reports summaries of top actors.
	// Errors reported will be indicated as junit test failure and will cause job runs to fail.
	WriteContentToStorage(ctx context.Context, storageDir, timeSuffix string, finalIntervals monitorapi.Intervals, finalResourceState monitorapi.ResourcesMap) error

	// Cleanup must be idempotent and it may be called multiple times in any scenario.  Multiple defers, multi-registered
	// abort handlers, abort handler running concurrent to planned shutdown.  Make your cleanup callable multiple times.
	// Errors reported will cause job runs to fail to ensure cleanup functions work reliably.
	Cleanup(ctx context.Context) error
}

type InvariantRegistry interface {
	AddRegistryOrDie(registry InvariantRegistry)

	// AddInvariant adds an invariant test with a particular name, the name will be used to create a testsuite.
	// The jira component will be forced into every JunitTestCase.
	AddInvariant(name, jiraComponent string, invariantTest InvariantTest) error

	AddInvariantOrDie(name, jiraComponent string, invariantTest InvariantTest)

	// StartCollection is responsible for setting up all resources required for collection of data on the cluster.
	// An error will not stop execution, but will cause a junit failure that will cause the job run to fail.
	// This allows us to know when setups fail.
	StartCollection(ctx context.Context, adminRESTConfig *rest.Config, recorder monitorapi.RecorderWriter) ([]*junitapi.JUnitTestCase, error)

	// CollectData will only be called once near the end of execution, before all Intervals are inspected.
	// Errors reported will be indicated as junit test failure and will cause job runs to fail.
	CollectData(ctx context.Context, storageDir string, beginning, end time.Time) (monitorapi.Intervals, []*junitapi.JUnitTestCase, error)

	// ConstructComputedIntervals is called after all InvariantTests have produced raw Intervals.
	// Order of ConstructComputedIntervals across different InvariantTests is not guaranteed.
	// Return *only* the constructed intervals.
	// Errors reported will be indicated as junit test failure and will cause job runs to fail.
	ConstructComputedIntervals(ctx context.Context, startingIntervals monitorapi.Intervals, recordedResources monitorapi.ResourcesMap, beginning, end time.Time) (monitorapi.Intervals, []*junitapi.JUnitTestCase, error)

	// EvaluateTestsFromConstructedIntervals is called after all Intervals are known and can produce
	// junit tests for reporting purposes.
	// Errors reported will be indicated as junit test failure and will cause job runs to fail.
	EvaluateTestsFromConstructedIntervals(ctx context.Context, finalIntervals monitorapi.Intervals) ([]*junitapi.JUnitTestCase, error)

	// WriteContentToStorage writes content to the storage directory that is collected by openshift CI.
	// Do not write.
	// 1. junits.  Those should be returned from EvaluateTestsFromConstructedIntervals
	// 2. intervals.  Those should be returned from CollectData and ConstructComputedIntervals
	// 3. tracked resources.  Those are written by some default invariantTests.
	// You *may* choose to store state in CollectData that you later persist via this method. An example might be
	// code that scans audit logs and reports summaries of top actors.
	WriteContentToStorage(ctx context.Context, storageDir, timeSuffix string, finalIntervals monitorapi.Intervals, finalResourceState monitorapi.ResourcesMap) ([]*junitapi.JUnitTestCase, error)

	// Cleanup must be idempotent and it may be called multiple times in any scenario.  Multiple defers, multi-registered
	// abort handlers, abort handler running concurrent to planned shutdown.  Make your cleanup callable multiple times.
	// Errors reported will cause job runs to fail to ensure cleanup functions work reliably.
	Cleanup(ctx context.Context) ([]*junitapi.JUnitTestCase, error)

	getInvariantTests() map[string]*invariantItem
}
