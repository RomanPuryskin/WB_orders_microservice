package postgres

type PostgresConfig struct {
	Host     string `env:"DB_HOST" envDefault:"db"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Name     string `env:"DB_NAME,required"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}
