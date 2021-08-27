# BitcoinWallet
Wallet that holds bitcoins. Wallet allows you to deposit, withdraw money and get balance. Wallet is safe for concurrent use.

##Download 
To download this package write `go get github.com/Dyleme/BitcoinWallet.git` in your cmd

##Usage

###Initialization
To initialize your wallet, in the arguments you should provide amount of BitCoins that you would like to have in your wallet
To do it you have to use struct BitCoin
```go
    var wallet = Wallet{23.4}
```

###Balance, Deposit, Withdraw
These functions are safe for concurrent usage.

####Balance
To get wallet's balance use function `Balance()`

```go
    balance := wallet.Balance()
```

####Deposit, Witdraw
In the functions `Deposit(money BitCoin)` and `Witdraw(moeny BitCoin)` you should provide positive amount of BitCoins,
else them will return not edit balance and error.
```go
    balance, err := wallet.Deposit(23.4) // Deposit 23.4 bitcoins to wallet

    balance, err := wallet.WithDraw(23.4) // Withdraw 23.4 BitCoins from wallet
```


