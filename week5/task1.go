package main

import (
	"errors"
	"fmt"
)

type Notifier interface {
	Notify(msg string) error
}

type Email struct {
	Address string
}

func (e Email) Notify(msg string) error {
	fmt.Printf("Email to %s: %s\n", e.Address, msg)
	return nil
}

type SMS struct {
	Phone string
}

func (s SMS) Notify(msg string) error {
	if len(msg) > 20 {
		return errors.New("SMS message is longer than 20 characters")
	}

	fmt.Printf("SMS to %s: %s\n", s.Phone, msg)
	return nil
}

type Telegram struct {
	Chat string
	Sent int
}

func (t *Telegram) Notify(msg string) error {
	t.Sent++
	fmt.Printf("Telegram to %s: %s\n", t.Chat, msg)
	return nil
}

// Compile-time checks
var _ Notifier = Email{}
var _ Notifier = SMS{}
var _ Notifier = (*Telegram)(nil)

func Broadcast(ns []Notifier, msg string) []error {
	var errs []error

	for _, n := range ns {
		if err := n.Notify(msg); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func main() {
	tg := &Telegram{Chat: "@student"}

	notifiers := []Notifier{
		Email{Address: "student@kbtu.kz"},
		SMS{Phone: "+77001234567"},
		tg,
	}

	errs := Broadcast(notifiers, "Hello from Go!")

	fmt.Println("Errors:", errs)
	fmt.Println("Telegram sent:", tg.Sent)

	//errs = Broadcast(notifiers, "This message is definitely longer than twenty characters")
	//fmt.Println("Errors:", errs)
	//fmt.Println("Telegram sent:", tg.Sent)
}