# JWT Auth Microservice

Created a simple jwt auth service that can handle registration, login and jwt token signing.

Deployed on aws. Data is stored in dynamodb, api gateway is used for routing/proxying and lamdas for receiving events (http request) and running code.


### aws resources are bundled into a single stack (Infrastructure as Code) which can be managed from gows.go
