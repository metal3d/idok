What is it ?
============

IDOK (kodi reversed name) is a Go program that allows to serve medias to your Kodi plateform (raspbmc, xbmc...)

You may be able to send video, images and music from your computer.

Idok has got two modes:

* your computer serve media from a standard port (default 8080)
* your computer dig a tunnel and serve media

Installation
============

You can clone repository and compile source code yourself:

	git clone http://git.develipsum.com/metal3d/idok.git
	cd idok
	go build -o idok main.go

Then you can put binary in your PATH:

	mkdir -p ~/.local/bin
	cp idok ~/.local/bin

You can now ask some help:

	idok -h
	Usage of idok:
	  -login="": jsonrpc login (configured in xbmc settings)
	  -password="": jsonrpc password (configured in xbmc settings)
	  -port=8080: local port (ignored if you use ssh option)
	  -ssh=false: Use SSH Tunnelling (need ssh user and password)
	  -sshpass="raspberry": ssh password
	  -sshuser="pi": ssh login
	  -target="": xbmc/kodi ip (raspbmc address, ip or hostname)


There is a build for x86_64 cpu (if you're computer is not too old, this is ok) that is available (TODO: give the url)

Stream your first media
=======================


## HTTP (default)

The HTTP way is not secured. While you're streaming to Kodi (or XBMC), the directory where is the media can be accessed by other computer in your network.

This solution need to open port on your firewall. 

By default, idok open 8080 port, but you can specify other port.

To open firewall port on you linux installation:

	firewall-cmd --add-port=8080/tcp

When you will reload firewall, or restart computer, the port will be closed. If you want to keep that port open:


	firewall-cmd --add-port=8080/tcp --permanent

Then, send media:

	idok -target=IP_OF_KODI_OR_XBMC /path/to/media.mp3

If you've opened other port, you can set it. For example for port 1234:

	idok -port=1234 -target=IP_OF_KODI_OR_XBMC /path/to/media.mp3


## SSH

The SSH way is the easier and more secured way. You don't have to open port on your computer and only the Kodi instance will be able to access your content. This is the default mode

Unfortunately, raspbmc has a little problem with ssh server. But it's not very hard to fix.

We will replace dropbear by openssh-server.

	$ ssh pi@YOUR_RASPBMC_IP
	$ sudo nano /etc/default/dropbear

Change DROPBEAR_PORT value to "2222" then save by pressing CTRL+X then Y

Now, install openssh-server:

	$ sudo apt-get update
	$ sudo apt-get -y openssh-server

After installation, you should restart your raspberry
	
	$ sudo reboot

The next time you want to ssh your raspbmc, your computer will complain about fingerprint changes. Just open your ~/.ssh/known_hosts and remove the corresponding line for your raspberry ip.

That's it, you have finished openssh installation. Now try to send data:

	idok -ssh -target=IP_OF_RASPBERRY /path/to/media.mp3

Your kodi should open the file.

Pressing CTRL+C should stop media stream and exit program.

Options
=======

There are other options that may be usefull:

* -target: kodi instance ip or hostname 
* -login : xbmc or kodi login configured on web interface settings
* -password : xbmc or kodi password configured on web interface settings
* -ssh : If set, idok will dig ssh tunnel to stream content
* -sshuser : if you don't user "pi" user
* -sshpass : if you changed standard password of "pi" user
* -port : local port for media stream if you don't use ssh tunneling, default is 8080

TODO
====

I should ask dropbear maintainer why the tunneling won't work. Installing openssh-server on raspbmc can be complicated for some users.

Add ssh port option.

Add kodi/xbmc port option (some users changed the default 80) 


