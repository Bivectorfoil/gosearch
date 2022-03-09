#!/bin/bash

# for local dev to set proxy
set -e
# bash exec command from working directory (pwd)

# change work directory to script directory
cd ./script

host_ip=$(cat /etc/resolv.conf |grep "nameserver" |cut -f 2 -d " ")
curProxyConf='PROXY_HOST='"http://$host_ip:10809"
oldProxyConf=`sed '3q;d' '../.env'`
# oldProxy=`awk -F "=" '{print $2}' <<< $oldProxyConfig`

if [ "$oldProxyConf" != "$curProxyConf" ]
then
    echo "old: $oldProxyConf, cur: $curProxyConf"
    sed -i "3s|$oldProxyConf|$curProxyConf|g" '../.env'
fi


