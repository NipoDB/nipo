#!/bin/bash

declare -A osInfo;
osInfo[/etc/debian_version]="apt-get install -y"
osInfo[/etc/alpine-release]="apk --update add"
osInfo[/etc/centos-release]="yum install -y"
osInfo[/etc/fedora-release]="dnf install -y"

for f in ${!osInfo[@]}
do
    if [[ -f $f ]];then
        package_manager=${osInfo[$f]}
    fi
done

package="golang"

command=$1

case $command in

  "install")
        ${package_manager} ${package}
        if [ -d "/tmp/NipoDB" ]
        then
            rm -fr /tmp/NipoDB
            echo remove old git repo 
        fi
        echo cloning the repository
        git clone https://github.com/NipoDB/nipo.git /tmp/NipoDB
        echo get nipolib 
        go get github.com/NipoDB/nipolib
        echo get yaml.v2 
        go get gopkg.in/yaml.v2
        cd /tmp/NipoDB/nipo
        echo building the go binary
        go build
        if [ -f "/usr/local/bin/nipo" ]
        then
            rm -f /usr/local/bin/nipo
            echo coping the binary to /usr/local/bin/
            cp /tmp/NipoDB/nipo/nipo /usr/local/bin/
        else
            echo coping the binary to /usr/local/bin/
            cp /tmp/NipoDB/nipo/nipo /usr/local/bin/
        fi
        
        if [ -d "/etc/nipo" ]
        then
            echo create config file
            rm -fr /etc/nipo
            mkdir /etc/nipo
            cp /tmp/NipoDB/files/config/nipo-cfg.yaml /etc/nipo/
        else
            echo create config file
            mkdir /etc/nipo
            cp /tmp/NipoDB/files/config/nipo-cfg.yaml /etc/nipo/
        fi
        if [ -f "/lib/systemd/system/nipo.service" ]
        then
            echo create service file
            rm -f /lib/systemd/system/nipo.service
            cp /tmp/NipoDB/files/config/nipo.service /lib/systemd/system/
        else
            echo create service file
            cp /tmp/NipoDB/files/config/nipo.service /lib/systemd/system/
        fi
        if [ -d "/var/log/nipo" ]
        then
            echo log directory exists
        else
            echo create log directory /var/log/nipo
            mkdir /var/log/nipo
            touch /var/log/nipo/nipo.log
        fi
        if [ -d "/tmp/NipoDB" ]
        then
            rm -fr /tmp/NipoDB
            echo remove old git repo 
        fi
    ;;

  "uninstall")
        if [ -d "~/go/src/nipo" ]
        then
            echo removing ~/go/src/nipo
            rm -fr ~/go/src/nipo
        fi
        if [ -f "/usr/local/bin/nipo" ]
        then
            echo removing /usr/local/bin/nipo
            rm -fr /usr/local/bin/nipo
        fi
        if [ -d "/tmp/NipoDB" ]
        then
            echo removing /tmp/NipoDB
            rm -fr /tmp/NipoDB
        fi
        if [ -d "/var/log/nipo" ]
        then
            echo removing /var/log/nipo
            rm -fr /var/log/nipo
        fi
        if [ -d "/etc/nipo" ]
        then
            echo removing /etc/nipo
            rm -fr /etc/nipo
        fi
        if [ -f "/lib/systemd/system/nipo.service" ]
        then
            echo removing /lib/systemd/system/nipo.service
            rm -fr /lib/systemd/system/nipo.service
        fi
    ;;
esac