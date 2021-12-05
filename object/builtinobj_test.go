package object

import "testing"

// test zero values of *PanObj type BuiltInObj
func TestBuiltInObjZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		obj      PanObject
		expected PanObject
	}{
		{"BuiltInIntObj", BuiltInIntObj, BuiltInZeroInt},
		{"BuiltInFloatObj", BuiltInFloatObj, zeroFloat},
		{"BuiltInNumObj", BuiltInNumObj, zeroObj},
		{"BuiltInNilObj", BuiltInNilObj, BuiltInNil},
		{"BuiltInStrObj", BuiltInStrObj, zeroStr},
		{"BuiltInArrObj", BuiltInArrObj, zeroArr},
		{"BuiltInRangeObj", BuiltInRangeObj, zeroRange},
		// {"", BuiltInFuncObj, zeroFunc},
		// {"", BuiltInIterObj, zeroIter},
		{"BuiltInIterableObj", BuiltInIterableObj, zeroObj},
		{"BuiltInComparableObj", BuiltInComparableObj, zeroObj},
		{"BuiltInWrappableObj", BuiltInWrappableObj, zeroObj},
		// {"", BuiltInMatchObj, zeroMatch},
		{"BuiltInObjObj", BuiltInObjObj, zeroObj},
		{"BuiltInBaseObj", BuiltInBaseObj, zeroObj},
		{"BuiltInMapObj", BuiltInMapObj, zeroMap},
		{"BuiltInDiamondObj", BuiltInDiamondObj, zeroObj},
		{"BuiltInKernelObj", BuiltInKernelObj, zeroObj},
		{"BuiltInJSONObj", BuiltInJSONObj, zeroObj},
		// {"", BuiltInEitherObj, zeroEither},
		// {"", BuiltInEitherValObj, zeroEitherVal},
		// {"", BuiltInEitherErrObj, zeroEitherErr},
		// {"", BuiltInErrObj, zeroErr},
		// (other errors...),
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			actual := tt.obj.Zero()
			if actual != tt.expected {
				t.Errorf("wrong zero value: expected %s, got %s",
					tt.expected.Repr(), actual.Repr())
			}
		})
	}
}
