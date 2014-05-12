#!/bin/bash

# if you don't know your own md5_pass please visit https://www.duiadns.net/account/account_info.html

##################### USE YOUR OWN CREDENTIALS HERE!!! #####################
host="example.duia.us" # replace with your own hostname                    #
md5_pass="000000000000000000000000000000" # replace with your own md5_pass # 
############################################################################

if [ $host == "example.duia.us" ] || [ $md5_pass == "000000000000000000000000000000" ] ; then
 echo "Edit this script and add your own hostname and md5 password!"
 exit
fi

### This script is using "wget". If you prefer "curl" instead, uncomment "curl" lines and comment (#) "wget" ones.
### IPv4 address update ###

ip4="`wget -qO- http://ipv4.duia.ro`"
#ip4="`curl -sG http://ipv4.duia.ro`"

echo $ip4 > duia4.last
if [ -f "duia4.cache" ] ; then
echo "The file duia4.cache exists"
n=`grep $ip4 duia4.cache | wc -l`
if [ $n -eq 0 ] ; then
echo "Your IPv4 address is not in duia4.cache file; initiate DNS update & duia4.cache file update!"
`wget -qO- --user-agent="duia-cunix-1.0.0.2" "http://ipv4.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip4=$ip4" > /dev/null`
#`curl -sG -A "duia-cunix-1.0.0.2" "http://ipv4.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip4=$ip4" > /dev/null`
cp duia4.last duia4.cache
rm -f duia4.last
else
echo "Your IPv4 address is in duia4.cache file; do nothing!"
rm -f duia4.last
fi
else
echo "The duia4.cache file does not exists; initiate DNS update & create duia4.cache file"
`wget -qO- --user-agent="duia-cunix-1.0.0.2" "http://ipv4.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip4=$ip4" > /dev/null`
#`curl -sG -A "duia-cunix-1.0.0.2" "http://ipv4.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip4=$ip4" > /dev/null`
cp duia4.last duia4.cache
rm -f duia4.last
fi

### IPv6 address update - delete this part of the script if you don't have an IPv6 address on your router ###

ip6="`wget -qO- http://ipv6.duia.ro`"
#ip6="`curl -sG http://ipv6.duia.ro`"

echo $ip6 > duia6.last
if [ -f "duia6.cache" ] ; then
echo "The file duia6.cache exists"
m=`grep $ip6 duia6.cache | wc -l`
if [ $m -eq 0 ] ; then
echo "Your IPv6 address is not in duia6.cache file; initiate DNS update & duia6.cache file update!"
`wget -qO- --user-agent="duia-cunix-1.0.0.2" "http://ipv6.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip6=$ip6" > /dev/null`
#`curl -sG -A "duia-cunix-1.0.0.2" "http://ipv6.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip6=$ip6" > /dev/null`
cp duia6.last duia6.cache
rm -f duia6.last
else
echo "Your IPv6 address is in duia6.cache file; do nothing!"
rm -f duia6.last
fi
else
echo "The duia6.cache file does not exists; initiate DNS update & create duia6.cache file"
`wget -qO- --user-agent="duia-cunix-1.0.0.2" "http://ipv6.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip6=$ip6" > /dev/null`
#`curl -sG -A "duia-cunix-1.0.0.2" "http://ipv6.duia.ro/dynamic.duia?host=$host&password=$md5_pass&ip6=$ip6" > /dev/null`
cp duia6.last duia6.cache
rm -f duia6.last
fi
