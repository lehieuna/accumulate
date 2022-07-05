package block

import (
	"fmt"

	"gitlab.com/accumulatenetwork/accumulate/config"
	"gitlab.com/accumulatenetwork/accumulate/internal/chain"
	"gitlab.com/accumulatenetwork/accumulate/internal/database"
	"gitlab.com/accumulatenetwork/accumulate/internal/errors"
	"gitlab.com/accumulatenetwork/accumulate/internal/logging"
	"gitlab.com/accumulatenetwork/accumulate/internal/url"
	"gitlab.com/accumulatenetwork/accumulate/protocol"
	"gitlab.com/accumulatenetwork/accumulate/smt/storage"
)

// ProcessTransaction processes a transaction. It will not return an error if
// the transaction fails - in that case the status code will be non zero. It
// only returns an error in cases like a database failure.
func (x *Executor) ProcessTransaction(batch *database.Batch, delivery *chain.Delivery) (*protocol.TransactionStatus, *chain.ProcessTransactionState, error) {
	r := x.BlockTimers.Start(BlockTimerTypeProcessTransaction)
	defer x.BlockTimers.Stop(r)
	// Load the status
	status, err := batch.Transaction(delivery.Transaction.GetHash()).GetStatus()
	if err != nil {
		return nil, nil, err
	}
	if status.Initiator == nil {
		// This should never happen
		return nil, nil, fmt.Errorf("transaction initiator is missing")
	}

	// Load the principal
	principal, err := batch.Account(delivery.Transaction.Header.Principal).GetState()
	switch {
	case err == nil, errors.Is(err, storage.ErrNotFound):
		// Ok
	default:
		err = errors.Format(errors.StatusUnknownError, "load principal: %w", err)
		return x.recordFailedTransaction(batch, delivery, err)
	}

	// Check if the transaction is ready to be executed
	ready, err := x.TransactionIsReady(batch, delivery, status, principal)
	if err != nil {
		return x.recordFailedTransaction(batch, delivery, err)
	}
	if !ready {
		return x.recordPendingTransaction(&x.Describe, batch, delivery)
	}

	if delivery.Transaction.Body.Type().IsSynthetic() {
		// Verify that the synthetic transaction has all the right signatures
		err = processSyntheticTransaction(batch, delivery.Transaction, status)
		if err != nil {
			return x.recordFailedTransaction(batch, delivery, err)
		}
	}

	// Set up the state manager
	var st *chain.StateManager
	if x.isGenesis {
		st = chain.NewStateManager(&x.Describe, nil, batch.Begin(true), principal, delivery.Transaction, x.logger.With("operation", "ProcessTransaction"))
	} else {
		st, err = chain.LoadStateManager(&x.Describe, &x.globals.Active, batch.Begin(true), principal, delivery.Transaction, status, x.logger.With("operation", "ProcessTransaction"))
		if err != nil {
			return x.recordFailedTransaction(batch, delivery, err)
		}
	}
	defer st.Discard()

	// Execute the transaction
	executor, ok := x.executors[delivery.Transaction.Body.Type()]
	if !ok {
		// An invalid transaction should not make it to this point
		err = errors.Format(errors.StatusInternalError, "missing executor for %v", delivery.Transaction.Body.Type())
		return x.recordFailedTransaction(batch, delivery, err)
	}

	r2 := x.BlockTimers.Start(executor.Type())
	result, err := executor.Execute(st, &chain.Delivery{Transaction: delivery.Transaction})
	x.BlockTimers.Stop(r2)
	if err != nil {
		err = errors.Wrap(errors.StatusUnknownError, err)
		return x.recordFailedTransaction(batch, delivery, err)
	}

	// Commit changes, queue state creates for synthetic transactions
	state, err := st.Commit()
	if err != nil {
		err = fmt.Errorf("commit: %w", err)
		return x.recordFailedTransaction(batch, delivery, err)
	}

	// Do extra processing for special network accounts
	err = x.processNetworkAccountUpdates(batch, delivery, principal)
	if err != nil {
		return x.recordFailedTransaction(batch, delivery, err)
	}

	return x.recordSuccessfulTransaction(batch, state, delivery, result)
}

func (x *Executor) TransactionIsReady(batch *database.Batch, delivery *chain.Delivery, status *protocol.TransactionStatus, principal protocol.Account) (bool, error) {
	var ready bool
	var err error
	typ := delivery.Transaction.Body.Type()
	switch {
	case typ.IsUser():
		ready, err = x.userTransactionIsReady(batch, delivery, status, principal)
	case typ.IsSynthetic():
		ready, err = x.synthTransactionIsReady(batch, delivery, status, principal)
	case typ.IsSystem():
		ready, err = x.systemTransactionIsReady(batch, delivery, status, principal)
	default:
		return false, errors.Format(errors.StatusInternalError, "unknown transaction type %v", typ)
	}
	return ready, errors.Wrap(errors.StatusUnknownError, err)
}

func (x *Executor) userTransactionIsReady(batch *database.Batch, delivery *chain.Delivery, status *protocol.TransactionStatus, principal protocol.Account) (bool, error) {
	// If the principal is missing, check if that's ok
	if principal == nil {
		val, ok := getValidator[chain.PrincipalValidator](x, delivery.Transaction.Body.Type())
		if !ok || !val.AllowMissingPrincipal(delivery.Transaction) {
			return false, errors.NotFound("missing principal: %v not found", delivery.Transaction.Header.Principal)
		}
	}

	// Internally produced transactions are always executed immediately
	if delivery.WasProducedInternally() {
		return true, nil
	}

	// UpdateKey transactions are always M=1 and always require a signature from
	// the initiator
	if delivery.Transaction.Body.Type() == protocol.TransactionTypeUpdateKey {
		if status.Initiator == nil {
			return false, fmt.Errorf("missing initiator")
		}

		initSigs, err := batch.Transaction(delivery.Transaction.GetHash()).ReadSignatures(status.Initiator)
		if err != nil {
			return false, fmt.Errorf("load initiator signatures: %w", err)
		}

		if initSigs.Count() == 0 {
			return false, fmt.Errorf("missing initiator signature")
		}

		return true, nil
	}

	// Delegate to the transaction executor?
	val, ok := getValidator[chain.SignerValidator](x, delivery.Transaction.Body.Type())
	if ok {
		ready, fallback, err := val.TransactionIsReady(x, batch, delivery.Transaction, status)
		if err != nil {
			return false, errors.Wrap(errors.StatusUnknownError, err)
		}
		if !fallback {
			return ready, nil
		}
	}

	// At this point we cannot continue without the principal
	if principal == nil {
		return false, errors.NotFound("missing principal: %v not found", delivery.Transaction.Header.Principal)
	}

	// Get the principal's account auth
	auth, err := x.GetAccountAuthoritySet(batch, principal)
	if err != nil {
		return false, fmt.Errorf("unable to load authority of %v: %w", delivery.Transaction.Header.Principal, err)
	}

	// For each authority
	authRequired := delivery.Transaction.Body.Type().RequireAuthorization()
	for _, entry := range auth.Authorities {
		// Do not check signers for disabled authorities
		if entry.Disabled && !authRequired {
			continue
		}

		// Check if any signer has reached its threshold
		ok, err := x.AuthorityIsSatisfied(batch, delivery.Transaction, status, entry.Url)
		if err != nil {
			return false, errors.Wrap(errors.StatusUnknownError, err)
		}
		if !ok {
			return false, nil
		}
	}

	// If every authority is disabled, at least one signature is required
	return len(status.Signers) > 0, nil
}

func (x *Executor) AuthorityIsSatisfied(batch *database.Batch, transaction *protocol.Transaction, status *protocol.TransactionStatus, authUrl *url.URL) (bool, error) {
	// Check if any signer has reached its threshold
	for _, signer := range status.FindSigners(authUrl) {
		ok, err := x.SignerIsSatisfied(batch, transaction, status, signer)
		if err != nil {
			return false, errors.Wrap(errors.StatusUnknownError, err)
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

func (x *Executor) SignerIsSatisfied(batch *database.Batch, transaction *protocol.Transaction, status *protocol.TransactionStatus, signer protocol.Signer) (bool, error) {
	// Load the signature set
	signatures, err := batch.Transaction(transaction.GetHash()).ReadSignaturesForSigner(signer)
	if err != nil {
		return false, fmt.Errorf("load signatures set %v: %w", signer.GetUrl(), err)
	}

	// Check if the signature set includes a completed set
	for _, e := range signatures.Entries() {
		if e.Type == protocol.SignatureTypeSet {
			return true, nil
		}
	}

	// Check if the threshold has been reached
	if uint64(signatures.Count()) >= signer.GetSignatureThreshold() {
		return true, nil
	}

	return false, nil
}

func (x *Executor) synthTransactionIsReady(batch *database.Batch, delivery *chain.Delivery, status *protocol.TransactionStatus, principal protocol.Account) (bool, error) {
	// Do not check the principal until the transaction is ready (see below). Do
	// not delegate "is ready?" to the transaction executor - synthetic
	// transactions _must_ be sequenced and proven before being executed.

	if status.Proof == nil {
		return false, nil
	}

	// Determine which anchor chain to load
	var partition string
	if x.Describe.NetworkType != config.Directory {
		partition = protocol.Directory
	} else {
		var ok bool
		partition, ok = protocol.ParsePartitionUrl(status.SourceNetwork)
		if !ok {
			return false, errors.Format(errors.StatusUnknownError, "%v is not a valid partition URL", status.SourceNetwork)
		}
	}

	// Load the anchor chain
	anchorChain, err := batch.Account(x.Describe.AnchorPool()).ReadChain(protocol.RootAnchorChain(partition))
	if err != nil {
		return false, errors.Format(errors.StatusUnknownError, "load %s intermediate anchor chain: %w", partition, err)
	}

	// Is the result a valid DN anchor?
	_, err = anchorChain.HeightOf(status.Proof.Anchor)
	switch {
	case err == nil:
		// Ready
	case errors.Is(err, storage.ErrNotFound):
		return false, nil
	default:
		return false, errors.Format(errors.StatusUnknownError, "get height of entry %X of %s intermediate anchor chain: %w", status.Proof.Anchor[:4], partition, err)
	}

	// Load the ledger
	var ledger *protocol.SyntheticLedger
	err = batch.Account(x.Describe.Synthetic()).GetStateAs(&ledger)
	if err != nil {
		return false, errors.Format(errors.StatusUnknownError, "load synthetic transaction ledger: %w", err)
	}

	// If the transaction is out of sequence, mark it pending
	partitionLedger := ledger.Partition(status.SourceNetwork)
	if partitionLedger.Delivered+1 != status.SequenceNumber {
		x.logger.Info("Out of sequence synthetic transaction",
			"hash", logging.AsHex(delivery.Transaction.GetHash()).Slice(0, 4),
			"seq-got", status.SequenceNumber,
			"seq-want", partitionLedger.Delivered+1,
			"source", status.SourceNetwork,
			"destination", status.DestinationNetwork,
			"type", delivery.Transaction.Body.Type(),
			"hash", logging.AsHex(delivery.Transaction.GetHash()).Slice(0, 4),
		)
		return false, nil
	}

	if principal != nil {
		return true, nil
	}

	// If the principal is required but missing, do not return an error unless
	// the transaction is ready to execute.
	// https://accumulate.atlassian.net/browse/AC-1704
	val, ok := getValidator[chain.PrincipalValidator](x, delivery.Transaction.Body.Type())
	if !ok || !val.AllowMissingPrincipal(delivery.Transaction) {
		return false, errors.NotFound("missing principal: %v not found", delivery.Transaction.Header.Principal)
	}

	return true, nil
}

func (x *Executor) systemTransactionIsReady(batch *database.Batch, delivery *chain.Delivery, status *protocol.TransactionStatus, principal protocol.Account) (bool, error) {
	// Do not check the principal until the transaction is ready (see below). Do
	// not delegate "is ready?" to the transaction executor - anchors _must_ be
	// sequenced.

	switch delivery.Transaction.Body.Type() {
	case protocol.TransactionTypeSystemGenesis, protocol.TransactionTypeSystemWriteData:
		// Do not check these
		return true, nil

	default:
		// Anchors must be sequenced
	}

	// Load the ledger
	var ledger *protocol.SyntheticLedger
	err := batch.Account(x.Describe.Synthetic()).GetStateAs(&ledger)
	if err != nil {
		return false, errors.Format(errors.StatusUnknownError, "load synthetic transaction ledger: %w", err)
	}

	// If the transaction is out of sequence, mark it pending
	partLedger := ledger.Anchor(delivery.SourceNetwork)
	if partLedger.Delivered+1 != status.SequenceNumber {
		x.logger.Info("Out of sequence anchor transaction",
			"hash", logging.AsHex(delivery.Transaction.GetHash()).Slice(0, 4),
			"seq-got", status.SequenceNumber,
			"seq-want", partLedger.Delivered+1,
			"source", status.SourceNetwork,
			"destination", status.DestinationNetwork,
			"type", delivery.Transaction.Body.Type(),
			"hash", logging.AsHex(delivery.Transaction.GetHash()).Slice(0, 4),
		)
		return false, nil
	}

	if principal != nil {
		return true, nil
	}

	// If the principal is required but missing, do not return an error unless
	// the transaction is ready to execute.
	// https://accumulate.atlassian.net/browse/AC-1704
	val, ok := getValidator[chain.PrincipalValidator](x, delivery.Transaction.Body.Type())
	if !ok || !val.AllowMissingPrincipal(delivery.Transaction) {
		return false, errors.NotFound("missing principal: %v not found", delivery.Transaction.Header.Principal)
	}

	return true, nil
}

func (x *Executor) recordTransaction(batch *database.Batch, delivery *chain.Delivery, updateStatus func(*protocol.TransactionStatus)) (*protocol.TransactionStatus, error) {
	// Store the transaction state (without signatures)
	db := batch.Transaction(delivery.Transaction.GetHash())
	err := db.PutState(&database.SigOrTxn{Transaction: delivery.Transaction})
	if err != nil {
		return nil, fmt.Errorf("store transaction: %w", err)
	}

	// Update the status
	status, err := db.GetStatus()
	if err != nil {
		return nil, fmt.Errorf("load transaction status: %w", err)
	}

	status.TxID = delivery.Transaction.ID()
	updateStatus(status)
	err = db.PutStatus(status)
	if err != nil {
		return nil, fmt.Errorf("store transaction status: %w", err)
	}

	// If the transaction is synthetic, update the synthetic ledger
	if delivery.Transaction.Body.Type().IsUser() {
		return status, nil
	}
	switch delivery.Transaction.Body.Type() {
	case protocol.TransactionTypeSystemGenesis, protocol.TransactionTypeSystemWriteData:
		return status, nil
	}

	// Update the synthetic ledger
	var ledger *protocol.SyntheticLedger
	err = batch.Account(x.Describe.Synthetic()).GetStateAs(&ledger)
	if err != nil {
		return nil, errors.Format(errors.StatusUnknownError, "load synthetic transaction ledger: %w", err)
	}

	var partLedger *protocol.PartitionSyntheticLedger
	if delivery.Transaction.Body.Type().IsSystem() {
		partLedger = ledger.Anchor(delivery.SourceNetwork)
	} else {
		partLedger = ledger.Partition(delivery.SourceNetwork)
	}

	// This should never happen, but if it does Add will panic
	if status.Pending() && delivery.SequenceNumber <= partLedger.Delivered {
		return nil, errors.Format(errors.StatusFatalError, "synthetic transactions executed out of order: delivered %d, executed %d", partLedger.Delivered, delivery.SequenceNumber)
	}

	// The ledger's Delivered number needs to be updated if the transaction
	// succeeds or fails
	if partLedger.Add(!status.Pending(), delivery.SequenceNumber, delivery.Transaction.ID()) {
		err = batch.Account(x.Describe.Synthetic()).PutState(ledger)
		if err != nil {
			return nil, errors.Format(errors.StatusUnknownError, "store synthetic transaction ledger: %w", err)
		}
	}

	return status, nil
}

func (x *Executor) recordPendingTransaction(net *config.Describe, batch *database.Batch, delivery *chain.Delivery) (*protocol.TransactionStatus, *chain.ProcessTransactionState, error) {
	// Record the transaction
	status, err := x.recordTransaction(batch, delivery, func(status *protocol.TransactionStatus) {
		status.Code = errors.StatusPending
	})
	if err != nil {
		return nil, nil, err
	}

	if delivery.Transaction.Body.Type().IsSystem() {
		return status, new(chain.ProcessTransactionState), nil
	}

	// Add the user transaction to the principal's list of pending transactions
	if delivery.Transaction.Body.Type().IsUser() {
		err = batch.Account(delivery.Transaction.Header.Principal).AddPending(delivery.Transaction.ID())
		if err != nil {
			return nil, nil, fmt.Errorf("store pending list: %w", err)
		}

		return status, new(chain.ProcessTransactionState), nil
	}

	if status.Proof == nil {
		x.logger.Error("Missing receipt for pending synthetic transaction", "hash", logging.AsHex(delivery.Transaction.GetHash()).Slice(0, 4), "type", delivery.Transaction.Body.Type())
		return status, new(chain.ProcessTransactionState), nil
	}

	x.logger.Debug("Pending synthetic transaction", "hash", logging.AsHex(delivery.Transaction.GetHash()).Slice(0, 4), "type", delivery.Transaction.Body.Type(), "anchor", logging.AsHex(status.Proof.Anchor).Slice(0, 4), "module", "synthetic")

	err = batch.Account(net.Ledger()).AddSyntheticForAnchor(*(*[32]byte)(status.Proof.Anchor), delivery.Transaction.ID())
	if err != nil {
		return nil, nil, errors.Wrap(errors.StatusUnknownError, err)
	}

	return status, new(chain.ProcessTransactionState), nil
}

func (x *Executor) recordSuccessfulTransaction(batch *database.Batch, state *chain.ProcessTransactionState, delivery *chain.Delivery, result protocol.TransactionResult) (*protocol.TransactionStatus, *chain.ProcessTransactionState, error) {
	// Record the transaction
	status, err := x.recordTransaction(batch, delivery, func(status *protocol.TransactionStatus) {
		status.Code = errors.StatusDelivered
		if result == nil {
			status.Result = new(protocol.EmptyResult)
		} else {
			status.Result = result
		}
	})
	if err != nil {
		return nil, nil, err
	}

	// Remove the transaction from the principal's list of pending transactions
	err = batch.Account(delivery.Transaction.Header.Principal).RemovePending(delivery.Transaction.ID())
	if err != nil {
		return nil, nil, fmt.Errorf("store pending list: %w", err)
	}

	// Add the transaction to the principal's main or scratch chain
	targetChain := selectTargetChain(delivery.Transaction.Body)

	err = state.ChainUpdates.AddChainEntry(batch, delivery.Transaction.Header.Principal, targetChain, protocol.ChainTypeTransaction, delivery.Transaction.GetHash(), 0, 0)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return nil, nil, fmt.Errorf("add to chain: %v", err)
	}

	typ := delivery.Transaction.Body.Type()
	if typ.IsUser() || typ.IsSystem() && !typ.IsAnchor() {
		return status, state, nil
	}

	// Check for pending synthetic/system transactions
	var ledger *protocol.SyntheticLedger
	err = batch.Account(x.Describe.Synthetic()).GetStateAs(&ledger)
	if err != nil {
		return nil, nil, errors.Format(errors.StatusUnknownError, "load synthetic transaction ledger: %w", err)
	}

	var partLedger *protocol.PartitionSyntheticLedger
	if typ.IsSynthetic() {
		partLedger = ledger.Partition(delivery.SourceNetwork)
	} else {
		partLedger = ledger.Anchor(delivery.SourceNetwork)
	}

	nextHash, ok := partLedger.Get(delivery.SequenceNumber + 1)
	if ok {
		state.ProcessAdditionalTransaction(delivery.NewSyntheticFromSequence(nextHash.Hash()))
	}

	return status, state, nil
}

func selectTargetChain(body protocol.TransactionBody) string {
	if writeData, ok := body.(*protocol.WriteData); ok {
		if writeData.Scratch {
			return protocol.ScratchChain
		}
	}
	return protocol.MainChain
}

func (x *Executor) recordFailedTransaction(batch *database.Batch, delivery *chain.Delivery, failure error) (*protocol.TransactionStatus, *chain.ProcessTransactionState, error) {
	// Record the transaction
	status, err := x.recordTransaction(batch, delivery, func(status *protocol.TransactionStatus) {
		status.Set(failure)
	})
	if err != nil {
		return nil, nil, err
	}

	// If this transaction is a synthetic transaction, send a refund
	state := new(chain.ProcessTransactionState)
	if swo, ok := delivery.Transaction.Body.(protocol.SynthTxnWithOrigin); ok {
		init, refundAmount := swo.GetRefund()
		if refundAmount > 0 {
			refund := new(protocol.SyntheticDepositCredits)
			refund.Amount = refundAmount.AsUInt64()
			state.DidProduceTxn(init, refund)
		}
	}

	// Execute the post-failure hook if the transaction executor defines one
	if val, ok := getValidator[chain.TransactionExecutorCleanup](x, delivery.Transaction.Body.Type()); ok {
		err = val.DidFail(state, delivery.Transaction)
		if err != nil {
			return nil, nil, err
		}
	}

	// Remove the transaction from the principal's list of pending transactions
	err = batch.Account(delivery.Transaction.Header.Principal).RemovePending(delivery.Transaction.ID())
	if err != nil {
		return nil, nil, fmt.Errorf("update pending list: %w", err)
	}

	// Refund the signer
	if status.Initiator == nil || !delivery.Transaction.Body.Type().IsUser() {
		return status, state, nil
	}

	// TODO Send a refund for a failed remotely initiated transaction
	if !delivery.Transaction.Header.Principal.LocalTo(status.Initiator) {
		return status, state, nil
	}

	// But only if the paid paid is larger than the max failure paid
	paid, err := protocol.ComputeTransactionFee(delivery.Transaction)
	if err != nil {
		return nil, nil, fmt.Errorf("compute fee: %w", err)
	}
	if paid <= protocol.FeeFailedMaximum {
		return status, state, nil
	}

	var signer protocol.Signer
	obj := batch.Account(status.Initiator)
	err = obj.GetStateAs(&signer)
	if err != nil {
		return nil, nil, fmt.Errorf("load initial signer: %w", err)
	}

	refund := paid - protocol.FeeFailedMaximum
	signer.CreditCredits(refund.AsUInt64())
	err = obj.PutState(signer)
	if err != nil {
		return nil, nil, fmt.Errorf("store initial signer: %w", err)
	}

	return status, state, nil
}
