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
