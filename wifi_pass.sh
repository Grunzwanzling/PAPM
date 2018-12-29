#!/bin/sh
wpa_cli -i wlo1 set_network '0' 'password' \"$(/home/max/cli -command wifi/eduroam)\"
