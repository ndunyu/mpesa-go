package examples

import (
	"encoding/json"
	"fmt"
	mpesa_go "github.com/ndunyu/mpesa-go"
	"log"
	"net/http"
	"time"
)
func MpesaExpressExample() {
	//NOTE it easy to create only one mpesa object per app
	//and use it though out you app
	//instead of creating a new object per request
	mpesa := mpesa_go.New("ConsumerKey_preferably_form_environmental_variables", "ConsumerSecret", false)
	mpesa.SetDefaultTimeOut(10 * time.Second)
	mpesa.SetDefaultPassKey("PASSKEY")
	mpesa.SetDefaultB2CShortCode("MPESASHORTCODE")

	//For this request I do not pass a passKey since I already
	//set a DefaultPassKey if you did not set it you should pass a passkey as the second
	//parameter as a second argument
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
