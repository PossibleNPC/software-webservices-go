package bank

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// This comes from the server stubs generated; we did the process backwards; we should have generated the proto file first
// and then generated the functions from that, with implementation logic; this will need to be refactored
type Service struct {
	users *Users
	UnimplementedBankServer
}

func NewBankService() *Service {
	return &Service{users: NewBankUsers()}
}

// Implementes the gRPC interface UnimplementedBankServer
func (b Service) TransferMoneyFromTo(ctx context.Context, request *TransferMoneyFromToRequest) (*TransferMoneyFromToResponse, error) {
	err := b.transferMoney(request.From, request.To, float64(request.Amount))
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

// Internal functions that might have a different interface than the gRPC interface
// I cannot put a finger on it, but something seems off about this with calling these functions on the Service struct
func (b Service) transferMoney(from, to string, amount float64) error {
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

// Helper funcs
func generateRandomBankAccountNumber() string {
	return uuid.New().String()
}
