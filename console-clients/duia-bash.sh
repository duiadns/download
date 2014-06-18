#!/bin/bash

# if you don't know your own md5_pass please visit https://www.duiadns.net/account/account_info.html

##################### USE YOUR OWN CREDENTIALS HERE!!! #####################
host="example.duia.us" # replace with your own hostname                    #
md5_pass="000000000000000000000000000000" # replace with your own md5_pass # 
############################################################################

use_ip_version=4
ip_cache_file="duia${use_ip_version}.cache"
tmp_ip_file="$( mktemp )"

if [ $host == "example.duia.us" ] || [ $md5_pass == "000000000000000000000000000000" ] ; then
 echo "Edit this script and add your own hostname and md5 password!" >&2
 exit
fi

### This script is using "wget". If you prefer "curl" instead, uncomment "curl" lines and comment (#) "wget" ones.

ip="`wget -qO- http://ipv${use_ip_version}.duia.ro`"
#ip="`curl -sG http://ipv${use_ip_version}.duia.ro`"

function set_ip_for_host() {
	`wget -qO- --user-agent="duia-cunix-1.0.0.2" "http://ipv${use_ip_version}.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip${use_ip_version}=$ip" > /dev/null`
	#`curl -sG -A "duia-cunix-1.0.0.2" "http://ipv${use_ip_version}.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip${use_ip_version}=$ip" > /dev/null`
	cp "${tmp_ip_file}" "${ip_cache_file}"
	rm -f "${tmp_ip_file}"
	
	ip_check=$( cat "${ip_cache_file}" )
	if [ "$ip_check" != "$ip" ] ; then
		echo "The cache file ${ip_cache_file} does not contain what expected." >&2
		exit 1
	fi
}

echo $ip > $tmp_ip_file
if [ -f "${ip_cache_file}" ] ; then
	echo "The file ${ip_cache_file} exists"
	n=`grep $ip "$ip_cache_file" | wc -l`
	if [ $n -eq 0 ] ; then
		echo "Your IPv${use_ip_version} address is not in ${ip_cache_file} file; initiate DNS update & ${ip_cache_file} file update!"
		set_ip_for_host
	else
		echo "Your IPv${use_ip_version} address is in ${ip_cache_file} file; do nothing!"
		rm -f "${tmp_ip_file}"
	fi
else
	echo "The ${ip_cache_file} file does not exists; initiate DNS update & create ${ip_cache_file} file"
	set_ip_for_host
fi
