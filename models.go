package mpesa_go




//Start of Mpesa express models<-----------------------------------------------

// StKPushRequestBody this is the body we will send when sending an
//mpesa express request
type StKPushRequestBody struct {
	BusinessShortCode string
	Password string
	Timestamp string
	///use only [ CustomerPayBillOnline ]
	TransactionType string
	Amount string
	//sender phone number
	PartyA string
	///receiver shortcode
	PartyB string
	////Sending funds
	PhoneNumber string
	///
	CallBackURL string
	///use this with paybill
	AccountReference string
	//
	TransactionDesc string
}

// StkPushResult is the  result returned
//when you send a Mpesa express  result
type StkPushResult struct {
	///
	CheckoutRequestID string
	CustomerMessage string
	MerchantRequestID string
	ResponseCode string
	ResponseDescription string

}

//End of Mpesa express models<-----------------------------------------------


//Start of Token Model-------------------------------------------

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

//End of access  token<-----------------------------------------------