package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type message struct {
	id          string
	eventType   string
	roomId      string
	parentId    string
	targetType  string
	targetValue float32
}

type workout struct {
	Id       int `json:"id"`
	Duration int `json:"duration"`
	RoomId   int `json:"roomId"`
}

type redisSessionEventConsumer struct {
	config      *models.Config
	simService  core.SimulationService
	logger      stages_logging.StagesLogger
	redisClient redis.UniversalClient
	healthScore float64
}

func NewSessionStreamEventConsumer(
	config *models.Config,
	simService core.SimulationService,
	redisClient redis.UniversalClient,
	logger stages_logging.StagesLogger,
) stages_services.BackgroundService {

	return &redisSessionEventConsumer{
		config:      config,
		simService:  simService,
		logger:      logger,
		redisClient: redisClient,
	}
}

// Health implements stages_services.BackgroundService
func (r *redisSessionEventConsumer) Health() (float64, error) {
	return r.healthScore, nil
}

// Start implements stages_services.BackgroundService
func (r *redisSessionEventConsumer) Start(ctx context.Context) error {
	var (
		err       error
		isExiting bool
	)

	go func() {
		<-ctx.Done()
		isExiting = true
	}()

	for !isExiting {
		err = stages_redis.TestAndWaitForConnection(
			ctx,
			r.redisClient,
			1*time.Second,
			func(n uint, err error) {
				r.healthScore = stages_services.StatusUnhealthy
				r.logger.Warningf("process=\"StreamEventConsumer.Start: TestAndWaitForConnection\" attempt=%d error=\"%s\"", n, err)
			},
		)
		if err != nil {
			return err
		}

		r.logger.Info("process=\"RedisEventConsumer.Start\" message=\"Redis connected\"")
		r.healthScore = stages_services.StatusHealthy

		err = r.doWork(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			r.logger.Errorf("process=\"RedisEventConsumer\" message=\"Redis connection lost\" error=\"%s\"", err)
		}

		r.healthScore = stages_services.StatusUnhealthy
	}

	return nil
}

func (r *redisSessionEventConsumer) doWork(ctx context.Context) error {
	err := stages_redis.CreateConsumerGroupIfNotExists(
		ctx,
		r.redisClient,
		r.config.Redis.Streams.SessionEvents,
		r.config.Redis.EventConsumerGroup,
	)
	if err != nil {
		return errors.Wrap(err, "RedisEventConsumer.Start: CreateConsumerGroupIfNotExists")
	}

	sopt := &stages_redis.GroupConsumeOptions{
		ConsumerName: r.config.Redis.ConsumerName,
		Group:        r.config.Redis.EventConsumerGroup,
		Stream:       r.config.Redis.Streams.SessionEvents,

		OnMessageReceived: func(ctx context.Context, messageId string, e map[string]interface{}) {
			msg := toMessage(e)

			if msg != nil {

				switch msg.eventType {
				case constants.EventKeyPlanSegmentChanged:
					r.targetChangeHandler(ctx, msg)
				case constants.EventKeySessionStarted:
					r.workoutStartedHandler(ctx, msg)
				case constants.EventKeySessionCompleted:
					r.workoutEndedHandler(ctx, msg)

				}
			}
		},
	}

	return stages_redis.GroupConsumeFromStream(ctx, r.redisClient, sopt)

}

func (r *redisSessionEventConsumer) targetChangeHandler(ctx context.Context, msg *message) {

	tsk := r.simService.GetFromRoomId(msg.roomId)
	if tsk != nil && msg.targetType != "" && msg.targetValue > 0 {
		r.logger.LogInfo(constants.TaskIdLabel, tsk.Id, constants.SimulationNameLabel, tsk.Simulation.Name, constants.MessageLabel, fmt.Sprintf("[CONSUMER] Target '%s' changed to %v", msg.targetType, msg.targetValue), constants.IterationsLabel, tsk.Iterations, msg.roomId)
		targets := &models.TaskTargets{
			Primary: &models.TaskTarget{
				Type:  strings.ToLower(msg.targetType),
				Value: msg.targetValue,
			},
		}
		tsk.Simulation.Targets = targets
		tsk.TargetUpdater <- targets
	}
}

func (r *redisSessionEventConsumer) workoutEndedHandler(ctx context.Context, msg *message) {
	sid, err := strconv.Atoi(msg.id)

	if err != nil {
		return
	}

	tsk := r.simService.GetFromSessionId(sid)
	if tsk != nil && tsk.IsRunning {
		r.logger.LogInfo(constants.TaskIdLabel, tsk.Id, constants.SimulationNameLabel, tsk.Simulation.Name, constants.MessageLabel, "Stopping simulation automatically based on session start command", constants.IterationsLabel, tsk.Iterations, "session_id", sid)

		tsk.Stopper <- true
	}
}

func (r *redisSessionEventConsumer) workoutStartedHandler(ctx context.Context, msg *message) {
	sid, err := strconv.Atoi(msg.id)

	if err != nil {
		return
	}

	tsk := r.simService.GetFromSessionId(sid)
	if tsk != nil && !tsk.IsRunning {
		r.logger.LogInfo(constants.TaskIdLabel, tsk.Id, constants.SimulationNameLabel, tsk.Simulation.Name, constants.MessageLabel, "Starting simulation automatically based on session start command", constants.IterationsLabel, tsk.Iterations, "session_id", sid)

		// if msg.duration > 0 {

		// }
		r.simService.Start(tsk)
	}
}

func (r *redisSessionEventConsumer) getWorkoutFromRoom(ctx context.Context, roomId string) *workout {
	var obj workout

	key := fmt.Sprintf(constants.CacheKeyWorkout, roomId)

	exists, _ := r.redisClient.HExists(ctx, key, "data").Result()

	if exists {
		data, _ := r.redisClient.HGet(ctx, key, "data").Result()
		json.Unmarshal([]byte(data), &obj)
	}

	return &obj
}
func toMessage(values map[string]interface{}) *message {
	msg := message{}

	if id, ok := values["session_id"]; ok && id != nil {
		msg.id = id.(string)
	}

	if t, ok := values["event_type"]; ok && t != nil {
		msg.eventType = t.(string)
	}

	if val, ok := values["room_id"]; ok && val != nil {
		msg.roomId = val.(string)
	}
	if val, ok := values["parent_id"]; ok && val != nil {
		msg.parentId = val.(string)
	}

	// if msg.parentId == "" {
	// 	if val, ok := values["group"]; ok && val != nil {
	// 		msg.parentId = val.(string)
	// 	}
	// }

	if targetVal, tvok := values["target_value"]; tvok && targetVal != nil {
		tv, _ := strconv.ParseFloat(targetVal.(string), 32)

		msg.targetValue = float32(tv)
	}

	if targetType, ttok := values["target_type"]; ttok && targetType != nil {
		msg.targetType = targetType.(string)
	}

	return &msg
}
