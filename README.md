# PagerDutyEventsLite
A very simple PagerDuty command line EVENTS caller written in Go

This util can be used to make calls to PagerDuty events API ( not the REST API , yes there is a differnece between REST and EVENTS API at PagerDuty! ).

Why "lite"?

* My mate wrote the proper full version that's in production use using Python. 
* I hate Python!
* I couldn't run his Python version on Windows2008 boxes but the Go compiled version worked perfect
* It used the same command flags as the original version to make it 100% compatible.
* ( Mine is 3 times faster than his 'cos it's written in Go and compiles to a single optimized binary! )
* Go is the best language ever. Period!
* Python is for data scientists and "pretend coders"! LOL!

```
PagerDuty Util Lite - 1.2

--routing_key    <string> - The primary routing key for the PD event rule or service
--keyname        <string> - Unique user defined key.
--event          <string> - Must be one of { trigger | acknowledge | resolve}
--severity       <string> - {info | critical | warning | error}
--msg            <string> - Primary message alert title.
--source         <string> - Source of the alert, advise use of hostname.
--details        <string> - Simple logging details for the alert.
--jdetails       <string> - JSON formatted structure declaring sets of key:value pairs with log information sets. ( OPT )
--proxy_server   <string> - Force specific proxy server ( default use HTTP_PROXY/HTTPS_PROXY from environment).
--JSONresult             Return result in JSON format.
--saveJSONresponse       Save the JSON result to a file names <keyname>.json.


Note :
  - If you need to use a proxy then set HTTP_PROXY or HTTPS_PROXY in the environment first ( set on Windows cmd line or export on Unix ).
  - When supplying <jdetails> param, be sure to escape the double-quotes. PowerShell : '{\"key01\":\"value01\"}'
  - When supplying <jdetails> param it will override <details> param, although both will be sent to PD if supplied.
```
