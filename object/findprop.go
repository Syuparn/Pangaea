package object

// FindPropAlongProtos traces proto chain and returns prop.
func FindPropAlongProtos(o PanObject, propHash SymHash) (PanObject, bool) {
	// trace prototype chains
	for obj := o; obj != nil; obj = obj.Proto() {
		prop, ok := findProp(obj, propHash)

		if ok {
			return prop, true
		}
	}
	return nil, false
}

// FindPropOwner traces proto chain and returns the prop's owner.
func FindPropOwner(o PanObject, propHash SymHash) (PanObject, bool) {
	// trace prototype chains
	for obj := o; obj != nil; obj = obj.Proto() {
		_, ok := findProp(obj, propHash)

		if ok {
			return obj, true
		}
	}
	return nil, false
}

func findProp(o PanObject, propHash SymHash) (PanObject, bool) {
	obj, ok := o.(*PanObj)
	if !ok {
		return nil, false
	}

	elem, ok := (*obj.Pairs)[propHash]

	if !ok {
		return nil, false
	}

	return elem.Value, true
}
