#!/bin/sh
wpa_cli -i wlo1 set_network '3' 'password' \"$(/home/max/git/YouShallNotPassword/cli -command 'get;wifi/KIT' -socket /home/max/git/YouShallNotPassword/socket)\"
