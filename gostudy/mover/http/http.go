package http

import (
	"app/tools/mover/info"
	"app/tools/mover/logging"
	"app/tools/mover/moverconfig"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	httpMutex        sync.Mutex
	httpBodyCacheMap map[string]moverconfig.CacheData
)

func HttpGetBody(url string) ([]byte, error) {
	httpMutex.Lock()
	defer httpMutex.Unlock()

	if nil == httpBodyCacheMap {
		httpBodyCacheMap = make(map[string]moverconfig.CacheData)
	}
	if v, ok := httpBodyCacheMap[url]; ok {
		if v.ExpireTime >= time.Now().Unix() {
			return v.Data.([]byte), nil
		}
		delete(httpBodyCacheMap, url)
	}

	// get url
	logging.LogDebug("Get url, url=%s", url)
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Get(url)
	if nil != err {
		return nil, err
	}

	// parse body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	var cacheData moverconfig.CacheData
	cacheData.Data = body
	cacheData.ExpireTime = time.Now().Unix() + int64(info.CacheTime)
	httpBodyCacheMap[url] = cacheData
	return body, nil
}
