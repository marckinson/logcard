/*
Dans cette première implémentation les LogCard correspondent aux Parts
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//===============================================================================================================================
//	 Participant types 
//==============================================================================================================================
const   SUPPLIER 	= "Suppliers" 								
const   MT_USER 	= "MRO_User" 								
const 	CUSTOMER 	= "customer" 								 
const 	CERTIFIER 	= "EASAA_FAA" 							
const	AH 			= "AirbusHelicopter" 										
const   SHIPPING 	= "shipping_company"							

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type SimpleChaincode struct {
}

//==============================================================================================================================
//	Part - Defines the structure for a part object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Part struct { // Part et eLogcard sont regroupés dans cette première version
	Id   		string  `json:"id"` 					// Concaténation des deux PN et SN
	PN			string 	`json:"pn"` 					// Part Number
	SN 			string 	`json:"sn"` 					// Serial Number 
	PType  		string  `json:"pType"` 					// Part Type (Voir excel Jean-Guillaume) voir où l'on doit décrire ces type là.
	Owner  		string  `json:"owner"` 					// Propriétaire de la pièce 
	Responsible string  `json:"responsible"` 			// Responsable à l'instant T de la pièce
	Signature	string 	`json:"signature"`				// Process de signature manuel ??
	Logs        []Log 	`json:"logs"` 					// Correspondent au Log 
}

//================================================================================================================================
//	Part - Defines the structure for a log object. It represents transactions for a part, states changes, maintenance tasks, etc..
//================================================================================================================================
type Log struct { // remplacement de transaction par Log
	PType  		string  `json:"pType"`
	Responsible string  `json:"responsible"`  
	Signature	string 	`json:"signature"`
	VDate 		string   `json:"vDate"` 
	Location  	string   `json:"location"` 
	LType 		string   `json:"ttype"` 
}

//================================================================================================================================
//	Part - Defines the structure for a AllParts object. It represents all the parts ..
//================================================================================================================================
type AllParts struct{  
	Parts []string `json:"parts"`
}

//================================================================================================================================
//	Part - Defines the structure for a AllPartsDetails object. It represents all the parts with their content.
//================================================================================================================================
type AllPartsDetails struct{ 
	Parts []Part `json:"parts"`
}

//================================================================================================================================
//	Part - Defines the structure for a Aircraft object. It represents the AirCraft.
//================================================================================================================================
type Aircraft struct{ 
	Id   	string `json:"id"` 
	AType  	string  `json:"aType"` 				
	Owner  	string  `json:"owner"`
	Parts 	[]Part `json:"parts"`
}
//================================================================================================================================
//	Part - Defines the structure for a AllAircraft object. It represents all the aircrafts.
//================================================================================================================================
type AllAircraft struct {
	Airfcrafts []string `json:"aircrafts"`
}

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    
	if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }
	var err error 												// déclaration de la variable err de type error             
	
	// Initialize the chaincode
	// Qu'est ce qu'on doit mettre ici ?????
	
	
	// Test du Network (Read, write)
	// Write the state to the ledger
	//err = stub.PutState("allParts", []byte(args[0]))  			//
	//if err != nil {
	//	return nil, err
	//}	
		
	var parts AllParts 											// array of string a la place de AllParts	?? // déclaration de la variable parts de type AllParts 
	jsonAsBytes, _ := json.Marshal(parts)   					// marshal de cet asset AllParts 
	err = stub.PutState("allParts", jsonAsBytes)  				//
	if err != nil {
		return nil, err
	}	
	
	
	
	// Fini 
	
	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke is our entry point to invoke a chaincode function
// Run - Our entry point _ Invoke is called when an invoke message is received
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "createPart" {									  
		return t.createPart(stub, args)
	} else if function == "transferPart_Responsility" {
		return t.transferPart_Responsility(stub, args)
	} else if function == "claimOwnershipOnPart" {
		return t.claimOwnershipOnPart(stub, args)
	} else if function == "write" { 						// à enlever plus tard 
		return t.write(stub, args)
	}
    
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation")
}

// ================================================================================================================================
// Functions Handled by Invoke
// ================================================================================================================================
//=================================================================================================================================
//	 Create Function
//=================================================================================================================================
// ================================================================================================================================
// Creation of the Part (creation of the eLogcard)
// ================================================================================================================================

func (t *SimpleChaincode) createPart(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Running createPart")

	if len(args) != 8 { 																// A verifier 
		fmt.Println("Incorrect number of arguments. Expecting 8")
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}

/*	
	if args[2] != SUPPLIER{     														// A vérifier la syntaxe 
		fmt.Println("You are not allowed to create a new part")
		return nil, errors.New("You are not allowed to create a new part") 
	}

	if args[2] != AH{     																// Chercher la syntaxe OU 
		fmt.Println("You are not allowed to create a new part")
		return nil, errors.New("You are not allowed to create a new part") 
	}
*/
	
	var err error
	
	var pt Part
	pt.Id 			= args[0]
	pt.PN			= args[1]
	pt.SN			= args[2]
	pt.PType		= args[3]
	pt.Owner		= args[4]
	pt.Responsible  = args[5]
	pt.Signature 	= ""

	var tx Log
	tx.VDate		= args[6]
	tx.Location 	= args[7]
	tx.LType 		= "CREATE"
	tx.PType 		= pt.PType
	tx.Responsible 	= pt.Responsible
	tx.Signature 	= pt.Signature 

	pt.Logs = append(pt.Logs, tx)
	
	//Commit part to ledger
	fmt.Println("createPart Commit Part To Ledger");
	ptAsBytes, _ := json.Marshal(pt)
	err = stub.PutState(pt.Id, ptAsBytes)	
	if err != nil {
		return nil, err
	}
	
	//Update All Parts Array
	allPAsBytes, err := stub.GetState("allParts")
	if err != nil {
		return nil, errors.New("Failed to get all Parts")
	}
	var allp AllParts
	err = json.Unmarshal(allPAsBytes, &allp)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal all Parts")
	}
	allp.Parts = append(allp.Parts,pt.Id)

	allPuAsBytes, _ := json.Marshal(allp)
	err = stub.PutState("allParts", allPuAsBytes)	
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

//=================================================================================================================================
//	 Transfer Functions
//=================================================================================================================================
// ================================================================================================================================
// Transfer a part = Transfert physique de la part & Transfert de la responsabilité
// ================================================================================================================================
func (t *SimpleChaincode) transferPart_Responsility(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running transferPart_Responsility")

	if len(args) != 6 {
		fmt.Println("Incorrect number of arguments. Expecting 6 (PartId, user, date, location, newResponsible, signature)")
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}

//	if args[1] != CUSTOMER { return nil, errors.New("You are not allowed to transfer a part") }
/*
	if args[1] != MT_USER { return nil, errors.New("You are not allowed to transfer a part") }
	if args[1] != SUPPLIER { return nil, errors.New("You are not allowed to transfer a part") }
	if args[1] != AH { return nil, errors.New("You are not allowed to transfer a part") }
	if args[1] != SHIPPING { return nil, errors.New("You are not allowed to transfer a part") }
*/

	//Update Part data
	pAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get Part #" + args[0])
	}
	var prt Part
	err = json.Unmarshal(pAsBytes, &prt)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal Part #" + args[0])
	}
	prt.Responsible = args[4]
	prt.Signature = args[5]

	var tx Log
	tx.VDate		= args[2]
	tx.Location 	= args[3]
	tx.LType 		= "TRANSFERT"
	tx.PType 		= prt.PType
	tx.Responsible  = prt.Responsible
	tx.Signature 	= prt.Signature 

	prt.Logs = append(prt.Logs, tx)

	//Commit updates part to ledger
	fmt.Println("transferPart Commit Updates To Ledger");
	ptAsBytes, _ := json.Marshal(prt)
	err = stub.PutState(prt.Id, ptAsBytes)	
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//=================================================================================================================================
//	 Claim Functions
//=================================================================================================================================
// ================================================================================================================================
// Claim the ownership on a part = Correspond au transfert de propriété
// ================================================================================================================================
func (t *SimpleChaincode)claimOwnershipOnPart(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running claimOwnershipOnPart")

	if len(args) != 4 {
		fmt.Println("Incorrect number of arguments. Expecting 4 (PartId, user, date, location)")
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// if args[1] != AH { return nil, errors.New("You are not allowed to claimOwnership on a Part") } 

/*
	if args[1] != SUPPLIER { return nil, errors.New("You are not allowed to claimOwnership on a Part") } 
	if args[1] != CUSTOMER { return nil, errors.New("You are not allowed to claimOwnership on a Part") } 
*/

	//Update Part owner
	bAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get Part #" + args[0])
	}
	var prt Part
	err = json.Unmarshal(bAsBytes, &prt)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal Part #" + args[0])
	}
	prt.Owner = args[1]

	var tx Log
	tx.VDate		= args[2]
	tx.Location 	= args[3]
	tx.LType 		= "CLAIM"
	tx.PType 		= prt.PType
	tx.Responsible  = prt.Responsible
	tx.Signature 	= prt.Signature 

	prt.Logs = append(prt.Logs, tx)

	//Commit updates part to ledger
	fmt.Println("claimOwnershipOnPart Commit Updates To Ledger");
	ptAsBytes, _ := json.Marshal(prt)
	err = stub.PutState(prt.Id, ptAsBytes)	
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ===========================================================================================================================================
// Query - read a variable from chaincode state - (aka read) _ As the name implies, Query is called whenever you query your chaincode's state. 
// Queries do not result in blocks being added to the chain. 
// You cannot use functions like PutState inside of Query or any helper functions it calls. 
// ============================================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    // Handle different functions
    if function == "getPart" {return t.getPart(stub, args[0])}
	if function == "getAllParts" { return t.getAllParts(stub, args[0]) }
	if function == "getAllPartsDetails" { return t.getAllPartsDetails(stub, args[0]) }
	if function == "read" {	return t.read(stub, args)}
	
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
}

// =============================================================================================================================================
// Functions Handled by Query
// =============================================================================================================================================


// =================
// Get Part Details
// =================
func (t *SimpleChaincode) getPart(stub shim.ChaincodeStubInterface, partId string)([]byte, error){
	
	fmt.Println("Start find Part")
	fmt.Println("Looking for Part #" + partId);

	//get the part index
	pAsBytes, err := stub.GetState(partId)
	if err != nil {
		return nil, errors.New("Failed to get Part #" + partId)
	}
	
	return pAsBytes, nil
	
}

// ===============
// Get All Parts 
// ===============
func (t *SimpleChaincode) getAllParts(stub shim.ChaincodeStubInterface, user string)([]byte, error){
	
	fmt.Println("Start find getAllParts ")
	fmt.Println("Looking for All Parts " + user);

	//get the AllParts index
	allBAsBytes, err := stub.GetState("allParts")
	if err != nil {
		return nil, errors.New("Failed to get all Parts")
	}

	var res AllParts
	err = json.Unmarshal(allBAsBytes, &res)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal all Parts")
	}

	var rab AllParts

	for i := range res.Parts{

		sbAsBytes, err := stub.GetState(res.Parts[i])
		if err != nil {
			return nil, errors.New("Failed to get Part")
		}
		var sb Part
		json.Unmarshal(sbAsBytes, &sb)

		if(sb.Owner == user || user == CERTIFIER || user == AH) {
			rab.Parts = append(rab.Parts,sb.Id); 
		}
	}
	rabAsBytes, _ := json.Marshal(rab)

	return rabAsBytes, nil
}


// ============================================
// Get All Parts Details for a specific user
// ============================================
func (t *SimpleChaincode) getAllPartsDetails(stub shim.ChaincodeStubInterface, user string)([]byte, error){
	
	fmt.Println("Start find getAllPartsDetails ")
	fmt.Println("Looking for All Parts Details " + user);

	//get the AllParts index
	allBAsBytes, err := stub.GetState("allParts")
	if err != nil {
		return nil, errors.New("Failed to get all Parts")
	}

	var res AllParts
	err = json.Unmarshal(allBAsBytes, &res)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal all Parts")
	}

	var rab AllPartsDetails

	for i := range res.Parts{

		sbAsBytes, err := stub.GetState(res.Parts[i])
		if err != nil {
			return nil, errors.New("Failed to get Part")
		}
		var sb Part
		json.Unmarshal(sbAsBytes, &sb)

		if(sb.Owner == user) {
			sb.Logs = nil
			sb.Signature = ""
			rab.Parts = append(rab.Parts,sb); 
		}

	}

	rabAsBytes, _ := json.Marshal(rab)

	return rabAsBytes, nil
}



// A enlever
// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// A enlever
// write - invoke function to write key/value pair

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] 												//rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) 					//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}