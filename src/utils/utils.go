package utils

type filter[V any] func(V) bool
type mapper[X, Y any] func(X) Y
type reducer[X any] func(X, X) X

func FilterValuesToMap[K comparable, V any](mp map[K]V, f filter[V]) map[K]V {
	res := map[K]V{}
	for k, v := range mp {
		if f(v) {
			res[k] = v
		}
	}
	return res
}

func FilterValuesToList[K comparable, V any](mp map[K]V, f filter[V]) []V {
	res := []V{}
	for _, v := range mp {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

func MapValuesToList[K comparable, V, W any](mp map[K]V, f mapper[V, W]) []W {
	res := []W{}
	for _, v := range mp {
		res = append(res, f(v))
	}
	return res
}

func ValuesToList[K comparable, V any](mp map[K]V) []V {
	res := make([]V, len(mp))
	ptr := 0
	for _, v := range mp {
		res[ptr] = v
		ptr++
	}
	return res
}

func Filter[V any](l []V, f filter[V]) []V {
	res := []V{}
	for _, v := range l {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

func Map[X, Y any](l []X, f mapper[X, Y]) []Y {
	res := make([]Y, len(l))
	for i, v := range l {
		res[i] = f(v)
	}
	return res
}

func ToMapMappingToKey[V any, K comparable](l []V, getKey mapper[V, K]) map[K]V {
	res := map[K]V{}
	for _, v := range l {
		res[getKey(v)] = v
	}
	return res
}

func Reduce[V any](l []V, f reducer[V]) *V {
	if len(l) == 0 {
		return nil
	}
	res := l[0]
	for i := 1; i < len(l); i++ {
		res = f(res, l[i])
	}
	return &res
}

func MapCount[V any, K comparable](l []V, f mapper[V, K]) map[K]int {
	res := map[K]int{}
	for _, v := range l {
		res[f(v)]++
	}
	return res
}

func GroupByMapReduce[K comparable, V, W any](l []W, toKey mapper[W, K], mp mapper[W, V], rd reducer[V]) map[K]V {
	res := map[K]V{}
	for _, w := range l {
		k := toKey(w)
		if v, exists := res[k]; exists {
			res[k] = rd(v, mp(w))
		} else {
			res[k] = mp(w)
		}
	}
	return res
}
