# Data Model: 越野跑App MVP

**Date**: 2026-03-30
**Feature**: 越野跑App MVP

---

## 1. Core Entities

### User (用户)

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | 用户唯一标识 |
| nickname | VARCHAR(50) | NOT NULL | 昵称 |
| email | VARCHAR(255) | UNIQUE, NOT NULL | 邮箱（登录用） |
| password_hash | VARCHAR(255) | NOT NULL | 密码哈希 |
| avatar_url | VARCHAR(500) | NULLABLE | 头像URL |
| phone | VARCHAR(20) | NULLABLE | 手机号 |
| itra_account | VARCHAR(100) | NULLABLE | ITRA账户 |
| created_at | TIMESTAMP | DEFAULT NOW | 创建时间 |
| updated_at | TIMESTAMP | AUTO UPDATE | 更新时间 |

### Race (赛事)

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | 赛事唯一标识 |
| name | VARCHAR(255) | NOT NULL | 赛事名称 |
| date | DATE | NOT NULL | 比赛日期 |
| location | VARCHAR(255) | NOT NULL | 地点 |
| province | VARCHAR(50) | NOT NULL | 省份 |
| city | VARCHAR(50) | NOT NULL | 城市 |
| distance_km | DECIMAL(6,2) | NOT NULL | 距离（公里） |
| elevation_m | INT | NOT NULL | 总爬升（米） |
| difficulty | ENUM('Easy','Medium','Hard','Extrem') | NOT NULL | 难度 |
| itra_points | DECIMAL(4,1) | NULLABLE | ITRA积分 |
| start_lat | DECIMAL(10,7) | NOT NULL | 起点纬度 |
| start_lng | DECIMAL(10,7) | NOT NULL | 起点经度 |
| end_lat | DECIMAL(10,7) | NOT NULL | 终点纬度 |
| end_lng | DECIMAL(10,7) | NOT NULL | 终点经度 |
| route_gpx_url | VARCHAR(500) | NULLABLE | 赛道GPX文件URL |
| weather_city_code | VARCHAR(20) | NULLABLE | 天气预报城市码 |
| status | ENUM('draft','published') | DEFAULT draft | 状态 |
| created_at | TIMESTAMP | DEFAULT NOW | 创建时间 |

### AidStation (补给站)

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | 补给站唯一标识 |
| race_id | UUID | FK → Race | 所属赛事 |
| name | VARCHAR(100) | NOT NULL | 站名 |
| distance_km | DECIMAL(6,2) | NOT NULL | 距起点距离 |
| elevation_m | INT | NOT NULL | 海拔 |
| lat | DECIMAL(10,7) | NOT NULL | 纬度 |
| lng | DECIMAL(10,7) | NOT NULL | 经度 |
| supplies | JSON | NOT NULL | 补给类型列表 |
| close_time | TIME | NOT NULL | 关门时间 |

**supplies JSON示例**: `["水", "运动饮料", "能量胶", "香蕉", "医疗服务"]`

### Equipment (装备清单项)

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | 装备唯一标识 |
| race_id | UUID | FK → Race | 所属赛事 |
| name | VARCHAR(100) | NOT NULL | 装备名称 |
| is_mandatory | BOOLEAN | DEFAULT false | 是否强制装备 |
| category | ENUM('穿着','装备','补给','安全','其他') | NOT NULL | 分类 |

### Result (完赛记录)

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | 记录唯一标识 |
| user_id | UUID | FK → User | 用户 |
| race_id | UUID | FK → Race | 赛事 |
| finish_time | TIME | NULLABLE | 完赛时间 |
| ranking | INT | NULLABLE | 总排名 |
| ranking_age_group | INT | NULLABLE | 年龄组排名 |
| gpx_url | VARCHAR(500) | NULLABLE | GPX轨迹文件URL |
| photos | JSON | NULLABLE | 照片URL列表 |
| generated_image_url | VARCHAR(500) | NULLABLE | 生成的轨迹图片URL |
| status | ENUM('registered','ongoing','finished') | DEFAULT registered | 状态 |
| created_at | TIMESTAMP | DEFAULT NOW | 创建时间 |

### EquipmentCheck (装备勾选状态)

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | 唯一标识 |
| user_id | UUID | FK → User | 用户 |
| race_id | UUID | FK → Race | 赛事 |
| equipment_id | UUID | FK → Equipment | 装备 |
| checked | BOOLEAN | DEFAULT false | 是否已勾选 |
| UNIQUE | | (user_id, race_id, equipment_id) | 组合唯一约束 |

### FavoriteRace (收藏赛事)

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | 唯一标识 |
| user_id | UUID | FK → User | 用户 |
| race_id | UUID | FK → Race | 赛事 |
| reminder_set | BOOLEAN | DEFAULT false | 是否设置提醒 |
| created_at | TIMESTAMP | DEFAULT NOW | 收藏时间 |
| UNIQUE | | (user_id, race_id) | 组合唯一约束 |

---

## 2. Relationships

```
User
├── 1:N → Result (用户的完赛记录)
├── 1:N → EquipmentCheck (用户的装备勾选)
└── 1:N → FavoriteRace (用户收藏的赛事)

Race
├── 1:N → AidStation (赛事的补给站)
├── 1:N → Equipment (赛事的装备清单)
├── 1:N → Result (赛事的完赛记录)
└── 1:N → FavoriteRace (被收藏)

Result
└── N:1 → User
└── N:1 → Race
```

---

## 3. Indexes

```sql
-- 高频查询优化
CREATE INDEX idx_race_date ON races(date);
CREATE INDEX idx_race_location ON races(province, city);
CREATE INDEX idx_result_user ON results(user_id);
CREATE INDEX idx_result_race ON results(race_id);
CREATE INDEX idx_favorite_user ON favorite_races(user_id);
CREATE INDEX idx_aid_station_race ON aid_stations(race_id);
```

---

## 4. Validation Rules

| Entity | Field | Rule |
|--------|-------|------|
| Race | distance_km | > 0 |
| Race | elevation_m | >= 0 |
| Race | date | >= TODAY |
| AidStation | distance_km | > 0 (within race distance) |
| AidStation | close_time | Valid TIME format |
| Equipment | name | NOT empty |
| Result | finish_time | > 0 (if status = finished) |

---

## 5. State Transitions

### Result Status

```
registered → ongoing → finished
     ↓
   (用户开始比赛) (用户完赛)
```

### Race Status

```
draft → published
  (审核通过后发布)
```
