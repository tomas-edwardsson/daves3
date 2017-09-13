# daves3 - A S3 backed DAV server (WebDav)

## Is this something for me?
Maybe, if you need
* a dav server with basic auth
* which redirects GET requests to a s3 bucket
* and proxies PUT request to a s3 bucket

## What is this used for?
I use it as a [caching](https://www.pantsbuild.org/setup_repo.html#build-cache) server for [pantsbuild](https://pantsbuild.github.io/)

## Install
* Fetch it via `git clone` and run `go build`

## Running
Configuration is done via environment variables, the following are required for daves3. You will also have to make sure that you have working AWS credentials via AWS_PROFILE or AWS access key/secret.

| Env | Value |
| --- | ----- |
| DAVES3_USERNAME | username for http access |
| DAVES3_PASSWORD | password for http access |
| DAVES3_BUCKET   | bucket to store objects  |

## What AWS IAM access is required to my S3 bucket
* s3:GetObject
* s3:PutObject
