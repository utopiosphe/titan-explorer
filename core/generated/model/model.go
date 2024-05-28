package model

import "time"

type Language string

const (
	LanguageEN = "en"
	LanguageCN = "cn"
)

type UserRole int

const (
	UserRoleDefault UserRole = iota
	UserRoleAdmin
	UserRoleKOL
)

var SupportLanguages = []Language{LanguageEN, LanguageCN}

type RewardEvent string

var (
	RewardEventInviteFrens RewardEvent = "invite_frens"
	RewardEventBindDevice  RewardEvent = "bind_device"
	RewardEventEarning     RewardEvent = "earning"
	RewardEventReferrals   RewardEvent = "referrals"
)

type Project struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type StorageStats struct {
	ID                         int64     `db:"id" json:"id"`
	Rank                       int64     `db:"s_rank" json:"s_rank"`
	ProjectId                  int64     `db:"project_id" json:"project_id"`
	ProjectName                string    `db:"project_name" json:"project_name"`
	TotalSize                  int64     `db:"total_size" json:"total_size"`
	UserCount                  int64     `db:"user_count" json:"user_count"`
	ProviderCount              int64     `db:"provider_count" json:"provider_count"`
	Expiration                 time.Time `db:"expiration" json:"expiration"`
	Time                       string    `db:"time" json:"time"`
	StorageChange24H           int64     `db:"storage_change_24h" json:"storage_change_24h"`
	StorageChangePercentage24H float64   `db:"storage_change_percentage_24h" json:"storage_change_percentage_24h"`
	Gas                        float64   `db:"gas" json:"gas"`
	Pledge                     float64   `db:"pledge" json:"pledge"`
	Locations                  string    `db:"locations" json:"locations"`
	CreatedAt                  time.Time `db:"created_at" json:"-"`
	UpdatedAt                  time.Time `db:"updated_at" json:"-"`
}

type StorageSummary struct {
	TotalSize         float64 `db:"total_size" json:"total_size"`
	Projects          int64   `db:"projects" json:"projects"`
	Users             int64   `db:"users" json:"users"`
	Pledges           float64 `db:"pledges" json:"pledges"`
	Gases             float64 `db:"gases" json:"gases"`
	Providers         int64   `db:"providers" json:"providers"`
	RetrievalProvider int64   `db:"retrieval_providers" json:"retrieval_providers"`
	StorageProvider   int64   `db:"storage_providers" json:"storage_providers"`
	LatestUpdateTime  string  `db:"-" json:"latest_update_time"`
}

type StorageProvider struct {
	ID          int64     `db:"id" json:"id"`
	ProviderID  string    `db:"provider_id" json:"provider_id"`
	IP          string    `db:"ip" json:"ip"`
	Location    string    `db:"location" json:"location"`
	Retrievable bool      `db:"retrievable" json:"retrievable"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

type InviteFrensRecord struct {
	Email      string    `db:"email" json:"email"`
	Status     int       `db:"status" json:"status"`
	BoundCount int       `db:"bound_count" json:"bound_count"`
	Reward     float64   `db:"reward" json:"reward"`
	Referrer   string    `db:"referrer" json:"referrer"`
	Time       time.Time `db:"time" json:"time"`
}

type SignInfo struct {
	MinerID      string `json:"miner_id" db:"miner_id"`
	Address      string `json:"address" db:"address"`
	Date         int64  `json:"date" db:"date"`
	SignedMsg    string `json:"signed_msg" db:"signed_msg"`
	MinerPower   string `json:"miner_power" db:"miner_power"`
	MinerBalance string `json:"miner_balance" db:"miner_balance"`
}

type DeviceDistribution struct {
	Country string `json:"country" db:"country"`
	Count   int    `json:"count" db:"count"`
}

type AppVersion struct {
	ID          int64     `db:"id" json:"-"`
	Version     string    `db:"version" json:"version"`
	MinVersion  string    `db:"min_version" json:"min_version"`
	Description string    `db:"description" json:"description"`
	Url         string    `db:"url" json:"url"`
	Cid         string    `db:"cid" json:"cid"`
	Size        int64     `db:"size" json:"size"`
	Platform    string    `db:"platform" json:"platform"`
	Lang        string    `db:"lang" json:"lang"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

type KOLLevelConf struct {
	ID                      int64     `db:"id" json:"-"`
	Level                   int       `json:"level" db:"level"`
	ParentCommissionPercent int       `db:"parent_commission_percent" json:"parent_commission_percent"`
	ChildrenBonusPercent    int       `db:"children_bonus_percent" json:"children_bonus_percent"`
	Status                  int       `db:"status" json:"status"`
	UserThreshold           int       `db:"user_threshold" json:"user_threshold"`
	DeviceThreshold         int       `db:"device_threshold" json:"device_threshold"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time `db:"updated_at" json:"updated_at"`
}

type KOL struct {
	ID        int64     `db:"id" json:"-"`
	UserId    string    `json:"user_id" db:"user_id"`
	Level     int       `json:"level" db:"level"`
	Comment   string    `json:"comment" db:"comment"`
	Status    int       `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type KOLLevel struct {
	UserId                  string `db:"user_id"`
	Level                   int    `json:"level" db:"level"`
	ParentCommissionPercent int    `db:"parent_commission_percent" json:"parent_commission_percent"`
	ChildrenBonusPercent    int    `db:"children_bonus_percent" json:"children_bonus_percent"`
	DeviceThreshold         int64  `db:"device_threshold" json:"device_threshold"`
}

type ReferralCode struct {
	ID        int64     `db:"id" json:"-"`
	UserId    string    `json:"user_id" db:"user_id"`
	Code      string    `json:"code" db:"code"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type ReferralCodeProfile struct {
	Code                string    `db:"code" json:"code"`
	ReferralUsers       int       `db:"referral_users" json:"referral_users"`
	ReferralNodes       int       `db:"referral_nodes" json:"referral_nodes"`
	ReferralOnlineNodes int       `db:"referral_online_nodes" json:"referral_online_nodes"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
}

type KolLevelUpInfo struct {
	CurrenLevel          int `json:"curren_level"`
	CommissionPercent    int `json:"commission_percent"`
	BonusPercent         int `json:"bonus_percent"`
	ReferralUsers        int `json:"referral_users"`
	ReferralNodes        int `json:"referral_nodes"`
	LevelUpReferralUsers int `json:"level_up_referral_users"`
	LevelUpReferralNodes int `json:"level_up_referral_nodes"`
}

type ReferralRewardDaily struct {
	ReferrerUserId string    `json:"referrer_user_id" db:"referrer_user_id"`
	UserId         string    `json:"user_id" db:"user_id"`
	OnlineCount    int64     `json:"online_count" db:"online_count"`
	ReferrerReward float64   `json:"referrer_reward" db:"referrer_reward"`
	RefereeReward  float64   `json:"referee_reward" db:"referee_reward"`
	Time           time.Time `db:"time" json:"time"`
}

type DataCollectionEvent int

const (
	DataCollectionEventReferralCodePV = iota + 1
)

type DataCollection struct {
	Event     DataCollectionEvent `json:"event" db:"event"`
	Url       string              `json:"url" db:"url"`
	Os        string              `json:"os" db:"os"`
	Value     string              `json:"value" db:"value"`
	IP        string              `json:"ip" db:"ip"`
	CreatedAt time.Time           `json:"created_at" db:"created_at"`
}

type DateValue struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type UserRewardDaily struct {
	UserId            string    `json:"user_id" db:"user_id"`
	CumulativeReward  float64   `json:"cumulative_reward" db:"cumulative_reward"`
	Reward            float64   `json:"reward" db:"reward"`
	AppReward         float64   `json:"app_reward" db:"app_reward"`
	CliReward         float64   `json:"cli_reward" db:"cli_reward"`
	KOLBonus          float64   `json:"kol_bonus" db:"kol_bonus"`
	ReferralReward    float64   `json:"referral_reward" db:"referral_reward"`
	ReferrerUserId    string    `json:"referrer_user_id" db:"referrer_user_id"`
	DeviceOnlineCount int64     `json:"device_online_count" db:"device_online_count"`
	TotalDeviceCount  int64     `json:"total_device_count" db:"total_device_count"`
	IsKOL             int64     `json:"is_kol" db:"is_kol"`
	IsReferrerKOL     int64     `json:"is_referrer_kol" db:"is_referrer_kol"`
	ReferrerReward    float64   `json:"referrer_reward" db:"referrer_reward"`
	CommissionPercent int64     `json:"commission_percent" db:"commission_percent"`
	KOLBonusPercent   int64     `json:"kol_bonus_percent" db:"kol_bonus_percent"`
	Time              time.Time `db:"time" json:"time"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"-"`
}

type UserReferralRecord struct {
	UserId            string    `json:"user_id" db:"user_id"`
	ReferrerUserId    string    `json:"referrer_user_id" db:"referrer_user_id"`
	DeviceOnlineCount int64     `json:"device_online_count" db:"device_online_count"`
	Reward            float64   `json:"reward" db:"reward"`
	ReferrerReward    float64   `json:"referrer_reward" db:"referrer_reward"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type ReferralCounter struct {
	ReferralUsers  int64   `json:"referral_users" db:"referral_users"`
	ReferralNodes  int64   `json:"referral_nodes" db:"referral_nodes"`
	ReferrerReward float64 `json:"referrer_reward" db:"referrer_reward"`
	RefereeReward  float64 `json:"referee_reward" db:"referee_reward"`
}
