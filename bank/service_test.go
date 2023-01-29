package bank

import (
	"context"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"testing"
)

// I don't know if this is a good convention for global test variables in Go or not, but I don't know how to really use
// context, set values in the context, and then I don't understand the "hooks" being offered from godog.
// TODO: Do we migrate the service to the context, or is this something that goes into an init?
var (
	bankService = NewBankService()
)

type bankUserCtxKey struct{}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func aUserProvidesTheirFirstNameLastNameSocialSecurityNumberAndBalance(ctx context.Context, arg1, arg2, arg3 string, arg4 int) (context.Context, error) {
	return context.WithValue(ctx, bankUserCtxKey{}, &User{firstName: arg1, lastName: arg2, socialSecurityNumber: arg3, balance: float64(arg4)}), nil
}

func theUserIsCreated(ctx context.Context) (context.Context, error) {
	// An interesting part of my research indicates that struct validation might need to be handled in a step where
	// the creation of the struct is done otherwise subtle bugs can be introduced.
	_, ok := ctx.Value(bankUserCtxKey{}).(*User)
	if !ok {
		return ctx, errors.New("expected user to be created, but was nil")
	}
	return ctx, nil
}

func hasABalanceOf(ctx context.Context, arg1 int) error {
	user, ok := ctx.Value(bankUserCtxKey{}).(*User)
	if !ok {
		return errors.New("there is no bank user available")
	}
	if user.GetBalance() != float64(arg1) {
		return fmt.Errorf("expected %.2f, got %.2f", float64(arg1), user.GetBalance())
	}
	return nil
}

func aUser(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, bankUserCtxKey{}, &User{firstName: "John", lastName: "Doe", socialSecurityNumber: "123-45-6789", balance: 100.00})
	// TODO: Research if this is a good way to do this, or do I need to do this kind of a test in a Unit Test versus Service Test?
	// I think this is going to be answered with "what is your intent of the test doing?"
	ctx = context.WithValue(ctx, "firstName", "John")
	ctx = context.WithValue(ctx, "lastName", "Doe")
	ctx = context.WithValue(ctx, "socialSecurityNumber", "123-45-6789")
	ctx = context.WithValue(ctx, "balance", 100.00)

	return ctx
}

func theUserIsAddedToTheBankUsers(ctx context.Context) (context.Context, error) {
	_, ok := ctx.Value(bankUserCtxKey{}).(*User)
	if !ok {
		return ctx, errors.New("there is no bank user available")
	}

	ctx = context.WithValue(ctx, "bankUsers", &Users{users: []*User{}})

	ctx.Value("bankUsers").(*Users).AddUser(ctx.Value(bankUserCtxKey{}).(*User))

	return ctx, nil
}

func theUserIsInTheBankUsers(ctx context.Context) error {
	users, ok := ctx.Value("bankUsers").(*Users)
	if !ok {
		return fmt.Errorf("expected bank users to be present, but none are")
	}

	for _, user := range users.users {
		if user.firstName == ctx.Value("firstName") && user.lastName == ctx.Value("lastName") && user.socialSecurityNumber == ctx.Value("socialSecurityNumber") && user.balance == ctx.Value("balance") {
			return nil
		}
	}
	return fmt.Errorf("expected user to be in bank users, but was not")
}

func theUserIsRemovedFromTheBankUsers(ctx context.Context) {
	ctx.Value("bankUsers").(*Users).RemoveUser(ctx.Value(bankUserCtxKey{}).(*User))
}

func theUserIsNotInTheBankUsers(ctx context.Context) error {
	users, ok := ctx.Value("bankUsers").(*Users)
	if !ok {
		return fmt.Errorf("expected bank users to be present, but none are")
	}
	for _, user := range users.users {
		if user.firstName == ctx.Value("firstName") && user.lastName == ctx.Value("lastName") && user.socialSecurityNumber == ctx.Value("socialSecurityNumber") && user.balance == ctx.Value("balance") {
			return fmt.Errorf("expected user to not be in bank users, but got %v", user)
		}
	}
	return nil
}

func aUserWithABalanceOf(ctx context.Context, arg1 int) (context.Context, error) {
	bankId, err := bankService.createUserWithBalance("John", "Doe", "123-45-6789", float64(arg1))
	if err != nil {
		return nil, errors.New("expected user to be created, but service returned nil")
	}
	ctx = context.WithValue(ctx, "bankAccount", bankId)
	return ctx, nil
}

func anotherUserWithABalanceOf(ctx context.Context, arg1 int) (context.Context, error) {
	bankId, err := bankService.createUserWithBalance("Jane", "Doe", "987-65-4321", float64(arg1))
	if err != nil {
		return nil, errors.New("expected user to be created, but service returned nil")
	}
	ctx = context.WithValue(ctx, "bankAccount2", bankId)
	return ctx, nil
}

func theFirstUserTransfersToTheSecondUser(ctx context.Context, arg1 int) (context.Context, error) {
	// TODO: split these into variables to easily communicate intent to the reader
	err := bankService.transferMoney(ctx.Value("bankAccount").(string), ctx.Value("bankAccount2").(string), float64(arg1))
	if err != nil {
		return nil, fmt.Errorf("expected transfer to be successful, but it was not. %s", err.Error())
	}
	return ctx, nil
}

func theFirstUserHasABalanceOf(ctx context.Context, arg1 int) (context.Context, error) {
	user, _ := bankService.findUserByBankAccountNumber(ctx.Value("bankAccount").(string))
	if user.GetBalance() != float64(arg1) {
		return nil, fmt.Errorf("expected %.2f, got %.2f", float64(arg1), user.GetBalance())
	}
	return ctx, nil
}

func theSecondUserHasABalanceOf(ctx context.Context, arg1 int) (context.Context, error) {
	user, _ := bankService.findUserByBankAccountNumber(ctx.Value("bankAccount2").(string))
	if user.GetBalance() != float64(arg1) {
		return nil, fmt.Errorf("expected %.2f, got %.2f", float64(arg1), user.GetBalance())
	}
	return ctx, nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a user provides their first name "([^"]*)", last name "([^"]*)", social security number "([^"]*)", and balance (\d+)$`, aUserProvidesTheirFirstNameLastNameSocialSecurityNumberAndBalance)
	ctx.Step(`^the user is created$`, theUserIsCreated)
	ctx.Step(`^has a balance of (\d+)$`, hasABalanceOf)
	ctx.Step(`^a user$`, aUser)
	ctx.Step(`^the user is added to the bank users$`, theUserIsAddedToTheBankUsers)
	ctx.Step(`^the user is in the bank users$`, theUserIsInTheBankUsers)
	ctx.Step(`^the user is removed from the bank users$`, theUserIsRemovedFromTheBankUsers)
	ctx.Step(`^the user is not in the bank users$`, theUserIsNotInTheBankUsers)
	ctx.Step(`^a user with a balance of (\d+)$`, aUserWithABalanceOf)
	ctx.Step(`^another user with a balance of (\d+)$`, anotherUserWithABalanceOf)
	ctx.Step(`^the first user transfers (\d+) to the second user$`, theFirstUserTransfersToTheSecondUser)
	ctx.Step(`^the first user has a balance of (\d+)$`, theFirstUserHasABalanceOf)
	ctx.Step(`^the second user has a balance of (\d+)$`, theSecondUserHasABalanceOf)
}
