package token

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestOne(t *testing.T) {
	token := NewToken("BTC")
	token.Mint("0xJasper", 10)
	fmt.Println(token.Balances)
	assert.EqualValues(t, 10, token.TotalSupply(), "error")

	token.Transfer("0xJasper", "0xTom", 5)
	fmt.Println(token.Balances)
	assert.EqualValues(t, 10, token.TotalSupply(), "error")
	assert.EqualValues(t, 5, token.Balances["0xJasper"], "error")
	assert.EqualValues(t, 5, token.Balances["0xTom"], "error")

	token.Burn("0xJasper", 5)
	fmt.Println(token.Balances)
	assert.EqualValues(t, 5, token.TotalSupply(), "error")
	assert.EqualValues(t, 0, token.Balances["0xJasper"], "error")
	assert.EqualValues(t, 5, token.Balances["0xTom"], "error")
}
