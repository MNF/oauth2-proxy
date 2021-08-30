#ensure that previous version uninstalled
docker-compose -f docker-compose.yaml -f docker-compose-nginx.yaml down
# Alternatively:#    make nginx-<command> (eg make nginx-up, make nginx-down)
start http://oauth2-proxy.oauth2-proxy.localhost:8088/
docker-compose -f docker-compose.yaml -f docker-compose-nginx.yaml up