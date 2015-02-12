package main




import (
 "crypto/md5"
 "crypto/tls"
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

const DuiaVersion = "2" + "." + "0" + "." + "0" + "." + "4"
const UserAgent = "github.go" + "-" + DuiaVersion

func getIpFromSite(version int) (s string, err error) {

 resp, err := doRequest("https://" + "ipv" + strconv.Itoa(version) + ".duia.ro")
 if err != nil {

  return "none", err
 }
 defer resp.Body.Close()
 ip, err := ioutil.ReadAll(resp.Body)

 return string(ip), nil
}

func doRequest(urlStr string) (r *http.Response, err error) {
 tr := &http.Transport{
  TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
 }
 client := &http.Client{Transport: tr}
 req, err := http.NewRequest("GET", urlStr, nil)
 if err != nil {
  return nil, err
 }

 req.Header.Set("User-Agent", UserAgent)
 resp, err := client.Do(req)

 return resp, err
}

func updateDNS4(host, password, ip4 string) (d int, err error) {

 resp, err := doRequest("https://ipv4.duia.ro/dynamic.duia?host="+host+"&password="+password+"&ip4="+ip4)
 if err != nil {
  return -1, nil
 }
 defer resp.Body.Close()
 return resp.StatusCode, nil
}

func updateDNS6(host, password, ip6 string) (d int, err error) {

 resp, err := doRequest("https://ipv6.duia.ro/dynamic.duia?host="+host+"&password="+password+"&ip6="+ip6)
 if err != nil {
  return -1, nil
 }
 defer resp.Body.Close()
 return resp.StatusCode, nil
}

func updateDNS(host, passwordAsMd5, ip4, ip6 string) (d int, err error) {

 resp, err := doRequest("https://ipv4.duia.ro/dynamic.duia?host="+host+"&password="+passwordAsMd5+"&ip4="+ip4+"&ip6="+ip6)
 if err != nil {
  return -1, nil
 }
 defer resp.Body.Close()
 return resp.StatusCode, nil
}

func readCache() (ip4, ip6 string, err error) {
 path, _ := os.Getwd()
 file, err := os.Open(path + string(filepath.Separator) + "duia.cache")
 if err != nil {
  return "", "", err
 }
 defer file.Close()

 fmt.Fscanf(file, "%s %s", &ip4, &ip6)
 return ip4, ip6, nil
}

func updateCache(ip4, ip6 string) (err error) {
 path, _ := os.Getwd()
 file, err := os.Create(path + string(filepath.Separator) + "duia.cache")
 if err != nil {
  return err
 }
 defer file.Close()


 fmt.Fprintf(file, "%s %s", ip4, ip6)
 return nil
}

func readCfg() (host, passwordAsMd5, update string, err error) {
 path, _ := os.Getwd()

 file, err := os.Open(path + string(filepath.Separator) + "duia.cfg")
 if err != nil {
  return "", "", "", err
 }
 defer file.Close()

 fmt.Fscanf(file, "%s %s %s", &host, &passwordAsMd5, &update)
 return host, passwordAsMd5, update, nil
}

func computeMd5(str string) (strAsMd5 string){

 h := md5.New()
 io.WriteString(h, str)
 strAsMd5 = fmt.Sprintf("%x", h.Sum(nil))

 return strAsMd5
}

func updateCfg(host, passwordAsMd5, update string) {

 path, _ := os.Getwd()
 file, err := os.Create(path + string(filepath.Separator) + "duia.cfg")
 if err != nil {
  return
 }
 defer file.Close()


 fmt.Fprintf(file, "%s %s %s", host, passwordAsMd5, update)
}

func updateLog(message string) {
 path, _ := os.Getwd()
 file, err := os.OpenFile(path+string(filepath.Separator)+"duia.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
 if err != nil {
  return
 }
 defer file.Close()

 t := time.Now().Format("2006-01-02 15:04:05")
 fmt.Fprintf(file, "%s", "\r\n"+t+" - "+message)
}

func duiaInit() (file string) {

 file, _ = exec.LookPath(os.Args[0])
 dir, _ := path.Split(filepath.ToSlash(file))


 os.Chdir(dir)


 return file
}

func duiaMain(file string) {

 host, passwordAsMd5, update, err := readCfg()
 if err != nil {
  var password string
  fmt.Println("Creating duia.cfg file ... please add:\n")
  fmt.Print("Hostname: ")
  fmt.Scan(&host)
  fmt.Print("Password: ")
  fmt.Scan(&password)
  passwordAsMd5 = computeMd5(password)
  fmt.Print("Update (ipv4/ipv6/both): ")
     fmt.Scan(&update)
     fmt.Println("")
  updateCfg(host, passwordAsMd5, update)
  fmt.Println("duia.cfg file created\n")
 }


 ip4, ip6, err := readCache()
 if err != nil {
  ip4 = "none"
  ip6 = "none"
  updateCache(ip4, ip6)
  fmt.Println("duia.cache file created\n")
 }
 needUpdate := false


 if (update == "ipv4"){
  ip4FromSite, err := getIpFromSite(4)
  if err == nil {

   if ip4FromSite != ip4 {
    ip4 = ip4FromSite
    needUpdate = true
   }
  }
  ip6="none"
  if needUpdate {
   http_code, err := updateDNS4(host, passwordAsMd5, ip4)
   if err == nil{
    if http_code == 200{
     if ip4 != "none" {
      fmt.Println("IPv4 DNS entry (" + host + ", " + ip4 + ") updated.")
      updateLog("IPv4 DNS entry (" + host + ", " + ip4 + ") updated.")
     }
     if ip4 == "none" {
      fmt.Println("IPv4 DNS entry for " + host + " deleted. IPv4 internet connection cannot be established!")
      updateLog("IPv4 DNS entry for " + host + " deleted. IPv4 internet connection cannot be established!")
     }
     updateCache(ip4, ip6)
     fmt.Println("duia.cache file updated.\n")
     fmt.Println("duia.log file updated.\n")
    }else if http_code == 401{
     log_err := "DNS update failed, invalid hostname or password!"
                                        fmt.Println(log_err)
     updateLog(log_err)
           }else{
                                        log_err := "DNS update failed, error code " + strconv.Itoa(http_code) + "!"
                                        fmt.Println(log_err)
                                        updateLog(log_err)
    }
   }
  } else {
   fmt.Println("No IPv4 DNS updates.")
  }
 }


 if (update == "ipv6"){
  ip6FromSite, err := getIpFromSite(6)
  if err == nil {

   if ip6FromSite != ip6 {
    ip6 = ip6FromSite
    needUpdate = true
   }
  }
  ip4="none"
  if needUpdate {
   http_code, err := updateDNS6(host, passwordAsMd5, ip6)
   if err == nil{
    if http_code == 200{
     if ip6 != "none" {
      fmt.Println("IPv6 DNS entry (" + host + ", " + ip6 + ") updated.")
      updateLog("IPv6 DNS entry (" + host + ", " + ip6 + ") updated.")
     }
     if ip6 == "none" {
      fmt.Println("IPv6 DNS entry for " + host + " deleted. IPv6 internet connection cannot be established!")
      updateLog("IPv6 DNS entry for " + host + " deleted. IPv6 internet connection cannot be established!")
     }
     updateCache(ip4, ip6)
     fmt.Println("duia.cache file updated.\n")
     fmt.Println("duia.log file updated.\n")
                                }else if http_code == 401{
                                        log_err := "DNS update failed, invalid hostname or password!"
                                        fmt.Println(log_err)
                                        updateLog(log_err)
    }else{
     log_err := "DNS update failed, error code " + strconv.Itoa(http_code) + "!"
     fmt.Println(log_err)
     updateLog(log_err)
    }
   }
  } else {
   fmt.Println("No IPv6 DNS updates.")
  }
 }


 if (update == "both"){
  ip4FromSite, err := getIpFromSite(4)
  if err == nil {

   if ip4FromSite != ip4 {
    ip4 = ip4FromSite
    needUpdate = true
   }
  }
  ip6FromSite, err := getIpFromSite(6)
  if err == nil {

   if ip6FromSite != ip6 {
    ip6 = ip6FromSite
    needUpdate = true
   }
  }
  if needUpdate {
   http_code, err := updateDNS6(host, passwordAsMd5, ip6)
   if err == nil{
    if http_code == 200{
     updateDNS(host, passwordAsMd5, ip4, ip6)
     if ip4 != "none" {
      fmt.Println("IPv4 DNS entry (" + host + ", " + ip4 + ") updated.")
      updateLog("IPv4 DNS entry (" + host + ", " + ip4 + ") updated.")
     }
     if ip4 == "none" {
      fmt.Println("IPv4 DNS entry for " + host + " deleted. IPv4 internet connection cannot be established!")
      updateLog("IPv4 DNS entry for " + host + " deleted. IPv4 internet connection cannot be established!")
     }
     if ip6 != "none" {
      fmt.Println("IPv6 DNS entry (" + host + ", " + ip6 + ") updated.")
      updateLog("IPv6 DNS entry (" + host + ", " + ip6 + ") updated.")
     }
     if ip6 == "none" {
      fmt.Println("IPv6 DNS entry for " + host + " deleted. IPv6 internet connection cannot be established!")
      updateLog("IPv6 DNS entry for " + host + " deleted. IPv6 internet connection cannot be established!")
     }
     updateCache(ip4, ip6)
     fmt.Println("duia.cache file updated.\n")
     fmt.Println("duia.log file updated.\n")
                                }else if http_code == 401{
                                        log_err := "DNS update failed, invalid hostname or password!"
                                        fmt.Println(log_err)
                                        updateLog(log_err)
    }else{
     log_err := "DNS update failed, error code " + strconv.Itoa(http_code) + "!"
     fmt.Println(log_err)
     updateLog(log_err)
    }
   }
  } else {
   fmt.Println("No IPv4/IPv6 DNS updates.")
  }
 }
}
