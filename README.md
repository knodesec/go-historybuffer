# History Buffer
History buffer is based on a circular buffer, but without reads changing the head pointer. This allows for a historical log of *n* bytes to be retrieved by the consumer at any point. This has a fairly narrow use case, as bytes are often only important in context.
  
This package is currently not safe for concurrent use.
#### Installation
```
go get -v github.com/knodesec/go-historybuffer
```

#### Usage

```go
package main

import (
    "log"
    hbuf "github.com/knodesec/go-historybuffer"
)

func main() {

    log.SetFlags(0)

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

    // Write and fill up the buffer using 0x00 to 0x05
    written, err := history.Write(data[0:5])
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Wrote %d bytes to history buffer\n", written)

    // Read the full 5 btyes buffer
    readOut := make([]byte, 5)
    read, err := history.Read(readOut)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Read %d bytes from the history buffer\n", read)
    log.Println(readOut)

    // Write 2 more bytes 0x06 and 0x07 to the buffer
    written, err = history.Write([]byte{0x06, 0x07})
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Wrote a further %d bytes\n", written)

    // Read the buffer again and see the last 5 bytes written to it.
    readOut2 := make([]byte, 5)
    read, err = history.Read(readOut2)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("After the last write, the buffer now reads\n")
    log.Println(readOut2)
}
```

The output from the above looks like the following
```
Wrote 5 bytes to hsitory buffer
Read 5 bytes from the history buffer
[1 2 3 4 5]
Wrote a further 2 bytes
After the last write, the buffer now reads
[3 4 5 6 7]
```

#### Notes
Currently Write implements a naive write, EG it just loops over the incoming byte slice and writes them to the buffer one by one. This could be improved drastically.