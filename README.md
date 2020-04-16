## MyHTTP  
This tool displays the MD5 hash of the response body from the given
list of URLs  

### Install  
```bash
~$ go get github.com/wesleyvicthor/myhttp
```
or clone
```bash
~$ git clone https://github.com/wesleyvicthor/myhttp
```

### Tests
```
myhttp/pkg~$ go test
```

### Usage  
```bash
~$ myhttp example.com
```
or multiple urls by space separated  
```bash
~$ myhttp adjust.com example.com https://httpbin.org
```
you can adjust the parallel processes using the flag
```bash
~$ myhttp -parallel=2 adjust.com example.com httpbin.org
``` 
the default is 10 process  