# Run locally

https://local.webjet.com.au/oauth2/start?rd=https://local.webjet.com.au/FlightSearch/

(Previous https://local.webjet.com.au/FlightSearch/SignIn/Index)

https://local.webjet.com.au/oauth2/start?rd=https://services.dev.webjet.com.au/web/customer/customerprofile/
https://local.webjet.com.au/oauth2/start?rd=https://services.dev.webjet.co.nz/web/customer/travellerprofile/

View AAD B2C logs: https://portal.azure.com/#blade/Microsoft_AAD_B2CAdmin/TenantManagementMenuBlade/manageAuditLogs

# Run remotely

## wcs

Dev NZ https://devci.webjet.co.nz/oauth2/start?rd=https://devCI.webjet.co.nz
Production https://flights.webjet.co.nz/oauth2/start?rd=https://flights.webjet.co.nz
logout link  
Production https://flights.webjet.co.nz/oauth2/sign_out?rd=https://webjetnz.b2clogin.com/webjetnz.onmicrosoft.com/B2C_1A_signin_TSA_Prod_NZ/oauth2/v2.0/logout?post_logout_redirect_uri=https://flights.webjet.co.nz/

## Micro-services

https://devCI.webjet.com.au/oauth2/start?rd=https://services.dev.webjet.com.au/web/customer/customerprofile/
https://devCI.webjet.com.au/oauth2/start?rd=https://services.dev.webjet.com.au/web/customer/mybookingdetails/?itineraryId=250746

Simulate "Invalid redirect" error in INCOGNITO window or logout before next login
https://devCI.webjet.co.nz/oauth2/start?ReturnUrl=https://services.dev.webjet.co.nz/web/customer/mybookingdetails/?itineraryId=79691

# Run remotely ota-b2c-oauth standalone

https://services.dev.webjet.co.nz/ota-b2c-oauth - returns to the same Sign In B2C page
https://services.dev.webjet.co.nz/ota-b2c-oauth/start?rd=https://services.dev.webjet.com.au/web/customer/customerprofile/ - returns to the same Sign In B2C page
https://services.dev.webjet.co.nz/ota-b2c-oauth-proxy/oauth2/start?rd=https://services.dev.webjet.com.au/web/customer/customerprofile/ - returns to the same Sign In B2C page
https://services.dev.webjet.com.au/ota-b2c-oauth-proxy/oauth2/start?rd=https://services.dev.webjet.com.au/web/customer/customerprofile/

Ensure that https://services.dev.webjet.co.nz/ota-b2c-oauth/callback is added to Redirect URIs in https://portal.azure.com/#blade/Microsoft_AAD_RegisteredApps/ApplicationMenuBlade/Authentication/appId/7c0fc11a-9e15-4c0b-8d6d-625acb879bb3/isMSAApp/

# Locally you can run pipeline.sh (pipeline\ota-b2c\pipeline.sh)

(ensure that in ota-b2c-oauth-proxy\pipeline\ota-b2c\oauth2_proxy.local.cfg secret values are copied FROM LastPass b2c-oauthproxy.cfg DEV)
and (after build) run
http://localhost:4180 or
http://localhost:4180/oauth2/start
or http://localhost:4180/oauth2/start?rd=https://jwt.ms //see C:\GitRepos\ota-b2c-oauth-proxy\contrib\local-environment\nginx.conf
or http://localhost:4180/oauth2/sign_in?rd=https://jwt.ms
All return after callback to the same https://webjetb2cdev.b2clogin.com/

To save oroginal url run wcs-reverse-b2c-oauth-pipeline.sh instead

http://localhost:4180/oauth2/userinfo currently returns Unauthorized

# contrib\local-environment

If install local environment from C:\GitRepos\ota-b2c-oauth-proxy\contrib\local-environment\docker-compose-nginx.yaml
run
oauth-nginx-restart.sh
or
docker-compose -f docker-compose.yaml -f docker-compose-nginx.yaml up

Then you can open in browser
http://oauth2-proxy.oauth2-proxy.localhost:8088/
(note that it may overwrite localhost:4180 from pipeline.sh )
and getting prompt for login (don't know how to create account. waiting for https://github.com/oauth2-proxy/oauth2-proxy/issues/1110 )

#wcs-reverse-proxy and b2c-oauth local
See ota-b2c-oauth-proxy\README.md#run and test locally
Run  
 wcs-reverse-b2c-oauth-pipeline.sh

#Test access to protected microservices
Use in INCOGNITO to avoid existing cookies
https://services.dev.webjet.co.nz/web/customer/travellerprofile/
https://services.dev.webjet.com.au/web/customer/travellerprofile/

https://services.dev.webjet.co.nz/web/customer/mybookings/

#Test internal(aak Admin) oauth2 proxy
https://services.dev.webjet.co.nz/oauth2
https://services.dev.webjet.com.au/oauth2

https://services.dev.webjet.com.au/admin/web/portal/itinerarylist/ItineraryList/Index

#local-environment -replace local-environment example with b2c-oauth POC not completed

docker-compose -f docker-compose-b2c-oauth.yaml -f docker-compose-nginx-b2c-oauth.yaml up (or down)

Access one of the following URLs to initiate a login flow: - doesn't work http://oauth2-proxy.localhost or http://oauth2-proxy.localhost:8181

- http://httpbin.oauth2-proxy.localhost

  The OAuth2 Proxy itself is hosted at http://oauth2-proxy.oauth2-proxy.localhost:8181/
