package amm

import (
	token "defi-primitives-go/token"
	"errors"
)

type Pool struct {
	Address string
	A       token.Token
	B       token.Token
	LpToken token.Token
}

var ErrConstantProduct = errors.New("a * b = k violated")

func (p *Pool) AddLiquidity(sender string, aAmount uint64, bAmount uint64) error {
	if p.LpToken.TotalSupply() == 0 {
		// first deposit
		if p.A.Transfer(sender, p.Address, aAmount) == token.ErrInsufficientBalance {
			return token.ErrInsufficientBalance
		}
		if p.B.Transfer(sender, p.Address, bAmount) == token.ErrInsufficientBalance {
			return token.ErrInsufficientBalance
		}
		p.LpToken.Mint(sender, 1)
		return nil
	}
	var addA uint64
	var addB uint64
	aBal := p.A.Balances[p.Address]
	bBal := p.B.Balances[p.Address]
	price := aBal / bBal

	if aAmount/price >= bAmount {
		// more a than b
		addA = price * bAmount
		addB = bAmount
	} else {
		addA = aAmount
		addB = aAmount / price
	}

	mint := addA / aBal * p.LpToken.TotalSupply()
	if p.A.Transfer(sender, p.Address, addA) == token.ErrInsufficientBalance {
		return token.ErrInsufficientBalance
	}
	if p.B.Transfer(sender, p.Address, addB) == token.ErrInsufficientBalance {
		return token.ErrInsufficientBalance
	}
	p.LpToken.Mint(sender, mint)
	return nil
}

func (p *Pool) RemoveLiquidity(sender string, lpTokenAmount uint64) error {
	lpShare := lpTokenAmount / p.LpToken.TotalSupply()
	removeAAmount := lpShare * p.A.Balances[p.Address]
	removeBAmount := lpShare * p.B.Balances[p.Address]
	if p.A.Transfer(p.Address, sender, removeAAmount) == token.ErrInsufficientBalance {
		return token.ErrInsufficientBalance
	}
	if p.A.Transfer(p.Address, sender, removeBAmount) == token.ErrInsufficientBalance {
		return token.ErrInsufficientBalance
	}
	if p.LpToken.Burn(sender, lpTokenAmount) == token.ErrInsufficientBalance {
		return token.ErrInsufficientBalance
	}
	return nil
}

func (p *Pool) SwapAToB(sender string, amountIn uint64) error {
	aBal := p.A.Balances[p.Address]
	bBal := p.B.Balances[p.Address]
	k := aBal * bBal
	newBBal := k / (aBal + amountIn)
	bAmountOut := bBal - newBBal
	if p.A.Transfer(sender, p.Address, amountIn) == token.ErrInsufficientBalance {
		return token.ErrInsufficientBalance
	}
	if p.B.Transfer(p.Address, sender, bAmountOut) == token.ErrInsufficientBalance {
		return token.ErrInsufficientBalance
	}
	if k < p.A.Balances[p.Address]*p.B.Balances[p.Address] {
		return ErrConstantProduct
	}
	return nil
}

func (p *Pool) SwapBToA(sender string, amountIn uint64) error {
	aBal := p.A.Balances[p.Address]
	bBal := p.B.Balances[p.Address]
	k := aBal * bBal
	newABal := k / (bBal + amountIn)
	aAmountOut := aBal - newABal
	if p.B.Transfer(sender, p.Address, amountIn) == token.ErrInsufficientBalance {
		return token.ErrInsufficientBalance
	}
	if p.A.Transfer(p.Address, sender, aAmountOut) == token.ErrInsufficientBalance {
		return token.ErrInsufficientBalance
	}
	if k < p.A.Balances[p.Address]*p.B.Balances[p.Address] {
		return ErrConstantProduct
	}
	return nil
}
