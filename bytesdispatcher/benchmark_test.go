package bytesdispatcher

import (
	"github.com/inverse-inc/packetfence/go/bytearraypool"
	"testing"
)

func handleBytes([]byte) {
}

var bytePool = bytearraypool.NewByteArrayPool(100, 1024)
var dispatcher = NewDispatcher(4, 100, BytesHandlerFunc(handleBytes), bytePool)

func BenchmarkBytesdispatcher(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dispatcher.SubmitJob(bytePool.Get())
	}
}

func init() {
	bytePool.Fill(100)
	dispatcher.Run()
}
