package ioutil

import (
	"io"
	"log"
)

type data struct {
	Data []byte
	Err  error
}

type RdChannels struct {
	Req  chan int
	Data chan data
}

func loop(r io.Reader, buf []byte, ch RdChannels) {
	var (
		d     data
		avail = buf[:0]
		n     int
	)

	for {
		if len(avail) == 0 {
			d.Data = avail
			//			log.Println("LOOP read")
			if nread, err := r.Read(buf); err == nil {
				avail = buf[:nread]
			} else {
				d.Err = err
			}
		}
		d.Data = avail

		select {
		case ch.Data <- d:
			n = <-ch.Req
		case n = <-ch.Req:
		}
		if n == 0 {
			break
		}
		avail = avail[n:]
	}
	close(ch.Data)
	log.Println("Exit loop")
}

// Create a wrapper around an io.Reader. A pair of channels is
// returned that can be used to request a number of bytes,
// and to receive a byte slice containing the bytes actually read.
func ChannelizeReader(r io.Reader, buf []byte) (ch RdChannels) {
	if buf == nil {
		buf = make([]byte, 4096)
	}
	ch.Req = make(chan int)
	ch.Data = make(chan data)

	go loop(r, buf, ch)

	return
}
