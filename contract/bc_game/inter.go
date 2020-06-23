package gametransfer

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"time"
)

type User struct {
	Id string
	Name string
	WinRate int
	LoseLate int
}

type GameTransfer interface {
	AddUser(shim.ChaincodeStubInterface, *User) error
	ReadUserInfo(shim.ChaincodeStubInterface) error
	CheckUser(shim.ChaincodeStubInterface, string) (bool, error)
	ListUsers(shim.ChaincodeStubInterface) ([]*User, error)

	CalRate(shim.ChaincodeStubInterface, string) error
}