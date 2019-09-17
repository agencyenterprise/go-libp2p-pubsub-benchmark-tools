package config

// Config is a struct to hold the config options
type Config struct {
	Host Host
}

// Host contains configs for the host
type Host struct {
	// Listen are addresses on which to listen
	Listen []string
}
