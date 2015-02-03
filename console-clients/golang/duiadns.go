package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

const UserAgent = "duia-go-1.0.0.2"

func getIpFromSite(version int) (s string, err error) {
	//get my ip from server (ipv4/ipv6 compatible)
	req, err := http.Get("http://" + "ipv" + strconv.Itoa(version) + ".duia.ro")
	if err != nil {
		//println("no ipv" + strconv.Itoa(version) + " connection")
		return "none", err
	}
	ip, err := ioutil.ReadAll(req.Body)
	//println("get ip" + strconv.Itoa(version) + " " + string(ip))
	return string(ip), nil
}

func updateDNS(host, password, ip4, ip6 string) (err error) {
	//connect to duia server and update both ipv4/ipv6 on the same request
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"http://ip.duia.ro/dynamic.duia?host="+host+"&password="+password+"&ip4="+ip4+"&ip6="+ip6, nil)
	if err != nil {
		return err
	}
	//tested with http://httpbin.org
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()
	return nil
}

func readCache() (ip4, ip6 string, err error) {
	path, _ := os.Getwd()
	file, err := os.Open(path + string(filepath.Separator) + "duia.cache")
	if err != nil {
		return "", "", err
	}
	fmt.Fscanf(file, "%s %s", &ip4, &ip6)
	return ip4, ip6, nil
}

func updateCache(ip4, ip6 string) (err error) {
	path, _ := os.Getwd()
	file, err := os.Create(path + string(filepath.Separator) + "duia.cache")
	if err != nil {
		return err
	}
	// write creditentials in unix clasic style format
	fmt.Fprintf(file, "%s %s", ip4, ip6)
	return nil
}

func readCfg() (host, password string, err error) {
	path, _ := os.Getwd()
	// get creditentials from file
	file, err := os.Open(path + string(filepath.Separator) + "duia.cfg")
	if err != nil {
		return "", "", err
	}
	fmt.Fscanf(file, "%s %s", &host, &password)
	return host, password, nil
}

func updateCfg(host, password string) {
	// md5 password encoding
	h := md5.New()
	io.WriteString(h, password)
	md5 := fmt.Sprintf("%x", h.Sum(nil))
	// create creditentials file
	path, _ := os.Getwd()
	file, _ := os.Create(path + string(filepath.Separator) + "duia.cfg")
	// write creditentials in unix clasic style format
	fmt.Fprintf(file, "%s %s", host, md5)
}

func updateLog(message string) {
	path, _ := os.Getwd()
	file, _ := os.OpenFile(path+string(filepath.Separator)+"duia.log", os.O_CREATE|os.O_APPEND, 0660)
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(file, "%s", "\n"+t+" - "+message)
}

func main() {
	// get current path
	file, _ := exec.LookPath(os.Args[0])
	dir, _ := path.Split(file)
	// change to curent path to avoid problems when the program
	// is launch from other location  without directory changed
	os.Chdir(dir)
	fmt.Println(dir)

	// if duia.cfg does not exist, get hostname and passowrd from command line
	host, password, err := readCfg()
	if err != nil {
		fmt.Println("Creating duia.cfg file ... please add:\n")
		fmt.Print("Hostname: ")
		fmt.Scan(&host)
		fmt.Print("Password: ")
		fmt.Scan(&password)
		fmt.Println("")
		updateCfg(host, password)
		fmt.Println("duia.cfg file created\n")
	}

	// if duia.cache does not exist, create it with "none" as ip addresses
	ip4, ip6, err := readCache()
	if err != nil {
		ip4 = "none"
		ip6 = "none"
		updateCache(ip4, ip6)
		fmt.Println("duia.cache file created\n")
	}
	needUpdate := false
	//ipv4 support
	ip4FromSite, err := getIpFromSite(4)
	if err == nil {
		// check if ip4 has changed to send DNS update
		if ip4FromSite != ip4 {
			ip4 = ip4FromSite
			needUpdate = true
		}
	}
	//ipv6 support
	ip6FromSite, err := getIpFromSite(6)
	if err == nil {
		// check if ip6 has changed to send DNS update
		if ip6FromSite != ip6 {
			ip6 = ip6FromSite
			needUpdate = true
		}
	}
	if needUpdate {
		updateDNS(host, password, ip4, ip6)
		if ip4 != "none" {
			fmt.Println("IPv4 DNS entry (" + host + ", " + ip4 + ") updated.")
			updateLog("IPv4 DNS entry (" + host + ", " + ip4 + ") updated.")
		}
		if ip4 == "none" {
			fmt.Println("IPv4 DNS entry for " + host + " deleted.")
			updateLog("IPv4 DNS entry for " + host + " deleted.")
		}
		if ip6 != "none" {
			fmt.Println("IPv6 DNS entry (" + host + ", " + ip6 + ") updated.\n")
			updateLog("IPv6 DNS entry (" + host + ", " + ip6 + ") updated.\n")
		}
		if ip6 == "none" {
			fmt.Println("IPv6 DNS entry for " + host + " deleted.\n")
			updateLog("IPv6 DNS entry for " + host + " deleted.\n")
		}
		updateCache(ip4, ip6)
		fmt.Println("duia.cache file updated.\n")
		fmt.Println("duia.log file updated.\n")
		//	fmt.Printf("cache updated: %02d:%02d.%02d\n", t.Hour(), t.Minute(), t.Second())
	} else {
		fmt.Println("No IPv4/IPv6 DNS updates.")
	}
}
