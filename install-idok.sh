#!/bin/bash

VERSION="@VERSION@"
echo "This will install idok on your computer."
echo "Select installation type"



select choice in \
	"Install for the current user $USER" \
	"Install for all user in /usr/local/bin (will use sudo)" \
	"Cancel"
do
	case $REPLY in
		1)	DIR=~/.local/bin
			PREFIXCMD="bash -c"
			break
			;;
		2)	DIR="/usr/local/bin"
			PREFIXCMD="sudo bash -c "
			break
			;;
		3) echo "Install cancelled, bye"; exit 0
			;;
		*) echo "Not valid answer"
			;;
	esac
done

ARCH=$(uname -i)

[ $ARCH == "amd64" ] && ARCH="x86_64"
[ $ARCH == "ia64" ] && ARCH="x86_64"
[ $ARCH == "ia32" ] && ARCH="i686"
[ $ARCH == "i386" ] && ARCH="i686"

URL="https://github.com/metal3d/idok/releases/download/$VERSION/idok-$ARCH.gz"
GET=$(which curl)
if [ $? == 0 ]; then
	CMD="curl -L $URL"
else
	CMD="wget $URL -qO -"
fi

echo $URL

COMMAND="$CMD | gunzip -c > $DIR/idok"
bash -c "$PREFIXCMD \"$COMMAND\""
sleep 1
bash -c "$PREFIXCMD \"chmod +x $DIR/idok\""
echo
echo "Installation ok, idok installed in $DIR"
exit 0
