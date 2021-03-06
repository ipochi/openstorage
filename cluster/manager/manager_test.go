/*
Copyright 2018 Portworx

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package manager

import (
	"testing"

	"github.com/libopenstorage/openstorage/config"
	"github.com/portworx/kvdb"
	"github.com/portworx/kvdb/mem"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	testClusterId   = "test-cluster-id"
	testClusterUuid = "test-cluster-uuid"
)

var (
	kv kvdb.Kvdb
)

func init() {
	var err error
	kv, err = kvdb.New(mem.Name, "manager_test/"+testClusterId, []string{}, nil, logrus.Panicf)
	if err != nil {
		logrus.Panicf("Failed to initialize KVDB")
	}
	if err := kvdb.SetInstance(kv); err != nil {
		logrus.Panicf("Failed to set KVDB instance")
	}
}

func TestClusterManagerUuid(t *testing.T) {
	oldInst := inst
	defer func() {
		inst = oldInst
	}()

	uuid := "uuid"
	id := "id"
	err := Init(config.ClusterConfig{
		ClusterId:   id,
		ClusterUuid: uuid,
	})
	assert.NoError(t, err)
	assert.Equal(t, uuid, inst.Uuid())
}

func TestUpdateSchedulerNodeName(t *testing.T) {
	nodeID := "node-alpha"
	Init(config.ClusterConfig{
		ClusterId:         testClusterId,
		ClusterUuid:       testClusterUuid,
		NodeId:            nodeID,
		SchedulerNodeName: "old-sched-name",
	})

	err := inst.Start(1, false, "1001")
	assert.NoError(t, err)

	node, err := inst.Inspect(nodeID)
	assert.NoError(t, err)
	assert.Equal(t, "old-sched-name", node.SchedulerNodeName)
	assert.Equal(t, node.GossipPort, "1001", "Expected gossip port to be updated in cluster database")

	err = inst.UpdateSchedulerNodeName("new-sched-name")
	assert.NoError(t, err)

	node, err = inst.Inspect(nodeID)
	assert.NoError(t, err)
	assert.Equal(t, "new-sched-name", node.SchedulerNodeName)
}
