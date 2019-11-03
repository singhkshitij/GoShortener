# GoShortener

A tiny url shortener written in GO having storage as Bbolt dB !

### Endpoints available :

- (GET)  / : Welcome Screen
- (POST) /api/v1/url/shorten : { URL : `URL`} Shortens url and returns 
- (GET) /u/{shortkey} : redirects to long url
- (GET) /api/v1/url/backup : Downloads current db state backup

### How to run ?

#### Local build : 

- ​	go build
- ​	.\GoShortener.exe

#### Docker build :

- ​	docker build -t go-shortener .
- ​	docker run -d -p 8080:8080 go-shortener

### //  TODO

1. Write tests for existing endpoints
2. Serving Concurrent request using goroutines
3. Adding more functionalities such as total urls, hit counts per url, generating QR Code, Custom Url provided by user, better analytics per url
4. Health check for service
5. Scheduled dB backup
6. Swagger UI Integration

