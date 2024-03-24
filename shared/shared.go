package shared

type CityStat struct {
	Min, Max, Sum float64
	Count         int64
}

type CityStatV2 struct {
	Min, Max, Sum int
	Count         int
}

type CityStatV3 struct {
	Name          string
	Min, Max, Sum int
	Count         int
}
