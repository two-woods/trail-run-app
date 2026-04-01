package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trail-run-app/internal/model"
	"trail-run-app/pkg/amap"
	"trail-run-app/pkg/utils"
)

// SearchHotels 搜索附近酒店（供赛程规划使用）
func SearchHotels(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	if latStr == "" || lngStr == "" {
		utils.RespondValidationError(c, "lat and lng are required")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		utils.RespondValidationError(c, "invalid lat")
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		utils.RespondValidationError(c, "invalid lng")
		return
	}

	// 默认搜索半径3公里
	radius := 3000
	radiusStr := c.Query("radius")
	if radiusStr != "" {
		if r, err := strconv.Atoi(radiusStr); err == nil && r > 0 && r <= 10000 {
			radius = r
		}
	}

	amapClient := amap.NewClient()
	pois, err := amapClient.SearchHotels(lat, lng, radius)
	if err != nil {
		utils.LogWarn("Failed to search hotels from API: %v", err)
		// 即使API失败也继续，返回空列表
	}

	var hotels []*model.HotelResponse
	for _, poi := range pois {
		var lat, lng float64
		if _, err := fmt.Sscanf(poi.Location, "%f,%f", &lng, &lat); err != nil {
			continue
		}

		distance := 0
		if poi.Distance != "" {
			if d, err := strconv.Atoi(poi.Distance); err == nil {
				distance = d
			}
		}

		hotel := &model.Hotel{
			ID:               uuid.New(),
			Name:             poi.Name,
			Lat:              lat,
			Lng:              lng,
			DistanceFromStart: distance,
			Address:          poi.Address,
			Phone:            poi.Tel,
			BusinessArea:     poi.Business,
			CityName:         poi.CityName,
		}
		hotels = append(hotels, hotel.ToResponse())
	}

	utils.RespondSuccess(c, gin.H{
		"hotels": hotels,
	})
}

// ListFavoriteHotels 获取用户收藏的酒店
func ListFavoriteHotels(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	raceIDStr := c.Query("race_id")
	if raceIDStr == "" {
		utils.RespondValidationError(c, "race_id is required")
		return
	}

	raceID, err := uuid.Parse(raceIDStr)
	if err != nil {
		utils.RespondValidationError(c, "invalid race_id")
		return
	}

	var favorites []model.FavoriteHotel
	if err := model.DB.Where("user_id = ? AND race_id = ?", userID, raceID).Find(&favorites).Error; err != nil {
		utils.RespondInternalError(c, "Failed to fetch favorite hotels")
		return
	}

	hotels := make([]*model.HotelResponse, len(favorites))
	for i := range favorites {
		hotels[i] = favorites[i].ToResponse()
	}

	utils.RespondSuccess(c, gin.H{
		"hotels": hotels,
	})
}

// CreateFavoriteHotel 添加收藏酒店
type CreateFavoriteHotelRequest struct {
	RaceID     string `json:"race_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Lat        float64 `json:"lat" binding:"required"`
	Lng        float64 `json:"lng" binding:"required"`
	Distance   int     `json:"distance"`
	Address    string  `json:"address"`
	Phone      string  `json:"phone"`
	BookingURL string  `json:"booking_url"`
	Notes      string  `json:"notes"`
}

func CreateFavoriteHotel(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req CreateFavoriteHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	raceID, err := uuid.Parse(req.RaceID)
	if err != nil {
		utils.RespondValidationError(c, "invalid race_id")
		return
	}

	favorite := model.FavoriteHotel{
		UserID:     userID,
		RaceID:     raceID,
		Name:       req.Name,
		Lat:        req.Lat,
		Lng:        req.Lng,
		Distance:   req.Distance,
		Address:    req.Address,
		Phone:      req.Phone,
		BookingURL: req.BookingURL,
		Notes:      req.Notes,
	}

	if err := model.DB.Create(&favorite).Error; err != nil {
		utils.RespondInternalError(c, "Failed to create favorite hotel")
		return
	}

	utils.RespondCreated(c, favorite.ToResponse())
}

// DeleteFavoriteHotel 删除收藏酒店
func DeleteFavoriteHotel(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondValidationError(c, "invalid id")
		return
	}

	result := model.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.FavoriteHotel{})
	if result.RowsAffected == 0 {
		utils.RespondNotFound(c, "Favorite hotel not found")
		return
	}

	utils.RespondSuccess(c, gin.H{"deleted": true})
}

// GetAllHotelsForRace 获取赛事的所有酒店（API + 收藏）
func GetAllHotelsForRace(c *gin.Context) {
	raceIDStr := c.Query("race_id")
	if raceIDStr == "" {
		utils.RespondValidationError(c, "race_id is required")
		return
	}

	raceID, err := uuid.Parse(raceIDStr)
	if err != nil {
		utils.RespondValidationError(c, "invalid race_id")
		return
	}

	// 获取赛事信息获取起点坐标
	var race model.Race
	if err := model.DB.First(&race, "id = ?", raceID).Error; err != nil {
		utils.RespondNotFound(c, "Race not found")
		return
	}

	// 1. 从高德API获取附近酒店
	amapClient := amap.NewClient()
	var apiHotels []*model.HotelResponse
	pois, err := amapClient.SearchHotels(race.StartLat, race.StartLng, 3000)
	if err == nil {
		for _, poi := range pois {
			var lat, lng float64
			if _, err := fmt.Sscanf(poi.Location, "%f,%f", &lng, &lat); err != nil {
				continue
			}
			distance := 0
			if poi.Distance != "" {
				if d, err := strconv.Atoi(poi.Distance); err == nil {
					distance = d
				}
			}
			hotel := &model.HotelResponse{
				ID:               uuid.New().String(),
				Name:             poi.Name,
				Lat:              lat,
				Lng:              lng,
				DistanceFromStart: distance,
				Address:          poi.Address,
				Phone:            poi.Tel,
				Source:           "api",
			}
			apiHotels = append(apiHotels, hotel)
		}
	}

	// 2. 获取用户收藏的酒店
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var favoriteHotels []*model.HotelResponse
	var favorites []model.FavoriteHotel
	if err := model.DB.Where("user_id = ? AND race_id = ?", userID, raceID).Find(&favorites).Error; err == nil {
		for i := range favorites {
			favoriteHotels = append(favoriteHotels, favorites[i].ToResponse())
		}
	}

	utils.RespondSuccess(c, gin.H{
		"api_hotels":     apiHotels,
		"favorite_hotels": favoriteHotels,
	})
}