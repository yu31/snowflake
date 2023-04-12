// Copyright (c) 2019, Yu Wu <yu.771991@gmail.com> All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

/*
	Signed(0) |     Millisecond Timestamp (41 bits)       | Instance (10 bits) | Sequence (12 bits)
	    0     | 00000000000000000000000000000000000000000 |     0000000000     |    000000000000
*/

const (
	sequenceBits  uint = 12
	instanceBits  uint = 10
	timestampBits uint = 41

	maxSequenceId = -1 ^ (-1 << sequenceBits)
	maxInstanceId = -1 ^ (-1 << instanceBits)
	maxTimestamp  = -1 ^ (-1 << timestampBits)

	instanceShift  = sequenceBits
	timestampShift = instanceShift + instanceBits

	maxNextIdsNum = 1024
)

const (
	originTime int64 = 1547417892000 // The default origin time 2019-01-14 06:18:12
)

// Snowflake for implements algorithm of snowflake
type Snowflake struct {
	mux            *sync.Mutex
	instanceId     int64
	lastTimestamp  int64
	lastSequenceId int64
}

// New return a new SnowFlake
func New(instanceId int64) (*Snowflake, error) {
	if instanceId < 0 {
		return nil, errors.New("snowflake: instanceID can't less than 0")
	}
	if instanceId > maxInstanceId {
		return nil, fmt.Errorf("snowflake: instanceId can't more than %d", maxInstanceId)
	}

	sf := &Snowflake{
		mux:            new(sync.Mutex),
		instanceId:     instanceId,
		lastTimestamp:  0,
		lastSequenceId: 0,
	}
	return sf, nil
}

// Batch return multiple ids at once
func (sf *Snowflake) Batch(num int) ([]int64, error) {
	if num > maxNextIdsNum || num < 0 {
		num = maxNextIdsNum
	}

	sf.mux.Lock()

	var err error
	ids := make([]int64, num)
	for i := 0; i < num; i++ {
		ids[i], err = sf.next()
		if err != nil {
			break
		}
	}

	sf.mux.Unlock()

	if err != nil {
		return nil, err
	}
	return ids, nil
}

// Next return a unique id with thread safe
func (sf *Snowflake) Next() (int64, error) {
	sf.mux.Lock()
	id, err := sf.next()
	sf.mux.Unlock()
	return id, err
}

// NextUnsafe return a unique id without thread safe
func (sf *Snowflake) NextUnsafe() (int64, error) {
	return sf.next()
}

// generate a unique id
func (sf *Snowflake) next() (int64, error) {
	var uniqueId int64
	var timestamp int64

	timestamp = sf.millTimestamp()
	if timestamp < sf.lastTimestamp {
		return 0, errors.New("snowflake; clock moved backwards")
	}

	for sf.lastSequenceId > maxSequenceId && sf.lastTimestamp == timestamp {
		time.Sleep(time.Millisecond)
		timestamp = sf.millTimestamp()
	}

	if (timestamp - originTime) >= maxTimestamp {
		return 0, errors.New("snowflake: over the time limit")
	}

	if sf.lastTimestamp == timestamp {
		sf.lastSequenceId++
	} else {
		sf.lastSequenceId = 0
	}

	sf.lastTimestamp = timestamp

	uniqueId = ((timestamp - originTime) << timestampShift) | (sf.instanceId << instanceShift) | sf.lastSequenceId
	return uniqueId, nil
}

// millTimestamp generate a unix millisecond
func (sf *Snowflake) millTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// Decompose decompose id to timestamp instance id and sequence id
func Decompose(id int64) (timestamp int64, instanceID int64, sequenceID int64) {
	timestamp = id>>timestampShift + originTime
	instanceID = id >> instanceShift & maxInstanceId
	sequenceID = id & maxSequenceId
	return
}
