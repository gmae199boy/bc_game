package main

import (
	"encoding/json"
	"fmt"
	// "strconv"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type User struct {
	Id string
	Name string
	WinRate int
	LoseLate int
}

func checkLen(logger *shim.ChaincodeLogger, expected int, args []string) error {
	if len(args) < expected {
		mes := fmt.Sprintf(
			"not enough number of arguments: %d given, %d expected",
			len(args),
			expected,
		)
		logger.Warning(mes)
		return errors.New(mes)
	}
	return nil
}

// GameCC example simple Chaincode implementation
type GameCC struct {
}

// Init initializes chaincode
// ===========================
func (g *GameCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger := shim.NewLogger("Game !")
	logger.Info("chaincode initalized")
	return shim.Success([]byte{})
}

// Invoke - Our entry point for Invocations
// ========================================
func (g *GameCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger := shim.NewLogger("Game !")
	function, args := stub.GetFunctionAndParameters()
	logger.Info("function name = %s", function)
	fmt.Println("invoke is running " + function)
	// Handle different functions
	switch function {
	case "addUser":
		return g.addUser(stub, args)
	case "readUserInfo":
		return g.readUserInfo(stub, args)
	case "readUserList":
		//create a new marble
		return g.readUserList(stub, args)
	default:
		//error
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

func (s *GameCC) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("fail!")
	}

	UserAsBytes, _ := json.Marshal(args[0])
	key, err := stub.CreateCompositeKey("User", []string{args[0]})
	if err != nil {
		return shim.Error("err")
	}

	err = stub.PutState(key, UserAsBytes)
	if err != nil {
		return shim.Error("err")
	}

	return shim.Success(nil)
}

func (g *GameCC) readUserInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	key, err := stub.CreateCompositeKey("User", []string{args[0]})
	if err != nil {
		return shim.Error("err")
	}
	UserAsBytes, _ := stub.GetState(key)

	return shim.Success(UserAsBytes)
}

func (g *GameCC) readUserList(stub shim.ChaincodeStubInterface) pb.Response {
	iter, err := stub.GetStateByPartialCompositeKey("User", []string{})
	if err != nil {
		return shim.Error("err")
	}

	defer iter.Close()

	users := []*User{}

	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return shim.Error("err")
		}
		user := new(User)
		err = json.Unmarshal(kv.Value, user)
		if err != nil {
			return shim.Error("err")
		}
		users = append(users, user)
	}

	UserAsBytes, err := json.Marshal(users)

	return shim.Success(UserAsBytes)
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(GameCC))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
