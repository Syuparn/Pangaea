package object

import (
	"testing"
)

func TestTraceProtoOfArr(t *testing.T) {
	tests := []struct {
		name     string
		obj      PanObject
		expected *PanArr
	}{
		{
			"arr returns itself",
			zeroArr,
			zeroArr,
		},
		{
			"Arr returns zero value [] so that Arr itself can be used as arr object",
			BuiltInArrObj,
			zeroArr,
		},
		{
			"child of Arr returns zero value []",
			NewPanObj(&map[SymHash]Pair{}, BuiltInArrObj),
			zeroArr,
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			actual, ok := TraceProtoOfArr(tt.obj)

			if !ok {
				t.Fatalf("ok must be true (obj=%s)", tt.obj.Repr())
			}

			if actual != tt.expected {
				t.Errorf("proto must be %s(%T). got=%s(%T)",
					tt.expected.Repr(), tt.expected, actual.Repr(), actual)
			}
		})
	}
}

func TestTraceProtoOfArrFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfArr(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfBool(t *testing.T) {
	proto := BuiltInFalse

	tests := []struct {
		obj      PanObject
		expected *PanBool
	}{
		// return proto
		{
			NewPanObj(&map[SymHash]Pair{}, proto),
			proto,
		},
		// return itself
		{
			proto,
			proto,
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfBool(tt.obj)

		if !ok {
			t.Errorf("ok must be true (obj=%v)", tt.obj)
		}

		if actual != tt.expected {
			t.Errorf("proto must be %+v(%T). got=%+v(%T)",
				tt.expected, tt.expected, actual, actual)
		}
	}
}

func TestTraceProtoOfBoolFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfBool(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfBuiltInFunc(t *testing.T) {
	proto := NewPanBuiltInFunc(newMockBuiltInFunc())

	tests := []struct {
		obj      PanObject
		expected *PanBuiltIn
	}{
		// return proto
		{
			NewPanObj(&map[SymHash]Pair{}, proto),
			proto,
		},
		// return itself
		{
			proto,
			proto,
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfBuiltInFunc(tt.obj)

		if !ok {
			t.Errorf("ok must be true (obj=%v)", tt.obj)
		}

		if actual != tt.expected {
			t.Errorf("proto must be %+v(%T). got=%+v(%T)",
				tt.expected, tt.expected, actual, actual)
		}
	}
}

func newMockBuiltInFunc() BuiltInFunc {
	return func(e *Env, Kwargs *PanObj, args ...PanObject) PanObject {
		return nil
	}
}

func TestTraceProtoOfBuiltInFuncFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfBuiltInFunc(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfBuiltInIter(t *testing.T) {
	proto := NewPanBuiltInIter(newMockBuiltInFunc(), NewEnv())

	tests := []struct {
		obj      PanObject
		expected *PanBuiltInIter
	}{
		// return proto
		{
			NewPanObj(&map[SymHash]Pair{}, proto),
			proto,
		},
		// return itself
		{
			proto,
			proto,
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfBuiltInIter(tt.obj)

		if !ok {
			t.Errorf("ok must be true (obj=%v)", tt.obj)
		}

		if actual != tt.expected {
			t.Errorf("proto must be %+v(%T). got=%+v(%T)",
				tt.expected, tt.expected, actual, actual)
		}
	}
}

func TestTraceProtoOfBuiltInIterFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfBuiltInIter(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfFloat(t *testing.T) {
	proto := NewPanFloat(0.0)

	tests := []struct {
		obj      PanObject
		expected *PanFloat
	}{
		// return proto
		{
			NewPanObj(&map[SymHash]Pair{}, proto),
			proto,
		},
		// return itself
		{
			proto,
			proto,
		},
		// Float returns zero value 0.0 so that Float itself can be used as float object
		{
			BuiltInFloatObj,
			zeroFloat,
		},
		// child of Float
		{
			NewPanObj(&map[SymHash]Pair{}, BuiltInFloatObj),
			zeroFloat,
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfFloat(tt.obj)

		if !ok {
			t.Errorf("ok must be true (obj=%v)", tt.obj)
		}

		if actual != tt.expected {
			t.Errorf("proto must be %+v(%T). got=%+v(%T)",
				tt.expected, tt.expected, actual, actual)
		}
	}
}

func TestTraceProtoOfFloatFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfFloat(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfFunc(t *testing.T) {
	proto := NewPanFunc(newMockFuncWrapper(), NewEnv())

	tests := []struct {
		obj      PanObject
		expected *PanFunc
	}{
		// return proto
		{
			NewPanObj(&map[SymHash]Pair{}, proto),
			proto,
		},
		// return itself
		{
			proto,
			proto,
		},
		// TODO: enable to use Func as func (besides 0 value, Func#call must handle PanObj Func)
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfFunc(tt.obj)

		if !ok {
			t.Errorf("ok must be true (obj=%v)", tt.obj)
		}

		if actual != tt.expected {
			t.Errorf("proto must be %+v(%T). got=%+v(%T)",
				tt.expected, tt.expected, actual, actual)
		}
	}
}

func TestTraceProtoOfFuncFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfFunc(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfInt(t *testing.T) {
	tests := []struct {
		name     string
		obj      PanObject
		expected *PanInt
	}{
		{
			"return proto",
			NewPanObj(&map[SymHash]Pair{}, BuiltInOneInt),
			BuiltInOneInt,
		},
		{
			"1 returns itself",
			BuiltInOneInt,
			BuiltInOneInt,
		},
		{
			"Int returns zero value 0 so that Int itself can be used as int object",
			BuiltInIntObj,
			BuiltInZeroInt,
		},
		// child of Int
		{
			"child of Int returns zero value 0",
			NewPanObj(&map[SymHash]Pair{}, BuiltInIntObj),
			BuiltInZeroInt,
		},
	}

	for _, tt := range tests {
		tt := tt // pin
		t.Run(tt.name, func(t *testing.T) {
			actual, ok := TraceProtoOfInt(tt.obj)

			if !ok {
				t.Errorf("ok must be true (obj=%v)", tt.obj)
			}

			if actual != tt.expected {
				t.Errorf("proto must be %+v(%T). got=%+v(%T)",
					tt.expected, tt.expected, actual, actual)
			}
		})
	}
}

func TestTraceProtoOfIntFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfInt(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfIO(t *testing.T) {
	proto := &PanIO{}

	tests := []struct {
		obj      PanObject
		expected *PanIO
	}{
		// return proto
		{
			NewPanObj(&map[SymHash]Pair{}, proto),
			proto,
		},
		// return itself
		{
			proto,
			proto,
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfIO(tt.obj)

		if !ok {
			t.Errorf("ok must be true (obj=%v)", tt.obj)
		}

		if actual != tt.expected {
			t.Errorf("proto must be %+v(%T). got=%+v(%T)",
				tt.expected, tt.expected, actual, actual)
		}
	}
}

func TestTraceProtoOfIOFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfIO(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfMap(t *testing.T) {
	proto := NewEmptyPanMap()

	tests := []struct {
		obj      PanObject
		expected *PanMap
	}{
		// return proto
		{
			NewPanObj(&map[SymHash]Pair{}, proto),
			proto,
		},
		// return itself
		{
			proto,
			proto,
		},
		// Map returns zero value %{} so that Map itself can be used as map object
		{
			BuiltInMapObj,
			zeroMap,
		},
		// child of Map
		{
			NewPanObj(&map[SymHash]Pair{}, BuiltInMapObj),
			zeroMap,
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfMap(tt.obj)

		if !ok {
			t.Errorf("ok must be true (obj=%v)", tt.obj)
		}

		if actual != tt.expected {
			t.Errorf("proto must be %+v(%T). got=%+v(%T)",
				tt.expected, tt.expected, actual, actual)
		}
	}
}

func TestTraceProtoOfMapFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfMap(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfMatch(t *testing.T) {
	proto := &PanMatch{}

	tests := []struct {
		obj      PanObject
		expected *PanMatch
	}{
		// return proto
		{
			NewPanObj(&map[SymHash]Pair{}, proto),
			proto,
		},
		// return itself
		{
			proto,
			proto,
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfMatch(tt.obj)

		if !ok {
			t.Errorf("ok must be true (obj=%v)", tt.obj)
		}

		if actual != tt.expected {
			t.Errorf("proto must be %+v(%T). got=%+v(%T)",
				tt.expected, tt.expected, actual, actual)
		}
	}
}

func TestTraceProtoOfMatchFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfMatch(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfNil(t *testing.T) {
	tests := []struct {
		name     string
		obj      PanObject
		expected *PanNil
	}{
		{
			"return proto",
			NewPanObj(&map[SymHash]Pair{}, BuiltInNil),
			BuiltInNil,
		},
		{
			"nil returns itself",
			BuiltInNil,
			BuiltInNil,
		},
		{
			"Nil returns zero value nil so that Nil itself can be used as nil object",
			BuiltInNilObj,
			BuiltInNil,
		},
		{
			"child of Nil returns zero value nil",
			NewPanObj(&map[SymHash]Pair{}, BuiltInNilObj),
			BuiltInNil,
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			actual, ok := TraceProtoOfNil(tt.obj)

			if !ok {
				t.Fatalf("ok must be true (obj=%s)", tt.obj.Repr())
			}

			if actual != tt.expected {
				t.Errorf("proto must be %s(%T). got=%s(%T)",
					tt.expected.Repr(), tt.expected, actual.Repr(), actual)
			}
		})
	}
}

func TestTraceProtoOfNilFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfNil(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfObj(t *testing.T) {
	proto := BuiltInIntObj

	tests := []struct {
		name     string
		obj      PanObject
		expected *PanObj
	}{
		{
			"return *PanObj from ancestors",
			BuiltInOneInt,
			proto,
		},
		{
			"return itself",
			proto,
			proto,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, ok := TraceProtoOfObj(tt.obj)

			if !ok {
				t.Errorf("ok must be true (obj=%v)", tt.obj)
			}

			if actual != tt.expected {
				t.Errorf("proto must be %+v(%T). got=%+v(%T)",
					tt.expected, tt.expected, actual, actual)
			}
		})
	}
}

func TestTraceProtoOfObjFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		// NOTE: all valid panobjects have ObjType proto BuiltInBaseObj.
		// In this test, Invaild err object (, whose proto is nil) is used.
		{
			&PanErr{},
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfObj(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfRange(t *testing.T) {
	proto := NewPanRange(NewPanNil(), NewPanNil(), NewPanNil())

	tests := []struct {
		obj      PanObject
		expected *PanRange
	}{
		// return proto
		{
			NewPanObj(&map[SymHash]Pair{}, proto),
			proto,
		},
		// return itself
		{
			proto,
			proto,
		},
		// Range returns zero value (nil:nil) so that Range itself can be used as range object
		{
			BuiltInRangeObj,
			zeroRange,
		},
		// child of Range
		{
			NewPanObj(&map[SymHash]Pair{}, BuiltInRangeObj),
			zeroRange,
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfRange(tt.obj)

		if !ok {
			t.Errorf("ok must be true (obj=%v)", tt.obj)
		}

		if actual != tt.expected {
			t.Errorf("proto must be %+v(%T). got=%+v(%T)",
				tt.expected, tt.expected, actual, actual)
		}
	}
}

func TestTraceProtoOfRangeFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfRange(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfStr(t *testing.T) {
	proto := NewPanStr("")

	tests := []struct {
		obj      PanObject
		expected *PanStr
	}{
		// return proto
		{
			NewPanObj(&map[SymHash]Pair{}, proto),
			proto,
		},
		// return itself
		{
			proto,
			proto,
		},
		// Str returns zero value "" so that Str itself can be used as str object
		{
			BuiltInStrObj,
			zeroStr,
		},
		// child of Str
		{
			NewPanObj(&map[SymHash]Pair{}, BuiltInStrObj),
			zeroStr,
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfStr(tt.obj)

		if !ok {
			t.Errorf("ok must be true (obj=%v)", tt.obj)
		}

		if actual != tt.expected {
			t.Errorf("proto must be %+v(%T). got=%+v(%T)",
				tt.expected, tt.expected, actual, actual)
		}
	}
}

func TestTraceProtoOfStrFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfStr(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}

func TestTraceProtoOfErrWrapper(t *testing.T) {
	proto := WrapErr(NewPanErr("error"))

	tests := []struct {
		obj      PanObject
		expected *PanErrWrapper
	}{
		// return proto
		{
			NewPanObj(&map[SymHash]Pair{}, proto),
			proto,
		},
		// return itself
		{
			proto,
			proto,
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfErrWrapper(tt.obj)

		if !ok {
			t.Errorf("ok must be true (obj=%v)", tt.obj)
		}

		if actual != tt.expected {
			t.Errorf("proto must be %+v(%T). got=%+v(%T)",
				tt.expected, tt.expected, actual, actual)
		}
	}
}

func TestTraceProtoOfErrWrapperFailed(t *testing.T) {
	tests := []struct {
		obj PanObject
	}{
		{
			PanObjInstancePtr(&map[SymHash]Pair{}),
		},
	}

	for _, tt := range tests {
		actual, ok := TraceProtoOfErrWrapper(tt.obj)

		if ok {
			t.Errorf("ok must be false (obj=%v)", tt.obj)
		}

		if actual != nil {
			t.Errorf("actual must be nil. got=%+v(%T)", actual, actual)
		}
	}
}
