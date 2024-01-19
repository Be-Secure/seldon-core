/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package experiment

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestValidateExperiment(t *testing.T) {
	g := NewGomegaWithT(t)

	type test struct {
		name       string
		store      *ExperimentStore
		experiment *Experiment
		err        error
	}

	getStrPtr := func(val string) *string { return &val }
	tests := []test{
		{
			name: "valid",
			store: &ExperimentStore{
				modelBaselines: map[string]*Experiment{},
				experiments:    map[string]*Experiment{},
			},
			experiment: &Experiment{
				Name:    "a",
				Default: getStrPtr("model1"),
				Candidates: []*Candidate{
					{
						Name: "model1",
					},
					{
						Name: "model2",
					},
				},
			},
		},
		{
			name: "duplicate candidate and mirror",
			store: &ExperimentStore{
				modelBaselines: map[string]*Experiment{},
				experiments:    map[string]*Experiment{},
			},
			experiment: &Experiment{
				Name:    "a",
				Default: getStrPtr("model1"),
				Candidates: []*Candidate{
					{
						Name: "model1",
					},
					{
						Name: "model2",
					},
				},
				Mirror: &Mirror{
					Name: "model2",
				},
			},
			err: &ExperimentNoDuplicates{experimentName: "a", resource: "model2"},
		},
		{
			name: "duplicate candidate",
			store: &ExperimentStore{
				modelBaselines: map[string]*Experiment{},
				experiments:    map[string]*Experiment{},
			},
			experiment: &Experiment{
				Name:    "a",
				Default: getStrPtr("model1"),
				Candidates: []*Candidate{
					{
						Name: "model1",
					},
					{
						Name: "model2",
					},
					{
						Name: "model2",
					},
				},
			},
			err: &ExperimentNoDuplicates{experimentName: "a", resource: "model2"},
		},
		{
			name: "baseline already exists",
			store: &ExperimentStore{
				modelBaselines: map[string]*Experiment{"model1": {Name: "b"}},
				experiments:    map[string]*Experiment{},
			},
			experiment: &Experiment{
				Name:    "a",
				Default: getStrPtr("model1"),
				Candidates: []*Candidate{
					{
						Name: "model1",
					},
					{
						Name: "model2",
					},
				},
			},
			err: &ExperimentBaselineExists{experimentName: "a", name: "model1"},
		},
		{
			name: "baseline already exists but its this model so ignore",
			store: &ExperimentStore{
				modelBaselines: map[string]*Experiment{"model1": {Name: "a"}},
				experiments:    map[string]*Experiment{},
			},
			experiment: &Experiment{
				Name:    "a",
				Default: getStrPtr("model1"),
				Candidates: []*Candidate{
					{
						Name: "model1",
					},
					{
						Name: "model2",
					},
				},
			},
		},
		{
			name: "No Canidadates or mirrors",
			store: &ExperimentStore{
				modelBaselines: map[string]*Experiment{},
				experiments:    map[string]*Experiment{},
			},
			experiment: &Experiment{
				Name: "a",
			},
			err: &ExperimentNoCandidatesOrMirrors{experimentName: "a"},
		},
		{
			name: "No Canidadates but mirror",
			store: &ExperimentStore{
				modelBaselines: map[string]*Experiment{},
				experiments:    map[string]*Experiment{},
			},
			experiment: &Experiment{
				Name: "a",
				Mirror: &Mirror{
					Name: "model",
				},
			},
		},
		{
			name: "Default model is not candidate",
			store: &ExperimentStore{
				modelBaselines: map[string]*Experiment{},
				experiments:    map[string]*Experiment{},
			},
			experiment: &Experiment{
				Name:    "a",
				Default: getStrPtr("model1"),
			},
			err: &ExperimentDefaultNotFound{experimentName: "a", defaultResource: "model1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.store.validate(test.experiment)
			if test.err != nil {
				g.Expect(err.Error()).To(Equal(test.err.Error()))
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
