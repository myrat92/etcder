#!/bin/bash

if [ $# -gt 0 ]; then
    version=$1
else
    version="0.0"
fi

## windows amd64
fyne-cross windows -arch=amd64 -app-id=etcder.myrat -app-version=${version} -icon=docs/logos/etcder.png -env GOPROXY=https://goproxy.cn
#
## mac amd64
fyne-cross darwin -arch=amd64 -app-id=etcder.myrat -app-version=${version} -icon=docs/logos/etcder.png -env GOPROXY=https://goproxy.cn
#
## mac arm64
fyne-cross darwin -arch=arm64 -app-id=etcder.myrat -app-version=${version} -icon=docs/logos/etcder.png -env GOPROXY=https://goproxy.cn

# linux amd64
fyne-cross linux -arch=amd64 -app-id=etcder.myrat -app-version=${version} -icon=docs/logos/etcder.png -env GOPROXY=https://goproxy.cn

# linux arm64
fyne-cross linux -arch=arm64 -app-id=etcder.myrat -app-version=${version} -icon=docs/logos/etcder.png -env GOPROXY=https://goproxy.cn