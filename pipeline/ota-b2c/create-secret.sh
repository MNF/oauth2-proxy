#1.Copy content from LastPass  b2c-oauthproxy.cfg Dev/ b2c-oauthproxy.cfg Prod 
#2.Uncomment one of curl at the bottom of the file(for appropriate secret) 
#3.Do NOT commit the changes of this file
#4 You may need re-deploy to apply the change (in devservices you can kill pods instead)

# bots port-forward needs to be run only once, for subsequent calls comment the line or ignore  "Unable to listen on port 10010:" 
#kubectl config use-context k8s-dev-workloads-syd
#kubectl -n bots port-forward svc/kubebot 10010:80 
PROXYCONFIG=$(echo -n '
# Insert THE BLOCK below From LastPass b2c-oauthproxy.cfg Dev or Prod
' | base64 | tr -d '\n')
echo "
kind: Secret
apiVersion: v1
metadata:
  name: '[variables('KUBE_BOT_SECRET_NAME')]'
  namespace: serviceportal
type: Opaque 
data:
  oauth2_proxy.cfg: $PROXYCONFIG
" | 
curl -X POST http://localhost:10010/secret/dev/b2c-oauthproxy-wjau --data-binary @- -H 'Content-Type: application/yaml'
#curl -X POST http://localhost:10010/secret/dev/b2c-oauthproxy-wjnz --data-binary @- -H 'Content-Type: application/yaml'
#curl -X POST http://localhost:10010/secret/prod/b2c-oauthproxy-wjau --data-binary @- -H 'Content-Type: application/yaml'
#curl -X POST http://localhost:10010/secret/prod/b2c-oauthproxy-wjnz --data-binary @- -H 'Content-Type: application/yaml'
read -p "Press [Enter] key ..."