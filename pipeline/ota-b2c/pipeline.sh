#!/bin/bash
#to test with wcs-reverse run wcs-reverse-b2c-oauth-pipeline.sh instead
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
WORKSPACE="$(dirname "$DIR")"

echo "~~~~~~~~~~~~~~WORKSPACE" + $WORKSPACE

function check_error_exit
{
    if [ "$?" = "0" ]; then
	    echo "Stage succeeded!"
    else
        echo "Stage Error!" 1>&2
		read -p "Press [Enter] key ..."
        exit 1
    fi
}
PROJECTNAME="b2c-oauthproxy"  #must be lowercase
IMAGETAG="local"
# DEPLOY_TO_PROD = "FALSE" ## deployed by jenkins no need to set

echo "login"
docker login -u $CONTAINERREGISTRYUSERNAME -p $CONTAINERREGISTRYPASSWORD $CONTAINERREGISTRY

##remove previous container (is it ok in prod?) - cause error  "docker rm" requires at least 1 argument. https://docs.docker.com/engine/reference/commandline/rm/#examples #docker rm -f $(docker ps -aq) all containers
docker rm -f $(docker ps -f name=$(PROJECTNAME) -q)

echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ building image"
cd ../..
docker build --build-arg REGISTRY=${CONTAINERREGISTRY} -t "${PROJECTNAME}":$IMAGETAG  . -f Dockerfile 
check_error_exit

#docker build --build-arg IMAGE_NAME="${PROJECTNAME}_builder" --build-arg IMAGE_VERSION=${IMAGETAG} -t "containerregistrydev.azurecr.io/webjet/${PROJECTNAME}:${IMAGETAG}"  .
#check_error_exit

echo "~~~~~~~~~~~~~~~~~~ Run application: ${PROJECTNAME}:${IMAGETAG} " 
#docker run -d -p 7000:80 bsp:local
#winpty docker run -it --rm --name ${PROJECTNAME} -p 4180:4180 containerregistrydev.azurecr.io/webjet/${PROJECTNAME}:${IMAGETAG}
#expect in dockerfile COPY pipeline/ota-b2c/oauth2_proxy.local.cfg /etc/config/oauth2_proxy.cfg
#--config=/etc/config/oauth2_proxy.cfg  C:/GitRepos/ota-b2c-oauth-proxy/pipeline/ota-b2c/oauth2_proxy.local.cfg
winpty docker run -it --rm --name ${PROJECTNAME} -p 4180:4180 ${PROJECTNAME}:${IMAGETAG}  --config=etc/config/oauth2_proxy.cfg
check_error_exit

# ##PUSH TO REGISTRY - UNCOMMENT IF NEED TO PUSH AND DEPLOY
# echo "Push Production Image ${PROJECTNAME}:${IMAGETAG}"
# docker push "containerregistrydev.azurecr.io/webjet/${PROJECTNAME}:${IMAGETAG}"
# check_error_exit

#./deploy-kube.sh "$DOCKERIMAGE" "$DOCKERIMAGETAG" "$DOCKERCONTAINERPORT" "$SERVICE_NAME" "$AZ_STORAGENAME" "$AZ_STORAGEKEY" "$KUBE_SECRET_WJAU" "$SERVICE_HOST_DEV" "$SERVICE_PATH" "$KUBE_NAMESPACE"
#check_error_exit

echo "before running ensure that in ota-b2c-oauth-proxy\pipeline\ota-b2c\oauth2_proxy.local.cfg secret values are copied FROM LastPass b2c-oauthproxy.cfg DEV"
#start http://localhost:4180/
start http://localhost:4180/oauth2/start
read -p "Press [Enter] key ..."