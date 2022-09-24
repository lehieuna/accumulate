package main

import (
	"bytes"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/accumulatenetwork/accumulate/internal/database"
	"gitlab.com/accumulatenetwork/accumulate/internal/database/snapshot"
	"gitlab.com/accumulatenetwork/accumulate/internal/encoding"
	"gitlab.com/accumulatenetwork/accumulate/pkg/url"
	"gitlab.com/accumulatenetwork/accumulate/protocol"
	"gitlab.com/accumulatenetwork/accumulate/smt/managed"
	"gitlab.com/accumulatenetwork/accumulate/smt/pmt"
	"gitlab.com/accumulatenetwork/accumulate/smt/storage"
	"gitlab.com/accumulatenetwork/accumulate/smt/storage/badger"
)

var extractCmd = &cobra.Command{
	Use:   "extract <database> <snapshot>",
	Short: "Extracts the non-system accounts from a snapshot for use as a Genesis snapshot of a new network",
	Args:  cobra.ExactArgs(2),
	Run:   extractSnapshot,
}

func init() { cmd.AddCommand(extractCmd) }

func extractSnapshot(_ *cobra.Command, args []string) {
	store1, err := badger.New(args[0], nil)
	check(err)
	stx := store1.Begin(false)
	defer stx.Discard()

	db1 := database.New(store1, nil)
	batch1 := db1.Begin(false)
	defer batch1.Discard()

	var accounts []*url.URL
	txnHashes := new(snapshot.HashSet)
	sigHashes := new(snapshot.HashSet)

	bpt := pmt.NewBPTManager(stx)
	place := pmt.FirstPossibleBptKey
	const window = 1000 //                                       Process this many BPT entries at a time
	var count int       //                                       Recalculate number of nodes
	for {
		bptVals, next := bpt.Bpt.GetRange(place, int(window)) // Read a thousand values from the BPT
		count += len(bptVals)
		if len(bptVals) == 0 { //                                If there are none left, we break out
			break
		}
		place = next                //                           We will get the next 1000 after the last 1000
		for _, v := range bptVals { //                           For all the key values we got (as many as 1000)
			b, err := stx.Get(storage.Key(v.Key).Append("Main"))
			checkf(err, "get %v", v.Key)

			// Every account starts with the account type and the URL
			r := encoding.NewReader(bytes.NewReader(b))
			r.ReadEnum(1, new(protocol.AccountType))
			u, ok := r.ReadUrl(2)
			if !ok {
				fatalf("get %v: URL is missing", v.Key)
			}
			_, err = r.Reset(nil)
			checkf(err, "get %v", v.Key)

			// Skip system accounts
			if protocol.AcmeUrl().LocalTo(u) ||
				protocol.FaucetUrl.LocalTo(u) {
				continue
			}
			if _, ok := protocol.ParsePartitionUrl(u); ok {
				continue
			}

			accounts = append(accounts, u)
			account := batch1.Account(u)
			pending, err := account.Pending().Get()
			checkf(err, "get %v pending", u)
			for _, txid := range pending {
				txnHashes.Add(txid.Hash())
			}

			err = txnHashes.CollectFromChain(account, account.MainChain())
			checkf(err, "get %v main chain", u)

			err = txnHashes.CollectFromChain(account, account.ScratchChain())
			checkf(err, "get %v scratch chain", u)

			err = sigHashes.CollectFromChain(account, account.SignatureChain())
			checkf(err, "get %v signature chain", u)
		}
	}

	db2 := database.OpenInMemory(nil)
	batch2 := db2.Begin(true)
	defer batch2.Discard()
	for _, hash := range txnHashes.Hashes {
		hash := hash
		c, err := snapshot.CollectTransaction(batch1.Transaction(hash[:]))
		checkf(err, "collect txn %x", hash)
		err = c.Restore(batch2)
		checkf(err, "restore txn %x", hash)
		for _, c := range c.SignatureSets {
			for _, c := range c.Entries {
				sigHashes.Add(c.SignatureHash)
			}
		}
	}
	for _, hash := range sigHashes.Hashes {
		hash := hash
		c, err := snapshot.CollectSignature(batch1.Transaction(hash[:]))
		checkf(err, "collect sig %x", hash)
		err = c.Restore(batch2)
		checkf(err, "restore sig %x", hash)
	}
	for _, u := range accounts {
		c, err := snapshot.CollectAccount(batch1.Account(u), true)
		checkf(err, "collect %v", u)
		err = c.Restore(batch2)
		checkf(err, "restore %v", u)
		for _, c := range c.Chains {
			if c.Type != managed.ChainTypeTransaction {
				continue // Exclude index and anchor chains
			}
			err = c.Restore(batch2.Account(u))
			checkf(err, "restore %v %s chain", u, c.Name)
		}
	}
	check(batch2.Commit())

	f, err := os.Create(args[1])
	check(err)
	defer f.Close()

	batch2 = db2.Begin(true)
	defer batch2.Discard()
	_, err = snapshot.Collect(batch2, new(snapshot.Header), f, func(account *database.Account) (bool, error) { return true, nil })
	check(err)
}