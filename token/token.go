package token

import (
	"errors"
)

var ErrInsufficientBalance = errors.New("insufficient balance")

type Token struct {
	Symbol   string
	Balances map[string]uint64
}

func NewToken(symbol string) *Token {
	t := Token{
		Symbol:   symbol,
		Balances: make(map[string]uint64),
	}
	return &t
}

func (t *Token) TotalSupply() uint64 {
	var sum uint64 = 0
	for _, v := range t.Balances {
		sum += v
	}
	return sum
}

func (t *Token) Transfer(src string, dst string, amount uint64) error {
	if t.Balances[src] < amount {
		return ErrInsufficientBalance
	}
	t.Balances[src] -= amount
	t.Balances[dst] += amount
	return nil
}

func (t *Token) Mint(addr string, amount uint64) {
	t.Balances[addr] += amount
}

func (t *Token) Burn(addr string, amount uint64) error {
	if t.Balances[addr] < amount {
		return ErrInsufficientBalance
	}
	t.Balances[addr] -= amount
	return nil
}
