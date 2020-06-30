package main

import (
	"fmt"
	"net/smtp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// smtpServer data to smtp server.
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server.
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}



func sendMail(stub shim.ChaincodeStubInterface, args []string) (string, error){

	if(len(args) != 1){
		return "", fmt.Errorf("Incorrect Argument, expecting an emailId")
	}

	value, err := stub.GetState(args[0]);
	if err!=nil {
		return "", fmt.Errorf("Counldn't retrieve the message: %s",err)
	}

	fmt.Println("Will Send Mail To: " + args[0] + " Message -- " + string(value) );

	// Sender data.
	from := "ratikjindal.dad99@gmail.com"
	password := "cc-c-_-"

	// Receiver email address.
	to := []string{args[0]}

	// smtp server configuration.
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}

	// Message.
	message := []byte(value)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	// Sending email.
	err = smtp.SendMail(smtpServer.Address(), auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("Counld not send the mail: %s",err)
	}

	fmt.Println("Email Sent!")

	return "Email Sent Successfully!!", nil
}





// HeroesServiceChaincode implementation of Chaincode
type sendmailChaincode struct {
}

// Init of the chaincodeerr
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *sendmailChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Chaincode for sending email - Instantiation")


	args := stub.GetStringArgs()

	fmt.Println("Got arguments at Initialization -- 1 -> " + args[0] + " 2 -> " + args[1]);

	if(len(args) != 2){
		return shim.Error("Incorrect Arguments, expecting an emailId along with a message");
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return a successful message
	return shim.Success(nil)
}


func (t *sendmailChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {


	// Get the function and arguments from the request
	fn, args := stub.GetFunctionAndParameters()

	fmt.Println("Invoke Function: " + fn + " Arguments -- " + args[0]);

	var result string
	var err error
	if fn == "Send Message" {
		result, err = sendMail(stub, args);
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(sendmailChaincode))
	if err != nil {
		fmt.Printf("Error starting sendmailChaincode: %s", err)
	}
}
