package object

// TraceProtoOfArr traces proto chain of obj and returns arr proto.
func TraceProtoOfArr(obj PanObject) (*PanArr, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanArr); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfBool traces proto chain of obj and returns bool proto.
func TraceProtoOfBool(obj PanObject) (*PanBool, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanBool); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfBuiltInFunc traces proto chain of obj and returns builtInFunc proto.
func TraceProtoOfBuiltInFunc(obj PanObject) (*PanBuiltIn, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanBuiltIn); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfBuiltInIter traces proto chain of obj and returns builtInIter proto.
func TraceProtoOfBuiltInIter(obj PanObject) (*PanBuiltInIter, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanBuiltInIter); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfFloat traces proto chain of obj and returns float proto.
func TraceProtoOfFloat(obj PanObject) (*PanFloat, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanFloat); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfFunc traces proto chain of obj and returns func proto.
func TraceProtoOfFunc(obj PanObject) (*PanFunc, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanFunc); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfInt traces proto chain of obj and returns int proto.
func TraceProtoOfInt(obj PanObject) (*PanInt, bool) {
	// HACK: proto of Int is zero value 0 so that Int itself can be used as int object
	if obj == BuiltInIntObj {
		return BuiltInZeroInt, true
	}

	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanInt); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfIO traces proto chain of obj and returns IO proto.
func TraceProtoOfIO(obj PanObject) (*PanIO, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanIO); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfMap traces proto chain of obj and returns map proto.
func TraceProtoOfMap(obj PanObject) (*PanMap, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanMap); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfMatch traces proto chain of obj and returns match proto.
func TraceProtoOfMatch(obj PanObject) (*PanMatch, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanMatch); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfNil traces proto chain of obj and returns nil proto.
func TraceProtoOfNil(obj PanObject) (*PanNil, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanNil); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfObj traces proto chain of obj and returns PanObj type proto.
func TraceProtoOfObj(obj PanObject) (*PanObj, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanObj); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfRange traces proto chain of obj and returns range proto.
func TraceProtoOfRange(obj PanObject) (*PanRange, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanRange); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfStr traces proto chain of obj and returns str proto.
func TraceProtoOfStr(obj PanObject) (*PanStr, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanStr); ok {
			return v, true
		}
	}
	return nil, false
}

// TraceProtoOfErrWrapper traces proto chain of obj and returns errWrapper proto.
func TraceProtoOfErrWrapper(obj PanObject) (*PanErrWrapper, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanErrWrapper); ok {
			return v, true
		}
	}
	return nil, false
}
