package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/extensions"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &DisableMonitoringPoller{}

// DisableMonitoringPoller polls until the Monitoring for the specified HDInsight Cluster has been disabled
// This works around an issue outlined in  https://github.com/hashicorp/go-azure-sdk/issues/518 where the API
// is a LRO which doesn't use `provisioningState`.
type DisableMonitoringPoller struct {
	client    *extensions.ExtensionsClient
	clusterId extensions.ClusterId
}

func NewDisableMonitoringPoller(client *extensions.ExtensionsClient, clusterId extensions.ClusterId) *DisableMonitoringPoller {
	return &DisableMonitoringPoller{
		client:    client,
		clusterId: clusterId,
	}
}

func (p *DisableMonitoringPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.GetMonitoringStatus(ctx, p.clusterId)
	if err != nil {
		return nil, fmt.Errorf("retrieving Monitoring Status for %s: %+v", p.clusterId, err)
	}
	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving Monitoring Status for %s: `model` was nil", p.clusterId)
	}
	if resp.Model.ClusterMonitoringEnabled == nil {
		return nil, fmt.Errorf("retrieving Monitoring Status for %s: `model.ClusterMonitoringEnabled` was nil", p.clusterId)
	}

	status := pollers.PollingStatusInProgress
	if !*resp.Model.ClusterMonitoringEnabled {
		status = pollers.PollingStatusSucceeded
	}

	return &pollers.PollResult{
		HttpResponse: &client.Response{
			OData:    resp.OData,
			Response: resp.HttpResponse,
		},
		PollInterval: 10 * time.Second,
		Status:       status,
	}, nil
}
