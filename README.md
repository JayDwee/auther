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
- [6750 - The OAuth 2.0 Authorization Framework: Bearer Token Usage](https://datatracker.ietf.org/doc/html/rfc6750)
- [7009 - Oauth 2.0 Token Revocation](https://datatracker.ietf.org/doc/html/rfc7009)
- [7515 - JSON Web Signature (JWS)](https://datatracker.ietf.org/doc/html/rfc7515)
- [7516 - JSON Web Encryption (JWE)](https://datatracker.ietf.org/doc/html/rfc7516)
- [7517 - JSON Web Key (JWK)](https://datatracker.ietf.org/doc/html/rfc7517)
- [7518 - JSON Web Algorithms (JWA)](https://datatracker.ietf.org/doc/html/rfc7518)
- [7519 - JSON Web Token (JWT)](https://datatracker.ietf.org/doc/html/rfc7519)
- [7523 - JSON Web Token (JWT) Profile for OAuth 2.0 Client Authentication and Authorization Grants](https://datatracker.ietf.org/doc/html/rfc7523)
- [7617 - The 'Basic' HTTP Authentication Scheme](https://datatracker.ietf.org/doc/html/rfc7617)
- [7636 - Proof Key for Code Exchange by OAuth Public Clients](https://datatracker.ietf.org/doc/html/rfc7636)
- [7638 - JSON Web Key (JWK) Thumbprint](https://datatracker.ietf.org/doc/html/rfc7638)
- [7662 - Oauth 2.0 Token Introspection](https://datatracker.ietf.org/doc/html/rfc7662)
- [7797 - JSON Web Signature (JWS) Unencoded Payload Option](https://datatracker.ietf.org/doc/html/rfc7797)
- [8252 - OAuth 2.0 for Native Apps](https://datatracker.ietf.org/doc/html/rfc8252)
- [8628 - OAuth 2.0 Device Authorization Grant](https://datatracker.ietf.org/doc/html/rfc8628)
- [8693 - OAuth 2.0 Token Exchange](https://datatracker.ietf.org/doc/html/rfc8693)
- [8725 - JSON Web Token Best Current Practices](https://datatracker.ietf.org/doc/html/rfc8725)

- [OpenID Connect Core 1.0 incorporating errata set 2](https://openid.net/specs/openid-connect-core-1_0.html)
- [OpenID Connect Discovery 1.0 incorporating errata set 2](https://openid.net/specs/openid-connect-discovery-1_0.html)
- [OpenID Connect Dynamic Client Registration 1.0 incorporating errata set 2](https://openid.net/specs/openid-connect-registration-1_0.html)
- [OpenID Connect Session Management 1.0](https://openid.net/specs/openid-connect-session-1_0.html)

## Infrastructure
- Namecheap
- Cloudflare
- API Gateway
- Lambda
- S3
- DynamoDB
