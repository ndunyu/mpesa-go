package mpesa_go

import (
	"bytes"
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
	//For sandbox use false and for production use true
	Live           bool
	ConsumerKey    string
	ConsumerSecret string
	DefaultTimeOut time.Duration
}

func New(ConsumerKey, ConsumerSecret string, live bool) Mpesa {
	return Mpesa{
		Live:           live,
		ConsumerKey:    ConsumerKey,
		ConsumerSecret: ConsumerSecret,
		DefaultTimeOut: 20 * time.Second,
	}

}

// SetDefaultTimeOut this will set the connection timeout to mpesa the defaul is
//20 seconds
func (m *Mpesa) SetDefaultTimeOut(timeOut time.Duration) {
	m.DefaultTimeOut = timeOut

}

// SetMode  changes from production to sandbox and viceversa
//at runtime.
func (m *Mpesa) SetMode(mode bool) {
	m.Live = mode

}

func (m *Mpesa) StkPushRequest(body StKPushRequestBody, passKey string) (*StkPushResult, error) {
	if body.Timestamp == "" {
		t := time.Now()
		fTime := t.Format("20060102150405")
		body.Timestamp = fTime
		body.Password = GeneratePassword(body.BusinessShortCode, passKey, fTime)
	}
	body.TransactionType = CustomerPayBillOnline
	var stkPushResult StkPushResult
	err := m.sendAndProcessStkPushRequest(m.getStkPush(), body, &stkPushResult, nil)
	return &stkPushResult, err
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
	log.Println(token.AccessToken)
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
	resp, err := postRequest(url, data, headers,m.DefaultTimeOut)
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

	return &token, nil
}

func postRequest(url string, data interface{}, headers map[string]string,timeOut time.Duration) (*http.Response, error) {

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

func getRequest(url string, headers map[string]string, queryParameters map[string]string,timeOut time.Duration) (*http.Response, error) {
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
