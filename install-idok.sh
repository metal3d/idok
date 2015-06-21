#!/bin/bash

echo "This will install idok on your computer."
echo "Select installation type"

DIR="~/.local/bin"
PREFIXCMD="bash -c"
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

[ $ARCH == "amd64" ]                                 && ARCH="x86_64"
[ $ARCH == "ia64" ]                                  && ARCH="x86_64"
[ $ARCH == "ia32" ]                                  && ARCH="i686"
[ $ARCH == "i386" ]                                  && ARCH="i686"
[ x"$OSTYPE" == x"darwin" ]                          && ARCH="darwin"
[ x"$OSTYPE" == x"freebsd" ] && [ $ARCH == "i686" ]  && ARCH="freebsd32"
[ x"$OSTYPE" == x"freebsd" ] && [$ARCH == "x86_64" ] && ARCH="freebsd64"

CMD="wget -q -O -"
[ -x $(which curl) ] && CMD="curl -# -X GET -L"

URL=$($CMD "https://api.github.com/repos/metal3d/idok/releases" 2>/dev/null | awk -NF":" '
    BEGIN{
        ok=0
    }
    {
        if (/"prerelease"\s*:\s*false/) {
            ok=1
        }
        if (/browser_download_url/ && /idok-'$ARCH'/ && ok == 1){
            print $2 ":" $3
            exit 0
        }
    }
')

COMMAND="$CMD $URL | gunzip -c > $DIR/idok"
bash -c "$PREFIXCMD \"$COMMAND\""
sleep 1
bash -c "$PREFIXCMD \"chmod +x $DIR/idok\""
echo
echo "Installation ok, idok installed in $DIR"
exit 0
