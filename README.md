What is it ?
============

IDOK (kodi reversed name) is a Go program that allows to serve medias to your Kodi plateform (raspbmc, xbmc...)

You may be able to send video, images and music from your computer.

Idok has got two modes:

* your computer serve media from a standard port (default 8080)
* your computer dig a tunnel and serve media

Installation
============

## Install distribution

Linux users can use the auto-install:

	bash <(wget https://github.com/metal3d/idok/releases/download/0.2.6/install-idok.sh -qO -)

Or with curl:

	bash <(curl -L https://github.com/metal3d/idok/releases/download/0.2.6/install-idok.sh)

Check that ~/.local/bin is in your PATH. Then try to call:

	idok -h

If you want to get yourself the packed file for Linux, here are the urls:

* https://github.com/metal3d/idok/releases/download/0.2.4/idok-i686.gz
* https://github.com/metal3d/idok/releases/download/0.2.4/idok-x86_64.gz



That command could not work with MacOSX. You can get the gziped binary there: 
https://github.com/metal3d/idok/releases/download/0.2.4/idok-darwin.gz

Then gunzip the binary:

	gunzip idok-darwin.gz

(I need help for Mac because I don't have one and cannot be sure of how to install the command at the right path...)

Windows users can get exe:
https://github.com/metal3d/idok/releases/download/0.2.4/idok.zip
The "idok.exe" file should be launched from command line (cmd command). 

Windows users (again) may know that there is no graphical interface for the idok tool. Maybe one day...

If you have troubles, please fill an issue. But keep in mind that I don't have any Windows or Mac OSX installation. 

Stream your first media
=======================

## Youtube URL

The simplier command is to open youtube url. Idok reoconize "youtu.be" and "www.youtube.com" hosts. Example:

	idok -target=raspbmc.local "https://www.youtube.com/watch?v=o5snlP8Y5GY"

As next command examples, pressing CTRL+C will stop video stream, and the command exits when video ends.

**Note: you must enable youtube addon on your kodi/XBMC installation.**


## HTTP (default)

The HTTP way is not secured. While you're streaming to Kodi (or XBMC), the media can be accessed by other computer in your network. That's not a big problem while you're not streaming important information (restricted video). 

This solution need to open port on your firewall. 

By default, idok opens 8080 port (http-alternative), but you can specify other port using "-port" option.

At first, you must be sure that the port is opened. To open firewall port on you linux installation:

	firewall-cmd --add-port=8080/tcp

When you will reload firewall, or restart computer, the port will be closed. If you want to keep that port opened:

	firewall-cmd --add-port=8080/tcp --permanent

Then, send media:

	idok -target=IP_OF_KODI_OR_XBMC /path/to/media.mp3

If you've opened other port, you can set it. For example for port 1234:

	idok -port=1234 -target=IP_OF_KODI_OR_XBMC /path/to/media.mp3


## SSH

The SSH way is the easier and more secured way. Easier because you don't have to open port on your computer and only the Kodi instance will be able to access your content.

	idok -ssh -target=IP_OF_RASPBERRY /path/to/media.mp3

Your kodi should open the file.

Pressing CTRL+C should stop media stream and exit program.

With SSH, idok tries to use your ssh key pair to authenticate. If it fails, it will use login/password to auth. So, there are 2 possibilities:

* copy you public key to the kodi/xbmc host
* set -sshuser (default is "pi") and -sshpass options

To copy you key, type this command:

	ssh-copy-id USER@KODI_HOST

Where USER is the ssh kodi user ("pi" on raspbmc) and KODI_HOST is ip or hotname of the kodi host. By default, raspbmc use "raspberry" as password.

Now, should should be able to stream media without the need of password.

**Note: If you compiled yourself, remember to patch go.crypto/ssh package as explained above. Dropbear on raspbmc + crypto package are not compatibles without my patch**

Install from source
===================

**WARNING - Because there is a problem with dropbear ssh server on raspbmc, you should patch go.crypto/ssh package with the patched I made. See:
https://code.google.com/p/go/issues/detail?id=8657**

You can clone repository and compile source code yourself:

	git clone http://git.develipsum.com/metal3d/idok.git
	cd idok
	go build idok.go

Then you can put binary in your PATH:

	mkdir -p ~/.local/bin
	cp idok ~/.local/bin

Options
=======

There are other options that may be usefull:

* -target: kodi instance ip or hostname 
* -login : xbmc or kodi login configured on web interface settings
* -password : xbmc or kodi password configured on web interface settings
* -ssh : If set, idok will dig ssh tunnel to stream content
* -sshuser : if you don't user "pi" user
* -sshpass : if you changed standard password of "pi" user
* -sshport : if you changed standard ssh port or to use other ssh server (default is 22)
* -port : local port for media stream if you don't use ssh tunneling, default is 8080

TODO
====

I should ask dropbear maintainer why the tunneling won't work. Installing openssh-server on raspbmc can be complicated for some users.

Add kodi/xbmc port option (some users changed the default 80) 

ChangeLog
=========

* 0.2.6
  - fix issue #1 (reported at http://forum.xbmc.org/showthread.php?tid=203834)
  - Now, if idok wait reponse to check playing status
  - Add DSA key managment
  - Add some information in -verion option

* 0.2.4
  - Add youtube url support

* 0.2.2
  - Make use of SSH key pair
  - Fixed standard ssh tunnel on dropbear server https://code.google.com/p/go/issues/detail?id=8657
  - Fix randomized port for dropbear based ssh server
