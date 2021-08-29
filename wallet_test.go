package bitcoinwallet_test

import (
	"errors"
	"sync"
	"testing"

	bw "github.com/Dyleme/BitcoinWallet"
)

func TestWallet_Deposit(t *testing.T) {
	var flagDepositTest = []struct {
		testName     string
		startMoney   bw.BitCoin
		addableMoney bw.BitCoin
		outMoney     bw.BitCoin
		outError     error
	}{
		{"standard", 12.5, 124, 136.5, nil},
		{"border", 0, 0.34, 0.34, nil},
		{"Not positive argument", 34, -1.4, 34, bw.ErrNotPositiveArgumentError},
	}

	for _, tt := range flagDepositTest {
		t.Run(tt.testName, func(t *testing.T) {
			w := bw.NewWallet(tt.startMoney)
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

func TestWallet_Withdraw(t *testing.T) {
	var flagWithdrawTest = []struct {
		testName    string
		startMoney  bw.BitCoin
		pickedMoney bw.BitCoin
		outMoney    bw.BitCoin
		outError    error
	}{
		{"standard", 136.5, 124, 12.5, nil},
		{"border", 0.34, 0.34, 0, nil},
		{"NotHaveEnoughFunds", 9.54, 43.54, 9.54, bw.ErrNotHaveEnoughFundsError},
		{"Not positive argument", 34.5, -3.3, 34.5, bw.ErrNotPositiveArgumentError},
	}

	for _, tt := range flagWithdrawTest {
		t.Run(tt.testName, func(t *testing.T) {
			w := bw.NewWallet(tt.startMoney)
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

func TestWallet_Balance(t *testing.T) {
	var flagBalanceTest = []struct {
		testName string
		inMoney  bw.BitCoin
		outMoney bw.BitCoin
	}{
		{"standard", 12.5, 12.5},
		{"border", 0.0, 0.0},
	}

	for _, tt := range flagBalanceTest {
		t.Run(tt.testName, func(t *testing.T) {
			w := bw.NewWallet(tt.inMoney)
			if result := w.Balance(); result != tt.outMoney {
				t.Errorf("Result want %v, result get %v", tt.outMoney, result)
			}
		})
	}
}

func TestRaceDeposit(t *testing.T) {
	wallet := bw.NewWallet(13.5)

	var wg sync.WaitGroup

	wg.Add(2000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()

			_, err := wallet.Deposit(3.2)
			if err != nil {
				t.Errorf("Error in RaceDeposit %v", err)
			}
		}()
		go func() {
			wallet.Balance()
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestRaceBalance(t *testing.T) {
	wallet := bw.NewWallet(123.5)

	var wg sync.WaitGroup

	wg.Add(2000)

	for i := 0; i < 1000; i++ {
		var balanceFunc = func() {
			wallet.Balance()
			wg.Done()
		}
		go balanceFunc()
		go balanceFunc()
	}
	wg.Wait()
}

func TestRaceWithdraw(t *testing.T) {
	wallet := bw.NewWallet(1000)

	var wg sync.WaitGroup

	wg.Add(2000)

	for i := 0; i < 1000; i++ {
		go func() {
			_, err := wallet.Withdraw(1)
			if err != nil {
				t.Errorf("Error in RaceWithDraw %v", err)
			}

			wg.Done()
		}()
		go func() {
			_, err := wallet.Deposit(1)
			if err != nil {
				t.Errorf("Error in RaceWitdDraw %v", err)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	if wallet.Balance() != 1000 {
		t.Errorf("Results don't converge")
	}
}
