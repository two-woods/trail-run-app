package service

import (
	"trail-run-app/internal/model"
)

// EquipmentRecommender 基于赛事信息推荐装备
type EquipmentRecommender struct{}

type EquipmentRecommendation struct {
	Equipment model.Equipment
	Reason    string
	Priority  int // 1-3, 1为最高
}

// RecommendEquipment 根据距离、爬升、天气推荐装备
func (r *EquipmentRecommender) RecommendEquipment(race *model.Race, weather *WeatherInfo) []EquipmentRecommendation {
	var recommendations []EquipmentRecommendation

	// 基于距离推荐
	if race.DistanceKM > 50 {
		recommendations = append(recommendations, EquipmentRecommendation{
			Equipment: model.Equipment{
				Name:        "备用袜子",
				IsMandatory: false,
				Category:    model.EquipmentCategory装备,
			},
			Reason:   "超过50km，建议备用袜子防止脚起泡",
			Priority: 2,
		})
	}

	if race.DistanceKM > 100 {
		recommendations = append(recommendations, EquipmentRecommendation{
			Equipment: model.Equipment{
				Name:        "头灯",
				IsMandatory: true,
				Category:    model.EquipmentCategory装备,
			},
			Reason:   "超过100km预计会跑夜路",
			Priority: 1,
		})
	}

	// 基于爬升推荐
	if race.ElevationM > 3000 {
		recommendations = append(recommendations, EquipmentRecommendation{
			Equipment: model.Equipment{
				Name:        "登山杖",
				IsMandatory: false,
				Category:    model.EquipmentCategory装备,
			},
			Reason:   "爬升超过3000m，登山杖可节省体力",
			Priority: 2,
		})
	}

	if race.ElevationM > 5000 {
		recommendations = append(recommendations, EquipmentRecommendation{
			Equipment: model.Equipment{
				Name:        "保暖层",
				IsMandatory: true,
				Category:    model.EquipmentCategory穿着,
			},
			Reason:   "高海拔赛段需要保暖",
			Priority: 1,
		})
	}

	// 基于天气推荐
	if weather != nil {
		if weather.Temperature < 10 {
			recommendations = append(recommendations, EquipmentRecommendation{
				Equipment: model.Equipment{
					Name:        "保暖手套",
					IsMandatory: false,
					Category:    model.EquipmentCategory穿着,
				},
				Reason:   "气温低于10度，手部保暖很重要",
				Priority: 2,
			})
		}

		if weather.Temperature < 5 {
			recommendations = append(recommendations, EquipmentRecommendation{
				Equipment: model.Equipment{
					Name:        "压缩保暖衣",
					IsMandatory: true,
					Category:    model.EquipmentCategory穿着,
				},
				Reason:   "气温低于5度，需要压缩保暖层",
				Priority: 1,
			})
		}

		if weather.RainProbability > 30 {
			recommendations = append(recommendations, EquipmentRecommendation{
				Equipment: model.Equipment{
					Name:        "防水外套",
					IsMandatory: true,
					Category:    model.EquipmentCategory装备,
				},
				Reason:   "降雨概率较高，防止失温",
				Priority: 1,
			})
		}

		if weather.Temperature > 25 {
			recommendations = append(recommendations, EquipmentRecommendation{
				Equipment: model.Equipment{
					Name:        "防晒霜",
					IsMandatory: false,
					Category:    model.EquipmentCategory补给,
				},
				Reason:   "气温较高，需注意防晒",
				Priority: 3,
			})
		}
	}

	return recommendations
}

// WeatherInfo 天气信息
type WeatherInfo struct {
	Temperature       float64 // 摄氏度
	RainProbability   int     // 降雨概率 0-100
	WindSpeed         float64 // 风速 m/s
	Weather           string  // 天气描述
}

// CheckMissingEquipment 检查缺失的强制装备
func (r *EquipmentRecommender) CheckMissingEquipment(
	equipment []model.Equipment,
	checks map[string]bool,
) []model.Equipment {
	var missing []model.Equipment

	for _, eq := range equipment {
		if eq.IsMandatory {
			if !checks[eq.ID.String()] {
				missing = append(missing, eq)
			}
		}
	}

	return missing
}
