package bitcoinwallet

import (
	"fmt"
	"strconv"
	"testing"
)

var flagDepositTest = []struct {
	startMoney   float64
	addableMoney float64
	out          float64
}{
	{12.5, 124, 136.5},
	{0, 0.34, 0.34},
	{-34, 43.54, 9.54},
}

func TestWallet_Deposit(t *testing.T) {
	for _, tt := range flagDepositTest {
		t.Run(fmt.Sprintf("Start %v, add %v", tt.startMoney, tt.addableMoney), func(t *testing.T) {
			w := Wallet{balance: tt.startMoney}
			if result := w.Deposit(tt.addableMoney); result != tt.out {
				t.Errorf("Result want %v, result get %v", tt.out, result)
			}
		})
	}
}

var flagWithdrawTest = []struct {
	startMoney  float64
	pickedMoney float64
	out         float64
}{
	{136.5, 124, 12.5},
	{0.34, 0.34, 0},
	{9.54, 43.54, -34},
}

func TestWallet_Withdraw(t *testing.T) {
	for _, tt := range flagWithdrawTest {
		t.Run(fmt.Sprintf("Start %v, add %v", tt.startMoney, tt.pickedMoney), func(t *testing.T) {
			w := Wallet{balance: tt.startMoney}
			if result := w.Withdraw(tt.pickedMoney); result != tt.out {
				t.Errorf("Result want %v, result get %v", tt.out, result)
			}
		})
	}
}

var flagBalanceTest = []struct {
	in  float64
	out float64
}{
	{12.5, 12.5},
	{0.34, 0.34},
	{-9.54, -9.54},
}

func TestWallet_Balance(t *testing.T) {
	for _, tt := range flagBalanceTest {
		t.Run(strconv.FormatFloat(tt.in, 'E', -1, 32), func(t *testing.T) {
			w := Wallet{balance: tt.in}
			if result := w.Balance(); result != tt.out {
				t.Errorf("Result want %v, result get %v", tt.out, result)
			}
		})
	}
}
