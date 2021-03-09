package snowflake

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	worker, err := New(123)
	if err != nil {
		t.Fatalf("create new worker fail %v", err)
	}
	if worker == nil {
		t.Fatal("worker is nil")
	}
	if worker.instanceId != 123 {
		t.Fatalf("mistake field worker.instanceID %d", worker.instanceId)
	}
	if worker.lastTimestamp != 0 {
		t.Fatalf("mistake field worker.lastTimestamp %d", worker.lastTimestamp)
	}
	if worker.lastSequenceId != 0 {
		t.Fatalf("mistake field worker.lastSequenceID %d", worker.lastSequenceId)
	}
	if worker.mux == nil {
		t.Fatalf("worker.mux not initialization")
	}
}

func TestSnowflake_Next(t *testing.T) {
	var id int64

	worker, err := New(0)
	if err != nil {
		t.Fatalf("new snowflake fail: %v", err)
	}

	number := 2 << 16
	ids := make([]int64, number)
	for i := 0; i < number; i++ {
		id, err = worker.Next()
		if err != nil {
			t.Fatalf("generate id fail: %v", err)
		}
		if id == 0 {
			t.Fatalf("invalid id")
		}
		ids[i] = id
	}

	for i := 0; i < number-1; i++ {
		if ids[i] >= ids[i+1] {
			t.Fatalf("invalid id <%d> and <%d>", ids[i], ids[i+1])
		}
	}
}

func TestSnowflake_Batch(t *testing.T) {
	worker, err := New(0)
	if err != nil {
		t.Fatalf("new snowflake fail: %v", err)
	}

	ids, err := worker.Batch(32)
	if err != nil {
		t.Fatalf("batch get fail: %v", err)
	}
	if len(ids) != 32 {
		t.Fatalf("unexpected number of result %d", len(ids))
	}

	for i := 0; i < len(ids)-1; i++ {
		if ids[i] >= ids[i+1] {
			t.Fatalf("invalid id <%d> and <%d>", ids[i], ids[i+1])
		}
	}
}

func TestDecompose(t *testing.T) {
	worker, err := New(317)
	if err != nil {
		t.Fatalf("new snowflake fail: %v", err)
	}

	begin := time.Now().Unix()

	time.Sleep(time.Millisecond * 10)

	id, err := worker.Next()
	if err != nil {
		t.Fatalf("generate id fail: %v", err)
	}

	timestamp, instanceID, sequenceID := Decompose(id)

	if timestamp < begin {
		t.Fatalf("mistake timestamp %d", timestamp)
	}
	if instanceID != 317 {
		t.Fatalf("mistake instanceID %d", instanceID)
	}
	if sequenceID != 0 {
		t.Fatalf("mistake sequenceID %d", sequenceID)
	}
}

func BenchmarkSnowFlake_Next(b *testing.B) {
	worker, err := New(0)
	if err != nil {
		b.Fatalf("new snowflake fail: %v", err)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id, err := worker.Next()
			if err != nil {
				b.Fatalf("get id fail: %v", err)
			}
			_ = id
		}
	})
}
