# Dominos-Explorer

## Using the Pizza-Voucher-Explorer (Netherlands)

There are never any codes when you need them. I've had enough. 
From now on, the Dominos-Voucher-Explorer(TM) will do the job of finding those elusive dominos codes. 

To use this tool, you're going to need to know
1. The city you reside in (shouldn't be hard).
2. The number of the Dominos store you want to order from.
3. The ASP.Net Session ID of an open connection you have to Dominos.

Number 3 is the most confusing to explain. You basically need to open
a browser, navigate to a dominos store session, open the web-developer,
and then check the outgoing network requests when you submit a voucher.
What you want to do is visit the Cookies of the request. In it you can
find the ASP.NET_SessionId. That will be your session. The store number
and city are also listed in the cookies.

### Results

This program outputs valid codes on standard-out. However, the last code is always
the last number it checked. If the program terminates for any reason, then you'll be
able to restart from where you left off.

**NOTE:** You will *probably* get banned from Dominos for this. Use a high poll delay (~1s).
