#!/bin/bash

echo "Installing wadman..."
echo "Creating .wadman directory..."
USER_HOME=$(getent passwd $SUDO_USER | cut -d: -f6)
BASEPATH="$USER_HOME/.wadman"
PKGLIST=$BASEPATH/.pkglist.json
CONFIG=$USER_HOME/.config/wadman-config.json
if [ ! -d "$BASEPATH" ]; 
then
    sudo mkdir $BASEPATH
    sudo touch $PKGLIST
else
    echo "directory exists....skipping...."
fi

echo "creating config file"
if [ ! -d "$CONFIG" ];
then
    sudo touch $CONFIG
else
    echo "config exists...skipping"
fi


echo "building binary and copying to /usr/bin/wadman"
CGO_ENABLED=0 GOOS=linux go build .
sudo cp wadman /usr/bin/wadman
sudo cp README.md $BASEPATH/README.md


echo "writing to config file"
CONTENTS=$(tail $CONFIG)
echo "CONTENTS $CONTENTS"
if [ "$CONTENTS" == "" ]
then
    sudo echo "{" >> $CONFIG
    sudo echo "  \"launcher\": \"gzdoom\"," >> $CONFIG
    sudo echo "  \"launchArgs\": []," >> $CONFIG
    sudo echo "  \"iwads\": {}," >> $CONFIG
    sudo echo "  \"installDir\": \"$BASEPATH/\"" >> $CONFIG
    sudo echo "}" >> $CONFIG
else
    echo "config file isn't empty...skipping..."
fi