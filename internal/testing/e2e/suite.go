package e2e

import (
	"crypto/ed25519"
	"encoding"
	"encoding/hex"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/AccumulateNetwork/accumulated/internal/api"
	"github.com/AccumulateNetwork/accumulated/internal/relay"
	"github.com/AccumulateNetwork/accumulated/internal/url"
	"github.com/AccumulateNetwork/accumulated/types/api/transactions"
	"github.com/AccumulateNetwork/accumulated/types/state"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"golang.org/x/exp/rand"
)

type Suite struct {
	suite.Suite
	start func(*Suite) *api.Query
	query *api.Query
	rand  *rand.Rand

	synthMu *sync.Mutex
	synthTx map[[32]byte]*url.URL
}

var _ suite.SetupTestSuite = (*Suite)(nil)

func NewSuite(start func(*Suite) *api.Query) *Suite {
	s := new(Suite)
	s.start = start
	return s
}

func (s *Suite) SetupTest() {
	s.query = s.start(s)
	s.rand = rand.New(rand.NewSource(0))
	s.synthMu = new(sync.Mutex)
	s.synthTx = map[[32]byte]*url.URL{}
}

func (s *Suite) generateKey() ed25519.PrivateKey {
	_, key, _ := ed25519.GenerateKey(s.rand)
	return key
}

func (s *Suite) generateTmKey() tmed25519.PrivKey {
	return tmed25519.PrivKey(s.generateKey())
}

func (s *Suite) newTx(sponsor *url.URL, key tmed25519.PrivKey, nonce uint64, body encoding.BinaryMarshaler) *transactions.GenTransaction {
	tx, err := transactions.New(sponsor.String(), func(hash []byte) (*transactions.ED25519Sig, error) {
		sig := new(transactions.ED25519Sig)
		return sig, sig.Sign(nonce, key, hash)
	}, body)
	s.Require().NoError(err)
	return tx
}

func (s *Suite) getChainAs(url string, obj encoding.BinaryUnmarshaler) {
	r, err := s.query.QueryByUrl(url)

	s.Require().NoError(err)
	s.Require().Zero(r.Response.Code, "Query failed: %v", r.Response.Info)
	so := state.Object{}
	s.Require().NoError(so.UnmarshalBinary(r.Response.Value))
	s.Require().NoError(obj.UnmarshalBinary(so.Entry))
}

func (s *Suite) sendTxAsync(tx *transactions.GenTransaction) func(relay.BatchedStatus) {
	done := make(chan abci.TxResult)
	ti, err := s.query.BroadcastTx(tx, done)
	s.Require().NoError(err)

	return func(bs relay.BatchedStatus) {
		r, err := bs.ResolveTransactionResponse(ti)
		s.Require().NoError(err)
		s.Require().Zero(r.Code, "TX failed: %s", r.Log)
		s.Require().Empty(r.MempoolError, "TX failed: %s", r.MempoolError)

		var timer *time.Timer
		if os.Getenv("CI") == "true" {
			timer = time.NewTimer(15 * time.Minute)
		} else {
			timer = time.NewTimer(1 * time.Minute)
		}
		defer timer.Stop()

		var txr abci.TxResult
		select {
		case txr = <-done:
			s.Require().Zerof(txr.Result.Code, "TX failed: %s", txr.Result.Log)
		case <-timer.C:
			s.T().Fatal("Timed out while waiting for TX repsonse")
		}

		for _, e := range txr.Result.Events {
			if e.Type != "accSyn" {
				continue
			}

			var id [32]byte
			var u *url.URL
			for _, a := range e.Attributes {
				switch a.Key {
				case "txRef":
					b, err := hex.DecodeString(a.Value)
					if s.NoError(err) {
						copy(id[:], b)
					}
				case "url":
					u, err = url.Parse(a.Value)
					s.NoError(err)
				}
			}

			if id != ([32]byte{}) && u != nil {
				s.synthMu.Lock()
				s.synthTx[id] = u
				s.synthMu.Unlock()
			}
		}
	}
}

func (s *Suite) waitForSynth() {
	for {
		s.synthMu.Lock()
		if len(s.synthTx) == 0 {
			s.synthMu.Unlock()
			return
		}

		var id [32]byte
		var u *url.URL
		for id, u = range s.synthTx {
		}
		delete(s.synthTx, id)
		s.synthMu.Unlock()

		// Poll for TX results. This is hacky, but it's a test.
		for {
			r, err := s.query.GetTx(u.Routing(), id)
			if err == nil {
				s.Require().Zero(r.TxResult.Code, "TX failed: %s", r.TxResult.Log)
				break
			}

			if !strings.Contains(err.Error(), "not found") {
				s.Require().NoError(err)
				break
			}

			time.Sleep(10 * time.Millisecond)
		}
	}
}