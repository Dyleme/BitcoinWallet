package bitcoinwallet

import "sync"

type Wallet struct {
	balance float64
	sync.RWMutex
}

func (w *Wallet) Withdraw(money float64) float64 {
	w.Lock()
	defer w.Unlock()
	w.balance -= money

	return w.balance
}

func (w *Wallet) Deposit(money float64) float64 {
	w.Lock()
	defer w.Unlock()
	w.balance += money

	return w.balance
}

func (w *Wallet) Balance() float64 {
	w.RLock()
	defer w.RUnlock()

	return w.balance
}
