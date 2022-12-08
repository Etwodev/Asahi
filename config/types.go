package config

type Config struct {
	Port         string `json:"port"`
	Address      string `json:"address"`
	Assets		 string	 `json:"assets"`
	Experimental bool	 `json:"experimental"`
}

func Port() string {
	return c.Port
}

func Address() string {
	return c.Address
}

func Experimental() bool {
	return c.Experimental
}

func Assets() string {
	return c.Assets
}