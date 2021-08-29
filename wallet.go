// Package bitcoinwallet provide functional to
// interact with Bitcoin wallets safe for concurrent use
package bitcoinwallet

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotHaveEnoughFundsError  = errors.New("not have enough funds")
	ErrNotPositiveArgumentError = errors.New("not positive error")
)

// BitCoin is type that holds bitcoins.
type BitCoin float64

func (p *BitCoin) String() string {
	return fmt.Sprintf("%v BitCoins", *p)
}

// Wallet is type that holds BitCoins and provides
// functional to use it.
type Wallet struct {
	balance BitCoin
	mx      sync.RWMutex
}

// Withdraw method provides method to withdraw BitCoins from the wallet
// It returns updated amount of BitCoins in the wallet and error, if it occurs.
// The argument should be positive, if its not method will return NonPositiveArgumentError
// If after withdraw balance is negative, method will return NotHaveEnoughMoneyError
func (w *Wallet) Withdraw(coins BitCoin) (BitCoin, error) {
	if coins <= 0 {
		return w.balance, fmt.Errorf("withdraw: %w, amount = %v", ErrNotPositiveArgumentError, coins)
	}

	w.mx.Lock()
	defer w.mx.Unlock()

	if newBalance := w.balance - coins; newBalance >= 0 {
		w.balance = newBalance

		return w.balance, nil
	}

	return w.balance, fmt.Errorf("withdraw: %w, funds amount %v, want to withdraw %v", ErrNotHaveEnoughFundsError, w.balance, coins)
}

// Deposit method provides method to deposit BitCoins from the wallet
// It returns updated amount of BitCoins in the wallet and error, if it occurs.
// The argument should be positive, if its not method will return NonPositiveArgumentError
func (w *Wallet) Deposit(coins BitCoin) (BitCoin, error) {
	if coins <= 0 {
		return w.balance, fmt.Errorf("depost: %w, funds amount = %v", ErrNotPositiveArgumentError, coins)
	}

	w.mx.Lock()
	defer w.mx.Unlock()

	w.balance += coins

	return w.balance, nil
}

// Balance method provides method to get Balance of the wallet
func (w *Wallet) Balance() BitCoin {
	w.mx.RLock()
	defer w.mx.RUnlock()

	return w.balance
}

// NewWallet make instance of wallet
func NewWallet(coins BitCoin) Wallet {
	return Wallet{balance: coins}
}
