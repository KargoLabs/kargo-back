package models

type EventRoute string

const (
	// EventRouteNatural is natural trip route
	EventRouteNatural EventRoute = "natural"
	// EventRouteDenial is denial trip route
	EventRouteDenial EventRoute = "denial"
	// EventRouteDenial is cancellation trip route
	EventRouteCancellation EventRoute = "cancellation"
	// EventRouteReport is report in the trip route
	EventRouteReport EventRoute = "report"
)

var (
	// ValidEventRoutes represent all the valid event routes
	ValidEventRoutes = map[EventRoute]bool{
		EventRouteNatural:      true,
		EventRouteDenial:       true,
		EventRouteCancellation: true,
		EventRouteReport:       true,
	}
)
