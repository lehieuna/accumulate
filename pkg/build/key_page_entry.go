package build

import "gitlab.com/accumulatenetwork/accumulate/protocol"

type keyPageEntryBuilderArg[T any] interface {
	addEntry(protocol.KeySpecParams, []error) T
}

type KeyPageEntryBuilder[T keyPageEntryBuilderArg[T]] struct {
	t T
	parser
	entry protocol.KeySpecParams
}

func (b KeyPageEntryBuilder[T]) Owner(owner any, path ...string) KeyPageEntryBuilder[T] {
	b.entry.Delegate = b.parseUrl(owner, path...)
	return b
}

func (b KeyPageEntryBuilder[T]) Hash(hash any) KeyPageEntryBuilder[T] {
	b.entry.KeyHash = b.parseHash(hash)
	return b
}

func (b KeyPageEntryBuilder[T]) Key(key any, typ protocol.SignatureType) KeyPageEntryBuilder[T] {
	b.entry.KeyHash = b.hashKey(b.parsePublicKey(key), typ)
	return b
}

func (b KeyPageEntryBuilder[T]) FinishEntry() T {
	return b.t.addEntry(b.entry, b.errs)
}