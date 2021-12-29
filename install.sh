BASEPATH=$1
CONFIGPATH="/usr/share/.wadmanConfig"

go build .
sudo cp wadman /usr/bin/wadman
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