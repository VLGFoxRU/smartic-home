package engine

import (
    "context"
    "log"
    "time"

    "github.com/VLGFoxRU/smartic-home/internal/repository"
    "github.com/VLGFoxRU/smartic-home/internal/service"
)

type SceneEngine struct {
    sceneService   *service.SceneService
    controlService *service.ControlService
    cache          repository.CacheRepository
}

func NewSceneEngine(
    sceneService *service.SceneService,
    controlService *service.ControlService,
    cache repository.CacheRepository,
) *SceneEngine {
    return &SceneEngine{
        sceneService:   sceneService,
        controlService: controlService,
        cache:          cache,
    }
}

func (e *SceneEngine) ProcessTelemetry(ctx context.Context, deviceID, telemetryType string, value float64, timestamp time.Time) {
    homeID := "b0000000-0000-0000-0000-000000000001" // фиксированный дом

    scenes, err := e.sceneService.ListScenes(ctx, homeID)
    if err != nil {
        log.Printf("SceneEngine: ошибка получения сценариев: %v", err)
        return
    }

    for _, scene := range scenes {
        if !scene.IsActive() {
            continue
        }
        if e.checkCondition(scene.Condition(), deviceID, telemetryType, value) {
            // Проверяем cooldown
            cooldownSec := 300 // по умолчанию 5 минут
            if cd, ok := scene.Condition()["cooldown_sec"].(float64); ok {
                cooldownSec = int(cd)
            }
            cooldownKey := "scene:" + scene.ID() + ":last_triggered"
            acquired, err := e.cache.SetNX(ctx, cooldownKey, "1", time.Duration(cooldownSec)*time.Second)
            if err != nil {
                log.Printf("SceneEngine: ошибка проверки cooldown: %v", err)
                continue
            }
            if !acquired {
                log.Printf("SceneEngine: сценарий %s пропущен (cooldown)", scene.ID())
                continue
            }

            action := scene.Action()
            cmd, _ := action["command"].(string)
            params, _ := action["params"].(map[string]interface{})
            targetDeviceID, _ := action["device_id"].(string)
            if cmd == "" || targetDeviceID == "" {
                log.Printf("SceneEngine: некорректное действие в сценарии %s", scene.ID())
                continue
            }

            err = e.controlService.SendCommand(ctx, targetDeviceID, cmd, params, "c0000000-0000-0000-0000-000000000001", "admin")
            if err != nil {
                log.Printf("SceneEngine: ошибка выполнения команды сценария %s: %v", scene.ID(), err)
            } else {
                log.Printf("SceneEngine: сценарий %s выполнен (команда %s -> %s)", scene.ID(), cmd, targetDeviceID)
            }
        }
    }
}

// checkCondition простая проверка "type = temperature, operator = '>', value"
func (e *SceneEngine) checkCondition(condition map[string]interface{}, deviceID, telemetryType string, value float64) bool {
    condType, _ := condition["type"].(string)
    if condType != telemetryType {
        return false
    }
    condDeviceID, _ := condition["device_id"].(string)
    if condDeviceID != "" && condDeviceID != deviceID {
        return false
    }
    operator, _ := condition["operator"].(string)
    condValue, ok := condition["value"].(float64)
    if !ok {
        return false
    }
    switch operator {
    case ">":
        return value > condValue
    case "<":
        return value < condValue
    case ">=":
        return value >= condValue
    case "<=":
        return value <= condValue
    case "==":
        return value == condValue
    default:
        return false
    }
}