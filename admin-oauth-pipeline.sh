#ensure that previous version uninstalled down
docker-compose -f admin-oauth-compose.yaml down
PROJECTNAME="ota-b2c-oauth-proxy"
docker rm -f $(docker ps -f name=$(PROJECTNAME) -q)
docker-compose  -f admin-oauth-compose.yaml  build 
# Alternatively:#    make nginx-<command> (eg make nginx-up, make nginx-down)
#start https://local.webjet.com.au/oauth2/start?rd=https://local.webjet.com.au/FlightSearch/
#start https://local.webjet.com.au/oauth2/start?rd=https://services.dev.webjet.com.au/web/customer/customerprofile/
#start https://local.webjet.com.au/FlightSearch/
docker-compose  -f admin-oauth-compose.yaml up
read -p "Press [Enter] key ..."
