// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package model

import (
	"encoding/json"
	ers "errors"
	"time"
)

type Application struct {
	ID                int64     `db:"id" json:"id"`
	UserID            string    `db:"user_id" json:"user_id"`
	Email             string    `db:"email" json:"email"`
	Num               int32     `db:"num" json:"num"`
	AreaID            string    `db:"area_id" json:"area_id"`
	IpCountry         string    `db:"ip_country" json:"ip_country"`
	IpCity            string    `db:"ip_city" json:"ip_city"`
	PublicKey         string    `db:"public_key" json:"public_key"`
	NodeType          int32     `db:"node_type" json:"node_type"`
	Amount            int32     `db:"amount" json:"amount"`
	UpstreamBandwidth float64   `db:"upstream_bandwidth" json:"upstream_bandwidth"`
	DiskSpace         float64   `db:"disk_space" json:"disk_space"`
	Ip                string    `db:"ip" json:"ip"`
	Status            int32     `db:"status" json:"status"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type ApplicationResult struct {
	ID            int64     `db:"id" json:"id"`
	ApplicationID int64     `db:"application_id" json:"application_id"`
	UserID        string    `db:"user_id" json:"user_id"`
	DeviceID      string    `db:"device_id" json:"device_id"`
	NodeType      int32     `db:"node_type" json:"node_type"`
	Secret        string    `db:"secret" json:"secret"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type CacheEvent struct {
	ID           int64     `db:"id" json:"id"`
	DeviceID     string    `db:"device_id" json:"device_id"`
	CarfileCid   string    `db:"carfile_cid" json:"carfile_cid"`
	BlockSize    float64   `db:"block_size" json:"block_size"`
	Blocks       int64     `db:"blocks" json:"blocks"`
	Time         time.Time `db:"time" json:"time"`
	Status       int32     `db:"status" json:"status"`
	ReplicaInfos int32     `db:"replicaInfos" json:"replicaInfos"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
type DeviceDynamicInfo struct {
	DeviceID         string    `db:"device_id" json:"device_id"`
	UpdateTime       time.Time `db:"update_time" json:"update_time"`
	DiskUsage        float64   `db:"disk_usage" json:"disk_usage"`
	DeviceStatus     string    `db:"device_status" json:"device_status"`
	OnlineTime       float64   `db:"online_time" json:"online_time"`
	CumulativeProfit float64   `db:"cumulative_profit" json:"cumulative_profit"`
	BandwidthUp      float64   `db:"bandwidth_up" json:"bandwidth_up"`
	BandwidthDown    float64   `db:"bandwidth_down" json:"bandwidth_down"`
	DownloadTraffic  float64   `db:"download_traffic" json:"download_traffic"`
	UploadTraffic    float64   `db:"upload_traffic" json:"upload_traffic"`
	CacheCount       int64     `db:"cache_count" json:"cache_count"`
	RetrievalCount   int64     `db:"retrieval_count" json:"retrieval_count"`
}

type DeviceStaticInfo struct {
	DeviceID string `db:"device_id" json:"device_id"`
	//首次加入titan的时间
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	DeviceName    string    `db:"device_name" json:"device_name"`
	NodeType      int64     `db:"node_type" json:"node_type"`
	DiskSpace     float64   `db:"disk_space" json:"disk_space"`
	SystemVersion string    `db:"system_version" json:"system_version"`
	ExternalIp    string    `db:"external_ip" json:"external_ip"`
	InternalIp    string    `db:"internal_ip" json:"internal_ip"`
	MacLocation   string    `db:"mac_location" json:"mac_location"`
	IoSystem      string    `db:"io_system" json:"io_system"`
}

type DeviceInfo struct {
	DeviceID         string    `db:"device_id" json:"device_id"`
	DiskUsage        float64   `db:"disk_usage" json:"disk_usage"`
	DeviceStatus     string    `db:"device_status" json:"device_status"`
	DeviceStatusCode int64     `db:"device_status_code" json:"device_status_code"`
	OnlineTime       float64   `db:"online_time" json:"online_time"`
	CumulativeProfit float64   `db:"cumulative_profit" json:"cumulative_profit"`
	BandwidthUp      float64   `db:"bandwidth_up" json:"bandwidth_up"`
	BandwidthDown    float64   `db:"bandwidth_down" json:"bandwidth_down"`
	DownloadTraffic  float64   `db:"download_traffic" json:"download_traffic"`
	UploadTraffic    float64   `db:"upload_traffic" json:"upload_traffic"`
	CacheCount       int64     `db:"cache_count" json:"cache_count"`
	RetrievalCount   int64     `db:"retrieval_count" json:"retrieval_count"`
	NATType          string    `db:"nat_type" json:"nat_type"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt        time.Time `db:"deleted_at" json:"deleted_at"`
	BoundAt          time.Time `db:"bound_at" json:"bound_at"`
	UserID           string    `db:"user_id" json:"-"`
	LastSeen         time.Time `db:"last_seen" json:"last_seen"`
	IsMobile         int64     `db:"is_mobile" json:"is_mobile"`

	NodeType        int64   `db:"node_type" json:"node_type"`
	DeviceRank      int64   `db:"device_rank" json:"device_rank"`
	DeviceName      string  `db:"device_name" json:"device_name"`
	SystemVersion   string  `db:"system_version" json:"system_version"`
	NetworkInfo     string  `db:"network_info" json:"network_info"`
	ExternalIp      string  `db:"external_ip" json:"external_ip"`
	InternalIp      string  `db:"internal_ip" json:"internal_ip"`
	IpLocation      string  `db:"ip_location" json:"ip_location"`
	IpCountry       string  `db:"ip_country" json:"ip_country"`
	IpProvince      string  `db:"ip_province" json:"ip_province"`
	IpCity          string  `db:"ip_city" json:"ip_city"`
	Latitude        float64 `db:"latitude" json:"latitude"`
	Longitude       float64 `db:"longitude" json:"longitude"`
	MacLocation     string  `db:"mac_location" json:"mac_location"`
	CpuUsage        float64 `db:"cpu_usage" json:"cpu_usage"`
	CpuCores        int64   `db:"cpu_cores" json:"cpu_cores"`
	CpuInfo         string  `db:"cpu_info" json:"cpu_info"`
	GpuInfo  		string  `db:"gpu_info" json:"gpu_info"`
	MemoryUsage     float64 `db:"memory_usage" json:"memory_usage"`
	Memory          float64 `db:"memory" json:"memory"`
	DiskSpace       float64 `db:"disk_space" json:"disk_space"`
	BindStatus      string  `db:"bind_status" json:"bind_status"`
	ActiveStatus    int64   `db:"active_status" json:"active_status"`
	DiskType        string  `db:"disk_type" json:"disk_type"`
	IoSystem        string  `db:"io_system" json:"io_system"`
	TodayOnlineTime float64 `db:"today_online_time" json:"today_online_time"`
	TodayProfit     float64 `db:"today_profit" json:"today_profit"`
	YesterdayProfit float64 `db:"yesterday_profit" json:"yesterday_profit"`
	SevenDaysProfit float64 `db:"seven_days_profit" json:"seven_days_profit"`
	MonthProfit     float64 `db:"month_profit" json:"month_profit"`
	AvailableProfit float64 `db:"available_profit" json:"available_profit"`
	DeactivateTime  int64   `db:"deactivate_time" json:"deactivate_time"`
	IncomeIncr      float64 `db:"income_incr" json:"income_incr"`
	AreaID          string  `db:"area_id" json:"area_id"`
	TitanDiskSpace  float64 `db:"titan_disk_space" json:"titan_disk_space"`
	TitanDiskUsage  float64 `db:"titan_disk_usage" json:"titan_disk_usage"`

	Location
}

type NodesInfo struct {
	Rank        string  `db:"rank" json:"rank"`
	NodeType    string  `db:"node_type" json:"node_type"`
	UserId      string  `db:"user_id" json:"user_id"`
	NodeCount   int64   `db:"node_count" json:"node_count"`
	DiskSpace   float64 `db:"disk_space" json:"disk_space"`
	BandwidthUp float64 `db:"bandwidth_up" json:"bandwidth_up"`
}

type DeviceInfoDaily struct {
	ID                int64     `db:"id" json:"id"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt         time.Time `db:"deleted_at" json:"deleted_at"`
	UserID            string    `db:"user_id" json:"user_id"`
	DeviceID          string    `db:"device_id" json:"device_id"`
	Time              time.Time `db:"time" json:"time"`
	Income            float64   `db:"income" json:"income"`
	OnlineTime        float64   `db:"online_time" json:"online_time"`
	PkgLossRatio      float64   `db:"pkg_loss_ratio" json:"pkg_loss_ratio"`
	Latency           float64   `db:"latency" json:"latency"`
	NatRatio          float64   `db:"nat_ratio" json:"nat_ratio"`
	DiskUsage         float64   `db:"disk_usage" json:"disk_usage"`
	DiskSpace         float64   `db:"disk_space" json:"disk_space"`
	BandwidthUp       float64   `db:"bandwidth_up" json:"bandwidth_up"`
	BandwidthDown     float64   `db:"bandwidth_down" json:"bandwidth_down"`
	UpstreamTraffic   float64   `db:"upstream_traffic" json:"upstream_traffic"`
	DownstreamTraffic float64   `db:"downstream_traffic" json:"downstream_traffic"`
	RetrievalCount    int64     `db:"retrieval_count" json:"retrieval_count"`
	BlockCount        int64     `db:"block_count" json:"block_count"`
}

type DeviceInfoHour struct {
	ID                int64     `db:"id" json:"id"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt         time.Time `db:"deleted_at" json:"deleted_at"`
	UserID            string    `db:"user_id" json:"user_id"`
	DeviceID          string    `db:"device_id" json:"device_id"`
	Time              time.Time `db:"time" json:"time"`
	HourIncome        float64   `db:"hour_income" json:"hour_income"`
	OnlineTime        float64   `db:"online_time" json:"online_time"`
	PkgLossRatio      float64   `db:"pkg_loss_ratio" json:"pkg_loss_ratio"`
	Latency           float64   `db:"latency" json:"latency"`
	NatRatio          float64   `db:"nat_ratio" json:"nat_ratio"`
	DiskUsage         float64   `db:"disk_usage" json:"disk_usage"`
	DiskSpace         float64   `db:"disk_space" json:"disk_space"`
	BandwidthUp       float64   `db:"bandwidth_up" json:"bandwidth_up"`
	BandwidthDown     float64   `db:"bandwidth_down" json:"bandwidth_down"`
	UpstreamTraffic   float64   `db:"upstream_traffic" json:"upstream_traffic"`
	DownstreamTraffic float64   `db:"downstream_traffic" json:"downstream_traffic"`
	RetrievalCount    int64     `db:"retrieval_count" json:"retrieval_count"`
	BlockCount        int64     `db:"block_count" json:"block_count"`
}

type FullNodeInfo struct {
	ID                       int64     `db:"id" json:"id"`
	Date                     string    `db:"date" json:"date" `
	TotalNodeCount           int32     `db:"total_node_count" json:"total_node_count"`
	OnlineNodeCount          int32     `db:"online_node_count" json:"online_node_count"`
	TNodeOnlineRatio         float64   `db:"t_node_online_ratio" json:"t_node_online_ratio"`
	TUpstreamFileCount       int64     `db:"t_upstream_file_count" json:"t_upstream_file_count"`
	TAverageReplica          float64   `db:"t_average_replica" json:"t_average_replica"`
	FBackupsFromTitan        float64   `db:"f_backups_from_titan" json:"f_backups_from_titan"`
	ValidatorCount           int32     `db:"validator_count" json:"validator_count"`
	OnlineValidateorCount    int32     `db:"online_validator_count" json:"online_validator_count"`
	CandidateCount           int32     `db:"candidate_count" json:"candidate_count"`
	OnlineCandidateCount     int32     `db:"online_candidate_count" json:"online_candidate_count"`
	EdgeCount                int32     `db:"edge_count" json:"edge_count"`
	OnlineEdgeCount          int32     `db:"online_edge_count" json:"online_edge_count"`
	TotalStorage             float64   `db:"total_storage" json:"total_storage"`
	StorageUsed              float64   `db:"storage_used" json:"storage_used"`
	TitanDiskSpace           float64   `db:"titan_disk_space" json:"titan_disk_space"`
	TitanDiskUsage           float64   `db:"titan_disk_usage" json:"titan_disk_usage"`
	StorageLeft              float64   `db:"storage_left" json:"storage_left"`
	TotalUpstreamBandwidth   float64   `db:"total_upstream_bandwidth" json:"total_upstream_bandwidth"`
	TotalDownstreamBandwidth float64   `db:"total_downstream_bandwidth" json:"total_downstream_bandwidth"`
	TotalCarfile             int64     `db:"total_carfile" json:"total_carfile"`
	TotalCarfileSize         float64   `db:"total_carfile_size" json:"total_carfile_size"`
	RetrievalCount           int64     `db:"retrieval_count" json:"retrieval_count"`
	NextElectionTime         time.Time `db:"next_election_time" json:"next_election_time"`
	FVMOrderCount            int64     `db:"fvm_order_count" json:"fvm_order_count"`
	FNodeCount               int64     `db:"f_node_count" json:"f_node_count"`
	FHigh                    int64     `db:"f_high" json:"f_high"`
	TNextElectionHigh        int64     `db:"t_next_election_high" json:"t_next_election_high"`
	Time                     time.Time `db:"time" json:"time"`
	CPUCores                 int64     `db:"cpu_cores" json:"cpu_cores"`
	Memory                   int64     `db:"memory" json:"memory"`
	IPCount                  int64     `db:"ip_count" json:"ip_count"`
	CreatedAt                time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                time.Time `db:"updated_at" json:"updated_at"`
}

type LoginLog struct {
	ID            int64     `db:"id" json:"id"`
	LoginUsername string    `db:"login_username" json:"login_username"`
	IpAddress     string    `db:"ip_address" json:"ip_address"`
	LoginLocation string    `db:"login_location" json:"login_location"`
	Browser       string    `db:"browser" json:"browser"`
	Os            string    `db:"os" json:"os"`
	Status        int32     `db:"status" json:"status"`
	Msg           string    `db:"msg" json:"msg"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type OperationLog struct {
	ID               int64     `db:"id" json:"id"`
	Title            string    `db:"title" json:"title"`
	BusinessType     int32     `db:"business_type" json:"business_type"`
	Method           string    `db:"method" json:"method"`
	RequestMethod    string    `db:"request_method" json:"request_method"`
	OperatorType     int32     `db:"operator_type" json:"operator_type"`
	OperatorUsername string    `db:"operator_username" json:"operator_username"`
	OperatorUrl      string    `db:"operator_url" json:"operator_url"`
	OperatorIp       string    `db:"operator_ip" json:"operator_ip"`
	OperatorLocation string    `db:"operator_location" json:"operator_location"`
	OperatorParam    string    `db:"operator_param" json:"operator_param"`
	JsonResult       string    `db:"json_result" json:"json_result"`
	Status           int32     `db:"status" json:"status"`
	ErrorMsg         string    `db:"error_msg" json:"error_msg"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type RetrievalEvent struct {
	DeviceID   string    `db:"device_id" json:"device_id"`
	ClientID   string    `db:"client_id" json:"client_id"`
	CarfileCid string    `db:"carfile_cid" json:"carfile_cid"`
	BlockSize  float64   `db:"block_size" json:"block_size"`
	Time       time.Time `db:"time" json:"time"`
	StartTime  int64     `db:"start_time" json:"start_time"`
	EndTime    int64     `db:"end_time" json:"end_time"`
	Expiration time.Time `db:"expiration" json:"expiration"`
	Status     int32     `db:"status" json:"status"`
}

type Scheduler struct {
	ID        int64     `db:"id" json:"id"`
	Uuid      string    `db:"uuid" json:"uuid"`
	Area      string    `db:"area" json:"area"`
	Address   string    `db:"address" json:"address"`
	Status    int32     `db:"status" json:"status"`
	Token     string    `db:"token" json:"token"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt time.Time `db:"deleted_at" json:"deleted_at"`
}

type SystemInfo struct {
	ID               int64     `db:"id" json:"id"`
	SchedulerUuid    string    `db:"scheduler_uuid" json:"scheduler_uuid"`
	CarFileCount     int64     `db:"car_file_count" json:"car_file_count"`
	DownloadCount    int64     `db:"download_count" json:"download_count"`
	NextElectionTime time.Time `db:"next_election_time" json:"next_election_time"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type User struct {
	ID                     int64     `db:"id" json:"id"`
	Uuid                   string    `db:"uuid" json:"uuid"`
	Avatar                 string    `db:"avatar" json:"avatar"`
	Username               string    `db:"username" json:"username"`
	PassHash               string    `db:"pass_hash" json:"-"`
	UserEmail              string    `db:"user_email" json:"user_email"`
	WalletAddress          string    `db:"wallet_address" json:"wallet_address"`
	Role                   int32     `db:"role" json:"role"`
	AllocateStorage        int       `db:"allocate_storage" json:"allocate_storage"`
	ProjectId              int64     `db:"project_id"`
	Referrer               string    `db:"referrer" json:"referrer"`
	ReferrerUserId         string    `db:"referrer_user_id" json:"-"`
	ReferralCode           string    `db:"referral_code" json:"referral_code"`
	Reward                 float64   `db:"reward" json:"reward"`
	RefereralReward        float64   `db:"referral_reward" json:"referral_reward"`
	Payout                 float64   `db:"payout" json:"payout"`
	FrozenReward           float64   `db:"frozen_reward" json:"frozen_reward"`
	ClosedTestReward       float64   `db:"closed_test_reward" json:"closed_test_reward"`
	HuygensReward          float64   `db:"huygens_reward" json:"huygens_reward"`
	HuygensReferralReward  float64   `db:"huygens_referral_reward" json:"huygens_referral_reward"`
	HerschelReward         float64   `db:"herschel_reward" json:"herschel_reward"`
	HerschelReferralReward float64   `db:"herschel_referral_reward" json:"herschel_referral_reward"`
	DeviceCount            int64     `db:"device_count" json:"device_count"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"-"`
	DeletedAt              time.Time `db:"deleted_at" json:"-"`
}

type Link struct {
	ID        int64     `db:"id" json:"id"`
	UserName  string    `db:"username" json:"username"`
	Cid       string    `db:"cid" json:"cid"`
	LongLink  string    `db:"long_link" json:"long_link"`
	ShortLink string    `db:"short_link" json:"short_link"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt time.Time `db:"deleted_at" json:"deleted_at"`
}

type ValidationEvent struct {
	ID              int64     `db:"id" json:"id"`
	DeviceID        string    `db:"device_id" json:"device_id"`
	ValidatorID     string    `db:"validator_id" json:"validator_id"`
	Blocks          int64     `db:"blocks" json:"blocks"`
	Status          int32     `db:"status" json:"status"`
	Time            time.Time `db:"time" json:"time"`
	Duration        int64     `db:"duration" json:"duration"`
	UpstreamTraffic float64   `db:"upstream_traffic" json:"upstream_traffic"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type Location struct {
	ID          int64     `db:"id" json:"id"`
	Ip          string    `db:"ip" json:"ip"`
	Continent   string    `db:"continent" json:"continent"`
	Province    string    `db:"province" json:"province"`
	City        string    `db:"city" json:"city"`
	Country     string    `db:"country" json:"country"`
	Latitude    string    `db:"latitude" json:"latitude"`
	Longitude   string    `db:"longitude" json:"longitude"`
	AreaCode    string    `db:"area_code" json:"area_code"`
	Isp         string    `db:"isp" json:"isp"`
	ZipCode     string    `db:"zip_code" json:"zip_code"`
	Elevation   string    `db:"elevation" json:"elevation"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	CountryCode string    `db:"-" json:"country_code"`
}

//type LocationEn struct {
//	ID        int64     `db:"id" json:"id"`
//	Ip        string    `db:"ip" json:"ip"`
//	Continent string    `db:"continent" json:"continent"`
//	Province  string    `db:"province" json:"province"`
//	City      string    `db:"city" json:"city"`
//	Country   string    `db:"country" json:"country"`
//	Latitude  string    `db:"latitude" json:"latitude"`
//	Longitude string    `db:"longitude" json:"longitude"`
//	AreaCode  string    `db:"area_code" json:"area_code"`
//	Isp       string    `db:"isp" json:"isp"`
//	ZipCode   string    `db:"zip_code" json:"zip_code"`
//	Elevation string    `db:"elevation" json:"elevation"`
//	CreatedAt time.Time `db:"created_at" json:"created_at"`
//	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
//}

type UserInfo struct {
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	Date           string    `json:"date" db:"date"`
	UserId         string    `db:"user_id"`
	TotalSize      int64     `db:"total_storage_size"`
	UsedSize       int64     `db:"used_storage_size"`
	TotalBandwidth int64     `db:"total_bandwidth"`
	PeakBandwidth  int64     `db:"peak_bandwidth"`
	DownloadCount  int64     `db:"download_count"`
	Time           time.Time `db:"time"`
}

type LotusRequest struct {
	Jsonrpc string     `json:"jsonrpc"`
	Method  string     `json:"method"`
	Params  rawMessage `json:"params"`
	ID      int        `json:"id"`
}

type rawMessage []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m rawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *rawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return ers.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

// Response defines a JSON RPC response from the spec
// http://www.jsonrpc.org/specification#response_object
type LotusResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  interface{}     `json:"result,omitempty"`
	ID      interface{}     `json:"id"`
	Error   *LotusRespError `json:"error,omitempty"`
}

type LotusRespError struct {
	Code    errorCode       `json:"code"`
	Message string          `json:"message"`
	Meta    json.RawMessage `json:"meta,omitempty"`
}
type errorCode int

type Asset struct {
	ID         int64     `db:"id" json:"id"`
	NodeID     string    `db:"node_id" json:"node_id"`
	Event      int64     `db:"event" json:"event"`
	Cid        string    `db:"cid" json:"cid"`
	Hash       string    `db:"hash" json:"hash"`
	TotalSize  int64     `db:"total_size" json:"total_size"`
	Path       string    `db:"path" json:"path"`
	EndTime    time.Time `db:"end_time" json:"end_time"`
	Expiration time.Time `db:"expiration" json:"expiration"`
	UserId     string    `db:"user_id" json:"user_id"`
	Type       string    `db:"type" json:"type"`
	Name       string    `db:"name" json:"name"`
	ProjectId  int64     `db:"project_id" json:"project_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt  time.Time `db:"deleted_at" json:"deleted_at"`
}

type FilStorage struct {
	ID          int64     `db:"id" json:"id"`
	Provider    string    `db:"provider" json:"provider"`
	SectorNum   string    `db:"sector_num" json:"sector_num"`
	IP          string    `db:"ip" json:"ip"`
	Location    string    `db:"location" json:"location"`
	Cost        float64   `db:"cost" json:"cost"`
	MessageCid  string    `db:"message_cid" json:"message_cid"`
	PieceCid    string    `db:"piece_cid" json:"piece_cid"`
	PayloadCid  string    `db:"payload_cid" json:"payload_cid"`
	DealID      string    `db:"deal_id" json:"deal_id"`
	Path        string    `db:"path" json:"path"`
	FIndex      int64     `db:"f_index" json:"f_index"`
	PieceSize   float64   `db:"piece_size" json:"piece_size"`
	Gas         float64   `db:"gas" json:"gas"`
	Pledge      float64   `db:"pledge" json:"pledge"`
	StartHeight int64     `db:"start_height" json:"start_height"`
	EndHeight   int64     `db:"end_height" json:"end_height"`
	StartTime   time.Time `db:"start_time" json:"start_time"`
	EndTime     time.Time `db:"end_time" json:"end_time"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt   time.Time `db:"deleted_at" json:"deleted_at"`
}

type UserSecret struct {
	ID        int64     `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	AppKey    string    `db:"app_key" json:"app_key"`
	AppSecret string    `db:"app_secret" json:"app_secret"`
	Status    int32     `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt time.Time `db:"deleted_at" json:"deleted_at"`
}

type RewardStatement struct {
	ID        int64       `db:"id" json:"id"`
	FromUser  string      `db:"from_user" json:"from_user"`
	Username  string      `db:"username" json:"username"`
	Amount    int64       `db:"amount" json:"amount"`
	Event     RewardEvent `db:"event" json:"event"`
	Status    int32       `db:"status" json:"status"`
	DeviceId  string      `db:"device_id" json:"device_id"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
}

type Withdraw struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Amount    int64     `db:"amount" json:"amount"`
	ToAddress string    `db:"to_address" json:"to_address"`
	Hash      string    `db:"hash" json:"hash"'`
	Status    int32     `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Subscription struct {
	ID                      int64     `db:"id" json:"id"`
	Company                 string    `db:"company" json:"company"`
	Name                    string    `db:"name" json:"name"`
	Email                   string    `db:"email" json:"email"`
	Telegram                string    `db:"telegram" json:"mail"`
	Wechat                  string    `db:"wechat" json:"wechat"`
	Location                string    `db:"location" json:"location"`
	Storage                 string    `db:"storage" json:"storage"`
	Calculation             string    `db:"calculation" json:"calculation"`
	IdleResourcePercentages string    `db:"idle_resource_percentages" json:"idle_resource_percentages"`
	Bandwidth               string    `db:"bandwidth" json:"bandwidth"`
	JoinTestnet             int       `db:"join_testnet" json:"join_testnet,string"`
	Subscribe               int       `db:"subscribe" json:"subscribe,string"`
	Source                  string    `db:"source" json:"source"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time `db:"updated_at" json:"updated_at"`
}

type Signature struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	NodeId    string    `db:"node_id" json:"node_id"`
	AreaId    string    `db:"area_id" json:"area_id"`
	Message   string    `db:"message" json:"message"`
	Hash      string    `db:"hash" json:"hash"`
	Signature string    `db:"signature" json:"signature"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
