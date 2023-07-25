package clusterinfo_serializer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/openshift/origin/pkg/monitor"
	"github.com/openshift/origin/pkg/synthetictests/platformidentification"

	"github.com/openshift/origin/pkg/invariants"
	"github.com/openshift/origin/pkg/monitor/monitorapi"
	"github.com/openshift/origin/pkg/test/ginkgo/junitapi"
	"k8s.io/client-go/rest"
)

type clusterInfoSerializer struct {
	adminRESTConfig *rest.Config
}

func NewClusterInfoSerializer() invariants.InvariantTest {
	return &clusterInfoSerializer{}
}

func (w *clusterInfoSerializer) StartCollection(ctx context.Context, adminRESTConfig *rest.Config, recorder monitorapi.RecorderWriter) error {
	w.adminRESTConfig = adminRESTConfig
	return nil
}

func (w *clusterInfoSerializer) CollectData(ctx context.Context, beginning, end time.Time) (monitorapi.Intervals, []*junitapi.JUnitTestCase, error) {
	// because we are sharing a recorder that we're streaming into, we don't need to have a separate data collection step.
	return nil, nil, nil
}

func (*clusterInfoSerializer) ConstructComputedIntervals(ctx context.Context, startingIntervals monitorapi.Intervals, recordedResources monitorapi.ResourcesMap, beginning, end time.Time) (monitorapi.Intervals, error) {
	return nil, nil
}

func (*clusterInfoSerializer) EvaluateTestsFromConstructedIntervals(ctx context.Context, finalIntervals monitorapi.Intervals) ([]*junitapi.JUnitTestCase, error) {
	return nil, nil
}

func (w *clusterInfoSerializer) WriteContentToStorage(ctx context.Context, storageDir, timeSuffix string, finalIntervals monitorapi.Intervals, finalResourceState monitorapi.ResourcesMap) error {
	return writeClusterData(
		filepath.Join(storageDir, fmt.Sprintf("cluster-data%s.json", timeSuffix)),
		w.collectClusterData(monitor.WasMasterNodeUpdated(finalIntervals)),
	)
}

func (*clusterInfoSerializer) Cleanup(ctx context.Context) error {
	// TODO wire up the start to a context we can kill here
	return nil
}

func writeClusterData(filename string, clusterData platformidentification.ClusterData) error {
	jsonContent, err := json.MarshalIndent(clusterData, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, jsonContent, 0644)
}

func (w *clusterInfoSerializer) collectClusterData(masterNodeUpdated string) platformidentification.ClusterData {
	return monitor.CollectClusterData(w.adminRESTConfig, masterNodeUpdated)
}
