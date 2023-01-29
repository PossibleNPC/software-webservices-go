package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"software-webservices-go/bank"
)

// TODO: Migrate this to a CLI tool to deal with interaction with the gRPC server
func main() {
	// This is where we are creating the client to talk to our gRPC server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(":8080", opts...)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// We are going to create the client
	client := bank.NewBankClient(conn)
	accountDetails, err := client.FindUserByBankAccountNumber(context.Background(), &bank.FindUserByBankAccountNumberRequest{
		BankAccountNumber: "92489a16-114d-42ea-9779-da43b9b4e958",
	})
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	log.Printf("User %s found successfully\n", accountDetails)
	//accountNumber, err := client.CreateUserWithBalance(context.Background(), &bank.CreateUserWithBalanceRequest{
	//	FirstName:            "John",
	//	LastName:             "Doe",
	//	SocialSecurityNumber: "123-45-6789",
	//	Balance:              100,
	//})
	//if err != nil {
	//	log.Fatalf("Failed to create user: %v", err)
	//}
	//log.Printf("User %s created successfully\n", accountNumber)
}
