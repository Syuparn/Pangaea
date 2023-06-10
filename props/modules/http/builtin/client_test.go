package builtin

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Syuparn/pangaea/object"
)

func TestRequest(t *testing.T) {
	tests := []struct {
		name string
		// HACK: httptest.Server endpoint is dynamic
		mapToObj func(string) *object.PanObj
		expected func(string) *http.Request
	}{
		{
			"get",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr(endpoint),
				})
			},
			func(endpoint string) *http.Request {
				return &http.Request{
					Method: "GET",
					URL:    must(url.Parse("/")),
					Host:   strings.TrimPrefix(endpoint, "http://"),
				}
			},
		},
		{
			"get with path",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr(endpoint + "/foo"),
				})
			},
			func(endpoint string) *http.Request {
				return &http.Request{
					Method: "GET",
					URL:    must(url.Parse("/foo")),
					Host:   strings.TrimPrefix(endpoint, "http://"),
				}
			},
		},
		{
			"get with queries",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr(endpoint + "/foo"),
					"queries": mapToObj(map[string]object.PanObject{
						"query1": object.NewPanStr("value1"),
						"query2": object.NewPanStr("value2"),
					}),
				})
			},
			func(endpoint string) *http.Request {
				return &http.Request{
					Method: "GET",
					URL:    must(url.Parse("/foo?query1=value1&query2=value2")),
					Host:   strings.TrimPrefix(endpoint, "http://"),
				}
			},
		},
		{
			"get with queries (url also has queries)",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr(endpoint + "/foo?query1=value1"),
					"queries": mapToObj(map[string]object.PanObject{
						"query2": object.NewPanStr("value2"),
					}),
				})
			},
			func(endpoint string) *http.Request {
				return &http.Request{
					Method: "GET",
					URL:    must(url.Parse("/foo?query1=value1&query2=value2")),
					Host:   strings.TrimPrefix(endpoint, "http://"),
				}
			},
		},
		// TODO: allow duplicated queries
		{
			"get with headers",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr(endpoint + "/foo"),
					"headers": mapToObj(map[string]object.PanObject{
						"Content-Type": object.NewPanStr("text/plain"),
						"Accept":       object.NewPanStr("application/json"),
					}),
				})
			},
			func(endpoint string) *http.Request {
				return &http.Request{
					Method: "GET",
					URL:    must(url.Parse("/foo")),
					Host:   strings.TrimPrefix(endpoint, "http://"),
					Header: map[string][]string{
						"Content-Type": {"text/plain"},
						"Accept":       {"application/json"},
					},
				}
			},
		},
		// TODO: allow duplicated headers
		{
			"post with body",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("POST"),
					"url":    object.NewPanStr(endpoint),
					"body":   object.NewPanStr("hello"),
				})
			},
			func(endpoint string) *http.Request {
				return &http.Request{
					Method: "POST",
					URL:    must(url.Parse("/")),
					Host:   strings.TrimPrefix(endpoint, "http://"),
					Body:   io.NopCloser(strings.NewReader("hello")),
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			var actual *http.Request
			var actualBody bytes.Buffer
			h := func(w http.ResponseWriter, r *http.Request) {
				actual = r
				if r.Method != "GET" && r.Body != nil {
					defer r.Body.Close()
					actualBody.ReadFrom(r.Body)
				}
				io.WriteString(w, "ok")
			}

			ts := httptest.NewServer(http.HandlerFunc(h))
			defer ts.Close()

			expected := tt.expected(ts.URL)
			mapToObj := tt.mapToObj(ts.URL)

			env := object.NewEnv()
			res := request(env, mapToObj)

			if res.Type() == object.ErrType {
				t.Fatalf("error must not be raised: %s", res.Inspect())
			}

			if expected.Method != actual.Method {
				t.Errorf("wrong Method. expected=%v, got=%v", expected.Method, actual.Method)
			}

			if expected.URL.String() != actual.URL.String() {
				t.Errorf("wrong URL. expected=%v, got=%v", expected.URL, actual.URL)
			}
			if expected.Host != actual.Host {
				t.Errorf("wrong Host. expected=%v, got=%v", expected.Host, actual.Host)
			}
			if expected.Host != actual.Host {
				t.Errorf("wrong Host. expected=%v, got=%v", expected.Host, actual.Host)
			}

			for k, v := range expected.Header {
				actualHeaderValues, ok := actual.Header[k]
				if !ok {
					t.Errorf("Header %s not found", k)
				}

				for i, s := range actualHeaderValues {
					if s != v[i] {
						t.Errorf("wrong Header value (key %q). expected=%s, got=%s", k, s, v[i])
					}
				}
			}

			if expected.Body != nil {
				var expectedBody bytes.Buffer
				expectedBody.ReadFrom(expected.Body)
				if expectedBody.String() != actualBody.String() {
					t.Errorf("wrong Body. expected=%v, got=%v", expectedBody.String(), actualBody.String())
				}
			}
		})
	}
}

func TestRequestError(t *testing.T) {
	tests := []struct {
		name string
		// HACK: httptest.Server endpoint is dynamic
		mapToObj func(string) *object.PanObj
		expected *object.PanErr
	}{
		{
			"method is missing",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"url": object.NewPanStr(endpoint),
				})
			},
			object.NewTypeErr("method must be specified"),
		},
		{
			"method is not str",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanInt(1),
					"url":    object.NewPanStr(endpoint),
				})
			},
			object.NewTypeErr("method `1` cannot be treated as str"),
		},
		{
			"url is missing",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
				})
			},
			object.NewTypeErr("url must be specified"),
		},
		{
			"method is not str",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanInt(1),
				})
			},
			object.NewTypeErr("url `1` cannot be treated as str"),
		},
		{
			"url cannot be parsed",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr("\x01"),
				})
			},
			object.NewValueErr("failed to parse url: `\x01`"),
		},
		{
			"query value is not str",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr(endpoint),
					"queries": mapToObj(map[string]object.PanObject{
						"query": object.NewPanInt(1),
					}),
				})
			},
			object.NewTypeErr("query value of `query` cannot be treated as str: `1`"),
		},
		{
			"header value is not str",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr(endpoint),
					"headers": mapToObj(map[string]object.PanObject{
						"Content-Type": object.NewPanInt(1),
					}),
				})
			},
			object.NewTypeErr("header value of `Content-Type` cannot be treated as str: `1`"),
		},
		{
			"body is not str",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr(endpoint),
					"body":   object.NewPanInt(1),
				})
			},
			object.NewTypeErr("body `1` cannot be treated as str"),
		},
		{
			"failed to create request",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("ダミー"),
					"url":    object.NewPanStr(endpoint),
				})
			},
			object.NewPanErr("net/http: invalid method \"ダミー\""),
		},
		{
			"failed to request",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr(""),
				})
			},
			object.NewPanErr("Get \"\": unsupported protocol scheme \"\""),
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			h := func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "ok")
			}

			ts := httptest.NewServer(http.HandlerFunc(h))
			defer ts.Close()
			mapToObj := tt.mapToObj(ts.URL)

			env := object.NewEnv()
			res := request(env, mapToObj)

			if res.Inspect() != tt.expected.Inspect() {
				t.Errorf("wrong value. expected=%v, got=%v", tt.expected.Inspect(), res.Inspect())
			}
		})
	}
}

func TestRequestResponse(t *testing.T) {
	tests := []struct {
		name string
		// HACK: httptest.Server endpoint is dynamic
		mapToObj func(string) *object.PanObj
		expected object.PanObject
		handler  func(http.ResponseWriter, *http.Request)
	}{
		{
			"simple reponse",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr(endpoint),
				})
			},
			mapToObj(map[string]object.PanObject{
				"status": object.NewPanInt(200),
				"body":   object.NewPanStr("ok"),
				"headers": mapToObj(map[string]object.PanObject{
					"Content-Length": object.NewPanArr(object.NewPanStr("2")),
					"Content-Type":   object.NewPanArr(object.NewPanStr("text/plain; charset=utf-8")),
					"Date":           object.NewPanArr(object.NewPanStr("Sun, 1 Jan 2023 01:23:45 GMT")),
				}),
			}),
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Date", "Sun, 1 Jan 2023 01:23:45 GMT") // HACK: fix time-dependent header for comparison
				io.WriteString(w, "ok")
			},
		},
		{
			"Internal Server Error",
			func(endpoint string) *object.PanObj {
				return mapToObj(map[string]object.PanObject{
					"method": object.NewPanStr("GET"),
					"url":    object.NewPanStr(endpoint),
				})
			},
			mapToObj(map[string]object.PanObject{
				"status": object.NewPanInt(500),
				"body":   object.NewPanStr("ng"),
				"headers": mapToObj(map[string]object.PanObject{
					"Content-Length": object.NewPanArr(object.NewPanStr("2")),
					"Content-Type":   object.NewPanArr(object.NewPanStr("text/plain; charset=utf-8")),
					"Date":           object.NewPanArr(object.NewPanStr("Sun, 1 Jan 2023 01:23:45 GMT")),
				}),
			}),
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Date", "Sun, 1 Jan 2023 01:23:45 GMT") // HACK: fix time-dependent header for comparison
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, "ng")
			},
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer ts.Close()

			mapToObj := tt.mapToObj(ts.URL)

			env := object.NewEnv()
			res := request(env, mapToObj)

			if res.Inspect() != tt.expected.Inspect() {
				t.Errorf("wrong value. expected=%v, got=%v", tt.expected.Inspect(), res.Inspect())
			}
		})
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
