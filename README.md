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

--routing_key            <string> - The primary routing key for the PD event rule or service
--keyname                <string> - Unique user defined key.
--event                  <string> - Must be one of { trigger | acknowledge | resolve}
--severity               <string> - {info | critical | warning | error}
--msg                    <string> - Primary message alert title.
--source                 <string> - Source of the alert, advise use of hostname. ( OPT )
--details                <string> - Simple logging details for the alert. ( OPT )
--jsondetailsfile        <string> - JSON formatted text file with sets of key/value pairs holding extra alerting info. ( OPT )
--proxy_server           <string> - Force specific proxy server ( default use HTTP_PROXY/HTTPS_PROXY from environment). ( OPT )

--jsonresult             Return result to STDOUT in JSON format. Useful for other apps that need to capture the result. ( OPT )
--savejsonresponse       Save the JSON result to a file. ( OPT )

Note :
  - If you need to use a proxy then HTTP_PROXY or HTTPS_PROXY are drawn from the environment by default.
```
