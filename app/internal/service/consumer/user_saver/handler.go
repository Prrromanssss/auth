package usersaver

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2/log"

	"github.com/Prrromanssss/auth/internal/converter"
	"github.com/Prrromanssss/auth/internal/model"
)

func (s *service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	user := &model.CreateUserKafkaParams{}
	err := json.Unmarshal(msg.Value, user)
	if err != nil {
		return err
	}

	userParams := converter.ConvertCreateUserKafkaParamsFromConsumerServiceToUserService(*user)

	resp, err := s.userRepository.CreateUser(ctx, userParams)
	if err != nil {
		return err
	}

	log.Infof("User with id %d created: %+v\n", resp.UserID, resp)

	return nil
}
