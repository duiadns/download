#!/bin/sh

# if you don't know your own md5_pass please visit https://www.duiadns.net/account/account_info.html

##################### USE YOUR OWN CREDENTIALS HERE!!! #####################
host="example.duia.us" # replace with your own hostname                    #
md5_pass="000000000000000000000000000000" # replace with your own md5_pass #
############################################################################

has() {
	type $1 > /dev/null 2>&1
}

use_ip_version=4
ip_cache_file="duia${use_ip_version}.cache"
tmp_ip_file="$( mktemp )"
user_agent="github.bash-1.0.0.4"


has curl || has wget || { echo "Either curl or wget is required, but none were found" 1>&2 && exit 1; }


if [ $host = "example.duia.us" ] || [ $md5_pass = "000000000000000000000000000000" ] ; then
	echo "Edit this script and add your own hostname and md5 password!" >&2
	exit
fi

ip_url="http://ipv${use_ip_version}.duia.ro"
if has wget; then
	ip=`wget -qO- ${ip_url}`
else
	ip=`curl -sG ${ip_url}`
fi


set_ip_for_host () {
	set_ip_url="http://ipv${use_ip_version}.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip${use_ip_version}=$ip"
	server_response=0
	if has wget; then
		server_response=$(wget -S -qO- --user-agent="$user_agent" $set_ip_url 2>&1 | egrep "HTTP/[0-9\.[0-9]" | awk '{ print $2}')
	else
		server_response=$(curl -sG -L -w "%{http_code}" -A "${user_agent}" $set_ip_url -o /dev/null )
	fi
	if [ $server_response -eq 200 ] ; then
		cp "${tmp_ip_file}" "${ip_cache_file}"
	else
		echo "Update of your IPv${use_ip_version} address failed! Server response was $server_response" >&2
		exit 1
	fi
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
