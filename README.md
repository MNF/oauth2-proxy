# Admin Webjet/oauth2_proxy

oauth2_proxy to force webjet employees to login via Webjet Group AAD  as fork of oauth2-proxy/oauth2-proxy

If you need to work with B2C oauth-proxy, see https://github.com/Webjet/ota-b2c-oauth-proxy

## Run and test locally

1. Clone repo Webjet/oauth2_proxy
2. Ensure that in oauth2-proxy\pipeline\oauth2_proxy.local.cfg secret values are copied FROM LastPass b2c-oauthproxy.cfg DEV.
3. You may need to stop iis to use 443 port. Run as administrator  
   iisreset /STOP
  
4.   The `client_id` and `client_secret` are configured in the application settings.
Generate a unique `client_secret` to encrypt the cookie.
to get client_id -see https://portal.azure.com/#blade/Microsoft_AAD_RegisteredApps/ApplicationMenuBlade/Overview/appId/2d9e4170-c131-46eb-bbb8-321442f55b73/isMSAApp/ 
https://community.microfocus.com/cyberres/idm/w/identity_mgr_tips/17052/creating-the-application-client-id-and-client-secret-from-microsoft-azure-new-portal---part-1

To generate a strong cookie secret use python -c 'import os,base64; print(base64.urlsafe_b64encode(os.urandom(16)).decode())' 
or https://www.tutorialspoint.com/execute_python_online.php
5. Run  admin-oauth-pipeline.sh

5. When prompt, enter email address and password valid for TSA site (e.g. for DEV NZ).

## Debug Golang container

Follow instructions in https://github.com/Webjet/land-docs/tree/debug-code-in-container/Debug-code-in-container

### Debug Callback

the final callback URL after few redirects is  
https://local.webjet.com.au/oauth2/callback?state=c337190505ab10ede31b9a5539bdd06f%3a%2f&code=...  
To invoke debugger I have to manually modify URL to  
http://localhost:13310/oauth2/callback?state=3f9257fb3af3e05785af9c66bc6e96fa%3a%2fFlightSearch%2fSignIn%2f...  
where port is dynamically changes after each restart.

## kubectl commands to check deployment
  - kubectl config use-context k8s-dev-syd-b  
  - kubectl get pods -n oauth  
  - kubectl describe pod  -n  oauth [oauthservice-pod name from list]
  - kubectl logs pod/[oauthservice-pod name from list]

## Reference to original oauth2-proxy/oauth2-proxy repository

Set upstream  
git remote add upstream git://github.com/oauth2-proxy/oauth2-proxy.git  
git remote set-url origin https://github.com/Webjet/ota-b2c-oauthproxy.git

Disable irrelevant workflows https://docs.github.com/en/actions/managing-workflow-runs/disabling-and-enabling-a-workflow  
Code scanning - action (codeql.yml)  
Mark stale issues and pull requests(stale.yml).

# Original oauth2-proxy/oauth2-proxy README
https://github.com/oauth2-proxy/oauth2-proxy/blob/master/README.md
![OAuth2 Proxy](/docs/static/img/logos/OAuth2_Proxy_horizontal.svg)

[![Build Status](https://secure.travis-ci.org/oauth2-proxy/oauth2-proxy.svg?branch=master)](http://travis-ci.org/oauth2-proxy/oauth2-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/oauth2-proxy/oauth2-proxy)](https://goreportcard.com/report/github.com/oauth2-proxy/oauth2-proxy)
[![GoDoc](https://godoc.org/github.com/oauth2-proxy/oauth2-proxy?status.svg)](https://godoc.org/github.com/oauth2-proxy/oauth2-proxy)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![Maintainability](https://api.codeclimate.com/v1/badges/a58ff79407212e2beacb/maintainability)](https://codeclimate.com/github/oauth2-proxy/oauth2-proxy/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/a58ff79407212e2beacb/test_coverage)](https://codeclimate.com/github/oauth2-proxy/oauth2-proxy/test_coverage)

A reverse proxy and static file server that provides authentication using Providers (Google, GitHub, and others)
to validate accounts by email, domain or group.

**Note:** This repository was forked from [bitly/OAuth2_Proxy](https://github.com/bitly/oauth2_proxy) on 27/11/2018.
Versions v3.0.0 and up are from this fork and will have diverged from any changes in the original fork.
A list of changes can be seen in the [CHANGELOG](CHANGELOG.md).

**Note:** This project was formerly hosted as `pusher/oauth2_proxy` but has been renamed as of 29/03/2020 to `oauth2-proxy/oauth2-proxy`.
Going forward, all images shall be available at `quay.io/oauth2-proxy/oauth2-proxy` and binaries will be named `oauth2-proxy`.

![Sign In Page](https://cloud.githubusercontent.com/assets/45028/4970624/7feb7dd8-6886-11e4-93e0-c9904af44ea8.png)

## Installation

1.  Choose how to deploy:

    a. Download [Prebuilt Binary](https://github.com/oauth2-proxy/oauth2-proxy/releases) (current release is `v7.1.3`)

    b. Build with `$ go get github.com/oauth2-proxy/oauth2-proxy/v7` which will put the binary in `$GOROOT/bin`

    c. Using the prebuilt docker image [quay.io/oauth2-proxy/oauth2-proxy](https://quay.io/oauth2-proxy/oauth2-proxy) (AMD64, ARMv6 and ARM64 tags available)

Prebuilt binaries can be validated by extracting the file and verifying it against the `sha256sum.txt` checksum file provided for each release starting with version `v3.0.0`.

```
sha256sum -c sha256sum.txt 2>&1 | grep OK
oauth2-proxy-x.y.z.linux-amd64: OK
```

2.  [Select a Provider and Register an OAuth Application with a Provider](https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/oauth_provider)
3.  [Configure OAuth2 Proxy using config file, command line options, or environment variables](https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/overview)
4.  [Configure SSL or Deploy behind a SSL endpoint](https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/tls) (example provided for Nginx)


## Security

If you are running a version older than v6.0.0 we **strongly recommend you please update** to a current version.
See [open redirect vulnerability](https://github.com/oauth2-proxy/oauth2-proxy/security/advisories/GHSA-5m6c-jp6f-2vcv) for details.

## Docs

Read the docs on our [Docs site](https://oauth2-proxy.github.io/oauth2-proxy/docs/).

![OAuth2 Proxy Architecture](https://cloud.githubusercontent.com/assets/45028/8027702/bd040b7a-0d6a-11e5-85b9-f8d953d04f39.png)

## Getting Involved

If you would like to reach out to the maintainers, come talk to us in the `#oauth2-proxy` channel in the [Gophers slack](http://gophers.slack.com/).

## Contributing

Please see our [Contributing](CONTRIBUTING.md) guidelines. For releasing see our [release creation guide](RELEASE.md).
