// Package bitcoinwallet provide functional to
// interact with Bitcoin wallets safe for concurrent use
package bitcoinwallet

import (
	"fmt"
	"sync"
)

// BitCoin is type that holds bitcoins
type BitCoin float64

func (p *BitCoin) String() string {
	return fmt.Sprintf("%v BitCoins", *p)
}

// Wallet is type that holds BitCoins and provides
// functional to use it
type Wallet struct {
	balance BitCoin
	mx      sync.RWMutex
}

type NotHaveEnoughMoneyError struct {
	youGot  BitCoin
	youWant BitCoin
}

type NonPositiveArgumentError struct {
	arg BitCoin
}

func (e NonPositiveArgumentError) Error() string {
	return fmt.Sprintf("Argument must be not positive, your argument %v", e.arg)
}

func (e NotHaveEnoughMoneyError) Error() string {
	return fmt.Sprintf("You got %v, you want to wittdraw %v", e.youGot, e.youWant)
}

// Withdraw method provides method to withdraw BitCoins from the wallet
// It returns updated amount of BitCoins in the wallet and error, if it occurs.
// The argument should be positive, if its not method will return NonPositiveArgumentError
// If after withdraw balance is negative, method will return NotHaveEnoughMoneyError
func (w *Wallet) Withdraw(money BitCoin) (BitCoin, error) {
	if money <= 0 {
		return w.balance, NonPositiveArgumentError{money}
	}

	w.mx.Lock()
	defer w.mx.Unlock()

	if newBalance := w.balance - money; newBalance >= 0 {
		w.balance = newBalance

		return w.balance, nil
	}

	return w.balance, NotHaveEnoughMoneyError{w.balance, money}
}

// Deposit method provides method to deposit BitCoins from the wallet
// It returns updated amount of BitCoins in the wallet and error, if it occurs.
// The argument should be positive, if its not method will return NonPositiveArgumentError
func (w *Wallet) Deposit(money BitCoin) (BitCoin, error) {
	if money <= 0 {
		return w.balance, NonPositiveArgumentError{money}
	}

	w.mx.Lock()
	defer w.mx.Unlock()

	w.balance += money

	return w.balance, nil
}

// Balance method provides method to get Balance of the wallet
func (w *Wallet) Balance() BitCoin {
	w.mx.RLock()
	defer w.mx.RUnlock()

	return w.balance
}
