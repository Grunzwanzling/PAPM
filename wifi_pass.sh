#!/bin/sh
wpa_cli -i wlo1 set_network '3' 'password' \"$(/home/max/cli -command wifi/KIT)\"
