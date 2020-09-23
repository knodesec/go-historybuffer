# History Buffer
This package is designed for a fairly specific use case, so who knows if anyone else will want to use it.

#### Installation
```
go get -v github.com/knodesec/go-historybuffer
```

#### Usage

````go
package main

import (
    "fmt"
    hbuf "github.com/knodesec/go-historybuffer"
)

func main() {


    // Create a new history buffer
    history := hbuf.New(5)

    // Define some data to play with
    var data = []byte{
        0x01,
        0x02,
        0x03,
        0x04,
        0x05,        
    }

    // Write and fill up the buffer
    written, err := history.Write(data[0:5])
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Wrote %d bytes to hsitory buffer\n", written)

    readOut := make([]byte, 5)
    read, err := history.Read(readOut)
    if err != nil P
        log.Fatal(err)
    }

    log.Printf("Read %d bytes from the history buffer\n", read)
    log.Println(readOut)

    written, err = history.Write([]byte{0x06, 0x07})



}