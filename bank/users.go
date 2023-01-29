package bank

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
