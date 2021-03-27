package helpers

import (
	"os"
	"sync"

	"github.com/cheggaaa/pb"
)

type ProgressReader struct {
	ProgBar *pb.ProgressBar
	Fp      *os.File
	Size    int64
	Reads   int64
	mux     sync.Mutex
}

func NewProgressReader(f *os.File, fileSize int64) *ProgressReader {

	progBar := ProgressReader{
		ProgBar: pb.New64(fileSize).SetUnits(pb.U_BYTES),
		Fp:      f,
		Size:    fileSize,
		Reads:   -fileSize,
	}

	return &progBar

}
func (r *ProgressReader) Read(p []byte) (int, error) {

	return r.Fp.Read(p)
}

func (r *ProgressReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.Fp.ReadAt(p, off)
	c := make(chan int64)
	var read int64
	go func(address int64, delta int64, c chan int64) {
		read = address + delta
		c <- read
	}(r.Reads, int64(n), c)
	r.Reads = <-c
	if r.Reads >= 0 {
		r.ProgBar.Set64(r.Reads)
	}
	return n, err
}

func (r *ProgressReader) Seek(offset int64, whence int) (int64, error) {
	return r.Fp.Seek(offset, whence)
}
