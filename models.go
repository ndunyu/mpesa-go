package mpesa_go

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//Start of Mpesa express models<-----------------------------------------------

// StKPushRequestBody this is the body we will send when sending an
//mpesa express request

type StkRequestFullBody struct {
	StKPushRequestBody
	Password  string
	Timestamp string
	///use only [ CustomerPayBillOnline ]
	TransactionType string
	PartyA          string
	///receiver shortcode
	PartyB string
}




type StKPushRequestBody struct {
	BusinessShortCode string
	Amount            string
	//sender phone number
	//i.e where we will show pop up
	//format should be (254 followed by 9 digit)
	//note do not include + at the beginning of the phone number
	PhoneNumber string
	///
	CallBackURL string
	///use this with paybill
	AccountReference string
	//
	TransactionDesc string
}

func (s *StKPushRequestBody) Validate() error {
	if IsEmpty(s.BusinessShortCode) {
		return errors.New("business short code is required")
	}
	if IsEmpty(s.Amount) {
		return errors.New("amount is required")
	}
	if IsEmpty(s.PhoneNumber) {
		return errors.New("phone number is required")

	}
	//validate amount
	i, err := strconv.Atoi(s.Amount)
	if err != nil || i < 1 {
		return errors.New("amount should be a string number that is greater than 0")
	}
	if CheckKenyaInternationalPhoneNumber(s.PhoneNumber) {
		return errors.New("the phone number should be in the format 254000000000 i.e(254 followed by 9 digits)")
	}
	return nil
}

// StkPushResult is the  result returned
//when you send a Mpesa express  result
type StkPushResult struct {
	///
	CheckoutRequestID   string
	CustomerMessage     string
	MerchantRequestID   string
	ResponseCode        string
	ResponseDescription string
}

// StkPushQueryRequestBody when querying for success/failure(i.e verification)
type StkPushQueryRequestBody struct {
	BusinessShortCode	string
	Password	string
	Timestamp	string
	CheckoutRequestID	string
}






// StkPushQueryResponseBody returned when yous send mpesa express verification request
type StkPushQueryResponseBody struct {
	MerchantRequestID	string
	CheckoutRequestID	string
	ResponseCode	string
	ResultDesc	string
	ResponseDescription	string
	ResultCode	string

}

type StkPushCallBackResponseBody struct {
	Body Body `json:"Body"`
}

type Body struct {
	StkCallback StkCallback `json:"stkCallback"`
}

type StkCallback struct {
	MerchantRequestID string           `json:"MerchantRequestID"`
	CheckoutRequestID string           `json:"CheckoutRequestID"`
	ResultCode        int              `json:"ResultCode"`
	ResultDesc        string           `json:"ResultDesc"`
	CallbackMetadata  CallbackMetadata `json:"CallbackMetadata"`
}

type CallbackMetadata struct {
	Item []Item `json:"Item"`
}
type Item struct {
	Name  string `json:"Name"`
	Value interface{}   `json:"Value,omitempty"`
}


//End of Mpesa express models<-----------------------------------------------

//Start of Token Model-------------------------------------------

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	ExpireTime   time.Time `json:"expire_time"`
}

//End of access  token<-----------------------------------------------

func IsEmpty(s string) bool {
	if len(strings.TrimSpace(s)) == 0 {
		return true

	}
	return false

}

func CheckKenyaInternationalPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`(^254)\d{9}$`)
	return re.MatchString(phone)
}