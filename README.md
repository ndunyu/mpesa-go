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
	//You can cache the access token instead of 
	//getting a new one for each request
	//this will increase your speed
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
		CacheAccessToken: true,
	}

	//
	mpesa2.ShouldCacheAccessToken(true)

}
```

## Example of using the Mpesa Express Api

For this api PassKey is optional if you already set a DefaultPassKey when initializing your Mpesa Object.

If you want to use another pass key for this request or you have not set a default one remember to pass a passkey as a
second argument.

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
		//should be a string int
		Amount:            "",
		//note the format used for phone number
		PhoneNumber:       "254701047658",
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

## How to process the received Mpesa Express callback data

As you can recall we sent a calback url when sending the mpesa express request.

Safaricom send the response to this endpoint.Below example shows how you can process the received data.

Note when you receive this response always use the verification api to verify the result.

The mpesa_go.StkPushCallBackResponseBody also has mpesaRef and other data that you should remember to process

```go

func ExampleProcessingMpesaExpressCallBack(w http.ResponseWriter, r *http.Request) {
	var stkPushResponseBody mpesa_go.StkPushCallBackResponseBody
	err := json.NewDecoder(r.Body).Decode(&stkPushResponseBody)
	if err != nil {
		log.Println(err)
		///sentry.CaptureException(err)
		http.Error(w, "something went wrong", 400)
		return
	}
	defer r.Body.Close()
	if stkPushResponseBody.Body.StkCallback.ResultCode != 0 {
		///this request has failed
		///mark it as failed in the database or something
		//like that
		w.WriteHeader(200)
		return
	}
	//otherwise Resultcode is 0 so it is a success
	//sample of processing received data
	for _, item := range stkPushResponseBody.Body.StkCallback.CallbackMetadata.Item {
		switch item.Name {
		case "Amount":
			amount, ok := item.Value.(float64)
			if !ok {
				log.Fatal("error")
			}
			///do something with amount
			log.Println(amount)
		case "MpesaReceiptNumber":
			//this is the mpesa transaction id sent to user
			//e.g MWYWWUWUWUW
			str := fmt.Sprintf("%v", item.Value)
			log.Println(str)
		case "TransactionDate":
			date, ok := item.Value.(float64)

			if !ok {
				log.Fatal("error")

			}
			log.Println(date)
		}

	}
	
	//It is always wise to send a verification request to
	//confirm that it is true that this request was actually a success
	//just to double check
	//for that you can use the verification api
	mpesa := mpesa_go.New("consumerkey", "consumersecret", true)
	//even for this request you dont need to pass key if you already set
	//a default passkey
	verification, err := mpesa.StkPushVerification(stkPushResponseBody.Body.StkCallback.CheckoutRequestID, "ShortCODE", "PASS_KEY")
	if err != nil {
		///if an error occured
		///you can retry it again
		//or store it and retry it later
		return
	}
	if verification.ResultCode == "0" {
		///this request was a success
		//and the user actually paid
		//so award you user the money here
	} else {
		//user did not pay
		//mark it as a failed transaction
	}

	w.WriteHeader(200)

}

```

