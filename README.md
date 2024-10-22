# Auther
[![Serverless Package & Upload 📦🪣](https://img.shields.io/github/actions/workflow/status/jaydwee/auther/serverless-package.yml?style=flat-square&logo=github&label=build)](https://github.com/jaydwee/auther/actions/workflows/serverless-package.yml)

Serverless authentication provider Developed in Go

## Features
### Resource Owner
- Create a Resource Owner
- Customise UI
- Create User Pools
- Login/Register
- Create Roles
- Configure scopes per role/user

### Authorization Server
- Create an Authorization Server
- Link to Resource Owner
- Create Clients
- Custom Claims
- Configurable Token
  - JWT
  - JWE
  - Opaque
- Blacklist

## RFCs Supported
- [6749 - The OAuth 2.0 Authorization Framework](https://datatracker.ietf.org/doc/html/rfc6749)
- [7515 - JSON Web Signature (JWS)](https://datatracker.ietf.org/doc/html/rfc7515)
- [7516 - JSON Web Encryption (JWE)](https://datatracker.ietf.org/doc/html/rfc7516)
- [7517 - JSON Web Key (JWK)](https://datatracker.ietf.org/doc/html/rfc7517)
- [7518 - JSON Web Algorithms (JWA)](https://datatracker.ietf.org/doc/html/rfc7518)
- [7519 - JSON Web Token (JWT)](https://datatracker.ietf.org/doc/html/rfc7519)
- [7617 - The 'Basic' HTTP Authentication Scheme](https://datatracker.ietf.org/doc/html/rfc7617)
- [7638 - JSON Web Key (JWK) Thumbprint](https://datatracker.ietf.org/doc/html/rfc7638)
- [7797 - JSON Web Signature (JWS) Unencoded Payload Option](https://datatracker.ietf.org/doc/html/rfc7797)

## Infrastructure
- Namecheap
- Cloudflare
- API Gateway
- Lambda
- S3
- DynamoDB