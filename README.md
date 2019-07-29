# SnowFlake

# Get Started

```
package main

import (
	"fmt"
	
	"github.com/Yu-33/snowflake"
)

func main() {
    // create a new worker
    instanceID := 1
    idWorker, err := snowflake.NewSnowFlake(instanceID)
    if err != nil {
        fmt.Printf("New snowflake fail: %v\n", err)
        return
    }
    
    for i := 0; i < 8; i++ {
        id, err := idWorker.NextID()
        if err != nil {
            fmt.Printf("Generate id fail: %v\n", err)
            return
        }
        fmt.Printf("New ID: %d\n", id)
    }
	
}
```

## Installing

go get github.com/Yu-33/snowflake

