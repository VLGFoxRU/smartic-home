package engine

import (
    "context"
    "log"
	"time"

    "github.com/VLGFoxRU/smartic-home/internal/service"
)

type SceneEngine struct {
    sceneService   *service.SceneService
    controlService *service.ControlService
}

func NewSceneEngine(sceneService *service.SceneService, controlService *service.ControlService) *SceneEngine {
    return &SceneEngine{
        sceneService:   sceneService,
        controlService: controlService,
    }
}

// ProcessTelemetry вызывается при каждом новом показании
func (e *SceneEngine) ProcessTelemetry(ctx context.Context, deviceID, telemetryType string, value float64, timestamp time.Time) {
    // Пока упрощённо: обрабатываем только temperature, ищем сценарии для всех домов.
    // В реальности нужно определить дом по deviceID, но у нас пока нет метода GetHomeByDevice.
    // Используем фиксированный дом (b0000000-...)
    homeID := "b0000000-0000-0000-0000-000000000001"

    // Получаем все активные сценарии для дома
    scenes, err := e.sceneService.ListScenes(ctx, homeID)
    if err != nil {
        log.Printf("SceneEngine: ошибка получения сценариев: %v", err)
        return
    }

    for _, scene := range scenes {
        if !scene.IsActive() {
            continue
        }
        // Проверяем условие
        if e.checkCondition(scene.Condition(), deviceID, telemetryType, value) {
            // Выполняем действие
            action := scene.Action()
            cmd, _ := action["command"].(string)
            params, _ := action["params"].(map[string]interface{})
            targetDeviceID, _ := action["device_id"].(string)

            if cmd == "" || targetDeviceID == "" {
                log.Printf("SceneEngine: некорректное действие в сценарии %s", scene.ID())
                continue
            }

            // Вызываем ControlService (роль system, userID – системный)
            err := e.controlService.SendCommand(ctx, targetDeviceID, cmd, params, "0000000-0000-0000-0000-000000000000", "admin")
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