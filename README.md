Idok
====

[![status](https://sourcegraph.com/api/repos/github.com/metal3d/idok/.badges/status.svg)](https://sourcegraph.com/github.com/metal3d/idok)
[![library users](https://sourcegraph.com/api/repos/github.com/metal3d/idok/.badges/library-users.svg)](https://sourcegraph.com/github.com/metal3d/idok)

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

### With "go get"

If you've got a go installation and $GOPATH configured, you can install idok with "go get":

	go get github.com/metal3d/idok

### Easy install

Linux users can use the auto-install:

	bash <(wget https://goo.gl/imm9jP -qO -)

Or with curl:

	bash <(curl -L https://goo.gl/imm9jP)

Check that ~/.local/bin is in your PATH. Then try to call:

	idok -h

### Other install method

If you want to get yourself the packed file for Linux, FreeBSD, MacOSX and Windows, here are the urls:

Visit the release page: https://github.com/metal3d/idok/releases and pick the version for you OS. Then unpack the binary and put it in your PATH.

Windows users may know that there is no graphical interface for the idok tool. Maybe one day...

If you have troubles, please fill an issue. But keep in mind that I don't have any Windows or Mac OSX installation. 

Stream medias
=============

## Youtube URL

Open a youtube url is simple

	idok -target=YOUR_KODI_IP "https://www.youtube.com/watch?v=o5snlP8Y5GY"

This will ask XBMC/Kodi to open this video. This doesn't stream video from your computer, so that's not use port opening and/or ssh tunnel.

**Note: you must enable youtube addon on your kodi/XBMC installation.**

## Distant medias (http, ftp, and so on...)

You can open http, rtsp, mms, rtpm... media. That doesn't make usage of ssh or local port. Kodi will connect itself to the stream:

	idok -target=YOUR_KODI_IP scheme://url

Where "scheme://url" can be "rtpm://...", "http://...", etc... 

For example, to open "Tears of Steel" movie:

	idok -target=YOUR_KODI_IP \
	http://ftp.halifax.rwth-aachen.de/blender/demo/movies/ToS/ToS-4k-1920.mov


Open Jamendo rock radio (creative commons musics):

	idok -target=YOUR_KODI_IP https://streaming.jamendo.com/JamRock


## Stream your local media through HTTP (default)

To open a media that resides on your computer:

	idok -target=IP_OF_KODI_OR_XBMC /path/to/media.mp3

That command open LOCAL port 8080 (http-alt) to stream media. If you want to use another port:

	idok -port=1234 -target=IP_OF_KODI_OR_XBMC \
	/path/to/media.mp3

If your kodi installation use another port for jsonrpc, you may change "target port":

	idok -target-port=80 -target=IP_OF_KODI_OR_XBMC \
	/path/to/media.mp3

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

Configuration File
==================

To not repeat options each time you want to run idok, you can create a configuation file to keep recurrent values.

Idok will check if config files exists in that order:

- ./idok.conf (current directory)
- $HOME/.config/idok/idok.conf
- $HOME/.local/etc/idok.conf
- /etc/idok.conf

If one of this file is found, the next will not be parsed.

The command line option can always override configuration file values.

To get an example of the configuration, type this command:

	$ idok -conf-example
	# blank value means default
	#
	# Idok checks if that file exists in that order:
	# - ./idok.conf
	# - $HOME/.config/idok/idok.conf
	# - $HOME/.local/etc/idok.conf
	# - /etc/idok.conf
	#
	# If option is given, it override the configuration file corresponding option
	#
	# IP or hostname of Kodi/XBMC
	# (-target)
	target = 

	# port to connect jsonrpc
	# (-targetport)
	targetport = 

	# Kodi/XBMC jsonrpc username and password
	# (-login -password)
	login = 
	password = 

	# SSH user/password (user is needed for ssh even if you sent keypair)
	# (-sshuser -sshpass)
	sshuser =
	sshpass =

	# if you changed ssh port from 22 to other
	# (-sshport)
	sshport = 

	# force ssh usage (true or false)
	# (-ssh)
	ssh = 

You can easilly prepare configuration:

	$ mkdir -p ~/.config/idok/ && idok -conf-example > ~/.config/idok/idok.conf

Then edit ~/.config/idok/idok.conf file to change values.

That way, you will be able to launch idok without giving target, port, sshuser, and so on...

Some other streams you can make
===============================

The -stdin option is a cool new functionnality that "ianux" (an user on DLFP pages) asked me... I took this feature as a challenge ;) And it works !

Since Idok can now use this option, I discovered that I'm able to make a lot of nice stream to my Kodi installation.

## Direct encoding and stream

You may have some files that won't be opened on Kodi. Sometimes because format is badly converted, or because you only have codecs on your local machine.

To test format converstion, you can use "ffmpeg". Following example encode a badly encoded avi file to matroska:

    ffmpeg -i BadFile.avi -f matroska - | idok -stdin
	
If you didn't made configuration file:

	ffmpeg -i BadFile.avi -f matroska - | \
	idok -target=YOUR_KODI_IP -stdin 
	
Last minus sign of ffmpeg means "stream to STDOUT". Idok will read from STDIN (piped) and stream data to Kodi/XBMC server (you may use -ssh -port as explained below).

## Gstreamer - screencast to kodi

Gstreamer can be used to stream medias to stdout using "fdsink" or "filesink location=/dev/stdout". 

If you're using operating system that can be able to launch gstreamer pipelines, here is a nice "screencast stream":

	gst-launch-1.0 -q ximagesrc remote=1 ! videoconvert ! \
	avenc_mpeg4 ! \
	mpegtsmux ! filesink location=/dev/stdout | \
	idok -stdin -ssh -target=YOUR_KODI_IP

Remove "remote=1" on "non fedora 20", this option is needed as far as I know on Fedora 20 (reported bug)

"mpegtsmux" is made to stream mpeg content. This is the better plugin to stream with minimal latency AFAIK.

## Livestreamer - ISS station from space from ustream

Livestreamer is a python tool that is able to fectch streams from some servers and is able to give an url. For some streams, it's impossible to get URL, but we can use "-O" option that dump stream to stdout. So...

	livestreamer http://www.ustream.tv/channel/iss-hdev-payload \
	480p -Q -O | \
	idok -stdin -ssh -target=YOUR_KODI_IP

That will launch the ISS live video from space (sometimes the image is black because ISS station is on the night side. Wait 5 minutes and you will see...)


Develop with me ?
=================

**WARNING - Because there is a problem with dropbear ssh server on raspbmc, we are using go.crypto/ssh package with the patched I made. The package is, at this time, located in ./tunnel package. Soon, if bug is fixed, we will go back to the standard ssh package. See:
https://code.google.com/p/go/issues/detail?id=8657**

To help me to improve idok, bug fixes or optimisations, please fork the repository then make pull-requests. 

I'm not developping inside "master" branch. I'm using "devel" branch and, sometimes other named branch for certain tasks. For example, I was using "refacto" branch to rewrite the code in packages. 

If you want to help, please send me issues for bugs, fix translation, fix README file, I will help you to choose the right branch if needed and you name will be inserted in AUTHORS files as "contributors" or "developpers".


Options
=======

There are other options that may be usefull:

* -check-release=false: check for new release
* -conf-example=false: print a configuration file example to STDOUT
* -disable-check-release=false: disable release check
* -login="": jsonrpc login (configured in xbmc settings)
* -nossh=false: force to not use SSH tunnel - usefull to override configuration file
* -password="": jsonrpc password (configured in xbmc settings)
* -port=8080: local port (ignored if you use ssh option)
* -ssh=false: use SSH Tunnelling (need ssh user and password)
* -sshpass="": ssh password
* -sshport=22: target ssh port
* -sshuser="pi": ssh login
* -stdin=false: read file from stdin to stream
* -target="": xbmc/kodi ip (raspbmc address, ip or hostname)
* -targetport=80: XBMC/Kodi jsonrpc port
* -version=false: Print the current version



TODO
====

- GUI (or not...)
- better lookup adresse for -target option
- Launch a list of streams, playlist...


ChangeLog
=========

* v1-alpha2
  - Fix some function names
  - Cleanup code and documentation

* v1-alpha1
  - Refactorisation to split in packages
  - Move go.crypto package in local sources to not interfer with developper go standard packages
  - Sorry but I change version system one more time... I will now use v1, v2, and so on. Alpha, Beta and RC will be named v2-alpha1, v2-alpha2, and so on.
  - Now support configuration file, -conf-example print a configuration sample

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
