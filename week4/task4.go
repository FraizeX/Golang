package main

import (
	"errors"
	"fmt"
)

type BankAccount struct {
	Owner   string
	balance float64
}

func (a *BankAccount) Deposit(amount float64) {
	a.balance += amount
}

func (a *BankAccount) Withdraw(amount float64) error {
	if amount > a.balance {
		return errors.New("not enough money")
	}

	a.balance -= amount
	return nil
}

func (a *BankAccount) Balance() float64 {
	return a.balance
}

func main() {
	accounts := []BankAccount{
		{Owner: "Patrick", balance: 1000},
		{Owner: "Dexter", balance: 2000},
		{Owner: "Tyler", balance: 3000},
	}

	accounts[0].Deposit(500)

	err := accounts[1].Withdraw(500)
	if err != nil {
		fmt.Println("Error:", err)
	}

	for i := range accounts {
		interest := accounts[i].Balance() * 0.05
		accounts[i].Deposit(interest)
	}

	for _, account := range accounts {
		fmt.Printf(
			"%s: %.2f\n",
			account.Owner,
			account.Balance(),
		)
	}
}