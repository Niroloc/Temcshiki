package utils

type filter[V any] func(V) bool
type mapper[X, Y any] func(X) Y

func FilterValuesToMap[V any](mp map[int]V, f filter[V]) map[int]V {
	res := map[int]V{}
	for k, v := range mp {
		if f(v) {
			res[k] = v
		}
	}
	return res
}

func FilterValues[V any](mp map[int]V, f filter[V]) []V {
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

func ListMap[X, Y any](l []X, f mapper[X, Y]) []Y {
	res := make([]Y, len(l))
	for i, v := range l {
		res[i] = f(v)
	}
	return res
}

func ListToMap[V any, K comparable](l []V, getKey mapper[V, K]) map[K]V {
	res := map[K]V{}
	for _, v := range l {
		res[getKey(v)] = v
	}
	return res
}
