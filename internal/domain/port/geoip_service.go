package port

type GeoIPService interface {
	GetCountryByIP(ip string) (string, error)
}
