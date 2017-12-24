package btcid

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	domain               = "https://vip.bitcoin.co.id"
	endpointTicker       = "/ticker"
	pairBTCIDR           = "/btc_idr"
	pairETHIDR           = "/eth_idr"
	prvMethodGetInfo     = "getInfo"
	prvMethodTransHist   = "transHistory"
	prvMethodTrade       = "trade"
	prvMethodTradeHist   = "tradeHistory"
	prvMethodOpenOrders  = "openOrders"
	prvMethodOrderHist   = "orderHistory"
	prvMethodGetOrder    = "getOrder"
	prvMethodCancelOrder = "cancelOrder"
	prvMethodWithdCoin   = "withdrawCoin"
)

var (
	publicAPIBaseURL     = fmt.Sprintf("%s/api", domain)
	privateAPIBaseURL    = fmt.Sprintf("%s/tapi", domain)
	endpointBTCIDRTicker = fmt.Sprintf("%s%s", pairBTCIDR, endpointTicker)
	endpointETHIDRTicker = fmt.Sprintf("%s%s", pairETHIDR, endpointTicker)
)

type Ticker struct {
	High string `json:"high"`
	Low  string `json:"low"`
	Last string `json:"last"`
	Buy  string `json:"buy"`
	Sell string `json:"sell"`
}

type Trade struct {
	Date   string `json:"date"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
	TID    string `json:"tid"`
	Type   string `json:"type"`
}

type Depth struct {
	Buy  [][]interface{} `json:"buy"`
	Sell [][]interface{} `json:"sell"`
}

type Client struct {
	apiKey     string
	secret     string
	httpClient *http.Client
}

type UserInfo struct {
	Balance        map[string]interface{} `json:"balance"`
	BalanceHold    map[string]interface{} `json:"balance_hold"`
	Address        map[string]string      `json:"address"`
	UserID         string                 `json:"user_id"`
	ProfilePicture string                 `json:"profile_picture"`
	Name           string                 `json:"name"`
	ServerTime     int                    `json:"server_time"`
	Email          string                 `json:"email"`
}

type InfoRes struct {
	Success int      `json:"success"`
	Return  UserInfo `json:"return"`
}

func (c *Client) newPrvReq(PrivateMethod string) ([]byte, error) {
	// Prepare variables for signing and sending
	nonce := strconv.FormatInt(time.Now().Unix(), 10)

	// Build URL query parameters
	q := url.Values{}
	q.Set("method", PrivateMethod)
	q.Set("nonce", nonce)
	queryString := q.Encode()

	// Setup Request
	req, err := http.NewRequest(http.MethodPost, privateAPIBaseURL, strings.NewReader(queryString))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Sign request
	hmac512 := hmac.New(sha512.New, []byte(c.secret))
	hmac512.Write([]byte(queryString))
	signature := hex.EncodeToString(hmac512.Sum(nil))

	// Set headers
	req.Header.Set("Key", c.apiKey)
	req.Header.Set("Sign", signature)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("Res error")
		return nil, err
	}

	// Read the response
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("Read error")
		return nil, err
	}

	return body, nil
}

func (c *Client) newPubReq(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", publicAPIBaseURL, endpoint)
	payload := bytes.NewBuffer([]byte{})
	req, err := http.NewRequest(http.MethodGet, url, payload)
	if err != nil {
		fmt.Println("Req error", err.Error())
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("Res error")
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("Read error")
		return nil, err
	}

	return body, nil
}

// GetTicker fetches the latest ticker data from the API
func (c *Client) GetTicker() (Ticker, error) {
	body, err := c.newPubReq("/btc_idr/ticker")
	if err != nil {
		fmt.Println("Req error")
	}
	ticker := struct {
		Ticker Ticker `json:"ticker"`
	}{}
	err = json.Unmarshal(body, &ticker)
	if err != nil {
		fmt.Println("JSON error", err)
		return Ticker{}, err
	}
	return ticker.Ticker, nil
}

// GetTrades fetches the latest market trade data from the API
func (c *Client) GetTrades() ([]Trade, error) {
	body, err := c.newPubReq("/btc_idr/trades")
	if err != nil {
		fmt.Println("Req error")
	}
	trades := []Trade{}
	err = json.Unmarshal(body, &trades)
	if err != nil {
		fmt.Println("Trade error", err)
	}
	return trades, nil
}

// GetDepth fetches the market cap data from the API
func (c *Client) GetDepth() (Depth, error) {
	body, err := c.newPubReq("/btc_idr/depth")
	if err != nil {
		fmt.Println("Req error")
	}
	depth := Depth{}
	err = json.Unmarshal(body, &depth)
	if err != nil {
		fmt.Println("Depth error", err)
	}
	return depth, nil
}

// GetInfo fetches an account's information details
func (c *Client) GetInfo() (UserInfo, error) {
	body, err := c.newPrvReq(prvMethodGetInfo)
	if err != nil {
		fmt.Println("Req error")
	}
	infoRes := InfoRes{}
	err = json.Unmarshal(body, &infoRes)
	if err != nil {
		fmt.Println("Info error", err)
	}
	return infoRes.Return, nil
}
