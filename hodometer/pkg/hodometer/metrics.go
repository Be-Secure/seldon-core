/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package hodometer

type UsageMetrics struct {
	CollectorMetrics
	ClusterMetrics
	ResourceMetrics
	FeatureMetrics
}

type CollectorMetrics struct {
	CollectorVersion   string `json:"collector_version"`
	CollectorGitCommit string `json:"collector_git_commit"`
}

type ClusterMetrics struct {
	ClusterId         string `json:"cluster_id"`
	SeldonCoreVersion string `json:"seldon_core_version"`
	KubernetesMetrics
}

type KubernetesMetrics struct {
	KubernetesVersion string `json:"kubernetes_version"`
}

type ResourceMetrics struct {
	ModelCount         uint `json:"model_count"`
	PipelineCount      uint `json:"pipeline_count"`
	ExperimentCount    uint `json:"experiment_count"`
	ServerCount        uint `json:"server_count"`
	ServerReplicaCount uint `json:"server_replica_count"`
}

type FeatureMetrics struct {
	MultimodelEnabledCount uint    `json:"multimodel_enabled_count"`
	OvercommitEnabledCount uint    `json:"overcommit_enabled_count"`
	GpuEnabledCount        uint    `json:"gpu_enabled_count"`
	ServerCpuCoresSum      float32 `json:"server_cpu_cores_sum"`
	ServerMemoryGbSum      float32 `json:"server_memory_gb_sum"`
}
