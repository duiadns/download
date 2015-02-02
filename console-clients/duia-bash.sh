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
duiadns_propagation_timeout=20
ip_cache_file="duia${use_ip_version}.cache"
tmp_ip_file="$( mktemp )"
user_agent="duia-unix-1.0.0.3"


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

dig_query_parm="A"
wget_query_parm="-4"
if [ $use_ip_version -eq 6 ]
then
	dig_query_parm="AAAA"
	wget_query_parm="-6"
fi

if has dig
then 
	duia_ip=`dig @ns1.duiadns.net $dig_query_parm +short $host`
	if [ -z $duia_ip ]
	then
		duia_ip=`dig @ns2.duiadns.net $dig_query_parm +short $host`
	fi
else
	duia_ip=`wget wget_query_parm -t1 -T3 --no-dns-cache --spider $host 2>&1 | grep Resolving | grep -v failed | awk '{ print $4}'`
fi

set_ip_for_host () {
	set_ip_url="http://ipv${use_ip_version}.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip${use_ip_version}=$ip" 
	if has wget; then
		wget -qO- --user-agent="$user_agent" $set_ip_url > /dev/null
	else
		curl -sG -A "${user_agent}" $set_ip_url > /dev/null
	fi
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
	if test `find "$ip_cache_file" -mmin +$duiadns_propagation_timeout` ; then
		if [ "$duia_ip" != "$ip" ] ; then
			echo "Your IPv${use_ip_version} address was updated $duiadns_propagation_timeout minute(s) ago, but DuiaDNS report different address. Update forced!"
			set_ip_for_host
			n=1
		fi
	fi

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
