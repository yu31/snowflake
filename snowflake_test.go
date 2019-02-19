package snowflake_test

import (
	"testing"

	"github.com/Yu-33/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestGenerateId(t *testing.T) {
	var err error
	var id int64

	idWorker, err := snowflake.NewSnowFlake(0)
	assert.Nil(t, err)

	number := 2 << 16
	ids := make([]int64, number)
	for i := 0; i < number; i++ {
		id, err = idWorker.NextID()
		assert.Nil(t, err)
		assert.NotEqual(t, id, 0)
		ids[i] = id
	}

	for i := 0; i < number-1; i++ {
		assert.True(t, ids[i] < ids[i+1])
	}
}

func BenchmarkSnowFlake_NextId(b *testing.B) {
	var err error
	var id int64

	idWorker, err := snowflake.NewSnowFlake(0)
	assert.Nil(b, err)

	for i := 0; i < b.N; i++ {
		id, err = idWorker.NextID()
	}

	_ = id
	_ = err
}
