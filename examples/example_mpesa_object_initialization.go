package examples

import (
	mpesa_go "github.com/ndunyu/mpesa-go"
	"time"
)

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
	mpesa.ShouldCacheAccessToken(true)
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

func MpesaInitializationExampleTwo() {
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
		CacheAccessToken:    true,
	}

	//
	mpesa2.SetLiveMode(true)

}
