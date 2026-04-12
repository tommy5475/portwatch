// Package daemon – barrier
//
// A barrier is a synchronisation primitive that blocks a fixed number of
// goroutines (participants) until every one of them has reached the same
// checkpoint.  Once the last participant arrives the barrier releases all
// waiting goroutines simultaneously and resets itself for the next round.
//
// # Typical use
//
//	b := newBarrier(workers)
//	for i := 0; i < workers; i++ {
//	    go func() {
//	        doWork()
//	        b.Wait() // rendezvous
//	        continueAfterSync()
//	    }()
//	}
//
// # Thread safety
//
// All methods are safe for concurrent use.
package daemon
