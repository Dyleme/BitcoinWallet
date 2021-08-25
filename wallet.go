package BitcoinWallet

import "sync"

// Walleter Bad name
type Walleter interface {
	Deposit(money float64)
	Withdraw(money float64)
	Balance() float64
}

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
