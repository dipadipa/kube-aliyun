package aliyun

import (
	"errors"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	log "github.com/golang/glog"
	origcloudprovider "k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/types"
)

func nodeName2IP(s types.NodeName) string {
	return string(s)
}

func ip2NodeName(s string) types.NodeName {
	return types.NodeName(s)
}

func (w *AliyunProvider) ListRoutes(clusterName string) (routes []*origcloudprovider.Route, err error) {
	id2ip, err := w.getInstanceID2IP()
	if err != nil {
		return
	}

	tables, _, err := w.client.DescribeRouteTables(&ecs.DescribeRouteTablesArgs{
		VRouterId:    w.vrouterID,
		RouteTableId: w.routeTable,
		Pagination: common.Pagination{
			PageNumber: 0,
			PageSize:   10,
		},
	})

	if err != nil {
		err = errors.New("Unable list routes:" + err.Error())
		return
	}

	for _, t := range tables {
		for _, r := range t.RouteEntrys.RouteEntry {
			ip, ok := id2ip[r.NextHopId]
			if !ok {
				log.Warningf("Unable to get ip of instance: %v", r.NextHopId)
				continue
			}

			routes = append(routes, &origcloudprovider.Route{
				TargetNode:      ip2NodeName(ip),
				DestinationCIDR: r.DestinationCidrBlock,
			})
		}
	}

	return
}

func (w *AliyunProvider) DeleteRoute(clusterName string, route *origcloudprovider.Route) (err error) {
	ip2id, err := w.getInstanceIP2ID()
	if err != nil {
		return
	}

	var (
		instanceId string
		ok         bool
	)

	ip := nodeName2IP(route.TargetNode)
	if instanceId, ok = ip2id[ip]; !ok {
		err = errors.New("Unable to get instance id of node:" + ip)
		return
	}

	args := &ecs.DeleteRouteEntryArgs{
		RouteTableId:         w.routeTable,
		DestinationCidrBlock: route.DestinationCIDR,
		NextHopId:            instanceId,
	}

	err = w.client.DeleteRouteEntry(args)
	if err != nil {
		log.Warningf("Unable to remove vpc route for %s (subnet: %s): %+v", instanceId, route.DestinationCIDR, err)
	}

	return
}

func (w *AliyunProvider) CreateRoute(clusterName string, nameHint string, route *origcloudprovider.Route) (err error) {
	ip2id, err := w.getInstanceIP2ID()
	if err != nil {
		return
	}

	var (
		instanceId string
		ok         bool
	)

	ip := nodeName2IP(route.TargetNode)
	if instanceId, ok = ip2id[ip]; !ok {
		err = errors.New("Unable to get instance id of node:" + ip)
		return
	}

	args := &ecs.CreateRouteEntryArgs{
		RouteTableId:         w.routeTable,
		DestinationCidrBlock: route.DestinationCIDR,
		NextHopId:            instanceId,
	}

	err = w.client.CreateRouteEntry(args)
	// TODO: already there
	if err != nil {
		log.Warningf("Unable to add vpc route for %s (subnet: %s): %+v", instanceId, route.DestinationCIDR, err)
	}
	return
}