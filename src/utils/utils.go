package utils

type filter[V any] func(V) bool

func FilterMap[V any](mp map[int]V, f filter[V]) map[int]V {
	res := map[int]V{}
	for k, v := range mp {
		if f(v) {
			res[k] = v
		}
	}
	return res
}

func FilterMapValues[V any](mp map[int]V, f filter[V]) []V {
	res := []V{}
	for _, v := range mp {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

func FilterList[V any](l []V, f filter[V]) []V {
	res := []V{}
	for _, v := range l {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}
