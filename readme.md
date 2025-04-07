


-- 
符合 AES-128 標準（適用於 HLS #EXT-X-KEY）
```
openssl rand 16 > stream1.key

```



--  Initialization Vector (IV)
HLS IV（初始化向量） 用於加密時的隨機數，可以生成 16-byte (128-bit) 隨機 IV：
```
openssl rand -hex 16
```


##### Get key

  ```
  curl -X GET "http://localhost:9090/api/v1/hls/key?key=stream1.key" -i



  HTTP/1.1 200 OK
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With, User-Agent, Pragma, Referer, X-Forwarded-For, X-Real-Ip, Accept-Language, utoken, x-key
Access-Control-Allow-Methods: OPTIONS, GET, POST, PUT, DELETE, PATCH
Access-Control-Allow-Origin: *
Access-Control-Expose-Headers: Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, utoken, x-key
Content-Type: application/octet-stream
Referrer-Policy: origin
Date: Wed, 12 Mar 2025 06:54:23 GMT
Content-Length: 16

(�XNR�Cg��e�%
  ```

##### JWT test
- get tonken
```
curl -X POST "http://localhost:9090/api/v1/auth/token" \
     -d "username=wwxhyuyusj" \
     -H "header-key: 6HdSWud6jkNUYEt8XrK6PuW"

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDE3NjY5MjMsImlhdCI6MTc0MTc2NjMyMywic3ViIjoid3d4aHl1eXVzaiJ9.PS0dg-bwngFTM4s6dw8NILOb0AfNJ-XptgBZG9m9b18"}%
```
- get key
```
  curl -X POST "http://localhost:9090/api/v1/hls/key" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJobHMta2V5LWFwaSIsImV4cCI6MTc0MTg1NTY0MywiaWF0IjoxNzQxODU1MDQzLCJpc3MiOiJobHMta2V5LXNlcnZlciIsInN1YiI6Ind3eGh5dXl1c2oifQ.mb1IaNspIq-SYN4CEROY86X802gX7SK0NjVbUcWIrEA" \
     -d "key=stream1.key"

```



#### JWT decoder 

![](./resource/1.png)
