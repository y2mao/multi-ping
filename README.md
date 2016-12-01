# multi-ping
ping multiple hosts and rank by lantency

# Installation
```go
go get -v githubcom/y2mao/multi-ping
```
# Usage
Create a host list file. e.g. host.txt
```
www.google.com
www.github.com
www.facebook.com
```
Execute following command:
```bash
multiping -h host.txt
```

# Options
Argument | Default Value | Description
--- | --- | ---
-h | host.txt | path of host list file
-t | 5 | ping count
-r | false | reserve result

