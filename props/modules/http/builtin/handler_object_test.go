package builtin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/parser"
	"github.com/labstack/echo/v4"
)

func dummyCallback(input string) *object.PanFunc {
	// HACK: inject props that are necessary for tests
	object.BuiltInArrObj.AddPairs(&map[object.SymHash]object.Pair{
		object.GetSymHash("at"): {Key: object.NewPanStr("at"), Value: evaluator.NewPropContainer()["Arr_at"]},
	})

	env := object.NewEnvWithConsts()
	node, err := parser.Parse(parser.NewReader(strings.NewReader(input), "<stdin>"))
	if err != nil {
		panic(err)
	}
	panObject := evaluator.Eval(node, env)
	return panObject.(*object.PanFunc)
}

func TestHandlerType(t *testing.T) {
	obj := newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`))
	if obj.Type() != handlerType {
		t.Fatalf("wrong type: expected=%s, got=%s", handlerType, obj.Type())
	}
}

func TestHandlerInspect(t *testing.T) {
	tests := []struct {
		obj      *panHandler
		expected string
	}{
		{
			newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`)),
			`[handler]`,
		},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestHandlerRepr(t *testing.T) {
	tests := []struct {
		obj      *panHandler
		expected string
	}{
		{
			newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`)),
			`[handler]`,
		},
	}

	for _, tt := range tests {
		if tt.obj.Repr() != tt.expected {
			t.Errorf("wrong output: expected=%s, got=%s",
				tt.expected, tt.obj.Inspect())
		}
	}
}

func TestHandlerProto(t *testing.T) {
	a := newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`))
	if a.Proto() != object.BuiltInObjObj {
		t.Fatalf("Proto is not object.BuiltInObjObj. got=%T (%+v)",
			a.Proto(), a.Proto())
	}
}

func TestHandlerZero(t *testing.T) {
	tests := []struct {
		obj *panHandler
	}{
		{newPanHandler(object.NewEnv(), "GET", "/", dummyCallback(`{|res| "ok"}`))},
	}

	for _, tt := range tests {
		tt := tt // pin

		actual := tt.obj.Zero()

		if actual != tt.obj {
			t.Errorf("zero must be itself (%#v). got=%s (%#v)",
				tt.obj, actual.Repr(), actual)
		}
	}
}

// checked by compiler (this function works nothing)
func testHandlerIsPanObject() {
	var _ object.PanObject = &panHandler{}
}

func TestToHandler(t *testing.T) {
	tests := []struct {
		name    string
		obj     *object.PanFunc
		status  int
		body    string
		headers map[string][]string
	}{
		{
			"simple response",
			dummyCallback(`{|req| {}}`),
			200,
			"",
			map[string][]string{
				"Content-Type": {"text/plain; charset=UTF-8"},
			},
		},
		{
			"with body",
			dummyCallback("{|req| {body: \"abc\"}}"),
			200,
			"abc",
			map[string][]string{
				"Content-Type": {"text/plain; charset=UTF-8"},
			},
		},
		{
			"with json body",
			dummyCallback("{|req| {body: \"{\\\"a\\\": 1}\", _isJSON: true}}"),
			200,
			`{"a": 1}`,
			map[string][]string{
				"Content-Type": {"application/json; charset=UTF-8"},
			},
		},
		{
			"with status code",
			dummyCallback("{|req| {status: 204}}"),
			204,
			``,
			map[string][]string{
				"Content-Type": {"text/plain; charset=UTF-8"},
			},
		},
		{
			"with headers (keys are capitalized)",
			dummyCallback("{|req| {headers: {foo: \"bar\", baz: \"quux\"}}}"),
			200,
			``,
			map[string][]string{
				"Content-Type": {"text/plain; charset=UTF-8"},
				"Foo":          {"bar"},
				"Baz":          {"quux"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := toHandler(object.NewEnvWithConsts(), tt.obj)

			err := h(c)
			if err != nil {
				t.Fatalf("error raised: %s", err)
			}

			if rec.Code != tt.status {
				t.Errorf("wrong status. expected=%v, got=%v", tt.status, rec.Code)
			}
			if string(rec.Body.Bytes()) != tt.body {
				t.Errorf("wrong body. expected=%v, got=%v", tt.body, string(rec.Body.Bytes()))
			}

			for k, v := range tt.headers {
				actualHeaderValues, ok := rec.Result().Header[k]
				if !ok {
					t.Errorf("Header %s not found", k)
				}

				for i, s := range actualHeaderValues {
					if s != v[i] {
						t.Errorf("wrong Header value (key %q). expected=%s, got=%s", k, s, v[i])
					}
				}
			}
		})
	}
}

func TestToHandlerRequest(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		url      string
		body     string
		header   http.Header
		obj      *object.PanFunc
		expected string
	}{
		{
			"request body",
			"POST",
			"/",
			"abc",
			map[string][]string{},
			dummyCallback(`{|req| {body: req.body}}`),
			"abc",
		},
		{
			"request method",
			"POST",
			"/",
			"abc",
			map[string][]string{},
			dummyCallback(`{|req| {body: req.method}}`),
			"POST",
		},
		{
			"request path",
			"POST",
			"/foo/bar",
			"abc",
			map[string][]string{},
			dummyCallback(`{|req| {body: req.url}}`),
			"/foo/bar",
		},
		{
			"request headers",
			"POST",
			"/",
			"",
			map[string][]string{
				"Foo": {"bar"},
			},
			// HACK: pick up only 1 element because `Arr#S` is not injected
			dummyCallback(`{|req| {body: req.headers.Foo[0]}}`),
			"bar",
		},
		{
			"request query parameters",
			"GET",
			"/users?foo=bar",
			"",
			map[string][]string{},
			// HACK: pick up only 1 element because `Arr#S` is not injected
			dummyCallback(`{|req| {body: req.queries.foo[0]}}`),
			"bar",
		},
		// TODO: path parameters
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			req.Header = tt.header
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := toHandler(object.NewEnvWithConsts(), tt.obj)

			err := h(c)
			if err != nil {
				t.Fatalf("error raised: %s", err)
			}

			if string(rec.Body.Bytes()) != tt.expected {
				t.Errorf("wrong body. expected=%v, got=%v", tt.expected, string(rec.Body.Bytes()))
			}
		})
	}
}

func TestToHandlerError(t *testing.T) {
	tests := []struct {
		name     string
		obj      *object.PanFunc
		expected string
	}{
		{
			"func raises an error",
			dummyCallback(`{|req| a}`),
			"NameErr: name `a` is not defined",
		},
		{
			"headers contains non-string value",
			dummyCallback(`{|req| {headers: {foo: 1}}}`),
			"TypeErr: `1` cannot be treated as str",
		},
		{
			"body is not str",
			dummyCallback(`{|req| {body: 1}}`),
			"TypeErr: `1` cannot be treated as str",
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := toHandler(object.NewEnvWithConsts(), tt.obj)

			err := h(c)
			if err == nil {
				t.Fatalf("error must be raised")
			}
			if err.Error() != tt.expected {
				t.Errorf("wrong body. expected=%v, got=%v", tt.expected, err.Error())
			}
		})
	}
}
