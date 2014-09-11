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

	bash <(wget https://github.com/metal3d/idok/releases/download/20140910-1/install-idok.sh -qO -)

Or with curl:

	bash <(curl -L https://github.com/metal3d/idok/releases/download/20140910-1/install-idok.sh)

Check that ~/.local/bin is in your PATH. Then try to call:

	idok -h

If you want to get yourself the packed file for Linux, FreeBSD, MacOSX and Windows, here are the urls:

* [Linux x64](https://github.com/metal3d/idok/releases/download/20140910-1/idok-x86_64.gz)
* [Linux i686](https://github.com/metal3d/idok/releases/download/20140910-1/idok-i686.gz)
* [FreeBSD 64 bits](https://github.com/metal3d/idok/releases/download/20140910-1/idok-freebsd64.gz)
* [FreeBSD 32 bits](https://github.com/metal3d/idok/releases/download/20140910-1/idok-freebsd32.gz)
* [MacOSX 64 bits](https://github.com/metal3d/idok/releases/download/20140910-1/idok-darwin.gz)
* [Windows 64 bits](https://github.com/metal3d/idok/releases/download/20140910-1/idok.zip)

(I need help for Mac because I don't have one and cannot be sure of how to install the command at the right path...)

Windows users may know that there is no graphical interface for the idok tool. Maybe one day...

If you have troubles, please fill an issue. But keep in mind that I don't have any Windows or Mac OSX installation. 

Stream medias
=============

## Youtube URL

Open a youtube url is simple

	idok -target=YOUR_KODI_IP "https://www.youtube.com/watch?v=o5snlP8Y5GY"

This will ask XBMC/Kodi to open this video. This doesn't stream video from your computer, so that's not use port opening and/or ssh tunnel.

**Note: you must enable youtube addon on your kodi/XBMC installation.**

## distant medias

You can open http, rtsp, mms, rtpm... media. That doesn't make usage of ssh or local port. Kodi will connect itself to the stream:

	idok -target=YOUR_KODI_IP scheme://url

Where "scheme://url" can be "rtpm://...", "http://...", etc... 

For example, to open "Tears of Steel" movie:

	idok -target=YOUR_KODI_IP http://ftp.halifax.rwth-aachen.de/blender/demo/movies/ToS/ToS-4k-1920.mov


Open Jamendo rock radio (creative commons musics):

	idok -target=YOUR_KODI_IP https://streaming.jamendo.com/JamRock


## Stream your local media through HTTP (default)

To open a media that resides on your computer:

	idok -target=IP_OF_KODI_OR_XBMC /path/to/media.mp3

That command open port 8080 (http-alt) to stream media. If you want to use another port:

	idok -port=1234 -target=IP_OF_KODI_OR_XBMC /path/to/media.mp3

**Note**

This solution needs to open port on your firewall. 

You must be sure that the port is opened. On Linux, to open firewall port on you linux installation:

	firewall-cmd --add-port=8080/tcp

When you will reload firewall, or restart computer, the port will be closed. If you want to keep that port opened:

	firewall-cmd --add-port=8080/tcp --permanent

## Stream your local media throught SSH Tunnel

Idok can stream media through ssh tunnel. That way, you don't need to configure firewall.

	idok -ssh -target=IP_OF_RASPBERRY /path/to/media.mp3

Your kodi should open the file.

Pressing CTRL+C should stop media stream and exit program.

**Note**

With SSH, idok tries to use your ssh key pair to authenticate. If it fails, it will use login/password to auth. So, there are 2 possibilities:

* copy you public key to the kodi/xbmc host (with ssh-copy-id for example)
* set -sshuser (if user is not "pi") and -sshpass options (if password is not "raspberry")

To copy you key, type this command:

	ssh-copy-id USER@KODI_HOST

Where USER is the ssh kodi user ("pi" on raspbmc, "root" for openelec) and KODI_HOST is ip or hotname of the kodi host. By default, raspbmc use "raspberry" as password, "password" for openelec.

Now, should should be able to stream media without the need of password.

Some other streams you can make
===============================

The -stdin option is a cool new functionnality that "ianux" (an user on DLFP pages) asked me... I took this feature as a challenge ;) And it works !

Since Idok can now use this option, I discovered that I'm able to make a lot of nice stream to my Kodi installation.

## Gstreamer - screencast to kodi

Gstreamer can be used to stream medias to stdout using "fdsink" or "filesink location=/dev/stdout". 

If you're using operating system that can be able to launch gstreamer pipelines, here is a nice "screencast stream":

	gst-launch-1.0 -q ximagesrc remote=1 ! videoconvert  ! avenc_mpeg4 ! mpegtsmux ! filesink location=/dev/stdout | idok -stdin -ssh -target=YOUR_KODI_IP

Remove "remote=1" on "non fedora 20", this option is needed as far as I know on Fedora 20 (reported bug)

## livestreamer - ISS station from space from ustream

Livestreamer is a python tool that is able to fectch streams from some servers and is able to give an url. For some streams, it's impossible to get URL, but we can use "-O" option that dump stream to stdout. So...

	livestreamer http://www.ustream.tv/channel/iss-hdev-payload 480p -Q -O | idok -stdin -ssh -target=YOUR_KODI_IP

That will launch the ISS live video from space (sometimes the image is black because ISS station is on the night side. Wait 5 minutes and you will see...)


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
* -targetport : Kodi jsonrpc port to connect, default is 80
* -login : xbmc or kodi login configured on web interface settings
* -password : xbmc or kodi password configured on web interface settings
* -ssh : If set, idok will dig ssh tunnel to stream content. Not used for youtube url or scheme url
* -sshuser : if you don't user "pi" user
* -sshpass : if you changed standard password of "pi" user
* -sshport : if you changed standard ssh port or to use other ssh server (default is 22)
* -port : local port for media stream if you don't use ssh tunneling, default is 8080
* -stdin: Read media from stdin and stream this to Kodi.

TODO
====

- Refactorisation to be more maintainable
- GUI (or not...)
- better lookup adresse for -target option


ChangeLog
=========

* 20140910-1
  - Can now open streams through stdin (livestreamer, gstreamer, etc...) (thanks to "anxt" on gstreamer irc channel, thanks to ianux on DLFP that gives the idea)
  - Fixes some needed options
  - prepared to be refactored
  - version number is now "modern" (thanks user "Baud" on DLFP)
  - Fixes documentation
  - Prepare examples

* 0.2.6
  - can now open http stream or video other that youtube
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
