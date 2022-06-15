package protocol

// GENERATED BY go run ./tools/cmd/gen-types. DO NOT EDIT.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"gitlab.com/accumulatenetwork/accumulate/internal/encoding"
)

// NewAccount creates a new Account for the specified AccountType.
func NewAccount(typ AccountType) (Account, error) {
	switch typ {
	case AccountTypeIdentity:
		return new(ADI), nil
	case AccountTypeAnchorLedger:
		return new(AnchorLedger), nil
	case AccountTypeDataAccount:
		return new(DataAccount), nil
	case AccountTypeKeyBook:
		return new(KeyBook), nil
	case AccountTypeKeyPage:
		return new(KeyPage), nil
	case AccountTypeLiteDataAccount:
		return new(LiteDataAccount), nil
	case AccountTypeLiteIdentity:
		return new(LiteIdentity), nil
	case AccountTypeLiteTokenAccount:
		return new(LiteTokenAccount), nil
	case AccountTypeSyntheticLedger:
		return new(SyntheticLedger), nil
	case AccountTypeSystemLedger:
		return new(SystemLedger), nil
	case AccountTypeTokenAccount:
		return new(TokenAccount), nil
	case AccountTypeTokenIssuer:
		return new(TokenIssuer), nil
	case AccountTypeUnknown:
		return new(UnknownAccount), nil
	case AccountTypeUnknownSigner:
		return new(UnknownSigner), nil
	default:
		return nil, fmt.Errorf("unknown account %v", typ)
	}
}

//EqualAccount is used to compare the values of the union
func EqualAccount(a, b Account) bool {
	if a == b {
		return true
	}
	switch a := a.(type) {
	case *ADI:
		b, ok := b.(*ADI)
		return ok && a.Equal(b)
	case *AnchorLedger:
		b, ok := b.(*AnchorLedger)
		return ok && a.Equal(b)
	case *DataAccount:
		b, ok := b.(*DataAccount)
		return ok && a.Equal(b)
	case *KeyBook:
		b, ok := b.(*KeyBook)
		return ok && a.Equal(b)
	case *KeyPage:
		b, ok := b.(*KeyPage)
		return ok && a.Equal(b)
	case *LiteDataAccount:
		b, ok := b.(*LiteDataAccount)
		return ok && a.Equal(b)
	case *LiteIdentity:
		b, ok := b.(*LiteIdentity)
		return ok && a.Equal(b)
	case *LiteTokenAccount:
		b, ok := b.(*LiteTokenAccount)
		return ok && a.Equal(b)
	case *SyntheticLedger:
		b, ok := b.(*SyntheticLedger)
		return ok && a.Equal(b)
	case *SystemLedger:
		b, ok := b.(*SystemLedger)
		return ok && a.Equal(b)
	case *TokenAccount:
		b, ok := b.(*TokenAccount)
		return ok && a.Equal(b)
	case *TokenIssuer:
		b, ok := b.(*TokenIssuer)
		return ok && a.Equal(b)
	case *UnknownAccount:
		b, ok := b.(*UnknownAccount)
		return ok && a.Equal(b)
	case *UnknownSigner:
		b, ok := b.(*UnknownSigner)
		return ok && a.Equal(b)
	default:
		return false
	}
}

// UnmarshalAccountType unmarshals the AccountType from the start of a Account.
func UnmarshalAccountType(r io.Reader) (AccountType, error) {
	var typ AccountType
	err := encoding.UnmarshalEnumType(r, &typ)
	return typ, err
}

// UnmarshalAccount unmarshals a Account.
func UnmarshalAccount(data []byte) (Account, error) {
	typ, err := UnmarshalAccountType(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	v, err := NewAccount(typ)
	if err != nil {
		return nil, err
	}

	err = v.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalAccountFrom unmarshals a Account.
func UnmarshalAccountFrom(rd io.ReadSeeker) (Account, error) {
	// Get the reader's current position
	pos, err := rd.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	// Read the type code
	typ, err := UnmarshalAccountType(rd)
	if err != nil {
		return nil, err
	}

	// Reset the reader's position
	_, err = rd.Seek(pos, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Create a new transaction result
	v, err := NewAccount(AccountType(typ))
	if err != nil {
		return nil, err
	}

	// Unmarshal the result
	err = v.UnmarshalBinaryFrom(rd)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalAccountJson unmarshals a Account.
func UnmarshalAccountJSON(data []byte) (Account, error) {
	var typ *struct{ Type AccountType }
	err := json.Unmarshal(data, &typ)
	if err != nil {
		return nil, err
	}

	if typ == nil {
		return nil, nil
	}

	acnt, err := NewAccount(typ.Type)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, acnt)
	if err != nil {
		return nil, err
	}

	return acnt, nil
}

// NewDataEntry creates a new DataEntry for the specified DataEntryType.
func NewDataEntry(typ DataEntryType) (DataEntry, error) {
	switch typ {
	case DataEntryTypeAccumulate:
		return new(AccumulateDataEntry), nil
	case DataEntryTypeFactom:
		return new(FactomDataEntry), nil
	default:
		return nil, fmt.Errorf("unknown data entry %v", typ)
	}
}

//EqualDataEntry is used to compare the values of the union
func EqualDataEntry(a, b DataEntry) bool {
	if a == b {
		return true
	}
	switch a := a.(type) {
	case *AccumulateDataEntry:
		b, ok := b.(*AccumulateDataEntry)
		return ok && a.Equal(b)
	case *FactomDataEntry:
		b, ok := b.(*FactomDataEntry)
		return ok && a.Equal(b)
	default:
		return false
	}
}

// UnmarshalDataEntryType unmarshals the DataEntryType from the start of a DataEntry.
func UnmarshalDataEntryType(r io.Reader) (DataEntryType, error) {
	var typ DataEntryType
	err := encoding.UnmarshalEnumType(r, &typ)
	return typ, err
}

// UnmarshalDataEntry unmarshals a DataEntry.
func UnmarshalDataEntry(data []byte) (DataEntry, error) {
	typ, err := UnmarshalDataEntryType(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	v, err := NewDataEntry(typ)
	if err != nil {
		return nil, err
	}

	err = v.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalDataEntryFrom unmarshals a DataEntry.
func UnmarshalDataEntryFrom(rd io.ReadSeeker) (DataEntry, error) {
	// Get the reader's current position
	pos, err := rd.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	// Read the type code
	typ, err := UnmarshalDataEntryType(rd)
	if err != nil {
		return nil, err
	}

	// Reset the reader's position
	_, err = rd.Seek(pos, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Create a new transaction result
	v, err := NewDataEntry(DataEntryType(typ))
	if err != nil {
		return nil, err
	}

	// Unmarshal the result
	err = v.UnmarshalBinaryFrom(rd)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalDataEntryJson unmarshals a DataEntry.
func UnmarshalDataEntryJSON(data []byte) (DataEntry, error) {
	var typ *struct{ Type DataEntryType }
	err := json.Unmarshal(data, &typ)
	if err != nil {
		return nil, err
	}

	if typ == nil {
		return nil, nil
	}

	acnt, err := NewDataEntry(typ.Type)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, acnt)
	if err != nil {
		return nil, err
	}

	return acnt, nil
}

// NewTransactionBody creates a new TransactionBody for the specified TransactionType.
func NewTransactionBody(typ TransactionType) (TransactionBody, error) {
	switch typ {
	case TransactionTypeAcmeFaucet:
		return new(AcmeFaucet), nil
	case TransactionTypeAddCredits:
		return new(AddCredits), nil
	case TransactionTypeBurnTokens:
		return new(BurnTokens), nil
	case TransactionTypeCreateDataAccount:
		return new(CreateDataAccount), nil
	case TransactionTypeCreateIdentity:
		return new(CreateIdentity), nil
	case TransactionTypeCreateKeyBook:
		return new(CreateKeyBook), nil
	case TransactionTypeCreateKeyPage:
		return new(CreateKeyPage), nil
	case TransactionTypeCreateToken:
		return new(CreateToken), nil
	case TransactionTypeCreateTokenAccount:
		return new(CreateTokenAccount), nil
	case TransactionTypeDirectoryAnchor:
		return new(DirectoryAnchor), nil
	case TransactionTypeIssueTokens:
		return new(IssueTokens), nil
	case TransactionTypePartitionAnchor:
		return new(PartitionAnchor), nil
	case TransactionTypeRemote:
		return new(RemoteTransaction), nil
	case TransactionTypeSendTokens:
		return new(SendTokens), nil
	case TransactionTypeSyntheticBurnTokens:
		return new(SyntheticBurnTokens), nil
	case TransactionTypeSyntheticCreateIdentity:
		return new(SyntheticCreateIdentity), nil
	case TransactionTypeSyntheticDepositCredits:
		return new(SyntheticDepositCredits), nil
	case TransactionTypeSyntheticDepositTokens:
		return new(SyntheticDepositTokens), nil
	case TransactionTypeSyntheticForwardTransaction:
		return new(SyntheticForwardTransaction), nil
	case TransactionTypeSyntheticWriteData:
		return new(SyntheticWriteData), nil
	case TransactionTypeSystemGenesis:
		return new(SystemGenesis), nil
	case TransactionTypeSystemWriteData:
		return new(SystemWriteData), nil
	case TransactionTypeUpdateAccountAuth:
		return new(UpdateAccountAuth), nil
	case TransactionTypeUpdateKey:
		return new(UpdateKey), nil
	case TransactionTypeUpdateKeyPage:
		return new(UpdateKeyPage), nil
	case TransactionTypeWriteData:
		return new(WriteData), nil
	case TransactionTypeWriteDataTo:
		return new(WriteDataTo), nil
	default:
		return nil, fmt.Errorf("unknown transaction %v", typ)
	}
}

//EqualTransactionBody is used to compare the values of the union
func EqualTransactionBody(a, b TransactionBody) bool {
	if a == b {
		return true
	}
	switch a := a.(type) {
	case *AcmeFaucet:
		b, ok := b.(*AcmeFaucet)
		return ok && a.Equal(b)
	case *AddCredits:
		b, ok := b.(*AddCredits)
		return ok && a.Equal(b)
	case *BurnTokens:
		b, ok := b.(*BurnTokens)
		return ok && a.Equal(b)
	case *CreateDataAccount:
		b, ok := b.(*CreateDataAccount)
		return ok && a.Equal(b)
	case *CreateIdentity:
		b, ok := b.(*CreateIdentity)
		return ok && a.Equal(b)
	case *CreateKeyBook:
		b, ok := b.(*CreateKeyBook)
		return ok && a.Equal(b)
	case *CreateKeyPage:
		b, ok := b.(*CreateKeyPage)
		return ok && a.Equal(b)
	case *CreateToken:
		b, ok := b.(*CreateToken)
		return ok && a.Equal(b)
	case *CreateTokenAccount:
		b, ok := b.(*CreateTokenAccount)
		return ok && a.Equal(b)
	case *DirectoryAnchor:
		b, ok := b.(*DirectoryAnchor)
		return ok && a.Equal(b)
	case *IssueTokens:
		b, ok := b.(*IssueTokens)
		return ok && a.Equal(b)
	case *PartitionAnchor:
		b, ok := b.(*PartitionAnchor)
		return ok && a.Equal(b)
	case *RemoteTransaction:
		b, ok := b.(*RemoteTransaction)
		return ok && a.Equal(b)
	case *SendTokens:
		b, ok := b.(*SendTokens)
		return ok && a.Equal(b)
	case *SyntheticBurnTokens:
		b, ok := b.(*SyntheticBurnTokens)
		return ok && a.Equal(b)
	case *SyntheticCreateIdentity:
		b, ok := b.(*SyntheticCreateIdentity)
		return ok && a.Equal(b)
	case *SyntheticDepositCredits:
		b, ok := b.(*SyntheticDepositCredits)
		return ok && a.Equal(b)
	case *SyntheticDepositTokens:
		b, ok := b.(*SyntheticDepositTokens)
		return ok && a.Equal(b)
	case *SyntheticForwardTransaction:
		b, ok := b.(*SyntheticForwardTransaction)
		return ok && a.Equal(b)
	case *SyntheticWriteData:
		b, ok := b.(*SyntheticWriteData)
		return ok && a.Equal(b)
	case *SystemGenesis:
		b, ok := b.(*SystemGenesis)
		return ok && a.Equal(b)
	case *SystemWriteData:
		b, ok := b.(*SystemWriteData)
		return ok && a.Equal(b)
	case *UpdateAccountAuth:
		b, ok := b.(*UpdateAccountAuth)
		return ok && a.Equal(b)
	case *UpdateKey:
		b, ok := b.(*UpdateKey)
		return ok && a.Equal(b)
	case *UpdateKeyPage:
		b, ok := b.(*UpdateKeyPage)
		return ok && a.Equal(b)
	case *WriteData:
		b, ok := b.(*WriteData)
		return ok && a.Equal(b)
	case *WriteDataTo:
		b, ok := b.(*WriteDataTo)
		return ok && a.Equal(b)
	default:
		return false
	}
}

// UnmarshalTransactionType unmarshals the TransactionType from the start of a TransactionBody.
func UnmarshalTransactionType(r io.Reader) (TransactionType, error) {
	var typ TransactionType
	err := encoding.UnmarshalEnumType(r, &typ)
	return typ, err
}

// UnmarshalTransactionBody unmarshals a TransactionBody.
func UnmarshalTransactionBody(data []byte) (TransactionBody, error) {
	typ, err := UnmarshalTransactionType(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	v, err := NewTransactionBody(typ)
	if err != nil {
		return nil, err
	}

	err = v.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalTransactionBodyFrom unmarshals a TransactionBody.
func UnmarshalTransactionBodyFrom(rd io.ReadSeeker) (TransactionBody, error) {
	// Get the reader's current position
	pos, err := rd.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	// Read the type code
	typ, err := UnmarshalTransactionType(rd)
	if err != nil {
		return nil, err
	}

	// Reset the reader's position
	_, err = rd.Seek(pos, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Create a new transaction result
	v, err := NewTransactionBody(TransactionType(typ))
	if err != nil {
		return nil, err
	}

	// Unmarshal the result
	err = v.UnmarshalBinaryFrom(rd)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalTransactionBodyJson unmarshals a TransactionBody.
func UnmarshalTransactionBodyJSON(data []byte) (TransactionBody, error) {
	var typ *struct{ Type TransactionType }
	err := json.Unmarshal(data, &typ)
	if err != nil {
		return nil, err
	}

	if typ == nil {
		return nil, nil
	}

	acnt, err := NewTransactionBody(typ.Type)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, acnt)
	if err != nil {
		return nil, err
	}

	return acnt, nil
}

// NewAccountAuthOperation creates a new AccountAuthOperation for the specified AccountAuthOperationType.
func NewAccountAuthOperation(typ AccountAuthOperationType) (AccountAuthOperation, error) {
	switch typ {
	case AccountAuthOperationTypeAddAuthority:
		return new(AddAccountAuthorityOperation), nil
	case AccountAuthOperationTypeDisable:
		return new(DisableAccountAuthOperation), nil
	case AccountAuthOperationTypeEnable:
		return new(EnableAccountAuthOperation), nil
	case AccountAuthOperationTypeRemoveAuthority:
		return new(RemoveAccountAuthorityOperation), nil
	default:
		return nil, fmt.Errorf("unknown account auth operation %v", typ)
	}
}

//EqualAccountAuthOperation is used to compare the values of the union
func EqualAccountAuthOperation(a, b AccountAuthOperation) bool {
	if a == b {
		return true
	}
	switch a := a.(type) {
	case *AddAccountAuthorityOperation:
		b, ok := b.(*AddAccountAuthorityOperation)
		return ok && a.Equal(b)
	case *DisableAccountAuthOperation:
		b, ok := b.(*DisableAccountAuthOperation)
		return ok && a.Equal(b)
	case *EnableAccountAuthOperation:
		b, ok := b.(*EnableAccountAuthOperation)
		return ok && a.Equal(b)
	case *RemoveAccountAuthorityOperation:
		b, ok := b.(*RemoveAccountAuthorityOperation)
		return ok && a.Equal(b)
	default:
		return false
	}
}

// UnmarshalAccountAuthOperationType unmarshals the AccountAuthOperationType from the start of a AccountAuthOperation.
func UnmarshalAccountAuthOperationType(r io.Reader) (AccountAuthOperationType, error) {
	var typ AccountAuthOperationType
	err := encoding.UnmarshalEnumType(r, &typ)
	return typ, err
}

// UnmarshalAccountAuthOperation unmarshals a AccountAuthOperation.
func UnmarshalAccountAuthOperation(data []byte) (AccountAuthOperation, error) {
	typ, err := UnmarshalAccountAuthOperationType(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	v, err := NewAccountAuthOperation(typ)
	if err != nil {
		return nil, err
	}

	err = v.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalAccountAuthOperationFrom unmarshals a AccountAuthOperation.
func UnmarshalAccountAuthOperationFrom(rd io.ReadSeeker) (AccountAuthOperation, error) {
	// Get the reader's current position
	pos, err := rd.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	// Read the type code
	typ, err := UnmarshalAccountAuthOperationType(rd)
	if err != nil {
		return nil, err
	}

	// Reset the reader's position
	_, err = rd.Seek(pos, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Create a new transaction result
	v, err := NewAccountAuthOperation(AccountAuthOperationType(typ))
	if err != nil {
		return nil, err
	}

	// Unmarshal the result
	err = v.UnmarshalBinaryFrom(rd)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalAccountAuthOperationJson unmarshals a AccountAuthOperation.
func UnmarshalAccountAuthOperationJSON(data []byte) (AccountAuthOperation, error) {
	var typ *struct{ Type AccountAuthOperationType }
	err := json.Unmarshal(data, &typ)
	if err != nil {
		return nil, err
	}

	if typ == nil {
		return nil, nil
	}

	acnt, err := NewAccountAuthOperation(typ.Type)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, acnt)
	if err != nil {
		return nil, err
	}

	return acnt, nil
}

// NewKeyPageOperation creates a new KeyPageOperation for the specified KeyPageOperationType.
func NewKeyPageOperation(typ KeyPageOperationType) (KeyPageOperation, error) {
	switch typ {
	case KeyPageOperationTypeAdd:
		return new(AddKeyOperation), nil
	case KeyPageOperationTypeRemove:
		return new(RemoveKeyOperation), nil
	case KeyPageOperationTypeSetThreshold:
		return new(SetThresholdKeyPageOperation), nil
	case KeyPageOperationTypeUpdateAllowed:
		return new(UpdateAllowedKeyPageOperation), nil
	case KeyPageOperationTypeUpdate:
		return new(UpdateKeyOperation), nil
	default:
		return nil, fmt.Errorf("unknown key page operation %v", typ)
	}
}

//EqualKeyPageOperation is used to compare the values of the union
func EqualKeyPageOperation(a, b KeyPageOperation) bool {
	if a == b {
		return true
	}
	switch a := a.(type) {
	case *AddKeyOperation:
		b, ok := b.(*AddKeyOperation)
		return ok && a.Equal(b)
	case *RemoveKeyOperation:
		b, ok := b.(*RemoveKeyOperation)
		return ok && a.Equal(b)
	case *SetThresholdKeyPageOperation:
		b, ok := b.(*SetThresholdKeyPageOperation)
		return ok && a.Equal(b)
	case *UpdateAllowedKeyPageOperation:
		b, ok := b.(*UpdateAllowedKeyPageOperation)
		return ok && a.Equal(b)
	case *UpdateKeyOperation:
		b, ok := b.(*UpdateKeyOperation)
		return ok && a.Equal(b)
	default:
		return false
	}
}

// UnmarshalKeyPageOperationType unmarshals the KeyPageOperationType from the start of a KeyPageOperation.
func UnmarshalKeyPageOperationType(r io.Reader) (KeyPageOperationType, error) {
	var typ KeyPageOperationType
	err := encoding.UnmarshalEnumType(r, &typ)
	return typ, err
}

// UnmarshalKeyPageOperation unmarshals a KeyPageOperation.
func UnmarshalKeyPageOperation(data []byte) (KeyPageOperation, error) {
	typ, err := UnmarshalKeyPageOperationType(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	v, err := NewKeyPageOperation(typ)
	if err != nil {
		return nil, err
	}

	err = v.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalKeyPageOperationFrom unmarshals a KeyPageOperation.
func UnmarshalKeyPageOperationFrom(rd io.ReadSeeker) (KeyPageOperation, error) {
	// Get the reader's current position
	pos, err := rd.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	// Read the type code
	typ, err := UnmarshalKeyPageOperationType(rd)
	if err != nil {
		return nil, err
	}

	// Reset the reader's position
	_, err = rd.Seek(pos, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Create a new transaction result
	v, err := NewKeyPageOperation(KeyPageOperationType(typ))
	if err != nil {
		return nil, err
	}

	// Unmarshal the result
	err = v.UnmarshalBinaryFrom(rd)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalKeyPageOperationJson unmarshals a KeyPageOperation.
func UnmarshalKeyPageOperationJSON(data []byte) (KeyPageOperation, error) {
	var typ *struct{ Type KeyPageOperationType }
	err := json.Unmarshal(data, &typ)
	if err != nil {
		return nil, err
	}

	if typ == nil {
		return nil, nil
	}

	acnt, err := NewKeyPageOperation(typ.Type)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, acnt)
	if err != nil {
		return nil, err
	}

	return acnt, nil
}

// NewSignature creates a new Signature for the specified SignatureType.
func NewSignature(typ SignatureType) (Signature, error) {
	switch typ {
	case SignatureTypeBTCLegacy:
		return new(BTCLegacySignature), nil
	case SignatureTypeBTC:
		return new(BTCSignature), nil
	case SignatureTypeDelegated:
		return new(DelegatedSignature), nil
	case SignatureTypeED25519:
		return new(ED25519Signature), nil
	case SignatureTypeETH:
		return new(ETHSignature), nil
	case SignatureTypeInternal:
		return new(InternalSignature), nil
	case SignatureTypeLegacyED25519:
		return new(LegacyED25519Signature), nil
	case SignatureTypeRCD1:
		return new(RCD1Signature), nil
	case SignatureTypeReceipt:
		return new(ReceiptSignature), nil
	case SignatureTypeRemote:
		return new(RemoteSignature), nil
	case SignatureTypeSet:
		return new(SignatureSet), nil
	case SignatureTypeSynthetic:
		return new(SyntheticSignature), nil
	default:
		return nil, fmt.Errorf("unknown signature %v", typ)
	}
}

//EqualSignature is used to compare the values of the union
func EqualSignature(a, b Signature) bool {
	if a == b {
		return true
	}
	switch a := a.(type) {
	case *BTCLegacySignature:
		b, ok := b.(*BTCLegacySignature)
		return ok && a.Equal(b)
	case *BTCSignature:
		b, ok := b.(*BTCSignature)
		return ok && a.Equal(b)
	case *DelegatedSignature:
		b, ok := b.(*DelegatedSignature)
		return ok && a.Equal(b)
	case *ED25519Signature:
		b, ok := b.(*ED25519Signature)
		return ok && a.Equal(b)
	case *ETHSignature:
		b, ok := b.(*ETHSignature)
		return ok && a.Equal(b)
	case *InternalSignature:
		b, ok := b.(*InternalSignature)
		return ok && a.Equal(b)
	case *LegacyED25519Signature:
		b, ok := b.(*LegacyED25519Signature)
		return ok && a.Equal(b)
	case *RCD1Signature:
		b, ok := b.(*RCD1Signature)
		return ok && a.Equal(b)
	case *ReceiptSignature:
		b, ok := b.(*ReceiptSignature)
		return ok && a.Equal(b)
	case *RemoteSignature:
		b, ok := b.(*RemoteSignature)
		return ok && a.Equal(b)
	case *SignatureSet:
		b, ok := b.(*SignatureSet)
		return ok && a.Equal(b)
	case *SyntheticSignature:
		b, ok := b.(*SyntheticSignature)
		return ok && a.Equal(b)
	default:
		return false
	}
}

// UnmarshalSignatureType unmarshals the SignatureType from the start of a Signature.
func UnmarshalSignatureType(r io.Reader) (SignatureType, error) {
	var typ SignatureType
	err := encoding.UnmarshalEnumType(r, &typ)
	return typ, err
}

// UnmarshalSignature unmarshals a Signature.
func UnmarshalSignature(data []byte) (Signature, error) {
	typ, err := UnmarshalSignatureType(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	v, err := NewSignature(typ)
	if err != nil {
		return nil, err
	}

	err = v.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalSignatureFrom unmarshals a Signature.
func UnmarshalSignatureFrom(rd io.ReadSeeker) (Signature, error) {
	// Get the reader's current position
	pos, err := rd.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	// Read the type code
	typ, err := UnmarshalSignatureType(rd)
	if err != nil {
		return nil, err
	}

	// Reset the reader's position
	_, err = rd.Seek(pos, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Create a new transaction result
	v, err := NewSignature(SignatureType(typ))
	if err != nil {
		return nil, err
	}

	// Unmarshal the result
	err = v.UnmarshalBinaryFrom(rd)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalSignatureJson unmarshals a Signature.
func UnmarshalSignatureJSON(data []byte) (Signature, error) {
	var typ *struct{ Type SignatureType }
	err := json.Unmarshal(data, &typ)
	if err != nil {
		return nil, err
	}

	if typ == nil {
		return nil, nil
	}

	acnt, err := NewSignature(typ.Type)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, acnt)
	if err != nil {
		return nil, err
	}

	return acnt, nil
}
