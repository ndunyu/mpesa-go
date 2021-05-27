# mpesa-go
A golang sdk for safaricom Mpesa

## Quickstart
```go

// MpesaInitialization  show how to create an mpesa object to use through out your app
//note for better security always read your consumerkey and consumer secret from environmental variables
//or from a config file
//otherwise you might commit them to git accidentally
//note i have set live to false to indicate sandbox
//for production set it to true
func MpesaInitialization() {
	mpesa := mpesa_go.New("ConsumerKey_preferably_form_environmental_variables", "ConsumerSecret", false)
	//you can change from sandbox to production at runtime like below
	mpesa.SetLiveMode(true)

	//you can set a default timeout for all http request this is optional though
	//if not set 20 seconds is the default
	//below I have set it to 10 as an example
	mpesa.SetDefaultTimeOut(10 * time.Second)
	//you can also set a default passkey to be used in every mpesaexpress
	// request that you dont provide a passkey
	//it is advisable to save this passkey in a config file
	mpesa.SetDefaultPassKey("PASSKEY")
	//you can pass a dfault mpesa paybill that will
	//be used as default for all request that you dont pass a paybill
	mpesa.SetDefaultB2CShortCode("MPESASHORTCODE")

}
```
You can also initialize an mpesa object like below.

```go
func MpesaInitializationExampleTwo(){
	//Example two
	//you can also initialize an mpesa object like below
	//the problem with this is that you might forget to pass a consumerKey
	//and passkey that are required
	mpesa2 := mpesa_go.Mpesa{
		Live:           false,
		ConsumerKey:    "consumerKeyHere",
		ConsumerSecret: "ConsumerSecret here",
		//values below this line are optional
		//but values above it are required
		DefaultTimeOut:      10 * time.Second,
		DefaultPassKey:      "",
		DefaultC2BShortCode: "",
	}

	//
	mpesa2.SetLiveMode(true)

}
```
## Example of using the Mpesa Express Api
For this api PassKey is optional if you already set a  DefaultPassKey when initializing your
Mpesa Object.

If you want to use another pass key for this request or you have not set a default one remember to pass a passkey as a second argument.
```go
func MpesaExpressExample() {
	//NOTE it easy to create only one mpesa object per app
	//and use it though out you app
	//instead of creating a new mpesa object per request
	mpesa := mpesa_go.New("ConsumerKey_preferably_form_environmental_variables", "ConsumerSecret", false)
	mpesa.SetDefaultTimeOut(10 * time.Second)
	mpesa.SetDefaultPassKey("PASSKEY")
	mpesa.SetDefaultB2CShortCode("MPESASHORTCODE")
	
	//For this request I do not pass a passKey since I already 
	//set a DefaultPassKey
	response, err := mpesa.StkPushRequest(mpesa_go.StKPushRequestBody{
		///NOTE if you already set a DefaultB2CShortCode you dont have to
		//pass the BusinessShortCode if it is empty the default will be used
		BusinessShortCode: "",
		Amount:            "",
		PhoneNumber:       "",
		//change this to your callback url
		//when a user pays or payment fails something happens you will receive the response here
		CallBackURL: "https://send/the/callback/here",
		//something you use to identify which user has paid
		//for example for something like KPLC this would be METER NUMBER
		AccountReference: "",
		TransactionDesc:  "",
	})
	if err != nil {
		//Deal with your error here
		log.Fatal(err)
	}
	//Do something with your response here
	//for example:
	if response.ResponseCode == "0" {
		//here the request is a success do something

	} else {
		//here your response has failed
		//do tell user about it for example
	}

}
```
