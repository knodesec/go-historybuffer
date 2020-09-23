package historybuf

import (
	"fmt"
	"testing"
)

func TestNewBuffer(t *testing.T) {

	for bufsize := 1024; bufsize <= 16384; bufsize = bufsize * 2 {

		testname := fmt.Sprintf("%d", bufsize)
		t.Run(testname, func(t *testing.T) {
			res := New(bufsize)
			if res.size != bufsize {
				t.Errorf("got size %d want %d", res.size, bufsize)
			}
		})
	}
}

func TestWriteBasic(t *testing.T) {

	var tests = []struct {
		testname     string
		bufsize      int
		writedata    []byte
		writeback    int
		expectedbuf  []byte
		expectederr  error
		expectedhead int
		expectedtail int
	}{
		{"valid-write-1", 1, []byte{0x01}, 1, []byte{0x01}, nil, 0, 0},
		{"valid-write-2", 2, []byte{0x01, 0x02}, 2, []byte{0x01, 0x02}, nil, 0, 1},
		{"valid-write-2", 3, []byte{0x01, 0x02, 0x03}, 3, []byte{0x01, 0x02, 0x3}, nil, 0, 2},
		{"over-write-1", 3, []byte{0x01, 0x02, 0x03, 0x04}, 4, []byte{0x04, 0x02, 0x03}, nil, 1, 0},
		{"over-write-2", 3, []byte{0x01, 0x02, 0x03, 0x04, 0x05}, 5, []byte{0x04, 0x05, 0x03}, nil, 2, 1},
		{"double-over-write-1", 2, []byte{0x01, 0x02, 0x03, 0x04, 0x05}, 5, []byte{0x05, 0x04}, nil, 1, 0},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s-size-%d-write-%d", tt.testname, tt.bufsize, len(tt.writedata))
		t.Run(testname, func(t *testing.T) {

			cbuf := New(tt.bufsize)
			written, err := cbuf.Write(tt.writedata)

			if err != tt.expectederr && written != tt.writeback {
				t.Errorf("wrote %d to %d buffer, got err \"%s\" wanted %T", len(tt.writedata), tt.bufsize, err.Error(), tt.expectederr)
				return
			}

			if written != tt.writeback {
				t.Errorf("wrote %d to %d buffer, got len %d wanted %d", len(tt.writedata), tt.bufsize, written, tt.writeback)
				return
			}

			// break here for invalid write cases
			if tt.expectederr != nil {
				return
			}

			// check buffer contents match expected
			for i, expectedByte := range tt.expectedbuf {
				if expectedByte != cbuf.buf[i] {
					t.Errorf("reading byte %d expected %#v got %#v", i, expectedByte, cbuf.buf[i])
				}
			}

			// check expected head is correct
			if tt.expectedhead != cbuf.head {
				t.Errorf("buffer head expected %d got %d", tt.expectedhead, cbuf.head)
			}

			// check expected tail is correct
			if tt.expectedtail != cbuf.tail {
				t.Errorf("buffer tail expected %d got %d", tt.expectedtail, cbuf.tail)
			}
		})
	}
}

func TestReadBasic(t *testing.T) {
	var tests = []struct {
		testname     string
		bufsize      int
		data         []byte
		forcehead    int
		forcetail    int
		expecteddata []byte
		expectederr  error
		expectedsize int
	}{
		{"valid-read-1", 1, []byte{0x01}, 0, 0, []byte{0x01}, nil, 1},
		{"valid-read-2", 2, []byte{0x01, 0x02}, 0, 1, []byte{0x01, 0x02}, nil, 2},
		{"valid-read-3", 5, []byte{0x01, 0x02, 0x03, 0x04, 0x05}, 1, 0, []byte{0x02, 0x03, 0x04, 0x05, 0x01}, nil, 5},
		{"valid-read-3", 5, []byte{0x01, 0x02, 0x03, 0x04, 0x05}, 2, 1, []byte{0x03, 0x04, 0x05, 0x01, 0x02}, nil, 5},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s-size-%d-read", tt.testname, tt.bufsize)
		t.Run(testname, func(t *testing.T) {

			// force head/tail to avoid modifications to write breaking read tests.
			testBuf := New(tt.bufsize)
			testBuf.head = tt.forcehead
			testBuf.tail = tt.forcetail
			testBuf.buf = tt.data
			readOut := make([]byte, tt.expectedsize)
			readSize, err := testBuf.Read(readOut)

			if tt.expectederr != err {
				t.Errorf("expected err %T, got %T", tt.expectederr, err)
				return
			}

			if tt.expectedsize != readSize {
				t.Errorf("expected read size %d, got %d", tt.expectedsize, readSize)
				return
			}

			if len(readOut) != len(tt.expecteddata) {
				t.Errorf("expected buf len %d, got %d", len(tt.expecteddata), len(readOut))
				return
			}
			for i, b := range tt.expecteddata {
				if readOut[i] != b {
					t.Errorf("reading byte %d expected %#v, got %#v", i, b, readOut[i])
				}
			}
		})
	}
}
