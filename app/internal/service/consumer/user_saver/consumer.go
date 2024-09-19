package usersaver

import (
	"context"

	"github.com/Prrromanssss/platform_common/pkg/kafka"

	"github.com/Prrromanssss/auth/config"
	"github.com/Prrromanssss/auth/internal/repository"
)

type service struct {
	cfg            *config.Config
	userRepository repository.UserRepository
	consumer       kafka.Consumer
}

func NewService(
	cfg *config.Config,
	userRepository repository.UserRepository,
	consumer kafka.Consumer,
) *service {
	return &service{
		cfg:            cfg,
		userRepository: userRepository,
		consumer:       consumer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, s.cfg.KafkaConsumer.UsersCreationTopicName, s.UserSaveHandler)
	}()

	return errChan
}
