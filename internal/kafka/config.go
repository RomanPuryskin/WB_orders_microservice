package kafka

type KafkaConfig struct {
	ExternalPort int    `env:"KAFKA_EXTERNAL_PORT" envDefault:"9092"`
	Topic        string `env:"KAFKA_TOPIC,required"`
	Group        string `env:"KAFKA_GROUP,required"`
	Address      string `env:"KAFKA_ADDRESS" envDefault:"kafka"`
}
