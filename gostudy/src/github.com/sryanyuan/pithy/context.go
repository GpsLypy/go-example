package pithy

import (
	"encoding/json"
	"net/http"

	"strings"

	"strconv"

	"io/ioutil"

	"github.com/gorilla/mux"
)

type HTTPContext struct {
	r    *http.Request
	w    http.ResponseWriter
	body []byte
}

type APIFunc func(*HTTPContext) *APIResult

func wrapHTTPHandler(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &HTTPContext{
			w: w,
			r: r,
		}
		result := fn(ctx)
		if nil == result {
			result = NewAPIResult(0, 0, "")
		}
		ctx.writeAPIResult(result)
	}
}

func (c *HTTPContext) GetRequest() *http.Request {
	return c.r
}

func (c *HTTPContext) SetRequest(r *http.Request) {
	c.r = r
	return
}

func (c *HTTPContext) GetResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *HTTPContext) writeAPIResult(result *APIResult) {
	if 0 != result.StatusCode {
		c.w.WriteHeader(result.StatusCode)
	}
	// If raw is not nil, directly send raw bytes rather than sending json bytes of the result
	if nil != result.RawBytes {
		c.w.Write(result.RawBytes)
		return
	}
	// Send json bytes of the result
	if jsonBytes, err := json.Marshal(result); nil != err {
		panic(err)
	} else {
		c.w.Write(jsonBytes)
	}
}

func (c *HTTPContext) GetPathVarString(key string) string {
	vars := mux.Vars(c.r)
	v, ok := vars[key]
	if !ok {
		return ""
	}
	return v
}

func (c *HTTPContext) GetPathVarInt64(key string, def int64) int64 {
	sval := c.GetPathVarString(key)
	if "" == sval {
		return def
	}
	ival, err := strconv.ParseInt(sval, 10, 64)
	if nil != err {
		return def
	}
	return ival
}

func (c *HTTPContext) GetPathVersion() string {
	return c.GetPathVarString("version")
}

func (c *HTTPContext) GetFormValueInt64(key string, def int64) int64 {
	sval := c.GetFormValueStringTrimBlank(key)
	if "" == sval {
		return def
	}
	ival, err := strconv.ParseInt(sval, 10, 64)
	if nil != err {
		return def
	}
	return ival
}

func (c *HTTPContext) GetFormValueString(key string) string {
	c.r.ParseForm()
	return c.r.FormValue(key)
}

func (c *HTTPContext) GetFormValueStringTrimBlank(key string) string {
	v := c.GetFormValueString(key)
	return strings.Trim(v, " ")
}

func (c *HTTPContext) ReadBody() ([]byte, error) {
	if nil != c.body {
		return c.body, nil
	}
	c.r.ParseForm()

	data, err := ioutil.ReadAll(c.r.Body)
	if nil != err {
		return nil, err
	}
	c.body = data
	return c.body, nil
}

func (c *HTTPContext) UnmarshalJSONFromBody(i interface{}) error {
	body, err := c.ReadBody()
	if nil != err {
		return err
	}
	if err = json.Unmarshal(body, i); nil != err {
		return err
	}
	return nil
}
