package bitcoin

import "strconv"

// A Bitcoind represents a Bitcoind client
type Bitcoind struct {
	client *rpcClient
}

// New return a new bitcoind
func New(host string, port int, user, passwd string, useSSL bool) (*Bitcoind, error) {
	rpcClient, err := newClient(host, port, user, passwd, useSSL)
	if err != nil {
		return nil, err
	}
	return &Bitcoind{rpcClient}, nil
}

// GetDifficulty returns the proof-of-work difficulty as a multiple of
// the minimum difficulty.
func (b *Bitcoind) GetDifficulty() (difficulty float64, err error) {
	r, err := b.client.call("getdifficulty", nil)
	if err != nil {
		return 0.0, err
	}

	difficulty, err = strconv.ParseFloat(string(r.Result), 64)
	return
}
