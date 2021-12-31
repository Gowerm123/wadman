BASEPATH="/usr/share/wadman/"
PKGLIST=$BASEPATH/.pkglist
CONFIG=$BASEPATH/.config

go build .
sudo cp wadman /usr/bin/wadman

sudo mkdir $BASEPATH
sudo touch $PKGLIST
sudo touch $CONFIG

echo "{" > $CONFIG
echo "  \"launcher\": \"gzdoom\"," >> $CONFIG
echo "  \"launchArgs\": []," >> $CONFIG
echo "  \"iwads\": {}," >> $CONFIG
echo "  \"installDir\": \"$BASEPATH\"" >> $CONFIG
echo "}" >> $CONFIG