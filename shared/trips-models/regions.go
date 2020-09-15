package models

import (
	"errors"

	"gonum.org/v1/gonum/graph/simple"
)

// Region is enum of all Dominican Republic regions ISO codes
type Region string

const (
	// CibaoNordesteRegion ISO code
	CibaoNordesteRegion Region = "DO-33"
	// CibaoNoroesteRegion ISO code
	CibaoNoroesteRegion Region = "DO-34"
	// CibaoNorteRegion ISO code
	CibaoNorteRegion Region = "DO-35"
	// CibaoSurRegion ISO code
	CibaoSurRegion Region = "DO-36"
	// ElValleRegion ISO code
	ElValleRegion Region = "DO-37"
	// EnriquilloRegion ISO code
	EnriquilloRegion Region = "DO-38"
	// HiguamoRegion ISO code
	HiguamoRegion Region = "DO-39"
	// OzamaRegion ISO code
	OzamaRegion Region = "DO-40"
	// ValdesiaRegion ISO code
	ValdesiaRegion Region = "DO-41"
	// YumaRegion ISO code
	YumaRegion Region = "DO-42"
)

var (
	errNoRegionForProvince = errors.New("no region available for given province")

	// ValidRegions is the map of all valid regions
	ValidRegions = map[Region]bool{
		CibaoNordesteRegion: true,
		CibaoNoroesteRegion: true,
		CibaoNorteRegion:    true,
		CibaoSurRegion:      true,
		ElValleRegion:       true,
		EnriquilloRegion:    true,
		HiguamoRegion:       true,
		OzamaRegion:         true,
		ValdesiaRegion:      true,
		YumaRegion:          true,
	}

	// RegionToProvinces is map of provinces contained in each region
	RegionToProvinces = map[Region][]Province{
		CibaoNordesteRegion: []Province{
			ProvinceDuarte,
			ProvinceHermanasMirabal,
			ProvinceMariaTrinidadSanchez,
			ProvinceSamana,
		},
		CibaoNoroesteRegion: []Province{
			ProvinceDajabon,
			ProvinceMonteCristi,
			ProvinceSantiagoRodriguez,
			ProvinceValverde,
		},
		CibaoNorteRegion: []Province{
			ProvinceEspaillat,
			ProvincePuertoPlata,
			ProvinceSantiago,
		},
		CibaoSurRegion: []Province{
			ProvinceLaVega,
			ProvinceMonsenorNouel,
			ProvinceSanchezRamirez,
		},
		ElValleRegion: []Province{
			ProvinceAzua,
			ProvinceEliasPina,
			ProvinceSanJuan,
		},
		EnriquilloRegion: []Province{
			ProvinceBahoruco,
			ProvinceBarahona,
			ProvinceIndependencia,
			ProvincePedernales,
		},
		HiguamoRegion: []Province{
			ProvinceHatoMayor,
			ProvinceSanPedroDeMacoris,
		},
		OzamaRegion: []Province{
			ProvinceDistritoNacional,
			ProvinceMontePlata,
			ProvinceSantoDomingo,
		},
		ValdesiaRegion: []Province{
			ProvincePeravia,
			ProvinceSanCristobal,
			ProvinceSanJoseDeOcoa,
		},
		YumaRegion: []Province{
			ProvinceElSeibo,
			ProvinceLaAltagracia,
			ProvinceLaRomana,
		},
	}

	// Indexes were chosen according to the ISO codes
	cibaoNordesteRegionNode = simple.Node(33)
	cibaoNoroesteRegionNode = simple.Node(34)
	cibaoNorteRegionNode    = simple.Node(35)
	cibaoSurRegionNode      = simple.Node(36)
	elValleRegionNode       = simple.Node(37)
	enriquilloRegionNode    = simple.Node(38)
	higuamoRegionNode       = simple.Node(39)
	ozamaRegionNode         = simple.Node(40)
	valdesiaRegionNode      = simple.Node(41)
	yumaRegionNode          = simple.Node(42)

	// RegionToRegionNode is the map to convert Region to RegionNode
	RegionToRegionNode = map[Region]simple.Node{
		CibaoNordesteRegion: cibaoNordesteRegionNode,
		CibaoNoroesteRegion: cibaoNoroesteRegionNode,
		CibaoNorteRegion:    cibaoNorteRegionNode,
		CibaoSurRegion:      cibaoSurRegionNode,
		ElValleRegion:       elValleRegionNode,
		EnriquilloRegion:    enriquilloRegionNode,
		HiguamoRegion:       higuamoRegionNode,
		OzamaRegion:         ozamaRegionNode,
		ValdesiaRegion:      valdesiaRegionNode,
		YumaRegion:          yumaRegionNode,
	}
)

// GetRegionFromProvince returns Region of given province
func GetRegionFromProvince(givenProvince Province) (Region, error) {
	for region := range RegionToProvinces {
		for _, provinceRegion := range RegionToProvinces[region] {
			if provinceRegion == givenProvince {
				return region, nil
			}
		}
	}

	return "", errNoRegionForProvince
}
