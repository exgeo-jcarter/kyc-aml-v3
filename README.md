Checks for matches against a blacklist.  

To set up and use this program:  

$ go get github.com/sajari/fuzzy  
$ go get github.com/dotcypress/phonetics  
$ git clone https://github.com/exgeo-jcarter/kyc-aml-v2.git  
$ cd kyc-aml-v2/kyc-aml-server  
$ ./kyc-aml-server&  

... wait for server to load (about 1 minute)  

$ cd ../kyc-aml-client  
$ ./kyc-aml-client "ali akbar mohummad"  
