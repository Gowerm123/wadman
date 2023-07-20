#!/bin/bash
if [ $(id -u) != "0" ]
then
    sudo "$0" "$@"
    exit $?
fi

echo "Installing wadman..."
echo "Creating .wadman directory..."

USER_HOME=$(getent passwd $SUDO_USER | cut -d: -f6)
BASEPATH="$USER_HOME/.wadman"
WADMANIFEST=$BASEPATH/wadmanifest.json
CONFIG_PATH=$USER_HOME/.config/wadman-config.json
BASE_MIRRORS="[\"http://mirrors.syringanetworks.net\",\"http://www.quaddicted.com\",\"http://ftpmirror1.infania.net\"]"

if [ ! -d "$BASEPATH" ]; 
then
    sudo mkdir $BASEPATH
    sudo touch $WADMANIFEST
else
    echo "directory exists....skipping...."
fi

echo "creating config file"
if [ ! -d "$CONFIG_PATH" ];
then
    sudo touch $CONFIG_PATH
else
    echo "config exists...skipping"
fi


echo "building binary and copying to /usr/bin/wadman"
CGO_ENABLED=0 GOOS=linux go build .
sudo cp wadman /usr/bin/wadman
sudo cp README.md $BASEPATH/README.md


echo "writing to config file"
CONTENTS=$(cat $CONFIG_PATH)
if [ "$CONTENTS" == "" ]
then
    sudo echo "{" >> $CONFIG_PATH
    sudo echo "  \"launcher\": \"gzdoom\"," >> $CONFIG_PATH
    sudo echo "  \"launchArgs\": []," >> $CONFIG_PATH
    sudo echo "  \"iwads\": {}," >> $CONFIG_PATH
    sudo echo "  \"mirrors\": $BASE_MIRRORS" >> $CONFIG_PATH
    sudo echo "}" >> $CONFIG_PATH
else
    echo "config file isn't empty...skipping..."
fi