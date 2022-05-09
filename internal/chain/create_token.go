package chain

import (
	"fmt"

	"gitlab.com/accumulatenetwork/accumulate/protocol"
)

type CreateToken struct{}

func (CreateToken) Type() protocol.TransactionType { return protocol.TransactionTypeCreateToken }

func (CreateToken) Execute(st *StateManager, tx *Delivery) (protocol.TransactionResult, error) {
	return (CreateToken{}).Validate(st, tx)
}

func (CreateToken) Validate(st *StateManager, tx *Delivery) (protocol.TransactionResult, error) {
	body, ok := tx.Transaction.Body.(*protocol.CreateToken)
	if !ok {
		return nil, fmt.Errorf("invalid payload: want %T, got %T", new(protocol.CreateToken), tx.Transaction.Body)
	}

	if body.Precision > 18 {
		return nil, fmt.Errorf("precision must be in range 0 to 18")
	}

	token := new(protocol.TokenIssuer)
	token.Url = body.Url
	token.Precision = body.Precision
	token.SupplyLimit = body.SupplyLimit
	token.Symbol = body.Symbol
	token.Properties = body.Properties

	err := st.SetAuth(token, body.Authorities)
	if err != nil {
		return nil, err
	}

	err = st.Create(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create %v: %w", token.Url, err)
	}
	return nil, nil
}
