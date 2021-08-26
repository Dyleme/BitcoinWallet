package bitcoinwallet

import (
	"errors"
	"testing"
)

var flagDepositTest = []struct {
	testName     string
	startMoney   BitCoin
	addableMoney BitCoin
	outMoney     BitCoin
	outError     error
}{
	{"standard", 12.5, 124, 136.5, nil},
	{"border", 0, 0.34, 0.34, nil},
	{"negativeArgument", 34, -1.4, 34, NonPositiveArgumentError{-1.4}},
}

func TestWallet_Deposit(t *testing.T) {
	for _, tt := range flagDepositTest {
		t.Run(tt.testName, func(t *testing.T) {
			w := Wallet{balance: tt.startMoney}
			result, err := w.Deposit(tt.addableMoney)
			if result != tt.outMoney {
				t.Errorf("Result want %v, result get %v", tt.outMoney, result)
			}
			if !errors.Is(err, tt.outError) {
				t.Errorf("Error want %v, error get %v", tt.outError, err)
			}
		})
	}
}

var flagWithdrawTest = []struct {
	testName    string
	startMoney  BitCoin
	pickedMoney BitCoin
	outMoney    BitCoin
	outError    error
}{
	{"standard", 136.5, 124, 12.5, nil},
	{"border", 0.34, 0.34, 0, nil},
	{"NotHaveEnoughMoney", 9.54, 43.54, 9.54, NotHaveEnoughMoneyError{9.54, 43.54}},
	{"Negative Argument", 34.5, -3.3, 34.5, NonPositiveArgumentError{-3.3}},
}

func TestWallet_Withdraw(t *testing.T) {
	for _, tt := range flagWithdrawTest {
		t.Run(tt.testName, func(t *testing.T) {
			w := Wallet{balance: tt.startMoney}
			result, err := w.Withdraw(tt.pickedMoney)
			if result != tt.outMoney {
				t.Errorf("Result want %v, result get %v", tt.outMoney, result)
			}
			if !errors.Is(err, tt.outError) {
				t.Errorf("Error want %v, error get %v", tt.outError, err)
			}
		})
	}
}

var flagBalanceTest = []struct {
	testName string
	inMoney  BitCoin
	outMoney BitCoin
}{
	{"standard", 12.5, 12.5},
	{"border", 0.0, 0.0},
}

func TestWallet_Balance(t *testing.T) {
	for _, tt := range flagBalanceTest {
		t.Run(tt.testName, func(t *testing.T) {
			w := Wallet{balance: tt.inMoney}
			if result := w.Balance(); result != tt.outMoney {
				t.Errorf("Result want %v, result get %v", tt.outMoney, result)
			}
		})
	}
}

func TestRaceDeposit(t *testing.T) {
	wallet := Wallet{balance: BitCoin(13.5)}

	for i := 0; i < 1000; i++ {
		go func() {
			_, err := wallet.Deposit(3.2)
			if err != nil {
				t.Errorf("Error in RaceDeposit %v", err)
			}
		}()
		go wallet.Balance()
	}
}

func TestRaceBalance(t *testing.T) {
	wallet := Wallet{balance: BitCoin(123.5)}
	for i := 0; i < 1000; i++ {
		go wallet.Balance()
		go wallet.Balance()
	}
}

func TestRaceWithdraw(t *testing.T) {
	wallet := Wallet{balance: BitCoin(13.5)}

	for i := 0; i < 1000; i++ {
		go func() {
			_, err := wallet.Withdraw(1.5)
			if err != nil {
				t.Errorf("Error in RaceWithDraw %v", err)
			}
		}()
		go func() {
			_, err := wallet.Deposit(3.5)
			if err != nil {
				t.Errorf("Error in RaceWitdDraw %v", err)
			}
		}()
	}
}
