/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package filters

import (
	"testing"

	. "github.com/onsi/gomega"

	pb "github.com/seldonio/seldon-core/apis/go/v2/mlops/scheduler"

	"github.com/seldonio/seldon-core/scheduler/v2/pkg/store"
)

func TestSharingFilter(t *testing.T) {
	g := NewGomegaWithT(t)

	type test struct {
		name     string
		model    *store.ModelVersion
		server   *store.ServerSnapshot
		expected bool
	}
	serverName := "server1"
	modelExplicitServer := store.NewModelVersion(
		&pb.Model{ModelSpec: &pb.ModelSpec{Server: &serverName}, DeploymentSpec: &pb.DeploymentSpec{Replicas: 1}},
		1,
		serverName,
		map[int]store.ReplicaStatus{3: {State: store.Loading}},
		false,
		store.ModelProgressing)
	modelSharedServer := store.NewModelVersion(
		&pb.Model{ModelSpec: &pb.ModelSpec{}, DeploymentSpec: &pb.DeploymentSpec{Replicas: 1}},
		1,
		serverName,
		map[int]store.ReplicaStatus{3: {State: store.Loading}},
		false,
		store.ModelProgressing)
	tests := []test{
		{name: "ModelAndServerMatchNotShared", model: modelExplicitServer, server: &store.ServerSnapshot{Name: serverName, Shared: false}, expected: true},
		{name: "ModelAndServerMatchShared", model: modelExplicitServer, server: &store.ServerSnapshot{Name: serverName, Shared: true}, expected: true},
		{name: "ModelAndServerDontMatch", model: modelExplicitServer, server: &store.ServerSnapshot{Name: "foo", Shared: true}, expected: false},
		{name: "SharedModelAnyServer", model: modelSharedServer, server: &store.ServerSnapshot{Name: "foo", Shared: true}, expected: true},
		{name: "SharedModelNotSharedServer", model: modelSharedServer, server: &store.ServerSnapshot{Name: "foo", Shared: false}, expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filter := SharingServerFilter{}
			ok := filter.Filter(test.model, test.server)
			g.Expect(ok).To(Equal(test.expected))
		})
	}
}
