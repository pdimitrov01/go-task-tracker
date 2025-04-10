package config

type Config struct {
	ApiPort         string
	DbConnectionUrl string
}

func FromEnv() (*Config, error) {
	var conf Config
	cb := newBuilder()

	conf.ApiPort = cb.getString("API_PORT")
	conf.DbConnectionUrl = cb.getString("DB_CONNECTION_URL")

	cbError := cb.getError()
	if cbError != nil {
		return nil, cbError
	}

	return &conf, nil
}
