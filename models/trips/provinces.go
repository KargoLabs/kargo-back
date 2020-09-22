package models

// Province is enum of all Dominican Republic provinces ISO codes
type Province string

const (
	// ProvinceDistritoNacional ISO code
	ProvinceDistritoNacional Province = "DO-01"
	// ProvinceAzua ISO code
	ProvinceAzua Province = "DO-02"
	// ProvinceBahoruco ISO code
	ProvinceBahoruco Province = "DO-03"
	// ProvinceBarahona ISO code
	ProvinceBarahona Province = "DO-04"
	// ProvinceDajabon ISO code
	ProvinceDajabon Province = "DO-05"
	// ProvinceDuarte ISO code
	ProvinceDuarte Province = "DO-06"
	// ProvinceEliasPina ISO code
	ProvinceEliasPina Province = "DO-07"
	// ProvinceElSeibo ISO code
	ProvinceElSeibo Province = "DO-08"
	// ProvinceEspaillat ISO code
	ProvinceEspaillat Province = "DO-09"
	// ProvinceIndependencia ISO code
	ProvinceIndependencia Province = "DO-10"
	// ProvinceLaAltagracia ISO code
	ProvinceLaAltagracia Province = "DO-11"
	// ProvinceLaRomana ISO code
	ProvinceLaRomana Province = "DO-12"
	// ProvinceLaVega ISO code
	ProvinceLaVega Province = "DO-13"
	// ProvinceMariaTrinidadSanchez ISO code
	ProvinceMariaTrinidadSanchez Province = "DO-14"
	// ProvinceMonteCristi ISO code
	ProvinceMonteCristi Province = "DO-15"
	// ProvincePedernales ISO code
	ProvincePedernales Province = "DO-16"
	// ProvincePeravia ISO code
	ProvincePeravia Province = "DO-17"
	// ProvincePuertoPlata ISO code
	ProvincePuertoPlata Province = "DO-18"
	// ProvinceHermanasMirabal ISO code
	ProvinceHermanasMirabal Province = "DO-19"
	// ProvinceSamana ISO code
	ProvinceSamana Province = "DO-20"
	// ProvinceSanCristobal ISO code
	ProvinceSanCristobal Province = "DO-21"
	// ProvinceSanJuan ISO code
	ProvinceSanJuan Province = "DO-22"
	// ProvinceSanPedroDeMacoris ISO code
	ProvinceSanPedroDeMacoris Province = "DO-23"
	// ProvinceSanchezRamirez ISO code
	ProvinceSanchezRamirez Province = "DO-24"
	// ProvinceSantiago ISO code
	ProvinceSantiago Province = "DO-25"
	// ProvinceSantiagoRodriguez ISO code
	ProvinceSantiagoRodriguez Province = "DO-26"
	// ProvinceValverde ISO code
	ProvinceValverde Province = "DO-27"
	// ProvinceMonsenorNouel ISO code
	ProvinceMonsenorNouel Province = "DO-28"
	// ProvinceMontePlata ISO code
	ProvinceMontePlata Province = "DO-29"
	// ProvinceHatoMayor ISO code
	ProvinceHatoMayor Province = "DO-30"
	// ProvinceSanJoseDeOcoa ISO code
	ProvinceSanJoseDeOcoa Province = "DO-31"
	// ProvinceSantoDomingo ISO code
	ProvinceSantoDomingo Province = "DO-32"
)

var (
	// ValidProvinces is the map of all valid provinces
	ValidProvinces = map[Province]bool{
		ProvinceDistritoNacional:     true,
		ProvinceAzua:                 true,
		ProvinceBahoruco:             true,
		ProvinceBarahona:             true,
		ProvinceDajabon:              true,
		ProvinceDuarte:               true,
		ProvinceEliasPina:            true,
		ProvinceElSeibo:              true,
		ProvinceEspaillat:            true,
		ProvinceIndependencia:        true,
		ProvinceLaAltagracia:         true,
		ProvinceLaRomana:             true,
		ProvinceLaVega:               true,
		ProvinceMariaTrinidadSanchez: true,
		ProvinceMonteCristi:          true,
		ProvincePedernales:           true,
		ProvincePeravia:              true,
		ProvincePuertoPlata:          true,
		ProvinceHermanasMirabal:      true,
		ProvinceSamana:               true,
		ProvinceSanCristobal:         true,
		ProvinceSanJuan:              true,
		ProvinceSanPedroDeMacoris:    true,
		ProvinceSanchezRamirez:       true,
		ProvinceSantiago:             true,
		ProvinceSantiagoRodriguez:    true,
		ProvinceValverde:             true,
		ProvinceMonsenorNouel:        true,
		ProvinceMontePlata:           true,
		ProvinceHatoMayor:            true,
		ProvinceSanJoseDeOcoa:        true,
		ProvinceSantoDomingo:         true,
	}
)
