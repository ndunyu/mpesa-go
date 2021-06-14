package mpesa_go

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

type Mpesa struct {
	//if set to true the access roken will be reused
	//until it expires
	//otherwise each request will always get a new token
	//which will slow down your requests
	//the default is true
	CacheAccessToken bool
	//For sandbox use false and for production use true
	Live           bool
	ConsumerKey    string
	ConsumerSecret string
	DefaultTimeOut time.Duration
	//for those using only one payBill
	//you can set a default passkey to be used
	//instead of passing a passkey evey time
	//when doing an stk push
	//use SetDefaultPassKey to change default
	//pass key at runtime.
	DefaultPassKey string
	//You can pass the Mpesa shortcode
	//you want In case you dont to pass the shortcode each time
	//you are sending a request
	//this is ideal for those using a single shortcode
	DefaultC2BShortCode string
	DefaultB2CShortCode string

	///for b2c
	DefaultInitiatorName string

	///for b2c
	DefaultSecurityCredential string

	cache map[string]AccessTokenResponse
}

func New(ConsumerKey, ConsumerSecret string, live bool) Mpesa {
	TokenCache := make(map[string]AccessTokenResponse)
	return Mpesa{
		Live:             live,
		ConsumerKey:      ConsumerKey,
		ConsumerSecret:   ConsumerSecret,
		DefaultTimeOut:   20 * time.Second,
		CacheAccessToken: true,
		cache:            TokenCache,
	}
}
func (m *Mpesa) ShouldCacheAccessToken(shouldCache bool) {
	m.CacheAccessToken = shouldCache
}

//SetDefaultB2CShortCode will set the default shortcode
// to use if you do not provide any
func (m *Mpesa) SetDefaultB2CShortCode(shortCode string) {
	m.DefaultC2BShortCode = shortCode

}

//SetDefaultPassKey You can set the default pass key
//Over here so that you dont have to pass it each time
//You are sending an StkRequest
func (m *Mpesa) SetDefaultPassKey(passKey string) {
	m.DefaultPassKey = passKey

}

// SetDefaultTimeOut this will set the connection timeout to Mpesa
//the default is 20 seconds when sending an http request
func (m *Mpesa) SetDefaultTimeOut(timeOut time.Duration) {
	m.DefaultTimeOut = timeOut
}

// SetLiveMode  changes from production to sandbox and viceversa
//at runtime.
func (m *Mpesa) SetLiveMode(mode bool) {
	m.Live = mode
}

// B2CRequest Sends Money from a business to the Customer
func (m *Mpesa) B2CRequest(b2c B2CRequestBody) (*MpesaResult, error) {
	if b2c.CommandID == "" {

		b2c.CommandID = BusinessPayment
	}
	if IsEmpty(b2c.PartyA) {
		b2c.PartyA=m.DefaultB2CShortCode

	}
	if IsEmpty(b2c.InitiatorName) {
		b2c.InitiatorName=m.DefaultInitiatorName
	}
	if IsEmpty(b2c.SecurityCredential) {
		b2c.SecurityCredential=m.DefaultSecurityCredential
	}
	err := b2c.Validate()
	if err != nil {
		return nil,err
	}
	var mpesaResult MpesaResult
	err = m.sendAndProcessStkPushRequest(m.getB2CUrl(), b2c, &mpesaResult, nil)

	return &mpesaResult,err

}

//StkPushRequest send an Mpesa express request
//note if you have already set a DefaultPassKey you don't have to pass
//a pass key here its optional
//If you also set DefaultC2BShortCode you dont have to pass BusinessShortCode to the StKPushRequestBody
//the default will be used if you don't pass it
func (m *Mpesa) StkPushRequest(body StKPushRequestBody, passKey ...string) (*StkPushResult, error) {
	var stkPassKey string
	if len(passKey) > 0 {
		stkPassKey = passKey[0]
	} else {
		stkPassKey = m.DefaultPassKey
	}
	if IsEmpty(stkPassKey) {
		return nil, errors.New("pass key is needed set a default pass key or pass it in ths function")
	}
	if IsEmpty(body.BusinessShortCode) {
		body.BusinessShortCode = m.DefaultC2BShortCode
	}
	err := body.Validate()
	if err != nil {
		return nil, err
	}
	t := time.Now()
	fTime := t.Format("20060102150405")
	requestBody := StkRequestFullBody{
		StKPushRequestBody: body,
		Password:           GeneratePassword(body.BusinessShortCode, stkPassKey, fTime),
		Timestamp:          fTime,
		TransactionType:    CustomerPayBillOnline,
		PartyA:             body.PhoneNumber,
		PartyB:             body.BusinessShortCode,
	}
	var stkPushResult StkPushResult
	err = m.sendAndProcessStkPushRequest(m.getStkPush(), requestBody, &stkPushResult, nil)
	return &stkPushResult, err
}

//StkPushVerification use this to confirm your stk push if it was a failure or success
//CheckoutRequestID is the CheckoutRequestID you got when you sent the StkPushRequest request
//you dont have to send a passkey if you have a DefaultPassKey set
func (m *Mpesa) StkPushVerification(CheckoutRequestID string, BusinessShortCode string, passKey ...string) (*StkPushQueryResponseBody, error) {
	var stkPassKey string
	if len(passKey) > 0 {
		stkPassKey = passKey[0]
	} else {
		stkPassKey = m.DefaultPassKey
	}
	if IsEmpty(stkPassKey) {
		return nil, errors.New("pass key is needed set a default pass key or pass it in ths function")
	}
	t := time.Now()
	fTime := t.Format("20060102150405")
	body := StkPushQueryRequestBody{
		BusinessShortCode: BusinessShortCode,
		Password:          GeneratePassword(BusinessShortCode, stkPassKey, fTime),
		Timestamp:         fTime,
		CheckoutRequestID: CheckoutRequestID,
	}
	var stkPushResult StkPushQueryResponseBody
	err := m.sendAndProcessStkPushRequest(m.getStkPushQuery(), body, &stkPushResult, nil)
	return &stkPushResult, err
}


func (m *Mpesa)StkPushQuery(body StkPushQueryRequestBody,passKey ...string)(*StkPushQueryResponseBody, error){
	var stkPassKey string
	if len(passKey) > 0 {
		stkPassKey = passKey[0]
	} else {
		stkPassKey = m.DefaultPassKey
	}
	if IsEmpty(stkPassKey) {
		return nil, errors.New("pass key is needed set a default pass key or pass it in ths function")
	}
	if body.Timestamp=="" {
		t := time.Now()
		fTime := t.Format("20060102150405")
		body.Timestamp= fTime
		body.Password=GeneratePassword(body.BusinessShortCode,stkPassKey, fTime)
	}

	var stkPushResult StkPushQueryResponseBody
	err:=m.sendAndProcessStkPushRequest(m.getStkPushQuery(),body,&stkPushResult,nil)
	return  &stkPushResult,err
}

func (m *Mpesa) sendAndProcessStkPushRequest(url string, data interface{}, respItem interface{}, extraHeader map[string]string) error {
	if reflect.ValueOf(respItem).Kind() != reflect.Ptr {
		log.Println("not a pointer")

		return errors.New("response should be a pointer")
	}
	token, err := m.GetAccessToken()
	if err != nil {

		return err
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	///auth := "Bearer " + token.AccessToken
	//headers["authorization"]= auth
	headers["Authorization"] = "Bearer " + token.AccessToken
	///add the extra headers
	//Get union of the headers
	for k, v := range extraHeader {
		headers[k] = v
	}
	resp, err := postRequest(url, data, headers, m.DefaultTimeOut)
	if err != nil {

		return err
	}
	defer resp.Body.Close()
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		b, _ := ioutil.ReadAll(resp.Body)

		return &RequestError{Message: string(b), StatusCode: resp.StatusCode}

	}

	///var respe map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(respItem); err != nil {

		return errors.New("error converting from json")
	}

	return nil

}

//GetAccessToken will get the token to be used to query data
func (m *Mpesa) GetAccessToken() (*AccessTokenResponse, error) {
	if m.cache == nil {
		m.cache = make(map[string]AccessTokenResponse)
	}
	//check if we have a cached Token and return it instead
	//if we have allowed caching
	if m.CacheAccessToken {
		if val, ok := m.cache[AccessToken]; ok {
			///we have a token so check if it is expired
			if val.ExpireTime.Sub(time.Now()).Minutes() > 0 {
				return &val, nil
			}
		}
	}
	req, err := http.NewRequest(http.MethodGet, m.getAccessTokenUrl(), nil)
	if err != nil {
		return nil, err
	}
	log.Println(req.URL.String())

	req.SetBasicAuth(m.ConsumerKey, m.ConsumerSecret)

	req.Header.Add("Accept", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{
		Timeout: m.DefaultTimeOut,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		b, _ := ioutil.ReadAll(resp.Body)
		if string(b) == "" {
			return nil, &RequestError{Message: "Error getting token", StatusCode: resp.StatusCode}

		}
		return nil, &RequestError{Message: string(b), StatusCode: resp.StatusCode}
	}
	var token AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {

		return nil, errors.New("error converting from json")
	}

	if m.cache == nil {
		m.cache = make(map[string]AccessTokenResponse)

	}
	if m.CacheAccessToken {
		//cache the token
		token.ExpireTime = time.Now().Add(time.Minute * 50)
		m.cache[AccessToken] = token
	}
	return &token, nil
}

func postRequest(url string, data interface{}, headers map[string]string, timeOut time.Duration) (*http.Response, error) {

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: timeOut,
	}
	return client.Do(req)

}

func getRequest(context context.Context, url string, headers map[string]string, queryParameters map[string]string, timeOut time.Duration) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range queryParameters {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()

	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{

		Timeout: timeOut,
	}
	return client.Do(req)

}

type RequestError struct {
	StatusCode int

	Message string
	Url     string
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("url is: %s \n status code is: %d \n  and body is : %s", r.Url, r.StatusCode, r.Message)

}
