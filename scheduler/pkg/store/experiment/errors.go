/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package experiment

import "fmt"

type ExperimentNotFound struct {
	experimentName string
}

func (enf *ExperimentNotFound) Is(tgt error) bool {
	_, ok := tgt.(*ExperimentNotFound)
	return ok
}

func (enf *ExperimentNotFound) Error() string {
	return fmt.Sprintf("Experiment not found %s", enf.experimentName)
}

type ExperimentBaselineExists struct {
	experimentName string
	name           string
}

func (ebe *ExperimentBaselineExists) Error() string {
	return fmt.Sprintf("Resource %s already in experiment %s as a baseline. A model or pipeline can only appear in one experiment as a baseline", ebe.name, ebe.experimentName)
}

type ExperimentNoCandidatesOrMirrors struct {
	experimentName string
}

func (enc *ExperimentNoCandidatesOrMirrors) Error() string {
	return fmt.Sprintf("experiment %s has no candidates or mirror", enc.experimentName)
}

type ExperimentDefaultNotFound struct {
	experimentName  string
	defaultResource string
}

func (enc *ExperimentDefaultNotFound) Is(tgt error) bool {
	_, ok := tgt.(*ExperimentDefaultNotFound)
	return ok
}

func (enc *ExperimentDefaultNotFound) Error() string {
	return fmt.Sprintf("default model/pipeline %s not found in experiment %s candidates", enc.defaultResource, enc.experimentName)
}

type ExperimentNoDuplicates struct {
	experimentName string
	resource       string
}

func (enc *ExperimentNoDuplicates) Is(tgt error) bool {
	_, ok := tgt.(*ExperimentNoDuplicates)
	return ok
}

func (enc *ExperimentNoDuplicates) Error() string {
	return fmt.Sprintf("each candidate and mirror must be unique but found resource %s duplicated in experiment %s", enc.resource, enc.experimentName)
}
