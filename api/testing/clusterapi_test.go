package testing

import (
	"fmt"
	"testing"
	"time"

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/api/server"
	"github.com/libopenstorage/openstorage/cluster"
	config "github.com/libopenstorage/openstorage/config"

	clusterclient "github.com/libopenstorage/openstorage/api/client/cluster"
	"github.com/stretchr/testify/assert"
	"go.pedge.io/dlog"
)

func init() {
	startCluster()
}

func startCluster() {

	if err := server.StartClusterAPI(
		cluster.APIBase,
		pluginPort,
	); err != nil {
		dlog.Errorf("Error starting the server")
	}

	// adding sleep to avoid race condition of connection refused.
	time.Sleep(1 * time.Second)
}
func TestClusterNodeStatusSuccess(t *testing.T) {

	ts := newMocks()
	defer ts.Stop()

	var err error
	ts.client, err = clusterclient.NewClusterClient("http://127.0.0.1:2377", "v1")

	assert.Nil(t, err)

	cfg := config.ClusterConfig{
		ClusterId:     "cluster-uid",
		DefaultDriver: "mock",
	}

	err = cluster.Init(cfg)
	/*
		cus, err := cluster.Inst()

		fmt.Printf("cus is ---%+v", cus)
		fmt.Println("err -- ", err)

	*/
	ts.MockCluster().
		EXPECT().
		NodeStatus().
		Return(api.Status_STATUS_OK, nil)

	// create client

	cs := clusterclient.ClusterManager(ts.client)

	res, err := cs.NodeStatus()

	fmt.Println("res -- ", res)
	assert.Nil(t, err)
	assert.EqualValues(t, api.Status_STATUS_NONE, res)
}

func TestClusterGossipState(t *testing.T) {

	ts := newMocks()
	defer ts.Stop()

	var err error
	ts.client, err = clusterclient.NewClusterClient("http://127.0.0.1:2377", "v1")

	assert.Nil(t, err)

	cfg := config.ClusterConfig{
		ClusterId:     "cluster-uid",
		DefaultDriver: "mock",
	}

	err = cluster.Init(cfg)
	cus, err := cluster.Inst()

	cus.Start(1, true)

	cluster := clusterclient.ClusterManager(ts.client)

	res := cluster.GetGossipState()

	fmt.Println("gossip state --- ", res)
}

func TestClusterStatus(t *testing.T) {

	ts := newMocks()
	defer ts.Stop()

	var err error
	ts.client, err = clusterclient.NewClusterClient("http://127.0.0.1:2377", "v1")

	assert.Nil(t, err)

	cfg := config.ClusterConfig{
		ClusterId:     "cluster-uid",
		DefaultDriver: "mock",
	}

	err = cluster.Init(cfg)

	ts.MockCluster().
		EXPECT().
		Enumerate().
		Return(api.Cluster{Id: "cluster-id"}, nil)

	cs := clusterclient.ClusterManager(ts.client)
	resp, err := cs.Enumerate()
	fmt.Println("id ---- ", resp.Id)
	assert.Nil(t, err)

}
