/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package interfaces

type CacheManager interface {
	// add a new node with specific id and priority/value
	Add(string, int64) error
	// add a new node with specific id and default priority/value
	AddDefault(string) error
	// update value for given id, which would reflect in order
	Update(id string, value int64) error
	// default bump value for given id, which would reflect in order
	UpdateDefault(string) error
	// check if value exists
	Exists(string) bool
	// get value/priority of given id
	Get(string) (int64, error)
	// delete item with id from cache
	Delete(id string) error
	// get a list of all keys / values
	GetItems() ([]string, []int64)
	// peek top of queue (no evict)
	Peek() (string, int64, error)
	// evict
	Evict() (string, int64, error)
}
