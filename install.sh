#!/bin/bash

echo "Installing wadman..."
echo "Creating /usr/share/wadman directory..."
BASEPATH="/usr/share/wadman/"
PKGLIST=$BASEPATH/.pkglist
CONFIG=$BASEPATH/.config
if [ ! -d "$BASEPATH" ]; 
then
    sudo mkdir $BASEPATH
    sudo touch $PKGLIST
    sudo touch $CONFIG
else
    echo "directory exists....skipping...."
fi


echo "building binary and copying to /usr/bin/wadman"
CGO_ENABLED=0 GOOS=linux go build .
sudo cp wadman /usr/bin/wadman
sudo cp README.md /usr/share/wadman/README.md


echo "writing to config file"
CONTENTS=$(tail $CONFIG)

if [ "$CONTENTS" == "" ]
then
    echo "{" > $CONFIG
    echo "  \"launcher\": \"gzdoom\"," >> $CONFIG
    echo "  \"launchArgs\": []," >> $CONFIG
    echo "  \"iwads\": {}," >> $CONFIG
    echo "  \"installDir\": \"$BASEPATH\"" >> $CONFIG
    echo "}" >> $CONFIG
else
    echo "config file isn't empty...skipping..."
fi