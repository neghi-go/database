package database

type E struct {
	key   string
	value interface{}
}

func (e E) Key() string {
	return e.key
}
func (e E) Value() interface{} {
	return e.value
}

type D []E

type Params func() E

func SetParams(params ...Params) D {
	res := D{}
	for _, p := range params {
		res = append(res, p())
	}
	return res
}

func SetFilter(key string, value interface{}) Params {
	return func() E {
		return E{key: key, value: value}
	}
}

type OrderType int

const (
	ASC OrderType = iota
	DESC
)

var orderTypeMap = map[OrderType]string{
	ASC:  "asc",
	DESC: "desc",
}

func (o OrderType) String() string {
	return orderTypeMap[o]
}

func SetOrder(key string, order OrderType) Params {
	return func() E {
		return E{key: key, value: order}
	}
}
