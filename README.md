# snowflake
==============  

snowflake is a distributed unique ID generator inspired by [Twitter's Snowflake](https://blog.twitter.com/2010/announcing-snowflake).

A snowflake ID is composed of

    41 bits for time in units of 10 msec
    10 bits for a instance id
    12 bits for a sequence number
    
# Quick Started

## Installation

```bash
$ go get -u github.com/yu31/snowflake
```

## Usage
```go
package man

import (
	"fmt"
	
	"github.com/yu31/snowflake"
)

func main() {
    // create a new worker
    instanceID := int64(1)
    worker, err := snowflake.New(instanceID)
    if err != nil {
        fmt.Printf("new snowflake fail: %v\n", err)
        return
    }
    
    // take id
    for i := 0; i < 8; i++ {
        id, err := worker.Next()
        if err != nil {
            fmt.Printf("generate id fail: %v\n", err)
            return
        }
        fmt.Printf("New ID: %d\n", id)
    }
	
    // take batch
    ids, err := worker.Batch(100)
    if err != nil {
        fmt.Printf("generate ids fail: %v\n", err)
        return
    }
    fmt.Println(ids)
}
```
