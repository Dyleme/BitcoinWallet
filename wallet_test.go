package bitcoinwallet_test

import (
	"errors"
	"sync"
	"testing"

	bw "github.com/Dyleme/BitcoinWallet"
)

func TestWallet_Deposit(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
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

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			w := bw.NewWallet(tc.startMoney)
			result, err := w.Deposit(tc.addableMoney)
			if result != tc.outMoney {
				t.Errorf("wallet: want result %v, get result %v", tc.outMoney, result)
			}
			if !errors.Is(err, tc.outError) {
				t.Errorf("wallet: want error %v, get error %v", tc.outError, err)
			}
		})
	}
}

func TestWallet_Withdraw(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
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

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			w := bw.NewWallet(tc.startMoney)
			result, err := w.Withdraw(tc.pickedMoney)
			if result != tc.outMoney {
				t.Errorf("wallet: want result %v, get result %v", tc.outMoney, result)
			}
			if !errors.Is(err, tc.outError) {
				t.Errorf("wallet: want error %v, get error %v", tc.outError, err)
			}
		})
	}
}

func TestWallet_Balance(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		testName string
		inMoney  bw.BitCoin
		outMoney bw.BitCoin
	}{
		{"standard", 12.5, 12.5},
		{"border", 0.0, 0.0},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			w := bw.NewWallet(tc.inMoney)
			if result := w.Balance(); result != tc.outMoney {
				t.Errorf("wallet: want result %v, get result %v", tc.outMoney, result)
			}
		})
	}
}

func TestRaceDeposit(t *testing.T) {
	t.Parallel()

	wallet := bw.NewWallet(13.5)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()

			_, err := wallet.Deposit(0.5)
			if err != nil {
				t.Errorf("wallet: error in racedeposit %v", err)
			}
		}()
		go func() {
			wallet.Balance()
			wg.Done()
		}()
	}
	wg.Wait()

	if wallet.Balance() != 513.5 {
		t.Errorf("wallet: concurrent results of deposit don't converge")
	}
}

func TestRaceWithdrawAndDeposit(t *testing.T) {
	t.Parallel()

	wallet := bw.NewWallet(1000)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(2)

		go func() {
			_, err := wallet.Withdraw(1)
			if err != nil {
				t.Errorf("wallet: error in racewithdrawanddepit %v", err)
			}

			wg.Done()
		}()
		go func() {
			_, err := wallet.Deposit(1)
			if err != nil {
				t.Errorf("wallet: error in racewitddrawanddeposit %v", err)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	if wallet.Balance() != 1000 {
		t.Errorf("wallet: concurrent results of withdraw and deposit don't converge")
	}
}
