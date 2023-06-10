package builtin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Syuparn/pangaea/object"
)

func request(env *object.Env, kwargs *object.PanObj, args ...object.PanObject) object.PanObject {
	urlPair, ok := (*kwargs.Pairs)[object.GetSymHash("url")]
	if !ok {
		return object.NewTypeErr("url must be specified")
	}
	urlStr, ok := object.TraceProtoOfStr(urlPair.Value)
	if !ok {
		return object.NewTypeErr(fmt.Sprintf("url `%s` cannot be treated as str", urlPair.Value.Inspect()))
	}
	u, err := url.Parse(urlStr.Value)
	if err != nil {
		return object.NewValueErr(fmt.Sprintf("failed to parse url: `%s`", urlStr.Value))
	}

	methodPair, ok := (*kwargs.Pairs)[object.GetSymHash("method")]
	if !ok {
		return object.NewTypeErr("method must be specified")
	}
	methodStr, ok := object.TraceProtoOfStr(methodPair.Value)
	if !ok {
		return object.NewTypeErr(fmt.Sprintf("method `%s` cannot be treated as str", methodPair.Value.Inspect()))
	}

	if queriesPair, ok := (*kwargs.Pairs)[object.GetSymHash("queries")]; ok {
		queriesObj, ok := object.TraceProtoOfObj(queriesPair.Value)
		if !ok {
			return object.NewTypeErr(fmt.Sprintf("queries `%s` cannot be treated as obj", queriesPair.Value.Inspect()))
		}

		if errObj := addQueries(u, queriesObj); errObj != nil {
			return errObj
		}
	}

	var body io.Reader
	if bodyPair, ok := (*kwargs.Pairs)[object.GetSymHash("body")]; ok {
		bodyStr, ok := object.TraceProtoOfStr(bodyPair.Value)
		if !ok {
			return object.NewTypeErr(fmt.Sprintf("body `%s` cannot be treated as str", bodyPair.Value.Inspect()))
		}
		body = strings.NewReader(bodyStr.Value)
	}

	req, err := http.NewRequest(methodStr.Value, u.String(), body)
	if err != nil {
		// TODO: create specific error type
		return object.NewPanErr(err.Error())
	}

	if headersPair, ok := (*kwargs.Pairs)[object.GetSymHash("headers")]; ok {
		headersObj, ok := object.TraceProtoOfObj(headersPair.Value)
		if !ok {
			return object.NewTypeErr(fmt.Sprintf("headers `%s` cannot be treated as obj", headersPair.Value.Inspect()))
		}

		if errObj := addHeaders(headersObj, req.Header); errObj != nil {
			return errObj
		}
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		// TODO: create specific error type
		return object.NewPanErr(err.Error())
	}

	return responseObj(res)
}

func addHeaders(headersObj *object.PanObj, headers http.Header) *object.PanErr {
	for k, v := range *headersObj.Pairs {
		keyObj, _ := object.SymHash2Str(k)
		key := keyObj.(*object.PanStr).Value

		valueStr, ok := object.TraceProtoOfStr(v.Value)
		if !ok {
			return object.NewTypeErr(fmt.Sprintf("header value of `%s` cannot be treated as str: `%s`", key, v.Value.Inspect()))
		}

		headers.Add(key, valueStr.Value)
	}
	return nil
}

func addQueries(u *url.URL, queryObj *object.PanObj) *object.PanErr {
	// queries cannot be added directly because u may already contain other queries
	q := u.Query()
	for k, v := range *queryObj.Pairs {
		keyObj, _ := object.SymHash2Str(k)
		key := keyObj.(*object.PanStr).Value

		valueStr, ok := object.TraceProtoOfStr(v.Value)
		if !ok {
			return object.NewTypeErr(fmt.Sprintf("query value of `%s` cannot be treated as str: `%s`", key, v.Value.Inspect()))
		}

		q.Add(key, valueStr.Value)
	}

	u.RawQuery = q.Encode()

	return nil
}

func responseObj(res *http.Response) object.PanObject {
	var b bytes.Buffer
	if res.Body != nil {
		b.ReadFrom(res.Body)
	}

	headers := map[string]object.PanObject{}
	for k, v := range res.Header {
		elems := []object.PanObject{}
		for _, h := range v {
			elems = append(elems, object.NewPanStr(h))
		}

		headers[k] = object.NewPanArr(elems...)
	}

	return mapToObj(map[string]object.PanObject{
		"status":  object.NewPanInt(int64(res.StatusCode)),
		"body":    object.NewPanStr(b.String()),
		"headers": mapToObj(headers),
	})
}
