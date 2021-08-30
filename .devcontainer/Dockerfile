# you can use same build image here or you can replace it with any newer image for test
# ensure that the version is consistent with main dockerfile golang version
FROM containerregistrydev.azurecr.io/library/golang:1.16-alpine As builder

RUN apk update && apk upgrade
# by default cgoenabled=1 to run go in linux and it requires adding gcc and g++
RUN apk add git gcc g++

# install dlv for debugging
RUN go get -v github.com/go-delve/delve/cmd/dlv

# other Go dev tools required by VS code
RUN go get -v github.com/ramya-rao-a/go-outline
# gopls (go please language server provides the language intelligence etc; it runs ok if installing in built container but not in docker at the moment)
RUN go get -v golang.org/x/tools/gopls

WORKDIR /workspace

COPY . .
