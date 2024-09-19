package yaml

import (
	"strings"

	"github.com/IBM/sarama"
)

// KafkaConsumer represents the configuration for KafkaConsumer.
type KafkaConsumer struct {
	UsersCreationTopicName string `validate:"required" yaml:"user_creation_topic_name"`
	Brokers                string `validate:"required" yaml:"brokers"`
	GroupID                string `validate:"required" yaml:"group_id"`
}

func (k KafkaConsumer) BrokersList() []string {
	return strings.Split(k.Brokers, ",")
}

func (k *KafkaConsumer) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
