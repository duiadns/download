#!/bin/sh

# if you don't know your own md5_pass please visit https://www.duiadns.net/account/account_info.html

##################### USE YOUR OWN CREDENTIALS HERE!!! #####################
host="example.duia.us" # replace with your own hostname			   #
md5_pass="000000000000000000000000000000" # replace with your own md5_pass #
############################################################################

has() {
	type $1 > /dev/null 2>&1
}

cleanup_and_exit() {
	if [ -n "$tmp_ip_file" ] && [ -f "$tmp_ip_file" ]; then
		rm -f $tmp_ip_file
	fi
	exit "$1"
}

die() {
	echo "$@" >&2
	cleanup_and_exit 1
}

use_ip_version=4
ip_cache_file="duia${use_ip_version}.cache"
tmp_ip_file=$(mktemp -t duia.XXXXXXXXXX)
user_agent="github.bash-1.0.0.4"

has curl || has wget || die "Either curl or wget is required, but none were found"

if [ $host = "example.duia.us" ] || [ $md5_pass = "000000000000000000000000000000" ] ; then
	die "Edit this script and add your own hostname and md5 password!"
fi

# make sure the ip_cache_file exists and is readable and writeable
[ -f "${ip_cache_file}" ] || touch "${ip_cache_file}" ||
	die "Cannot create cache file ${ip_cache_file}"
[ -r "${ip_cache_file}" ] && [ -w "${ip_cache_file}" ] ||
	die "Cannot read and write cache file ${ip_cache_file}"

old_ip=$(cat ${ip_cache_file})

ip_url="http://ipv${use_ip_version}.duia.ro"
if has wget; then
	ip=`wget -qO- ${ip_url}`
else
	ip=`curl -sG ${ip_url}`
fi
echo $ip > $tmp_ip_file

size=${#ip}
if [ "$size" -gt "5" ]; then
	set_ip_for_host () {
		set_ip_url="http://ipv${use_ip_version}.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip${use_ip_version}=$ip"
		server_response=0
		if has wget; then
			server_response=$(wget -S -qO- -U "$user_agent" $set_ip_url 2>&1 | egrep "HTTP/[0-9\.[0-9]" | awk '{ print $2}')
		else
			server_response=$(curl -sG -L -w "%{http_code}" -A "${user_agent}" $set_ip_url -o /dev/null )
		fi
		if [ $server_response -eq 200 ] ; then
			cp "${tmp_ip_file}" "${ip_cache_file}"
		else
			die "Update of your IPv${use_ip_version} address failed! Server response was $server_response"
		fi
		ip_check=$( cat "${ip_cache_file}" )
		if [ "$ip_check" != "$ip" ] ; then
			die "The cache file ${ip_cache_file} does not contain what expected."
		fi
	}
	if [ "${old_ip}" != "${ip}" ] ; then
		echo "Your IPv${use_ip_version} address is not in ${ip_cache_file} file; initiate DNS update & ${ip_cache_file} file update!"
		set_ip_for_host
	else
		echo "Your IPv${use_ip_version} address is in ${ip_cache_file} file; do nothing!"
	fi
else
	die "The IP address returned by ipv${use_ip_version}.duia.ro is NULL. This is really odd, do nothing this time!"
fi

cleanup_and_exit 0
