![example workflow](https://github.com/serainville/bitbucket-webhooks/actions/workflows/go.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/serainville/bitbucket-webhooks)](https://goreportcard.com/report/github.com/serainville/bitbucket-webhooks)
# Go Bitbucket Webhook Module
## Overview
This module is used for handling Bitbucket Server and Bitbucket Cloud Webhook requests. Requests are detected and can be validated. Also, HMAC signatures can be verified to ensure the authenticity of the message itself.

Use this module to validate incoming Bitbucket Webhooks and to ensure the requests are authentic, by checking the HMAC signature when webhook secrets are used.

