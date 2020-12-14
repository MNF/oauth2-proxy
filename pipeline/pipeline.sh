
#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
WORKSPACE="$(dirname "$DIR")"

DOCKERIMAGE="containerregistrydev.azurecr.io/webjet/oauthservice"
DOCKERIMAGETAG="LOCAL"

function check_error_exit
{
    if [ "$?" = "0" ]; then
	    echo "Stage succeeded!"
    else
        echo "Stage Error!" 1>&2
        exit 1
    fi

}

PROJECTNAME="oauthservice"
IMAGETAG="02"

echo "loging into container registry"
docker login -u $CONTAINERREGISTRYUSERNAME -p $CONTAINERREGISTRYPASSWORD $CONTAINERREGISTRY

echo "building source image"
cd $WORKSPACE
docker build  -t "${PROJECTNAME}_src":$IMAGETAG  .
check_error_exit

echo "building release image"
cd $WORKSPACE/pipeline
docker build --build-arg IMAGE_NAME="${PROJECTNAME}_src" --build-arg IMAGE_VERSION=${IMAGETAG} -t "containerregistrydev.azurecr.io/webjet/${PROJECTNAME}:${IMAGETAG}"  .
check_error_exit

# echo "run locally"
# docker run -it --rm -p 10011:10011 containerregistrydev.azurecr.io/webjet/${PROJECTNAME}:${IMAGETAG}
# docker run -it containerregistrydev.azurecr.io/webjet/${PROJECTNAME}:${IMAGETAG} bash

# echo "pushing release image"
# docker push "containerregistrydev.azurecr.io/webjet/${PROJECTNAME}:${IMAGETAG}"
# check_error_exit

# echo "deploy to cluster in dev" 
# curl -X POST http://localhost:10010/deploy/dev/oauth/${PROJECTNAME}/${IMAGETAG} --data-binary "@$WORKSPACE/pipeline/dev-wjau.yaml" -H 'Content-Type: application/yaml'
# check_error_exit

#