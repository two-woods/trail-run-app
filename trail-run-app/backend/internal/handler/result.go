package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trail-run-app/internal/model"
	"trail-run-app/internal/service"
	"trail-run-app/pkg/utils"
)

// ListResults returns the current user's race results
func ListResults(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var results []model.Result
	if err := model.DB.Preload("Race").Where("user_id = ?", userID).
		Order("created_at desc").Find(&results).Error; err != nil {
		utils.RespondInternalError(c, "Failed to fetch results")
		return
	}

	items := make([]*model.ResultListItem, len(results))
	for i := range results {
		r := &results[i]
		item := &model.ResultListItem{
			ID:                r.ID,
			RaceName:          r.Race.Name,
			RaceDate:          r.Race.Date,
			DistanceKM:        r.Race.DistanceKM,
			ElevationM:        r.Race.ElevationM,
			FinishTime:        r.FinishTime,
			Ranking:           r.Ranking,
			RankingAgeGroup:   r.RankingAgeGroup,
			Status:            r.Status,
			GeneratedImageURL: r.GeneratedImageURL,
		}
		items[i] = item
	}

	utils.RespondSuccess(c, gin.H{
		"results": items,
	})
}

// CreateResultRequest for creating a new result
type CreateResultRequest struct {
	RaceID          string   `json:"race_id" binding:"required"`
	FinishTime      *string  `json:"finish_time"`
	Ranking         *int     `json:"ranking"`
	RankingAgeGroup *int     `json:"ranking_age_group"`
	GPXURL          *string  `json:"gpx_url"`
	Photos          []string `json:"photos"`
}

// CreateResult creates a new race result
func CreateResult(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req CreateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	raceID, err := uuid.Parse(req.RaceID)
	if err != nil {
		utils.RespondValidationError(c, "Invalid race ID")
		return
	}

	// Verify race exists
	var race model.Race
	if err := model.DB.First(&race, "id = ?", raceID).Error; err != nil {
		utils.RespondNotFound(c, "Race not found")
		return
	}

	result := model.Result{
		UserID:   userID,
		RaceID:   raceID,
		Status:   model.ResultStatusRegistered,
	}

	if req.FinishTime != nil {
		result.FinishTime = req.FinishTime
	}
	if req.Ranking != nil {
		result.Ranking = req.Ranking
	}
	if req.RankingAgeGroup != nil {
		result.RankingAgeGroup = req.RankingAgeGroup
	}
	if req.GPXURL != nil {
		result.GPXURL = req.GPXURL
	}
	if len(req.Photos) > 0 {
		result.SetPhotos(req.Photos)
	}

	if result.FinishTime != nil {
		result.Status = model.ResultStatusFinished
	}

	if err := model.DB.Create(&result).Error; err != nil {
		utils.RespondInternalError(c, "Failed to create result")
		return
	}

	utils.RespondCreated(c, gin.H{
		"id":      result.ID,
		"status": result.Status,
	})
}

// GetResult returns a specific result
func GetResult(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondValidationError(c, "Invalid result ID")
		return
	}

	var result model.Result
	if err := model.DB.Preload("Race").First(&result, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		utils.RespondNotFound(c, "Result not found")
		return
	}

	photos, _ := result.GetPhotos()

	detail := model.ResultDetail{
		ID:                result.ID,
		Race:              result.Race,
		FinishTime:        result.FinishTime,
		Ranking:           result.Ranking,
		RankingAgeGroup:   result.RankingAgeGroup,
		GPXURL:            result.GPXURL,
		Photos:            photos,
		GeneratedImageURL: result.GeneratedImageURL,
		CreatedAt:         result.CreatedAt,
	}

	// Parse GPX and generate GeoJSON if GPXURL exists
	if result.GPXURL != nil && *result.GPXURL != "" {
		stats, geojson, err := service.ParseGPXFromURL(*result.GPXURL)
		if err == nil {
			detail.GeoJSON = geojson
			detail.Stats = stats
		}
	}

	utils.RespondSuccess(c, detail)
}

// GenerateImageRequest for triggering image generation
type GenerateImageRequest struct {
	RaceName   string `json:"race_name" json:"race_name"`
	Date       string `json:"date"`
	GPXData    string `json:"gpx_data"`
	DistanceKM float64 `json:"distance_km"`
	ElevationM int     `json:"elevation_m"`
	FinishTime string `json:"finish_time"`
}

// GenerateImage triggers image generation for a result
func GenerateImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondValidationError(c, "Invalid result ID")
		return
	}

	var req GenerateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	var result model.Result
	if err := model.DB.First(&result, "id = ?", id).Error; err != nil {
		utils.RespondNotFound(c, "Result not found")
		return
	}

	// Call Python image service
	imageServiceURL := os.Getenv("IMAGE_SERVICE_URL")
	if imageServiceURL == "" {
		imageServiceURL = "http://localhost:8082"
	}
	imageServiceURL += "/generate"

	payload := map[string]interface{}{
		"race_name":   req.RaceName,
		"date":        req.Date,
		"gpx_data":    req.GPXData,
		"distance_km": req.DistanceKM,
		"elevation_m": req.ElevationM,
		"finish_time": req.FinishTime,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(imageServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != http.StatusOK {
		utils.RespondInternalError(c, "Image service unavailable")
		return
	}
	defer resp.Body.Close()

	var imgResp struct {
		ImageURL  string `json:"image_url"`
		ImageData string `json:"image_data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&imgResp); err != nil {
		utils.RespondInternalError(c, "Failed to parse image response")
		return
	}

	// Save the image URL to the result
	result.GeneratedImageURL = &imgResp.ImageURL
	if err := model.DB.Save(&result).Error; err != nil {
		utils.LogWarn("Failed to save image URL to result: %v", err)
		// Continue anyway - image was generated successfully
	}

	utils.RespondSuccess(c, gin.H{
		"status":     "completed",
		"image_url":  imgResp.ImageURL,
		"image_data": imgResp.ImageData, // Base64 encoded image
	})
}

// Helper function to marshal to JSON
func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
