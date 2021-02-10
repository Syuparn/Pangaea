package object

// ReprStr generate prittified string description of obj.
func ReprStr(obj PanObject) string {
	// if self has '_name and it is str, use it
	// NOTE: only refer _name in self (NOT PROTO)!,
	// otherwise all children are shown as _name
	if name, ok := extractName(obj); ok {
		return name.Value
	}

	return obj.Inspect()
}

func extractName(obj PanObject) (*PanStr, bool) {
	o, ok := obj.(*PanObj)
	if !ok {
		return nil, false
	}
	pair, ok := (*o.Pairs)[GetSymHash("_name")]
	if !ok {
		return nil, false
	}

	name, ok := TraceProtoOfStr(pair.Value)
	if !ok {
		return nil, false
	}

	return name, true
}
