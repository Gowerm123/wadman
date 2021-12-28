BASEPATH=$1

go build .
sudo cp dwpm /usr/bin/dwpm
sudo rm -rf $BASEPATH

sudo mkdir $BASEPATH
sudo touch $BASEPATH/.pkglist
sudo touch $BASEPATH/.config