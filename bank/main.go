package bank

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type BankServicer interface {
	transferMoneyFromTo(from, to string, amount float64) error
	createUserWithBalance(firstName, lastName, socialSecurityNumber string, balance float64) (string, error)
	removeUserByBankAccountNumber(bankAccountNumber string) error
	findUserByBankAccountNumber(bankAccountNumber string) (*User, error)
}

// This comes from the server stubs generated; we did the process backwards; we should have generated the proto file first
// and then generated the functions from that, with implementation logic; this will need to be refactored
type Service struct {
	users *Users
	UnimplementedBankServer
}

func NewBankService() *Service {
	return &Service{users: NewBankUsers()}
}

func (b Service) TransferMoneyFromTo(ctx context.Context, request *TransferMoneyFromToRequest) (*TransferMoneyFromToResponse, error) {
	err := b.transferMoneyFromTo(request.From, request.To, float64(request.Amount))
	if err != nil {
		return nil, err
	}
	return &TransferMoneyFromToResponse{}, nil
}

func (b Service) CreateUserWithBalance(ctx context.Context, request *CreateUserWithBalanceRequest) (*CreateUserWithBalanceResponse, error) {
	bankAccountNumber, err := b.createUserWithBalance(request.FirstName, request.LastName, request.SocialSecurityNumber, float64(request.Balance))
	if err != nil {
		return nil, err
	}
	return &CreateUserWithBalanceResponse{BankAccountNumber: bankAccountNumber}, nil
}

func (b Service) RemoveUserByBankAccountNumber(ctx context.Context, request *RemoveUserByBankAccountNumberRequest) (*RemoveUserByBankAccountNumberResponse, error) {
	err := b.removeUserByBankAccountNumber(request.BankAccountNumber)
	if err != nil {
		return nil, err
	}
	return &RemoveUserByBankAccountNumberResponse{}, nil
}

func (b Service) FindUserByBankAccountNumber(ctx context.Context, request *FindUserByBankAccountNumberRequest) (*FindUserByBankAccountNumberResponse, error) {
	user, err := b.findUserByBankAccountNumber(request.BankAccountNumber)
	if err != nil {
		return nil, err
	}
	return &FindUserByBankAccountNumberResponse{FirstName: user.firstName, LastName: user.lastName, SocialSecurityNumber: user.socialSecurityNumber, Balance: float32(int32(user.balance))}, nil
}

// I cannot put a finger on it, but something seems off about this
func (b Service) transferMoneyFromTo(from, to string, amount float64) error {
	// first we have to check that the users exist
	fromUser := b.users.FindUserByBankAccountNumber(from)
	toUser := b.users.FindUserByBankAccountNumber(to)
	if fromUser == nil || toUser == nil {
		return fmt.Errorf("user not found. Aborting transfer")
	}
	// now we can check that the from user has enough money
	if fromUser.GetBalance() < amount {
		return fmt.Errorf("not enough money in account. Aborting transfer")
	}
	// now we can transfer the money
	fromUser.balance -= amount
	toUser.balance += amount
	return nil
}

func (b Service) createUserWithBalance(firstName, lastName, socialSecurityNumber string, balance float64) (string, error) {
	b.users.AddUser(NewUser(firstName, lastName, socialSecurityNumber, balance))
	// This will return the Bank ID generated for the user added to the list
	return b.users.users[len(b.users.users)-1].bankAccountNumber, nil
}

//func (b Service) createUserWithBalance(firstName, lastName, socialSecurityNumber string, balance float64) (string, error) {
//	b.users.AddUser(NewUser(firstName, lastName, socialSecurityNumber, balance))
//	// This will return the Bank ID generated for the user added to the list
//	return b.users.users[len(b.users.users)-1].bankAccountNumber, nil
//}

func (b Service) removeUserByBankAccountNumber(bankAccountNumber string) error {
	user := b.users.FindUserByBankAccountNumber(bankAccountNumber)
	if user == nil {
		return fmt.Errorf("user not found")
	}
	b.users.RemoveUser(user)
	return nil
}

func (b Service) findUserByBankAccountNumber(bankAccountNumber string) (*User, error) {
	user := b.users.FindUserByBankAccountNumber(bankAccountNumber)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

type User struct {
	firstName, lastName, socialSecurityNumber, bankAccountNumber string
	balance                                                      float64
}

func NewUser(firstName, lastName, socialSecurityNumber string, balance float64) *User {
	// TODO: Add a check to make sure balance is not negative
	bankID := generateRandomBankAccountNumber()
	return &User{firstName: firstName, lastName: lastName, socialSecurityNumber: socialSecurityNumber, bankAccountNumber: bankID, balance: balance}
}

func (u *User) GetBalance() float64 {
	return u.balance
}

type Users struct {
	users []*User
}

func NewBankUsers() *Users {
	return &Users{users: []*User{}}
}

func (b *Users) AddUser(user *User) {
	b.users = append(b.users, user)
}

func (b *Users) RemoveUser(user *User) {
	// I think there is a bug here with the append [i+1:] part
	for i, u := range b.users {
		if u.bankAccountNumber == user.bankAccountNumber {
			b.users = append(b.users[:i], b.users[i+1:]...)
		}
	}
}

func (b *Users) FindUserByBankAccountNumber(bankAccountNumber string) *User {
	for _, user := range b.users {
		if user.bankAccountNumber == bankAccountNumber {
			return user
		}
	}
	return nil
}

// Helper funcs
func generateRandomBankAccountNumber() string {
	return uuid.New().String()
}
