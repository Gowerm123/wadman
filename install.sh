BASEPATH=$1
CONFIGPATH='/usr/share/.dwpmConfig'

go build .
sudo cp dwpm /usr/bin/dwpm
sudo rm -rf $BASEPATH

sudo mkdir $BASEPATH
sudo touch $BASEPATH/.pkglist
sudo touch $CONFIGPATH

echo "{" > $CONFIGPATH
echo "  \"launcher\": \"gzdoom\"," >> $CONFIGPATH
echo "  \"launchArgs\": []," >> $CONFIGPATH
echo "  \"iwads\": {}," >> $CONFIGPATH
echo "  \"installDir\": \"$BASEPATH\"" >> $CONFIGPATH
echo "}" >> $CONFIGPATH