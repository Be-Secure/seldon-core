/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package hodometer

import (
	"fmt"
	"strings"
)

type MetricsLevel int

const (
	metricsLevelCluster MetricsLevel = iota
	metricsLevelResource
	metricsLevelFeature //nolint:varcheck
)

var supportedMetricsLevels = [...]string{
	"CLUSTER",
	"RESOURCE",
	"FEATURE",
}

func MetricsLevelFrom(level string) (MetricsLevel, error) {
	asUppercase := strings.ToUpper(level)
	for idx, sml := range supportedMetricsLevels {
		if sml == asUppercase {
			return MetricsLevel(idx), nil
		}
	}
	return -1, fmt.Errorf("level %s not recognised", level)
}

func (ml *MetricsLevel) String() string {
	if ml == nil {
		return "UNKNOWN"
	}
	return supportedMetricsLevels[int(*ml)]
}
