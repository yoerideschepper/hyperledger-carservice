/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
"fmt"
"encoding/json"
"strconv"

"github.com/hyperledger/fabric-chaincode-go/shim"
sc "github.com/hyperledger/fabric-protos-go/peer"

)

// Chaincode is the definition of the chaincode structure.
type Chaincode struct {
}

type User struct{
	ObjectType  	string `json:"docType"`
	UserID 		string `json:"userId"`
	Password 		string `json:"password"`
	CompanyID 		string `json:"companyId"`
	Email 			string `json:"email"`
	Firstname 		string `json:"firstname"`
	Lastname 		string `json:"lastname"`
}

type Ad struct {
	ObjectType string    `json:"docType"`
	AdID       string    `json:"adId"`
	UserID     string    `json:"userId"`
	CreatedOn  string    `json:"createdOn"`
	Title      string    `json:"title"`
	Category   string    `json:"category"`
	From       string    `json:"from"`
	To         string    `json:"to"`
	Text       string    `json:"text"`
	Comments   []Comment `json:"comments"`
}

type Comment struct {
	UserID      string `json:"userId"`
	CommentText string `json:"CommentText"`
}

type Car struct {
	ObjectType  	string 	`json:"docType"`
	LicencePlate 	string 	`json:"LicencePlate"`		//using the plate as UID
	AvailableSeats int 	`json:"AvailableSeats"`
	Make 			string 	`json:"Make"`
	Color 			string 	`json:"Color"`
	Owner			string	`json:"Owner"`
}

type CarpoolRide struct {
	ObjectType  		string 	`json:"docType"`
	CarpoolRideId 		string 	`json:"CarpoolRideId"`
	Car 				string 	`json:"Car"`
	Driver 				string 	`json:"Driver"`
	Destination 		string 	`json:"Destination"`
	DepartureHour 		string 	`json:"DepartureHour"`
	ApprovedPassangers 	[]User 	`json:"ApprovedPassangers"`
}

type ApplicationForRide struct {
	//ApplicationID	string	`json:"ApplicationID"`
	userID 			string	`json:"userID"`
	CarpoolRideId	string	`json:"CarpoolRideId"`
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		panic(err)
	}
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Init()", fcn, params)
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error

	if fn=="query" {
		result, err = query(stub,args)
	} else if fn  == "createUser" {
		result, err = createUser(stub,args)
	} else if fn == "registerAd" {
		result, err = registerAd(stub,args)
	} else if fn == "addComment" {
		result, err = addComment(stub,args)
	} else if fn == "addCarForUser" {
		result, err = addCarForUser(stub,args)
	} else if fn == "removeAd" {
		result, err = removeAd(stub,args)
	} else if fn == "removeCar" {
		result, err = removeCar(stub,args)
	} else if fn == "setCandidateForRide" {
		result, err = setCandidateForRide(stub,args)
	} else if fn == "denyCandidateForRide" {
		result, err = denyCandidateForRide(stub,args)
	} else if fn == "addCarpoolRide" {
		result, err = addCarpoolRide(stub,args)
	} else if fn == "acceptCandidateForRide" {
		result, err = acceptCandidateForRide(stub,args)
	}










	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

func denyCandidateForRide(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments, expecting 1")
	}

	applicationID := args[0]

	deniedCandidate, err := stub.GetState(applicationID)
	if err != nil {
		return "", fmt.Errorf(" failed with error: %s", err)
	}
	if deniedCandidate == nil {
		return "", fmt.Errorf("failed to get application : %s", applicationID)
	}

	error := stub.DelState(applicationID)
	if error != nil {
		return "", fmt.Errorf("failed with error: %s", error)
	}

	return "", nil



}

//{"carpool001", "USER001", "BE-000-001"}
func acceptCandidateForRide(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments, expecting 2")
	}
	CarpoolRideId := args[0]
	candidate 		:= args[1]
	var carpool CarpoolRide
	var car Car
	var strArray[]string
	approvedPassengersArray := carpool.ApprovedPassangers
	avSeats := car.AvailableSeats

	carpoolAsBytes, err := stub.GetState(CarpoolRideId)
	if err != nil {
		return "", fmt.Errorf("failed with error: %s", err)
	}
	if len(approvedPassengersArray) <= avSeats {
		strArray = append(strArray, candidate)
		carpoolAsBytes, _ = json.Marshal(strArray)
		err = stub.PutState(CarpoolRideId, carpoolAsBytes)
		if err != nil {
			return "", fmt.Errorf("failed to put state for ride: %s", CarpoolRideId)
		}
		return string(carpoolAsBytes), nil
	}
	return "", nil
}

func query(stub shim.ChaincodeStubInterface, args []string) (string,error){
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments, expecting key")
	}

	value, err := stub.GetState(args[0])

	if err != nil {
		return "", fmt.Errorf("Failed to get: %s", args[0])
	}

	if value == nil {
		return "", fmt.Errorf("Error: :s", args[0])
	}

	return string(value), nil
}

func createUser(stub shim.ChaincodeStubInterface, args []string) (string,error){
	if len(args) != 6 {
		return "", fmt.Errorf("Incorrect arguments, expecting 6")
	}

	userID := args[0]
	password := args[1]
	companyID := args[2]
	email := args[3]
	firstname := args[4]
	lastname := args[5]

	objectType := "user"
	user := &User{objectType, userID, password, companyID, email, firstname, lastname}
	userJSONasBytes, err := json.Marshal(user)

	if err != nil  {
		return "", fmt.Errorf("Failed to set user: %s", args[0])
	}

	err = stub.PutState(userID, userJSONasBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to put state of user: %s", args[0])
	}

	return string(userJSONasBytes), nil
}

func registerAd(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 8 {
		return "", fmt.Errorf("Incorrect arguments, expecting 8")
	}

	var comments []Comment

	adID := args[0]
	userID := args[1]
	createdOn := args[2]
	title := args[3]
	category := args[4]
	from := args[5]
	to := args[6]
	text := args[7]

	// ==== Check if user with userId already exists ====
	userAsBytes, err := stub.GetState(args[1])
	if err != nil {
		return "", fmt.Errorf("failed with error: %s", err)
	}

	if userAsBytes == nil {
		return "", fmt.Errorf("failed to get user: %s", args[1])
	}

	// ==== Create ad object and marshal to JSON ====
	objectType := "ad"
	ad := &Ad{objectType, adID, userID, createdOn, title, category, from, to, text, comments}
	adJSONasBytes, err := json.Marshal(ad)

	if err != nil {
		return "", fmt.Errorf("failed to set ad: %s", args[0])
	}

	err = stub.PutState(adID, adJSONasBytes)

	if err != nil {
		return "", fmt.Errorf("failed to put state of ad: %s", args[0])
	}

	return string(adJSONasBytes), nil
}

func addComment(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect arguments, expecting 3 (AdID + UserID + Text)")
	}

	adID := args[0]
	userID := args[1]
	commentText := args[2]

	// ==== Check if user with userId already exists ====
	userAsBytes, err := stub.GetState(args[1])
	if err != nil {
		return "", fmt.Errorf("(USER) failed with error: %s", err)
	}

	if userAsBytes == nil {
		return "", fmt.Errorf("failed to get user: %s", args[1])
	}

	// ==== Create comment && add text ====
	var comment Comment
	comment.UserID = userID
	comment.CommentText = commentText

	// ==== Check if ad with adId already exists ====
	adAsBytes, err := stub.GetState(args[0])
	ad := Ad{}
	if err != nil {
		return "", fmt.Errorf("(AD) failed with error: %s", err)
	}

	if adAsBytes == nil {
		return "", fmt.Errorf("failed to get ad: %s", args[0])
	}

	json.Unmarshal(adAsBytes, &ad)
	ad.Comments = append(ad.Comments, comment)

	adAsBytes, _ = json.Marshal(ad)
	err = stub.PutState(adID, adAsBytes)

	if err != nil {
		return "", fmt.Errorf("failed to put state for ad: %s", adID)
	}

	return string(adAsBytes), nil
}

//{"BE-000-001" , "3" , "Audi", " Black"," User001"}
func addCarForUser(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 5 {
		return "", fmt.Errorf("Incorrect arguments, expecting cardetails and ownerId")
	}

	LicencePlate 	:=args[0]
	AvailableSeatsTussen:=args[1]
	AvailableSeats, _ := strconv.Atoi(AvailableSeatsTussen)
	Make 			:=args[2]
	Color 			:=args[3]
	Owner			:=args[4]

	//LicencePlate Check
	carAsBytes, err := stub.GetState(LicencePlate)
	if err != nil {
		return "", fmt.Errorf("failed with error: %s", err)
	}
	if carAsBytes == nil {
		return "", fmt.Errorf("LicencePlate is not registred yet: %s", args[0])
	}

	//New car with owner
	objectType := "car"
	car := &Car{objectType, LicencePlate, AvailableSeats, Make, Color, Owner}
	carJSONasBytes, err := json.Marshal(car)

	if err != nil  {
		return "", fmt.Errorf("Failed to set car with plate: %s", args[0])
	}

	err = stub.PutState(LicencePlate, carJSONasBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to put state of car: %s", args[0])
	}

	return string(carJSONasBytes), nil

}

//{"carpppo001", "BE-000-001" , "UserMarc" , "Ghent", " 11:30"," [Maggie,Frederic]"}
func addCarpoolRide(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 6 {
		return "", fmt.Errorf("Incorrect arguments, expecting 6")
	}
	var user[] User

	CarpoolRideId	:=args[0]
	Car				:=args[1]
	Driver			:=args[2]
	Destination		:=args[3]
	DepartureHour	:=args[4]
	Passangers   := user

	CarpoolRideIdAsBytes, err := stub.GetState(CarpoolRideId)
	if err != nil {
		return "", fmt.Errorf("failed with error: %s", err)
	}
	if CarpoolRideIdAsBytes == nil {
		return "", fmt.Errorf("no carpool found with number: %s", CarpoolRideId)
	}

	objectType := "carpoolride"
	carpoolRide := &CarpoolRide{objectType, CarpoolRideId, Car,Driver,Destination,DepartureHour,	Passangers}
	carpoolJSONasBytes, err := json.Marshal(carpoolRide)

	if err != nil {
		return "", fmt.Errorf("failed to set ad: %s", args[0])
	}

	err = stub.PutState(CarpoolRideId, carpoolJSONasBytes)

	if err != nil {
		return "", fmt.Errorf("failed to put state of carpoolid: %s", CarpoolRideId)
	}

	return string(carpoolJSONasBytes), nil
}


//{"carpool001", "USER001"}
func setCandidateForRide(stub shim.ChaincodeStubInterface, args []string) (string, error)  {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect, only insert your account and your selected ride")
	}

	CarpoolRideId := args[0]
	userID 			:= args[1]



	//check if carpoolRide is listed
	value, err := stub.GetState(CarpoolRideId)

	if err != nil {
		return "", fmt.Errorf("Failed to get: %s", CarpoolRideId)
	}
	if value == nil {
		//creates new application
		var appForRide ApplicationForRide
		appForRide.userID = userID
		appForRide.CarpoolRideId = CarpoolRideId

		applicationforride := &ApplicationForRide{appForRide.userID,appForRide.CarpoolRideId}
		applicationforrideJSONasBytes, err := json.Marshal(applicationforride)
		err = stub.PutState(userID, applicationforrideJSONasBytes)

		if err != nil {
			return "", fmt.Errorf("Failed to set carpool with: %s", CarpoolRideId)
		}

		return string(applicationforrideJSONasBytes), nil
	}

	return "", nil
}

func removeAd(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect, only insert the adID")
	}

	adID :=args[0]

	//Check if ad exists
	adAsBytes, err := stub.GetState(adID)
	if err != nil {
		return "", fmt.Errorf("failed with error: %s", err)
	}
	if adAsBytes == nil {
		return "", fmt.Errorf("No ad was found: %s", adID)
	}

	// delete state of the ad
	error := stub.DelState(adID)
	if error != nil {
		return "", fmt.Errorf("failed with error: %s", error)
	}

	return "", nil
}

func removeCar(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect, only insert the licensePlate")
	}

	licensePlate :=args[0]

	//Check if car exists
	adAsBytes, err := stub.GetState(licensePlate)
	if err != nil {
		return "", fmt.Errorf("failed with error: %s", err)
	}
	if adAsBytes == nil {
		return "", fmt.Errorf("No ad was found: %s", licensePlate)
	}

	// delate state of the car
	error := stub.DelState(licensePlate)
	if error != nil {
		return "", fmt.Errorf("failed with error: %s", error)
	}

	return "", nil
}






