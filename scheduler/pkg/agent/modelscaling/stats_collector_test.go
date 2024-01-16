/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package modelscaling

import (
	"sync"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestStatsCollectorSmoke(t *testing.T) {
	g := NewGomegaWithT(t)
	dummyModel := "model_0"

	lags := NewModelReplicaLagsKeeper()
	lastUsed := NewModelReplicaLastUsedKeeper()

	collector := NewDataPlaneStatsCollector(lags, lastUsed)

	var wg sync.WaitGroup
	wg.Add(1)

	err := collector.ScalingMetricsSetup(&wg, dummyModel)
	g.Expect(err).To(BeNil())

	lagCount, _ := collector.ModelLagStats.Get(dummyModel)
	lastUsedCount, _ := collector.ModelLastUsedStats.Get(dummyModel)

	g.Expect(lagCount).To(Equal(uint32(1)))
	g.Expect(lastUsedCount).Should(BeNumerically("<=", time.Now().Unix()))

	err = collector.ScalingMetricsTearDown(&wg, dummyModel)
	g.Expect(err).To(BeNil())

	lagCount, _ = collector.ModelLagStats.Get(dummyModel)
	g.Expect(lagCount).To(Equal(uint32(0)))

}
