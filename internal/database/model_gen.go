package database

// GENERATED BY go run ./tools/cmd/gen-model. DO NOT EDIT.

//lint:file-ignore S1008,U1000 generated code

import (
	"gitlab.com/accumulatenetwork/accumulate/internal/database/record"
	"gitlab.com/accumulatenetwork/accumulate/internal/errors"
	"gitlab.com/accumulatenetwork/accumulate/internal/logging"
	"gitlab.com/accumulatenetwork/accumulate/internal/url"
	"gitlab.com/accumulatenetwork/accumulate/protocol"
	"gitlab.com/accumulatenetwork/accumulate/smt/storage"
)

type Batch struct {
	logger      logging.OptionalLogger
	store       record.Store
	done        bool
	writable    bool
	id          int64
	nextChildId int64
	parent      *Batch
	kvstore     storage.KeyValueTxn
	bptEntries  map[storage.Key][32]byte

	account           map[storage.Key]*Account
	transaction       map[storage.Key]*Transaction
	blockChainUpdates map[storage.Key]*record.List[*ChainUpdate]
	blockState        map[storage.Key]*record.Set[*BlockStateSynthTxnEntry]
}

func (c *Batch) Account(url *url.URL) *Account {
	return getOrCreateMap(&c.account, record.Key{}.Append("Account", url), func() *Account {
		v := new(Account)
		v.logger = c.logger
		v.store = c.store
		v.key = record.Key{}.Append("Account", url)
		v.parent = c
		v.label = "account %[2]v"
		return v
	})
}

func (c *Batch) getTransaction(hash [32]byte) *Transaction {
	return getOrCreateMap(&c.transaction, record.Key{}.Append("Transaction", hash), func() *Transaction {
		v := new(Transaction)
		v.logger = c.logger
		v.store = c.store
		v.key = record.Key{}.Append("Transaction", hash)
		v.parent = c
		v.label = "transaction %[2]x"
		return v
	})
}

func (c *Batch) BlockChainUpdates(partition *url.URL, index uint64) *record.List[*ChainUpdate] {
	return getOrCreateMap(&c.blockChainUpdates, record.Key{}.Append("BlockChainUpdates", partition, index), func() *record.List[*ChainUpdate] {
		return record.NewList(c.logger.L, c.store, record.Key{}.Append("BlockChainUpdates", partition, index), "block chain updates %[2]v %[3]v", record.Struct[ChainUpdate]())
	})
}

func (c *Batch) BlockState(partition *url.URL) *record.Set[*BlockStateSynthTxnEntry] {
	return getOrCreateMap(&c.blockState, record.Key{}.Append("BlockState", partition), func() *record.Set[*BlockStateSynthTxnEntry] {
		return record.NewSet(c.logger.L, c.store, record.Key{}.Append("BlockState", partition), "block state %[2]v", record.Struct[BlockStateSynthTxnEntry](), func(u, v *BlockStateSynthTxnEntry) int { return u.Compare(v) })
	})
}

func (c *Batch) Resolve(key record.Key) (record.Record, record.Key, error) {
	switch key[0] {
	case "Account":
		if len(key) < 2 {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for batch")
		}
		url, okUrl := key[1].(*url.URL)
		if !okUrl {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for batch")
		}
		v := c.Account(url)
		return v, key[2:], nil
	case "Transaction":
		if len(key) < 2 {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for batch")
		}
		hash, okHash := key[1].([32]byte)
		if !okHash {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for batch")
		}
		v := c.getTransaction(hash)
		return v, key[2:], nil
	case "BlockChainUpdates":
		if len(key) < 3 {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for batch")
		}
		partition, okPartition := key[1].(*url.URL)
		index, okIndex := key[2].(uint64)
		if !okPartition || !okIndex {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for batch")
		}
		v := c.BlockChainUpdates(partition, index)
		return v, key[3:], nil
	case "BlockState":
		if len(key) < 2 {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for batch")
		}
		partition, okPartition := key[1].(*url.URL)
		if !okPartition {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for batch")
		}
		v := c.BlockState(partition)
		return v, key[2:], nil
	default:
		return nil, nil, errors.New(errors.StatusInternalError, "bad key for batch")
	}
}

func (c *Batch) IsDirty() bool {
	if c == nil {
		return false
	}

	for _, v := range c.account {
		if v.IsDirty() {
			return true
		}
	}
	for _, v := range c.transaction {
		if v.IsDirty() {
			return true
		}
	}
	for _, v := range c.blockChainUpdates {
		if v.IsDirty() {
			return true
		}
	}
	for _, v := range c.blockState {
		if v.IsDirty() {
			return true
		}
	}

	return false
}

func (c *Batch) baseCommit() error {
	if c == nil {
		return nil
	}

	var err error
	for _, v := range c.account {
		commitField(&err, v)
	}
	for _, v := range c.transaction {
		commitField(&err, v)
	}
	for _, v := range c.blockChainUpdates {
		commitField(&err, v)
	}
	for _, v := range c.blockState {
		commitField(&err, v)
	}

	return nil
}

type Account struct {
	logger logging.OptionalLogger
	store  record.Store
	key    record.Key
	label  string
	parent *Batch

	main                   *record.Value[protocol.Account]
	pending                *record.Set[*url.TxID]
	syntheticForAnchor     map[storage.Key]*record.Set[*url.TxID]
	directory              *record.Set[*url.URL]
	mainChain              *Chain2
	scratchChain           *Chain2
	signatureChain         *Chain2
	rootChain              *Chain2
	anchorSequenceChain    *Chain2
	majorBlockChain        *Chain2
	syntheticSequenceChain map[storage.Key]*Chain2
	anchorChain            map[storage.Key]*AccountAnchorChain
	chains                 *record.Set[*protocol.ChainMetadata]
	syntheticAnchors       *record.Set[[32]byte]
	data                   *AccountData
}

func (c *Account) Main() *record.Value[protocol.Account] {
	return getOrCreateField(&c.main, func() *record.Value[protocol.Account] {
		return record.NewValue(c.logger.L, c.store, c.key.Append("Main"), c.label+" main", false, record.Union(protocol.UnmarshalAccount))
	})
}

func (c *Account) Pending() *record.Set[*url.TxID] {
	return getOrCreateField(&c.pending, func() *record.Set[*url.TxID] {
		return record.NewSet(c.logger.L, c.store, c.key.Append("Pending"), c.label+" pending", record.Wrapped(record.TxidWrapper), record.CompareTxid)
	})
}

func (c *Account) SyntheticForAnchor(anchor [32]byte) *record.Set[*url.TxID] {
	return getOrCreateMap(&c.syntheticForAnchor, c.key.Append("SyntheticForAnchor", anchor), func() *record.Set[*url.TxID] {
		return record.NewSet(c.logger.L, c.store, c.key.Append("SyntheticForAnchor", anchor), c.label+" synthetic for anchor %[4]x", record.Wrapped(record.TxidWrapper), record.CompareTxid)
	})
}

func (c *Account) Directory() *record.Set[*url.URL] {
	return getOrCreateField(&c.directory, func() *record.Set[*url.URL] {
		return record.NewSet(c.logger.L, c.store, c.key.Append("Directory"), c.label+" directory", record.Wrapped(record.UrlWrapper), record.CompareUrl)
	})
}

func (c *Account) MainChain() *Chain2 {
	return getOrCreateField(&c.mainChain, func() *Chain2 {
		return newChain2(c, c.logger.L, c.store, c.key.Append("MainChain"), "main", c.label+" main chain")
	})
}

func (c *Account) ScratchChain() *Chain2 {
	return getOrCreateField(&c.scratchChain, func() *Chain2 {
		return newChain2(c, c.logger.L, c.store, c.key.Append("ScratchChain"), "scratch", c.label+" scratch chain")
	})
}

func (c *Account) SignatureChain() *Chain2 {
	return getOrCreateField(&c.signatureChain, func() *Chain2 {
		return newChain2(c, c.logger.L, c.store, c.key.Append("SignatureChain"), "signature", c.label+" signature chain")
	})
}

func (c *Account) RootChain() *Chain2 {
	return getOrCreateField(&c.rootChain, func() *Chain2 {
		return newChain2(c, c.logger.L, c.store, c.key.Append("RootChain"), "root", c.label+" root chain")
	})
}

func (c *Account) AnchorSequenceChain() *Chain2 {
	return getOrCreateField(&c.anchorSequenceChain, func() *Chain2 {
		return newChain2(c, c.logger.L, c.store, c.key.Append("AnchorSequenceChain"), "anchor-sequence", c.label+" anchor sequence chain")
	})
}

func (c *Account) MajorBlockChain() *Chain2 {
	return getOrCreateField(&c.majorBlockChain, func() *Chain2 {
		return newChain2(c, c.logger.L, c.store, c.key.Append("MajorBlockChain"), "major-block", c.label+" major block chain")
	})
}

func (c *Account) getSyntheticSequenceChain(partition string) *Chain2 {
	return getOrCreateMap(&c.syntheticSequenceChain, c.key.Append("SyntheticSequenceChain", partition), func() *Chain2 {
		return newChain2(c, c.logger.L, c.store, c.key.Append("SyntheticSequenceChain", partition), "synthetic-sequence(%[4]v)", c.label+" synthetic sequence chain %[4]v")
	})
}

func (c *Account) getAnchorChain(partition string) *AccountAnchorChain {
	return getOrCreateMap(&c.anchorChain, c.key.Append("AnchorChain", partition), func() *AccountAnchorChain {
		v := new(AccountAnchorChain)
		v.logger = c.logger
		v.store = c.store
		v.key = c.key.Append("AnchorChain", partition)
		v.parent = c
		v.label = c.label + " anchor chain %[4]v"
		return v
	})
}

func (c *Account) Chains() *record.Set[*protocol.ChainMetadata] {
	return getOrCreateField(&c.chains, func() *record.Set[*protocol.ChainMetadata] {
		return record.NewSet(c.logger.L, c.store, c.key.Append("Chains"), c.label+" chains", record.Struct[protocol.ChainMetadata](), func(u, v *protocol.ChainMetadata) int { return u.Compare(v) })
	})
}

func (c *Account) SyntheticAnchors() *record.Set[[32]byte] {
	return getOrCreateField(&c.syntheticAnchors, func() *record.Set[[32]byte] {
		return record.NewSet(c.logger.L, c.store, c.key.Append("SyntheticAnchors"), c.label+" synthetic anchors", record.Wrapped(record.HashWrapper), record.CompareHash)
	})
}

func (c *Account) Data() *AccountData {
	return getOrCreateField(&c.data, func() *AccountData {
		v := new(AccountData)
		v.logger = c.logger
		v.store = c.store
		v.key = c.key.Append("Data")
		v.parent = c
		v.label = c.label + " data"
		return v
	})
}

func (c *Account) Resolve(key record.Key) (record.Record, record.Key, error) {
	switch key[0] {
	case "Main":
		return c.Main(), key[1:], nil
	case "Pending":
		return c.Pending(), key[1:], nil
	case "SyntheticForAnchor":
		if len(key) < 2 {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for account")
		}
		anchor, okAnchor := key[1].([32]byte)
		if !okAnchor {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for account")
		}
		v := c.SyntheticForAnchor(anchor)
		return v, key[2:], nil
	case "Directory":
		return c.Directory(), key[1:], nil
	case "MainChain":
		return c.MainChain(), key[1:], nil
	case "ScratchChain":
		return c.ScratchChain(), key[1:], nil
	case "SignatureChain":
		return c.SignatureChain(), key[1:], nil
	case "RootChain":
		return c.RootChain(), key[1:], nil
	case "AnchorSequenceChain":
		return c.AnchorSequenceChain(), key[1:], nil
	case "MajorBlockChain":
		return c.MajorBlockChain(), key[1:], nil
	case "SyntheticSequenceChain":
		if len(key) < 2 {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for account")
		}
		partition, okPartition := key[1].(string)
		if !okPartition {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for account")
		}
		v := c.getSyntheticSequenceChain(partition)
		return v, key[2:], nil
	case "AnchorChain":
		if len(key) < 2 {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for account")
		}
		partition, okPartition := key[1].(string)
		if !okPartition {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for account")
		}
		v := c.getAnchorChain(partition)
		return v, key[2:], nil
	case "Chains":
		return c.Chains(), key[1:], nil
	case "SyntheticAnchors":
		return c.SyntheticAnchors(), key[1:], nil
	case "Data":
		return c.Data(), key[1:], nil
	default:
		return nil, nil, errors.New(errors.StatusInternalError, "bad key for account")
	}
}

func (c *Account) IsDirty() bool {
	if c == nil {
		return false
	}

	if fieldIsDirty(c.main) {
		return true
	}
	if fieldIsDirty(c.pending) {
		return true
	}
	for _, v := range c.syntheticForAnchor {
		if v.IsDirty() {
			return true
		}
	}
	if fieldIsDirty(c.directory) {
		return true
	}
	if fieldIsDirty(c.mainChain) {
		return true
	}
	if fieldIsDirty(c.scratchChain) {
		return true
	}
	if fieldIsDirty(c.signatureChain) {
		return true
	}
	if fieldIsDirty(c.rootChain) {
		return true
	}
	if fieldIsDirty(c.anchorSequenceChain) {
		return true
	}
	if fieldIsDirty(c.majorBlockChain) {
		return true
	}
	for _, v := range c.syntheticSequenceChain {
		if v.IsDirty() {
			return true
		}
	}
	for _, v := range c.anchorChain {
		if v.IsDirty() {
			return true
		}
	}
	if fieldIsDirty(c.chains) {
		return true
	}
	if fieldIsDirty(c.syntheticAnchors) {
		return true
	}
	if fieldIsDirty(c.data) {
		return true
	}

	return false
}

func (c *Account) baseCommit() error {
	if c == nil {
		return nil
	}

	var err error
	commitField(&err, c.main)
	commitField(&err, c.pending)
	for _, v := range c.syntheticForAnchor {
		commitField(&err, v)
	}
	commitField(&err, c.directory)
	commitField(&err, c.mainChain)
	commitField(&err, c.scratchChain)
	commitField(&err, c.signatureChain)
	commitField(&err, c.rootChain)
	commitField(&err, c.anchorSequenceChain)
	commitField(&err, c.majorBlockChain)
	for _, v := range c.syntheticSequenceChain {
		commitField(&err, v)
	}
	for _, v := range c.anchorChain {
		commitField(&err, v)
	}
	commitField(&err, c.chains)
	commitField(&err, c.syntheticAnchors)
	commitField(&err, c.data)

	return nil
}

type AccountAnchorChain struct {
	logger logging.OptionalLogger
	store  record.Store
	key    record.Key
	label  string
	parent *Account

	root *Chain2
	bpt  *Chain2
}

func (c *AccountAnchorChain) Root() *Chain2 {
	return getOrCreateField(&c.root, func() *Chain2 {
		return newChain2(c, c.logger.L, c.store, c.key.Append("Root"), "anchor(%[4]v)-root", c.label+" root")
	})
}

func (c *AccountAnchorChain) BPT() *Chain2 {
	return getOrCreateField(&c.bpt, func() *Chain2 {
		return newChain2(c, c.logger.L, c.store, c.key.Append("BPT"), "anchor(%[4]v)-bpt", c.label+" bpt")
	})
}

func (c *AccountAnchorChain) Resolve(key record.Key) (record.Record, record.Key, error) {
	switch key[0] {
	case "Root":
		return c.Root(), key[1:], nil
	case "BPT":
		return c.BPT(), key[1:], nil
	default:
		return nil, nil, errors.New(errors.StatusInternalError, "bad key for anchor chain")
	}
}

func (c *AccountAnchorChain) IsDirty() bool {
	if c == nil {
		return false
	}

	if fieldIsDirty(c.root) {
		return true
	}
	if fieldIsDirty(c.bpt) {
		return true
	}

	return false
}

func (c *AccountAnchorChain) Commit() error {
	if c == nil {
		return nil
	}

	var err error
	commitField(&err, c.root)
	commitField(&err, c.bpt)

	return nil
}

type AccountData struct {
	logger logging.OptionalLogger
	store  record.Store
	key    record.Key
	label  string
	parent *Account

	entry       *record.Counted[[32]byte]
	transaction map[storage.Key]*record.Value[[32]byte]
}

func (c *AccountData) Entry() *record.Counted[[32]byte] {
	return getOrCreateField(&c.entry, func() *record.Counted[[32]byte] {
		return record.NewCounted(c.logger.L, c.store, c.key.Append("Entry"), c.label+" entry", record.WrappedFactory(record.HashWrapper))
	})
}

func (c *AccountData) Transaction(entryHash [32]byte) *record.Value[[32]byte] {
	return getOrCreateMap(&c.transaction, c.key.Append("Transaction", entryHash), func() *record.Value[[32]byte] {
		return record.NewValue(c.logger.L, c.store, c.key.Append("Transaction", entryHash), c.label+" transaction %[5]x", false, record.Wrapped(record.HashWrapper))
	})
}

func (c *AccountData) Resolve(key record.Key) (record.Record, record.Key, error) {
	switch key[0] {
	case "Entry":
		return c.Entry(), key[1:], nil
	case "Transaction":
		if len(key) < 2 {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for data")
		}
		entryHash, okEntryHash := key[1].([32]byte)
		if !okEntryHash {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for data")
		}
		v := c.Transaction(entryHash)
		return v, key[2:], nil
	default:
		return nil, nil, errors.New(errors.StatusInternalError, "bad key for data")
	}
}

func (c *AccountData) IsDirty() bool {
	if c == nil {
		return false
	}

	if fieldIsDirty(c.entry) {
		return true
	}
	for _, v := range c.transaction {
		if v.IsDirty() {
			return true
		}
	}

	return false
}

func (c *AccountData) Commit() error {
	if c == nil {
		return nil
	}

	var err error
	commitField(&err, c.entry)
	for _, v := range c.transaction {
		commitField(&err, v)
	}

	return nil
}

type Transaction struct {
	logger logging.OptionalLogger
	store  record.Store
	key    record.Key
	label  string
	parent *Batch

	main       *record.Value[*SigOrTxn]
	status     *record.Value[*protocol.TransactionStatus]
	produced   *record.Set[*url.TxID]
	signatures map[storage.Key]*record.Value[*sigSetData]
	chains     *record.Set[*TransactionChainEntry]
}

func (c *Transaction) Main() *record.Value[*SigOrTxn] {
	return getOrCreateField(&c.main, func() *record.Value[*SigOrTxn] {
		return record.NewValue(c.logger.L, c.store, c.key.Append("Main"), c.label+" main", false, record.Struct[SigOrTxn]())
	})
}

func (c *Transaction) Status() *record.Value[*protocol.TransactionStatus] {
	return getOrCreateField(&c.status, func() *record.Value[*protocol.TransactionStatus] {
		return record.NewValue(c.logger.L, c.store, c.key.Append("Status"), c.label+" status", true, record.Struct[protocol.TransactionStatus]())
	})
}

func (c *Transaction) Produced() *record.Set[*url.TxID] {
	return getOrCreateField(&c.produced, func() *record.Set[*url.TxID] {
		return record.NewSet(c.logger.L, c.store, c.key.Append("Produced"), c.label+" produced", record.Wrapped(record.TxidWrapper), record.CompareTxid)
	})
}

func (c *Transaction) getSignatures(signer *url.URL) *record.Value[*sigSetData] {
	return getOrCreateMap(&c.signatures, c.key.Append("Signatures", signer), func() *record.Value[*sigSetData] {
		return record.NewValue(c.logger.L, c.store, c.key.Append("Signatures", signer), c.label+" signatures %[4]v", true, record.Struct[sigSetData]())
	})
}

func (c *Transaction) Chains() *record.Set[*TransactionChainEntry] {
	return getOrCreateField(&c.chains, func() *record.Set[*TransactionChainEntry] {
		return record.NewSet(c.logger.L, c.store, c.key.Append("Chains"), c.label+" chains", record.Struct[TransactionChainEntry](), func(u, v *TransactionChainEntry) int { return u.Compare(v) })
	})
}

func (c *Transaction) Resolve(key record.Key) (record.Record, record.Key, error) {
	switch key[0] {
	case "Main":
		return c.Main(), key[1:], nil
	case "Status":
		return c.Status(), key[1:], nil
	case "Produced":
		return c.Produced(), key[1:], nil
	case "Signatures":
		if len(key) < 2 {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for transaction")
		}
		signer, okSigner := key[1].(*url.URL)
		if !okSigner {
			return nil, nil, errors.New(errors.StatusInternalError, "bad key for transaction")
		}
		v := c.getSignatures(signer)
		return v, key[2:], nil
	case "Chains":
		return c.Chains(), key[1:], nil
	default:
		return nil, nil, errors.New(errors.StatusInternalError, "bad key for transaction")
	}
}

func (c *Transaction) IsDirty() bool {
	if c == nil {
		return false
	}

	if fieldIsDirty(c.main) {
		return true
	}
	if fieldIsDirty(c.status) {
		return true
	}
	if fieldIsDirty(c.produced) {
		return true
	}
	for _, v := range c.signatures {
		if v.IsDirty() {
			return true
		}
	}
	if fieldIsDirty(c.chains) {
		return true
	}

	return false
}

func (c *Transaction) Commit() error {
	if c == nil {
		return nil
	}

	var err error
	commitField(&err, c.main)
	commitField(&err, c.status)
	commitField(&err, c.produced)
	for _, v := range c.signatures {
		commitField(&err, v)
	}
	commitField(&err, c.chains)

	return nil
}

func getOrCreateField[T any](ptr **T, create func() *T) *T {
	if *ptr != nil {
		return *ptr
	}

	*ptr = create()
	return *ptr
}

func getOrCreateMap[T any](ptr *map[storage.Key]T, key record.Key, create func() T) T {
	if *ptr == nil {
		*ptr = map[storage.Key]T{}
	}

	k := key.Hash()
	if v, ok := (*ptr)[k]; ok {
		return v
	}

	v := create()
	(*ptr)[k] = v
	return v
}

func commitField[T any, PT record.RecordPtr[T]](lastErr *error, field PT) {
	if *lastErr != nil || field == nil {
		return
	}

	*lastErr = field.Commit()
}

func fieldIsDirty[T any, PT record.RecordPtr[T]](field PT) bool {
	return field != nil && field.IsDirty()
}