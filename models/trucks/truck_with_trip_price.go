package models

import trips "kargo-back/models/trips"

type TruckWithTripPrice struct {
	*Truck
	TripPrice float64 `json:"trip_price"`
}

func NewTruckWithTripPrice(origin, destination trips.Region, truck *Truck) *TruckWithTripPrice {
	distance := trips.GetDistance(origin, destination)
	tripPrice := (distance * truck.PerRegionPrice) + truck.BasePrice

	return &TruckWithTripPrice{
		Truck:     truck,
		TripPrice: tripPrice,
	}
}
