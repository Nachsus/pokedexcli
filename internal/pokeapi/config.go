package pokeapi

type config struct {
	mapBaseUrl string
	mapNextUrl string
	mapPrevUrl string
}

var Conf = config{
	mapBaseUrl: "https://pokeapi.co/api/v2/location-area/",
	mapNextUrl: "",
	mapPrevUrl: "",
}
