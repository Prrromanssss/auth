package user

import (
	"context"
	"strconv"

	cacheClient "github.com/Prrromanssss/platform_common/pkg/cache"
	redigo "github.com/gomodule/redigo/redis"

	"github.com/Prrromanssss/auth/internal/cache"
	"github.com/Prrromanssss/auth/internal/cache/user/converter"
	modelCache "github.com/Prrromanssss/auth/internal/cache/user/model"
	"github.com/Prrromanssss/auth/internal/model"
)

type userRedis struct {
	cache cacheClient.RedisClient
}

func NewCache(cache cacheClient.RedisClient) cache.UserCache {
	return &userRedis{
		cache: cache,
	}
}

func (c *userRedis) Create(ctx context.Context, params model.User) (err error) {
	paramsCache := converter.ConvertUserFromServiceToCache(params)
	userIDString := strconv.FormatInt(paramsCache.UserID, 10)

	err = c.cache.HashSet(ctx, userIDString, paramsCache)
	if err != nil {
		return err
	}

	return nil
}

func (c *userRedis) Get(ctx context.Context, params model.GetUserParams) (resp model.GetUserResponse, err error) {
	paramsCache := converter.ConvertGetUserParamsFromServiceToCache(params)

	userIDString := strconv.FormatInt(paramsCache.UserID, 10)

	exists, err := c.cache.Exists(ctx, userIDString)
	if err != nil {
		return
	}

	if !exists {
		err = modelCache.ErrUserNotFound
		return
	}

	values, err := c.cache.HGetAll(ctx, userIDString)
	if err != nil {
		return
	}

	if len(values) == 0 {
		err = modelCache.ErrUserNotFound
		return
	}

	var user modelCache.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return
	}

	return converter.ConvertGetUserResponseFromCacheToService(user), nil
}

func (c *userRedis) Delete(ctx context.Context, params model.DeleteUserParams) (err error) {
	paramsCache := converter.ConvertDeleteUserParamsFromServiceToCache(params)
	userIDString := strconv.FormatInt(paramsCache.UserID, 10)

	err = c.cache.Del(ctx, userIDString)
	if err != nil {
		return err
	}

	return nil
}
