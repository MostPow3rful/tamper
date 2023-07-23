# Tamper
Simple & Fast Tool To Test Verb Tampering


## Run :
```shell
# 1
go run main.go
```
```shell
# 2
go build -ldflags "-w -s" main.go
./main
```

## Arguments :
```shell
./main -h
```
```yaml
-c : Set Cookie in Header
  type   : string
  defult : ""
  status : optional
  
-d : Set URL Of Your Target
  type    : string
  default : ""
  status : required

-fc : Don't Match Response Code [use ',' To Split]
  type    : string
  default : ""
  status : optional

-mc : Match Response Code [use ',' To Split]
  type    : string
  default : "*"
  status : optional

-o : Set Output File
  type    : string
  default : ""
  status : optional

-x : FUZZ Extra HTTP Methods
  type    : bool
  default : false
  status : optional
```

## Example :
```shell
./main -d https://google.com/
```
```shell
./main -d https://google.com/ -fc 500
```
```shell
./main -d https://google.com/ -x -mc 200 -c username=MostPow3rful
```
```shell
./main -d https://google.com/ -x -mc 200,403 -o result.txt
```