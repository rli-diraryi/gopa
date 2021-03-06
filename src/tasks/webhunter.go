/**
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午5:42 
 */
package tasks

import (
	log "logging"
	. "net/url"
//	"regexp"
	"strings"
//	. "github.com/PuerkitoBio/purell"
//	. "github.com/zeebo/sbloom"
//	util "util"
//	. "types"
	net "net"
	"io/ioutil"
	"net/http"
	"time"
	"io"
	"compress/gzip"
	"bytes"
)

//parse to get url root
func getRootUrl(source *URL) string {
	if strings.HasSuffix(source.Path, "/") {
		return source.Host + source.Path
	} else {
		index := strings.LastIndex(source.Path, "/")
		if index > 0 {
			path := source.Path[0:index]
			return source.Host + path
		} else {
			return source.Host + "/"
		}
	}
	return ""
}

//format url,prepare for bloom filter
func formatUrlForFilter(url []byte) []byte {
	src := string(url)
	log.Debug("start to normalize url:", src)
	if strings.HasSuffix(src, "/") {
		src = strings.TrimRight(src, "/")
	}
	src = strings.TrimSpace(src)
	src = strings.ToLower(src)
	return []byte(src)
}

////func ExtractLinksFromTaskResponse(bloomFilter *Filter,curl chan []byte, task Task, siteConfig *SiteConfig) {
//func ExtractLinksFromTaskResponse(bloomFilter *Filter, broker *kafka.BrokerPublisher, task Task, siteConfig *TaskConfig) {
//	siteUrlStr := string(task.Url)
//	log.Debug("enter links extract,", siteUrlStr)
//	if siteConfig.SkipPageParsePattern.Match(task.Url) {
//		log.Debug("hit SkipPageParsePattern pattern,", siteUrlStr)
//		return
//	}
//
//	log.Debug("parsing external links:", siteUrlStr, ",using:", siteConfig.LinkUrlExtractRegex)
//	if siteConfig.LinkUrlExtractRegex == nil {
//		siteConfig.LinkUrlExtractRegex = regexp.MustCompile("src=\"(?<url1>.*?)\"|href=\"(?<url2>.*?)\"")
//		log.Debug("use default linkUrlExtractRegex,", siteConfig.LinkUrlExtractRegex)
//	}
//
//	matches := siteConfig.LinkUrlExtractRegex.FindAllSubmatch(task.Response, -1)
//	log.Debug("extract links with pattern:", len(matches), " match result")
//	xIndex := 0
//	for _, match := range matches {
//		log.Debug("dealing with match result,", xIndex)
//		xIndex = xIndex + 1
//		url := match[siteConfig.LinkUrlExtractRegexGroupIndex]
//		filterUrl := formatUrlForFilter(url)
//		log.Debug("url clean result:", string(filterUrl), ",original url:", string(url))
//		filteredUrl := string(filterUrl)
//
//		//filter error link
//		if filteredUrl == "" {
//			log.Debug("filteredUrl is empty,continue")
//			continue
//		}
//
//		result1 := strings.HasPrefix(filteredUrl, "#")
//		if result1 {
//			log.Debug("filteredUrl started with: # ,continue")
//			continue
//		}
//
//		result2 := strings.HasPrefix(filteredUrl, "javascript:")
//		if result2 {
//			log.Debug("filteredUrl started with: javascript: ,continue")
//			continue
//		}
//
//		//判断是否满足FetchRule
//		if(siteConfig.FetchUrlPattern!=nil){
//			 if(!siteConfig.FetchUrlPattern.Match(filterUrl)){
//				 log.Debug("filteredUrl does not match fetchUrlPattern,continue")
//				 continue
//			 }
//		}
//
//		hit := false
//
//		//		l.Lock();
//		//		defer l.Unlock();
//
//		if bloomFilter.Lookup(filterUrl) {
//			log.Debug("hit bloomFilter,continue")
//			hit = true
//			continue
//		}
//
//		if !hit {
//			currentUrlStr := string(url)
//			currentUrlStr = strings.Trim(currentUrlStr, " ")
//
//			seedUrlStr := siteUrlStr
//			seedURI, err := ParseRequestURI(seedUrlStr)
//
//			if err != nil {
//				log.Error("ParseSeedURI failed!: ", seedUrlStr, " , ", err)
//				continue
//			}
//
//			currentURI1, err := ParseRequestURI(currentUrlStr)
//			currentURI := currentURI1
//			if err != nil {
//				if strings.Contains(err.Error(), "invalid URI for request") {
//					log.Warn("invalid URI for request,fix relative url,original:", currentUrlStr)
//					log.Debug("old relatived url,", currentUrlStr)
//					//page based relative urls
//
//					currentUrlStr = "http://" + seedURI.Host + "/" + currentUrlStr
//					currentURI1, err = ParseRequestURI(currentUrlStr)
//					currentURI = currentURI1
//					if err != nil {
//						log.Error("ParseCurrentURI internal failed!: ", currentUrlStr, " , ", err)
//						continue
//					}
//
//					log.Debug("new relatived url,", currentUrlStr)
//
//				} else {
//					log.Error("ParseCurrentURI failed!: ", currentUrlStr, " , ", err)
//					continue
//				}
//			}
//
//			//relative links
//			if currentURI == nil || currentURI.Host == "" {
//				if strings.HasPrefix(currentURI.Path, "/") {
//					//root based relative urls
//					log.Debug("old relatived url,", currentUrlStr)
//					currentUrlStr = "http://" + seedURI.Host + currentUrlStr
//					log.Debug("new relatived url,", currentUrlStr)
//				} else {
//					log.Debug("old relatived url,", currentUrlStr)
//					//page based relative urls
//					urlPath := getRootUrl(currentURI)
//					currentUrlStr = "http://" + urlPath + currentUrlStr
//					log.Debug("new relatived url,", currentUrlStr)
//				}
//			} else {
//				log.Debug("host:", currentURI.Host, " ", currentURI.Host == "")
//
//				//resolve domain specific filter
//				if siteConfig.FollowSameDomain {
//
//					if siteConfig.FollowSubDomain {
//
//						//TODO handler com.cn and .com,using a TLC-domain list
//
//					} else if seedURI.Host != currentURI.Host {
//						log.Debug("domain mismatch,", seedURI.Host, " vs ", currentURI.Host)
//						continue
//					}
//				}
//			}
//
//			if len(siteConfig.LinkUrlMustContain) > 0 {
//				if !util.ContainStr(currentUrlStr, siteConfig.LinkUrlMustContain) {
//					log.Debug("link does not hit must-contain,ignore,", currentUrlStr, " , ", siteConfig.LinkUrlMustNotContain)
//					continue
//				}
//			}
//
//			if len(siteConfig.LinkUrlMustNotContain) > 0 {
//				if util.ContainStr(currentUrlStr, siteConfig.LinkUrlMustNotContain) {
//					log.Debug("link hit must-not-contain,ignore,", currentUrlStr, " , ", siteConfig.LinkUrlMustNotContain)
//					continue
//				}
//			}
//
//			//normalize url
//			currentUrlStr = MustNormalizeURLString(currentUrlStr, FlagLowercaseScheme|FlagLowercaseHost|FlagUppercaseEscapes|
//				FlagRemoveUnnecessaryHostDots|FlagRemoveDuplicateSlashes|FlagRemoveFragment)
//			log.Debug("normalized url:", currentUrlStr)
//			currentUrlByte := []byte(currentUrlStr)
//			if !bloomFilter.Lookup(currentUrlByte) {
//
//				//				if(CheckIgnore(currentUrlStr)){}
//
//				log.Debug("enqueue:", currentUrlStr)
//
//				//TODO 如果使用分布式队列，则不使用go的channel，抽象出接口
//				//				curl <- currentUrlByte
//
//				broker.Publish(kafka.NewMessage(currentUrlByte))
//
//				//				bloomFilter.Add(currentUrlByte)
//			}
//			//			bloomFilter.Add([]byte(filterUrl))
//		} else {
//			log.Debug("hit bloom filter,ignore,", string(url))
//		}
//		log.Debug("exit links extract,", siteUrlStr)
//
//	}
//
//	log.Info("all links within ", siteUrlStr, " is done")
//}
//

func get(url string,cookie string) []byte{

	log.Debug("let's get :"+url)

	client := &http.Client{
		CheckRedirect: nil,
	}
	reqest, _ := http.NewRequest("GET", url, nil)

	reqest.Header.Set("User-Agent"," Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	reqest.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Set("Accept-Charset","GBK,utf-8;q=0.7,*;q=0.3")
	reqest.Header.Set("Accept-Encoding","gzip,deflate,sdch")
	reqest.Header.Set("Accept-Language","zh-CN,zh;q=0.8")
	reqest.Header.Set("Cache-Control","max-age=0")
	reqest.Header.Set("Connection","keep-alive")
	reqest.Header.Set("Referer", url)



	if(len(cookie) > 0){
		log.Debug("dealing with cookie:"+cookie)
		array:=strings.Split(cookie,";")
		for item:= range array{
			array2:=strings.Split(array[item],"=")
			if(len(array2)==2){
				cookieObj:= http.Cookie{}
				cookieObj.Name=array2[0]
				cookieObj.Value=array2[1]
				reqest.AddCookie(&cookieObj)
			}else{
				log.Info("error,index out of range:"+array[item])
			}
		}
	}

	resp, err := client.Do(reqest)

	if err != nil {
		log.Error(url,err)
		return nil
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Error(url,err)
			return nil
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}


	if(reader!=nil){
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Error(url,err)
			return nil
		}
		return body

	}
	return nil
}

func post(url string,cookie string,postStr string) []byte{

	log.Debug("let's post :"+url)

	client := &http.Client{
		CheckRedirect: nil,
	}

	postBytesReader := bytes.NewReader([]byte(postStr))
	reqest, _ := http.NewRequest("POST", url, postBytesReader)

	reqest.Header.Set("User-Agent"," Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	reqest.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Set("Accept-Charset","GBK,utf-8;q=0.7,*;q=0.3")
	reqest.Header.Set("Accept-Encoding","gzip,deflate,sdch")
//	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	reqest.Header.Set("Accept-Language","zh-CN,zh;q=0.8")
	reqest.Header.Set("Cache-Control","max-age=0")
	reqest.Header.Set("Connection","keep-alive")
	reqest.Header.Set("Referer", url)



	if(len(cookie) > 0){
		log.Debug("dealing with cookie:"+cookie)
		array:=strings.Split(cookie,";")
		for item:= range array{
			array2:=strings.Split(array[item],"=")
			if(len(array2)==2){
				cookieObj:= http.Cookie{}
				cookieObj.Name=array2[0]
				cookieObj.Value=array2[1]
				reqest.AddCookie(&cookieObj)
			}else{
				log.Info("error,index out of range:"+array[item])
			}
		}
	}

	resp, err := client.Do(reqest)

	if err != nil {
		log.Error(url,err)
		return nil
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Error(url,err)
			return nil
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}


	if(reader!=nil){
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Error(url,err)
			return nil
		}
		return body
	}
	return nil
}


func HttpGetWithCookie(resource string,cookie string)(msg []byte,err error){

   out := get(resource,cookie)
   return out,nil
}

func HttpGet(resource string)(msg []byte,err error){

//	//validate url
//	host, err := ParseRequestURI(resource)
//	if err != nil {
//		log.Error(resource,err)
//		return nil,err
//	}

//	//check domain
//	_, err =net.LookupIP(host.Host)
//	if err != nil {
//		log.Error(resource,err)
//		return nil,err
//	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(10 * time.Second)
				c, err := net.DialTimeout(netw, addr, 5*time.Second) //连接超时时间
				if err != nil {
					log.Error(resource,err)
					return nil, err
				}

				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	req, err := http.NewRequest("GET", resource, nil)

	if err != nil {
		log.Error(resource,err)
		return nil,err
	}

	//support gzip
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; gopa/0.1; +http://infinitbyte.com/gopa)")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(req)
	if err != nil {
		log.Error(resource,err)
		return nil,err
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Error(resource,err)
			return nil,err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}
	if(reader!=nil){
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Error(resource,err)
			return nil,err
		}
		return body,nil
	}
	return nil,http.ErrNotSupported
}
