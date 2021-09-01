# Run remotely

## Micro-services
When testing changes, consider to use INCOGNITO or clear cookies
https://services.dev.webjet.com.au/admin/web/portal/itinerarylist/ItineraryList/
https://services.dev.webjet.com.au/admin/web/packages/serenityservice/
https://services.dev.webjet.com.au/admin/web/ui/pricelookup/
https://services.dev.webjet.com.au/admin/web/packages/ai-controller/
https://services.dev.webjet.com.au/admin/web/serviceportal/bsp-refund/
https://services.dev.webjet.com.au/admin/api/serviceportal/bsp-refund/


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


