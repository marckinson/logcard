package main

import (
	"errors"
	"fmt"
	"encoding/json"
	

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const   SUPPLIER = "SUPPLIER" 								
const   MT_USER = "MT_USER" 								
const 	CUSTOMER = "CUSTOMER" 								 
const 	CERTIFIER = "CERTIFIER" 							
const	AH = "AH" 										
const   SHIPPING = "SHIPPING"							


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

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

type Log struct { // remplacement de transaction par Log
	PType  		string  `json:"pType"`
	Responsible string  `json:"responsible"`  
	Signature	string 	`json:"signature"`
	VDate 		string   `json:"vDate"` 
	Location  	string   `json:"location"` 
	LType 		string   `json:"ttype"` 
}

type AllParts struct{  
	Parts []string `json:"parts"`
}

type AllPartsDetails struct{ 
	Parts []Part `json:"parts"`
}

type Aircraft struct{ 
	Id   	string `json:"id"` 
	AType  	string  `json:"aType"` 				
	Owner  	string  `json:"owner"`
	Parts 	[]Part `json:"parts"`
}

type AllAircraft struct {
	Airfcrafts []string `json:"aircrafts"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

	var err error
	var parts AllParts
	jsonAsBytes, _ := json.Marshal(parts)
	err = stub.PutState("allParts", jsonAsBytes)
	if err != nil {
		return nil, err
	}	
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
// Run - Our entry point _ Invoke is called when an invoke message is received
// ============================================================================================================================

// Invoke is our entry point to invoke a chaincode function
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
	}
    
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Create new Part of Items _ de façon simplifier c'est la création de la logCard
// ============================================================================================================================

func (t *SimpleChaincode) createPart(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running createPart")

	if len(args) != 8 { // A verifier 
		fmt.Println("Incorrect number of arguments. Expecting 8")
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}

	if args[2] != SUPPLIER{     // A vérifier la syntaxe // RAJOUTER DES CONDITIONS 
		fmt.Println("You are not allowed to create a new part")
		return nil, errors.New("You are not allowed to create a new part") 
	}
	
	if args[2] != AH{     // A vérifier la syntaxe // RAJOUTER DES CONDITIONS 
		fmt.Println("You are not allowed to create a new part")
		return nil, errors.New("You are not allowed to create a new part") 
	}
	
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

// ============================================================================================================================
// Transfer a part = Transfert physique de la part & Transfert de la responsabilité
// ============================================================================================================================
func (t *SimpleChaincode) transferPart_Responsility(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running transferPart_Responsility")

	if len(args) != 6 {
		fmt.Println("Incorrect number of arguments. Expecting 6 (PartId, user, date, location, newResponsible, signature)")
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}

	if args[1] != CUSTOMER { return nil, errors.New("You are not allowed to transfer a part") }
	if args[1] != MT_USER { return nil, errors.New("You are not allowed to transfer a part") }
	if args[1] != SUPPLIER { return nil, errors.New("You are not allowed to transfer a part") }
	if args[1] != AH { return nil, errors.New("You are not allowed to transfer a part") }
	if args[1] != SHIPPING { return nil, errors.New("You are not allowed to transfer a part") }


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

// ============================================================================================================================
// Claim the ownership on a part = Correspond au transfert de propriété
// ============================================================================================================================
func (t *SimpleChaincode)claimOwnershipOnPart(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running claimOwnershipOnPart")

	if len(args) != 4 {
		fmt.Println("Incorrect number of arguments. Expecting 4 (PartId, user, date, location)")
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	if args[1] != AH { return nil, errors.New("You are not allowed to claimOwnership on a Part") } 
	if args[1] != SUPPLIER { return nil, errors.New("You are not allowed to claimOwnership on a Part") } 
	if args[1] != CUSTOMER { return nil, errors.New("You are not allowed to claimOwnership on a Part") } 


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


// ============================================================================================================================
// Query - read a variable from chaincode state - (aka read) _ As the name implies, Query is called whenever you query your chaincode's state. 
// Queries do not result in blocks being added to the chain, and you cannot use functions like PutState inside of Query or any helper functions it calls. 
// You will use Query to read the value of your chaincode state's key/value pairs.
// ============================================================================================================================


// Query is our entry point for queries Marckinson

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    // Handle different functions
    if function == "getPart" {return t.getPart(stub, args[0])}
	if function == "getAllParts" { return t.getAllParts(stub, args[0]) }
	if function == "getAllPartsDetails" { return t.getAllPartsDetails(stub, args[0]) }
	
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
}



// ============================================================================================================================
// Get Part Details
// ============================================================================================================================
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

// ============================================================================================================================
// Get All Parts 
// ============================================================================================================================
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


// ============================================================================================================================
// Get All Parts Details for a specific user
// ============================================================================================================================
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


