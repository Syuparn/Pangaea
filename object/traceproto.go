package object

func TraceProtoOfArr(obj PanObject) (*PanArr, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanArr); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfBool(obj PanObject) (*PanBool, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanBool); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfBuiltInFunc(obj PanObject) (*PanBuiltIn, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanBuiltIn); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfBuiltInIter(obj PanObject) (*PanBuiltInIter, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanBuiltInIter); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfFloat(obj PanObject) (*PanFloat, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanFloat); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfFunc(obj PanObject) (*PanFunc, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanFunc); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfInt(obj PanObject) (*PanInt, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanInt); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfIO(obj PanObject) (*PanIO, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanIO); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfMap(obj PanObject) (*PanMap, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanMap); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfMatch(obj PanObject) (*PanMatch, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanMatch); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfNil(obj PanObject) (*PanNil, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanNil); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfObj(obj PanObject) (*PanObj, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanObj); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfRange(obj PanObject) (*PanRange, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanRange); ok {
			return v, true
		}
	}
	return nil, false
}

func TraceProtoOfStr(obj PanObject) (*PanStr, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if v, ok := o.(*PanStr); ok {
			return v, true
		}
	}
	return nil, false
}
