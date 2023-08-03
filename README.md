# Tamper
# is not complete . . .
Simple & Fast Tool To Test Verb Tampering


## Run :
```shell
# Way 1
go run main.go
```
```shell
# Way 2
go build -ldflags "-w -s" main.go
./main
```

## Options :
```yaml
    -D string
    URL Of Your Target Do You Want To Test
```
```yaml
    -C string
    Set Value Of Cookie Header
```
```yaml
    -FC string
    Don't Match Response Code [use ',' To Split]
```
```yaml
    -MC string
    Match Response Code [use ',' To Split]
```
```yaml
    -O string
    Name Of File To Set Result in it
```
```yaml
    -X bool
    FUZZ Extra HTTP Methods
```
```yaml
    -H string
    Set Custom Headers To Test
```
