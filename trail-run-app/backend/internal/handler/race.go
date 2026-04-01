package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trail-run-app/internal/model"
	"trail-run-app/pkg/amap"
	"trail-run-app/pkg/utils"
)

// ListRaces returns a paginated list of races
func ListRaces(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var races []model.Race
	var total int64

	query := model.DB.Model(&model.Race{}).Where("status = ?", model.RaceStatusPublished)

	// Apply filters
	if province := c.Query("province"); province != "" {
		query = query.Where("province = ?", province)
	}
	if city := c.Query("city"); city != "" {
		query = query.Where("city = ?", city)
	}
	if difficulty := c.Query("difficulty"); difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		query = query.Where("date >= ?", dateFrom)
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		query = query.Where("date <= ?", dateTo)
	}

	// Get total count
	query.Count(&total)

	// Get paginated results
	if err := query.Offset(offset).Limit(pageSize).Order("date desc").Find(&races).Error; err != nil {
		utils.RespondInternalError(c, "Failed to fetch races")
		return
	}

	// Convert to list items
	items := make([]*model.RaceListItem, len(races))
	for i := range races {
		items[i] = races[i].ToListItem()
	}

	utils.RespondSuccess(c, gin.H{
		"races":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetRace returns race details with aid stations and equipment
func GetRace(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondValidationError(c, "Invalid race ID")
		return
	}

	var race model.Race
	if err := model.DB.Preload("AidStations").First(&race, "id = ?", id).Error; err != nil {
		utils.RespondNotFound(c, "Race not found")
		return
	}

	utils.RespondSuccess(c, race)
}

// RacePlanResponse represents the race day plan
type RacePlanResponse struct {
	Race struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		Date         time.Time `json:"date"`
		DistanceKM   float64   `json:"distance_km"`
		ElevationM   int       `json:"elevation_m"`
		StartLat     float64   `json:"start_lat"`
		StartLng     float64   `json:"start_lng"`
	} `json:"race"`
	Commute struct {
		StartPoint struct {
			Lat  float64 `json:"lat"`
			Lng  float64 `json:"lng"`
			Name string  `json:"name"`
		} `json:"start_point"`
		Parking []struct {
			Name      string  `json:"name"`
			DistanceM int     `json:"distance_m"`
			Lat       float64 `json:"lat"`
			Lng       float64 `json:"lng"`
		} `json:"parking"`
	} `json:"commute"`
	AidStations interface{} `json:"aid_stations"`
	Weather      interface{} `json:"weather"`
}

// GetRacePlan returns race day planning data
func GetRacePlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondValidationError(c, "Invalid race ID")
		return
	}

	var race model.Race
	if err := model.DB.Preload("AidStations").First(&race, "id = ?", id).Error; err != nil {
		utils.RespondNotFound(c, "Race not found")
		return
	}

	plan := RacePlanResponse{}
	plan.Race.ID = race.ID
	plan.Race.Name = race.Name
	plan.Race.Date = race.Date
	plan.Race.DistanceKM = race.DistanceKM
	plan.Race.ElevationM = race.ElevationM
	plan.Race.StartLat = race.StartLat
	plan.Race.StartLng = race.StartLng
	plan.Commute.StartPoint.Lat = race.StartLat
	plan.Commute.StartPoint.Lng = race.StartLng
	plan.Commute.StartPoint.Name = "起点"

	// Get parking info from Amap API
	amapClient := amap.NewClient()
	pois, err := amapClient.GetNearbyPOIs(race.StartLat, race.StartLng, "停车场")
	if err == nil && len(pois) > 0 {
		for _, poi := range pois {
			var lat, lng float64
			fmt.Sscanf(poi.Location, "%f,%f", &lng, &lat)
			var dist int
			fmt.Sscanf(poi.Distance, "%d", &dist)
			plan.Commute.Parking = append(plan.Commute.Parking, struct {
				Name      string  `json:"name"`
				DistanceM int     `json:"distance_m"`
				Lat       float64 `json:"lat"`
				Lng       float64 `json:"lng"`
			}{
				Name:      poi.Name,
				DistanceM: dist,
				Lat:       lat,
				Lng:       lng,
			})
		}
	}

	// Convert aid stations
	aidStations := make([]*model.AidStationResponse, len(race.AidStations))
	for i := range race.AidStations {
		aidStations[i], _ = race.AidStations[i].ToResponse()
	}
	plan.AidStations = aidStations

	// Get weather from Amap API
	if race.WeatherCityCode != nil && *race.WeatherCityCode != "" {
		weather, err := amapClient.GetWeather(*race.WeatherCityCode)
		if err == nil {
			var temp float64
			fmt.Sscanf(weather.Temperature, "%f", &temp)
			plan.Weather = map[string]interface{}{
				"date":        time.Now().Format("2006-01-02"),
				"weather":     weather.Weather,
				"temperature": weather.Temperature,
				"wind":        weather.WindDirection + weather.WindPower,
				"humidity":    weather.Humidity + "%",
			}
		}
	}

	utils.RespondSuccess(c, plan)
}

// GetEquipment returns equipment checklist for a race
func GetEquipment(c *gin.Context) {
	idStr := c.Param("id")
	raceID, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondValidationError(c, "Invalid race ID")
		return
	}

	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var equipment []model.Equipment
	if err := model.DB.Where("race_id = ?", raceID).Find(&equipment).Error; err != nil {
		utils.RespondInternalError(c, "Failed to fetch equipment")
		return
	}

	// Get user's check states
	var checks []model.EquipmentCheck
	model.DB.Where("user_id = ? AND race_id = ?", userID, raceID).Find(&checks)

	checkMap := make(map[uuid.UUID]bool)
	for _, check := range checks {
		checkMap[check.EquipmentID] = check.Checked
	}

	// Separate mandatory and recommended
	type EquipmentWithCheck struct {
		model.EquipmentResponse
		Checked bool `json:"checked"`
	}

	var mandatory []EquipmentWithCheck
	var recommended []EquipmentWithCheck

	for _, e := range equipment {
		eq := EquipmentWithCheck{
			EquipmentResponse: *e.ToResponse(),
			Checked:            checkMap[e.ID],
		}
		if e.IsMandatory {
			mandatory = append(mandatory, eq)
		} else {
			recommended = append(recommended, eq)
		}
	}

	utils.RespondSuccess(c, gin.H{
		"mandatory":   mandatory,
		"recommended": recommended,
	})
}

// CheckEquipment updates equipment check state
func CheckEquipment(c *gin.Context) {
	idStr := c.Param("id")
	raceID, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondValidationError(c, "Invalid race ID")
		return
	}

	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req struct {
		EquipmentID string `json:"equipment_id" binding:"required"`
		Checked     bool   `json:"checked"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	equipmentID, err := uuid.Parse(req.EquipmentID)
	if err != nil {
		utils.RespondValidationError(c, "Invalid equipment ID")
		return
	}

	// Upsert check state
	check := model.EquipmentCheck{
		UserID:      userID,
		RaceID:     raceID,
		EquipmentID: equipmentID,
		Checked:     req.Checked,
	}

	result := model.DB.Where("user_id = ? AND race_id = ? AND equipment_id = ?",
		userID, raceID, equipmentID).Assign(check).FirstOrCreate(&check)

	if result.Error != nil {
		utils.RespondInternalError(c, "Failed to update check state")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"equipment_id": equipmentID.String(),
		"checked":     req.Checked,
	})
}
