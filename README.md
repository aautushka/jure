# jure

# Redis protocol (official doc [here](https://redis.io/command)) 
## Read non-existant key
``` 
get key\r\n
$-1\r\n
```
## Set key value 
```
set key value\r\n
+OK\r\n
```
## Read existing key 
```
get key\r\n
$5\r\n
value\r\n
```
## Delete non-existant key 
```
del key\r\n
:0\r\n
```
## Delete existing key 
```
del key\r\n
:1\r\n
```
## Create a new key-value pair (alternative to a SET) 
```
append key value\r\n
:5\r\n
```
## Append existing value 
```
append key value\r\n
:10\r\n
get key\r\n
$10\r\n
valuevalue\r\n
```
## Check existing key 
```
exists key\r\n
$1\r\n
```
## Check non-existing key 
``` 
exists key\r\n
$0\r\n
```
## Remove existing key after 10 seconds 
```
expire key 10\r\n
$:1\r\n
```
## Try to expire non-existant key 
```
expire key 10\r\n
$:0\r\n
```
## Dump all keys (3 keys in this example) 
```
keys *\r\n
*3r\n
$4\r\n
key3\r\n
$4\r\n
key2\r\n
$4\r\n
key1\r\n
```
## No keys to dump 
```
keys *\r\n
*0\r\n
```
## Get key TTL 
``` 
ttl key\r\n
:80\r\n
```
## Get non-existant key TTL
```
ttl key\r\n
:-2\r\n
```




