package api

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	constant "github.com/TestsLing/aj-captcha-go/const"

	"github.com/Masterminds/squirrel"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-explorer/config"
	"github.com/gnasnik/titan-explorer/core/dao"
	"github.com/gnasnik/titan-explorer/core/errors"
	"github.com/gnasnik/titan-explorer/core/generated/model"
	"github.com/gnasnik/titan-explorer/pkg/random"
	"github.com/gnasnik/titan-explorer/pkg/rsa"
	"github.com/go-redis/redis/v9"
	"github.com/golang-module/carbon/v2"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type NonceStringType string

const (
	NonceStringTypeRegister  NonceStringType = "1"
	NonceStringTypeLogin     NonceStringType = "2"
	NonceStringTypeReset     NonceStringType = "3"
	NonceStringTypeSignature NonceStringType = "4"
	NonceStringTypeDeactive  NonceStringType = "5"
)

const (
	deactivePre = "user_deactive_"
)

var defaultNonceExpiration = 5 * time.Minute

func GetUserInfoHandler(c *gin.Context) {
	quest := new(dao.UserAndQuest)
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)
	user, err := dao.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.UserNotFound, c))
		return
	}

	codes, err := dao.GetUserReferCodes(c.Request.Context(), username)
	if err != nil {
		log.Errorf("get user referral codes: %v", err)
	}

	if len(codes) > 0 {
		user.ReferralCode = codes[0].Code
	}

	user.CassiniReward = user.Reward + user.OnlineIncentiveReward
	user.CassiniReferralReward = user.ReferralReward

	copier.Copy(quest, user)
	quest.HerschelCredits, _ = dao.GetCreditByUn(c.Request.Context(), username)
	quest.HerschelInviteCredits, _ = dao.GetInviteCreditByUn(c.Request.Context(), username)

	c.JSON(http.StatusOK, respJSON(quest))
}

type registerParams struct {
	Username   string `json:"username"`
	Referrer   string `json:"referrer"`
	VerifyCode string `json:"verify_code"`
	Password   string `json:"password"`
}

func UserRegister(c *gin.Context) {
	var params registerParams
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	userInfo := &model.User{
		Username:  params.Username,
		UserEmail: params.Username,
		Referrer:  params.Referrer,
		CreatedAt: time.Now(),
	}

	verifyCode := params.VerifyCode
	passwd := params.Password
	if userInfo.Username == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	_, err := dao.GetUserByUsername(c.Request.Context(), userInfo.Username)
	if err == nil {
		c.JSON(http.StatusOK, respErrorCode(errors.UserEmailExists, c))
		return
	}

	if err != nil && err != sql.ErrNoRows {
		log.Errorf("GetUserByUsername: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	var referrer *model.User

	if userInfo.Referrer != "" {
		referrer, err = dao.GetUserByRefCode(c.Request.Context(), userInfo.Referrer)
		if err != nil {
			log.Errorf("GetUserByRefCode: %v", err)
			c.JSON(http.StatusOK, respErrorCode(errors.InvalidReferralCode, c))
			return
		}
		userInfo.ReferrerUserId = referrer.Username
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.PassWordNotAllowed, c))
		return
	}
	userInfo.PassHash = string(passHash)

	testCode := os.Getenv("TEST_ENV_VERIFY_CODE")
	if testCode != "" && testCode == verifyCode {
		err = dao.CreateUser(c.Request.Context(), userInfo)
		if err != nil {
			log.Errorf("create user : %v", err)
			c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
			return
		}

		c.JSON(http.StatusOK, respJSON(JsonObject{
			"msg": "success",
		}))
		return
	}

	nonce, err := getNonceFromCache(c.Request.Context(), userInfo.Username, NonceStringTypeRegister)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if nonce == "" || verifyCode == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidVerifyCode, c))
		return
	}

	if nonce != verifyCode {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidVerifyCode, c))
		return
	}

	err = dao.CreateUser(c.Request.Context(), userInfo)
	if err != nil {
		log.Errorf("create user : %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"msg": "success",
	}))
}

type resetParams struct {
	Username   string `json:"username"`
	VerifyCode string `json:"verify_code"`
	Password   string `json:"password"`
}

func PasswordRest(c *gin.Context) {
	//username := c.Query("username")
	//verifyCode := c.Query("verify_code")
	//passwd := c.Query("password")
	var params resetParams
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	_, err := dao.GetUserByUsername(c.Request.Context(), params.Username)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, respErrorCode(errors.NameNotExists, c))
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.PassWordNotAllowed, c))
		return
	}

	nonce, err := getNonceFromCache(c.Request.Context(), params.Username, NonceStringTypeReset)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.Unknown, c))
		return
	}

	if nonce == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.VerifyCodeExpired, c))
		return
	}

	if params.VerifyCode == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidVerifyCode, c))
		return
	}

	if nonce != params.VerifyCode && os.Getenv("TEST_ENV_VERIFY_CODE") != params.VerifyCode {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidVerifyCode, c))
		return
	}

	err = dao.ResetPassword(c.Request.Context(), string(passHash), params.Username)
	if err != nil {
		log.Errorf("update user : %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}
	c.JSON(http.StatusOK, respJSON(JsonObject{
		"msg": "success",
	}))
}

func GetNonceStringHandler(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	nonce, err := generateNonceString(c.Request.Context(), getRedisNonceSignatureKey(username))
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	info, err := dao.GetUserByUsername(c.Request.Context(), username)
	switch err {
	case sql.ErrNoRows:
		user := &model.User{
			Username:         username,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			TotalStorageSize: 100 * 1024 * 1024,
			ReferralCode:     random.GenerateRandomString(6),
		}
		err = dao.CreateUser(c.Request.Context(), user)
		if err != nil {
			c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
			return
		}
	case nil:
		if info.TotalStorageSize == 0 {
			err = dao.UpdateUserTotalSize(c.Request.Context(), info.Username, 100*1024*1024)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	default:
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"code": nonce,
	}))
}

func generateNonceString(ctx context.Context, key string) (string, error) {
	rand := random.GenerateRandomNumber(6)
	verifyCode := "TitanNetWork(" + rand + ")"
	bytes, err := json.Marshal(verifyCode)
	if err != nil {
		return "", err
	}

	_, err = dao.RedisCache.Set(ctx, key, bytes, defaultNonceExpiration).Result()
	if err != nil {
		log.Errorf("%v:", err)
		return "", err
	}

	return verifyCode, nil
}

// GetBlockCaptcha 滑块验证
func GetBlockCaptcha(c *gin.Context) {
	data, _ := factory.GetService(constant.BlockPuzzleCaptcha).Get()
	//输出json结果给调用方
	c.JSON(200, data)
}

type (
	// VerifyCodeReq 获取邮箱验证码
	VerifyCodeReq struct {
		Username  string `json:"username"`
		Token     string `json:"token"`
		PointJSON string `json:"pointJson"`
		Type      int64  `json:"type"`
	}
)

func GetNumericVerifyCodeHandler(c *gin.Context) {
	userInfo := &model.User{}
	req := &VerifyCodeReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}
	userInfo.Username = req.Username
	verifyType := strconv.FormatInt(req.Type, 10)
	lang := c.GetHeader("Lang")
	userInfo.UserEmail = userInfo.Username

	var key string
	switch NonceStringType(verifyType) {
	case NonceStringTypeRegister:
		key = getRedisNonceRegisterKey(userInfo.Username)
	case NonceStringTypeLogin:
		key = getRedisNonceLoginKey(userInfo.Username)
	case NonceStringTypeReset:
		key = getRedisNonceResetKey(userInfo.Username)
	case NonceStringTypeSignature:
		key = getRedisNonceSignatureKey(userInfo.Username)
	case NonceStringTypeDeactive:
		key = getRedisNonceDeactiveKey(userInfo.Username)
	default:
		c.JSON(http.StatusOK, respErrorCode(errors.UnsupportedVerifyCodeType, c))
		return
	}

	nonce, err := getNonceFromCache(c.Request.Context(), userInfo.Username, NonceStringType(verifyType))
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if nonce != "" {
		c.JSON(http.StatusOK, respErrorCode(errors.GetVCFrequently, c))
		return
	}

	verifyCode := random.GenerateRandomNumber(6)

	if err = sendEmail(userInfo.Username, verifyCode, lang); err != nil {
		log.Errorf("send email: %v", err)
		if strings.Contains(err.Error(), "timed out") {
			c.JSON(http.StatusOK, respErrorCode(errors.TimeoutCode, c))
			return
		}
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if err = cacheVerifyCode(c.Request.Context(), key, verifyCode); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"msg": "success",
	}))
}

func DeviceBindingHandler(c *gin.Context) {
	var params model.Signature

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	if params.Signature == "" || params.NodeId == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	sign, err := dao.GetSignatureByHash(c.Request.Context(), params.Hash)
	if err == dao.ErrNoRow {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidSignature, c))
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	deviceInfo, err := dao.GetDeviceInfo(c.Request.Context(), params.NodeId)
	if err == dao.ErrNoRow {
		device, err := getDeviceInfoFromSchedulerAndInsert(c.Request.Context(), params.NodeId, params.AreaId)
		if err != nil {
			c.JSON(http.StatusOK, respErrorCode(errors.DeviceNotExists, c))
			return
		}

		deviceInfo = device
	}

	if deviceInfo == nil {
		c.JSON(http.StatusOK, respErrorCode(errors.DeviceNotExists, c))
		return
	}

	if deviceInfo.UserID != "" {
		c.JSON(http.StatusOK, respErrorCode(errors.DeviceBound, c))
		return
	}

	//if params.AreaId == "" {
	//	params.AreaId = dao.GetAreaID(c.Request.Context(), sign.Username)
	//}

	schedulerClient, err := getSchedulerClient(c.Request.Context(), params.AreaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.NoSchedulerFound, c))
		return
	}

	pubKeyString, err := schedulerClient.GetNodePublicKey(c.Request.Context(), params.NodeId)
	if err != nil {
		log.Errorf("api get node public key: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	pubicKey, err := rsa.Pem2PublicKey([]byte(pubKeyString))
	if err != nil {
		log.Errorf("pem 2 publicKey: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	signature, err := hex.DecodeString(params.Signature)
	if err != nil {
		log.Errorf("hex decode: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidSignature, c))
		return
	}

	err = rsa.VerifySHA256Sign(pubicKey, signature, []byte(params.Hash))
	if err != nil {
		log.Errorf("pem 2 publicKey: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidSignature, c))
		return
	}

	if params.Remark == "" {
		params.Remark = deviceInfo.DeviceName
	}

	if err = dao.UpdateUserDeviceInfo(c.Request.Context(), &model.DeviceInfo{
		UserID:     sign.Username,
		DeviceID:   params.NodeId,
		BindStatus: "binding",
		DeviceName: params.Remark,
	}); err != nil {
		log.Errorf("update device binding status: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if err := dao.UpdateDeviceInfoDailyUser(c.Request.Context(), params.NodeId, sign.Username); err != nil {
		log.Errorf("binding update device info daily: %v", err)
	}

	if sign.Signature == "" {
		err = dao.UpdateSignature(c.Request.Context(), params.Signature, params.NodeId, params.AreaId, params.Hash)
		if err != nil {
			log.Errorf("update signature: %v", err)
			c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
			return
		}
	} else {
		params.Username = sign.Username
		err = dao.AddSignature(c.Request.Context(), &params)
		if err != nil {
			log.Errorf("add signature: %v", err)
			c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
			return
		}
	}

	c.JSON(http.StatusOK, respJSON(nil))

}

func DeviceUnBindingHandlerOld(c *gin.Context) {
	deviceInfo := &model.DeviceInfo{}
	deviceInfo.DeviceID = c.Query("device_id")
	UserID := c.Query("user_id")
	deviceInfo.BindStatus = "unbinding"
	deviceInfo.ActiveStatus = 2

	old, err := dao.GetDeviceInfoByID(c.Request.Context(), deviceInfo.DeviceID)
	if err != nil {
		log.Errorf("get user device: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if old == nil {
		c.JSON(http.StatusOK, respErrorCode(errors.DeviceNotExists, c))
		return
	}

	if old.UserID != UserID {
		c.JSON(http.StatusOK, respErrorCode(errors.UnbindingNotAllowed, c))
		return
	}

	err = dao.UpdateUserDeviceInfo(c.Request.Context(), deviceInfo)
	if err != nil {
		log.Errorf("update user device: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"msg": "success",
	}))
}

func DeviceUpdateHandler(c *gin.Context) {
	deviceId := c.Query("device_id")
	deviceName := c.Query("device_name")

	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	old, err := dao.GetDeviceInfoByID(c.Request.Context(), deviceId)
	if err != nil {
		log.Errorf("get user device: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}
	if old == nil {
		c.JSON(http.StatusOK, respErrorCode(errors.DeviceNotExists, c))
		return
	}

	if old.UserID != username {
		c.JSON(http.StatusOK, respErrorCode(errors.PermissionNotAllowed, c))
		return
	}

	err = dao.UpdateDeviceName(c.Request.Context(), &model.DeviceInfo{DeviceID: deviceId, DeviceName: deviceName})
	if err != nil {
		log.Errorf("update user device: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"msg": "success",
	}))
}

func cacheVerifyCode(ctx context.Context, key, verifyCode string) error {
	bytes, err := json.Marshal(verifyCode)
	if err != nil {
		return err
	}

	_, err = dao.RedisCache.Set(ctx, key, bytes, defaultNonceExpiration).Result()
	if err != nil {
		return err
	}

	return nil
}

func SetPeakBandwidth(userId string) {
	peakBandwidth, err := dao.GetPeakBandwidth(context.Background(), userId)
	if err != nil {
		log.Errorf("get peak bandwidth: %v", err)
		return
	}
	var expireTime time.Duration
	expireTime = time.Hour
	_ = SetUserInfo(context.Background(), userId, peakBandwidth, expireTime)
	return
}

func SetUserInfo(ctx context.Context, key string, peakBandwidth int64, expireTime time.Duration) error {
	bytes, err := json.Marshal(peakBandwidth)
	vc := GetUserInfo(ctx, key)
	if vc != 0 {
		_, err := dao.RedisCache.Expire(ctx, key, expireTime).Result()
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	_, err = dao.RedisCache.Set(ctx, key, bytes, expireTime).Result()
	if err != nil {
		return err
	}
	return nil
}

func getRedisNonceSignatureKey(username string) string {
	return fmt.Sprintf("TITAN::SIGN::%s", username)
}

func getRedisNonceRegisterKey(username string) string {
	return fmt.Sprintf("TITAN::REG::%s", username)
}

func getRedisNonceLoginKey(username string) string {
	return fmt.Sprintf("TITAN::LOGIN::%s", username)
}

func getRedisNonceResetKey(username string) string {
	return fmt.Sprintf("TITAN::RESET::%s", username)
}

func getRedisNonceDeactiveKey(username string) string {
	return fmt.Sprintf("TITAN::DEACTIVE::%s", username)
}

func getNonceFromCache(ctx context.Context, username string, t NonceStringType) (string, error) {
	var key string

	switch t {
	case NonceStringTypeRegister:
		key = getRedisNonceRegisterKey(username)
	case NonceStringTypeLogin:
		key = getRedisNonceLoginKey(username)
	case NonceStringTypeReset:
		key = getRedisNonceResetKey(username)
	case NonceStringTypeSignature:
		key = getRedisNonceSignatureKey(username)
	case NonceStringTypeDeactive:
		key = getRedisNonceDeactiveKey(username)
	default:
		return "", fmt.Errorf("unsupported nonce string type")
	}

	bytes, err := dao.RedisCache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	var verifyCode string
	err = json.Unmarshal(bytes, &verifyCode)
	if err != nil {
		return "", err
	}

	return verifyCode, nil
}

func VerifyMessage(message string, signedMessage string) (string, error) {
	// Hash the unsigned message using EIP-191
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message)
	hash := crypto.Keccak256Hash(hashedMessage)
	// Get the bytes of the signed message
	decodedMessage := hexutil.MustDecode(signedMessage)
	// Handles cases where EIP-115 is not implemented (most wallets don't implement it)
	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
		decodedMessage[64] -= 27
	}
	// Recover a public key from the signed message
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
	if sigPublicKeyECDSA == nil {
		log.Errorf("Could not get a public get from the message signature")
	}
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(*sigPublicKeyECDSA).String(), nil
}

func GetUserInfo(ctx context.Context, key string) int64 {
	bytes, err := dao.RedisCache.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return 0
	}
	if err == redis.Nil {
		return 0
	}
	var peakBandwidth int64
	err = json.Unmarshal(bytes, &peakBandwidth)
	if err != nil {
		return 0
	}
	return peakBandwidth
}

func BindWalletHandler(c *gin.Context) {
	type bindParams struct {
		Username   string `json:"username"`
		VerifyCode string `json:"verify_code"`
		Sign       string `json:"sign"`
		Address    string `json:"address"`
	}

	var param bindParams
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	nonce, err := getNonceFromCache(c.Request.Context(), param.Username, NonceStringTypeSignature)
	if err != nil {
		log.Errorf("query nonce string: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if nonce == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.VerifyCodeExpired, c))
		return
	}

	recoverAddress, err := VerifyMessage(nonce, param.Sign)
	if strings.ToUpper(recoverAddress) != strings.ToUpper(param.Address) {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidSignature, c))
		return
	}

	user, err := dao.GetUserByUsername(c.Request.Context(), param.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusOK, respErrorCode(errors.UserNotFound, c))
		return
	}

	if user.WalletAddress != "" {
		c.JSON(http.StatusOK, respErrorCode(errors.WalletBound, c))
		return
	}

	if err := dao.UpdateUserWalletAddress(context.Background(), param.Username, recoverAddress); err != nil {
		log.Errorf("update user wallet address: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func UnBindWalletHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	ctx := context.Background()
	user, err := dao.GetUserByUsername(ctx, username)
	if err != nil {
		log.Errorf("get user by username: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if user == nil {
		log.Errorf("user not found: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if err := dao.UpdateUserWalletAddress(context.Background(), user.Username, ""); err != nil {
		log.Errorf("update user wallet address: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func maskEmail(email string) string {
	words := strings.Split(email, ".")
	if len(words) <= 1 {
		return email
	}

	prefix, suffix := words[0], words[1]

	if len(prefix) > 5 {
		return prefix[:3] + "****" + prefix[len(prefix)-2:] + "." + suffix
	}

	return prefix[:3] + "****" + "." + suffix
}

func GetReferralListHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	page, _ := strconv.Atoi(c.Query("page"))
	order := c.Query("order")
	orderField := c.Query("order_field")
	option := dao.QueryOption{
		Page:       page,
		PageSize:   pageSize,
		Order:      order,
		OrderField: orderField,
	}

	total, referList, err := dao.GetReferralList(c.Request.Context(), username, option)
	if err != nil {
		log.Errorf("get referral list: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	var userIds []string
	for _, refer := range referList {
		rw, err := dao.GetReferralReward(c.Request.Context(), username, refer.Email)
		if err == nil && rw != nil {
			refer.Reward = rw.Reward
		}
		userIds = append(userIds, refer.Email)
		refer.Email = maskEmail(refer.Email)
	}

	user, err := dao.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		log.Errorf("get user: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.UserNotFound, c))
		return
	}

	referralCodes, err := dao.GetReferralCodeProfileByUserId(c.Request.Context(), username)
	if err != nil {
		log.Errorf("GetUserReferCodes: %v", err)
	}

	kolLevel, err := dao.GetKOLByUserId(c.Request.Context(), username)
	if err != nil {
		log.Errorf("GetKOLByUserId: %v", err)
	}

	var currentLevel int
	if kolLevel != nil {
		currentLevel = kolLevel.Level
	}

	level, err := dao.GetKOLLevelByLevel(c.Request.Context(), currentLevel)
	if err != nil {
		log.Errorf("GetKOLLevelByLevel: %v", err)
	}

	var (
		sumReferralNode int
	)

	for _, item := range referralCodes {
		sumReferralNode += item.EligibleNodes
	}

	levelUpInfo := &model.KolLevelUpInfo{
		CurrenLevel:             currentLevel,
		CommissionPercent:       level.CommissionPercent,
		ParentCommissionPercent: level.ParentCommissionPercent,
		ReferralNodes:           sumReferralNode,
		LevelUpReferralNodes:    level.DeviceThreshold,
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"list":           referList,
		"total":          total,
		"total_reward":   user.ReferralReward,
		"referral_codes": referralCodes,
		"kol_level":      levelUpInfo,
	}))
}

func AddReferralCodeHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	codes, err := dao.GetUserReferCodes(c.Request.Context(), username)
	if err != nil {
		log.Errorf("GetUserReferCodes: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if len(codes) >= 5 {
		c.JSON(http.StatusOK, respErrorCode(errors.ExceedReferralCodeNumbers, c))
		return
	}

	referralCode := &model.ReferralCode{
		UserId:    username,
		Code:      random.GenerateRandomString(6),
		CreatedAt: time.Now(),
	}

	err = dao.AddNewReferralCode(c.Request.Context(), referralCode)
	if err != nil {
		log.Errorf("AddNewReferralCode: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"code": referralCode.Code,
	}))
}

func GetReferralCodeDetailHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)
	code := c.Query("code")
	from := c.Query("from")
	to := c.Query("to")

	option := dao.QueryOption{}

	if from != "" {
		option.StartTime = carbon.Parse(from).StartOfDay().String()
	}

	if to != "" {
		option.EndTime = carbon.Parse(to).EndOfDay().String()
	}

	user, err := dao.GetUserByRefCode(c.Request.Context(), code)
	if err != nil {
		log.Errorf("GetUserByRefCode: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidReferralCode, c))
		return
	}

	if user.Username != username {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidReferralCode, c))
		return
	}

	ipCount, pvCount, err := dao.CountPageViewByEvent(c.Request.Context(), model.DataCollectionEventReferralCodePV, code, option)
	if err != nil {
		log.Errorf("CountPageViewByEvent: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidReferralCode, c))
		return
	}

	referralUsers, referralNodes, err := dao.CountReferralUsersByCode(c.Request.Context(), code, option)
	if err != nil {
		log.Errorf("CountReferralUsersByCode: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidReferralCode, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"code":           code,
		"ip_count":       ipCount,
		"pv_count":       pvCount,
		"referral_users": referralUsers,
		"referral_nodes": referralNodes,
	}))

}

func GetReferralCodeStatHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	t := c.Query("type")
	startTime := c.Query("from")
	endTime := c.Query("to")
	code := c.Query("code")

	user, err := dao.GetUserByRefCode(c.Request.Context(), code)
	if err != nil {
		log.Errorf("GetUserByRefCode: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidReferralCode, c))
		return
	}

	if user.Username != username {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidReferralCode, c))
		return
	}

	option := dao.QueryOption{
		StartTime: startTime,
		EndTime:   endTime,
	}

	if startTime == "" {
		option.StartTime = carbon.Now().SubDays(14).StartOfDay().String()
	} else {
		option.StartTime = carbon.Parse(startTime).StartOfDay().String()
	}

	if endTime == "" {
		option.EndTime = carbon.Now().EndOfDay().String()
	} else {
		option.EndTime = carbon.Parse(endTime).EndOfDay().String()
	}

	var out []*model.DateValue

	switch t {
	case "referral_users":
		out, err = dao.GetUserReferrerUsersDailyStat(c.Request.Context(), code, option)
	case "referral_nodes":
		out, err = dao.GetUserReferrerNodesDailyStat(c.Request.Context(), code, option)
	case "ip_count":
		out, err = dao.GetPageViewIPCountDailyStat(c.Request.Context(), model.DataCollectionEventReferralCodePV, code, option)
	case "pv_count":
		out, err = dao.GetPageViewCountDailyStat(c.Request.Context(), model.DataCollectionEventReferralCodePV, code, option)
	}

	if err != nil {
		log.Errorf("GetReferralCodeStatHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"list": out,
	}))
}

func GetKOLReferralCodeInfoHandler(c *gin.Context) {
	code := c.Query("code")

	user, err := dao.GetUserByRefCode(c.Request.Context(), code)
	if err != nil {
		log.Errorf("GetUserByRefCode: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if user == nil || model.UserRole(user.Role) != model.UserRoleKOL {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidReferralCode, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"code":        code,
		"kol_user_id": user.Username,
	}))
}

func GetBannersHandler(c *gin.Context) {
	lang := c.GetHeader("Lang")
	platform := c.Query("platform")

	if lang != "cn" && lang != "en" {
		c.JSON(http.StatusOK, respErrorCode(errors.AdsLangNotExist, c))
		return
	}

	pfm, _ := strconv.Atoi(platform)
	if pfm != dao.AdsPlatformAPP && pfm != dao.AdsPlatformPC {
		c.JSON(http.StatusOK, respErrorCode(errors.AdsPlatformNotExist, c))
		return
	}

	list, err := dao.ListBannersCtx(c.Request.Context(), int64(pfm), lang)
	if err != nil {
		log.Errorf("GetBannersHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.AdsFetchFailed, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"list": list,
	}))
}

func GetNoticesHandler(c *gin.Context) {
	lang := c.GetHeader("Lang")
	platform := c.Query("platform")

	if lang != "cn" && lang != "en" {
		c.JSON(http.StatusOK, respErrorCode(errors.AdsLangNotExist, c))
		return
	}

	pfm, _ := strconv.Atoi(platform)
	if pfm != dao.AdsPlatformAPP && pfm != dao.AdsPlatformPC {
		c.JSON(http.StatusOK, respErrorCode(errors.AdsPlatformNotExist, c))
		return
	}

	notices, err := dao.ListNoticesCtx(c.Request.Context(), int64(pfm), lang)
	if err != nil {
		log.Errorf("GetNoticesHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.AdsFetchFailed, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"list": notices,
	}))
}

func GetAdsHistoryHandler(c *gin.Context) {
	size, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	platrom := c.Query("platform")
	lang := c.GetHeader("Lang")

	sb := squirrel.Select()

	if platrom != "" {
		sb = sb.Where("platform = ?", platrom)
	}

	if lang != "" {
		sb = sb.Where("lang = ?", lang)
	}

	sb = sb.Where("ads_type = ?", dao.AdsTypeNotice).OrderBy("created_at DESC")

	list, n, err := dao.AdsListPageCtx(c, page, size, sb)
	if err != nil {
		log.Errorf("GetAdsHistoryHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"list":  list,
		"total": n,
	}))
}

func AdsClickIncrHandler(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ads, err := dao.AdsFindOne(c.Request.Context(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "ads not found"})
			return
		} else {
			log.Errorf("Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
	}
	ads.Hits++
	if err := dao.AdsUpdateCtx(c.Request.Context(), ads); err != nil {
		log.Errorf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func BugReportHandler(c *gin.Context) {
	var bug model.Bug
	if err := c.BindJSON(&bug); err != nil {
		log.Errorf("BugReportHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	bug.CreatedAt = time.Now()
	bug.UpdatedAt = time.Now()

	// claims := jwt.ExtractClaims(c)
	// username := claims[identityKey].(string)

	if bug.Code == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}
	info, err := dao.GetSignatureByHash(c.Request.Context(), bug.Code)
	if err != nil {
		log.Errorf("BugReportHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidSignature, c))
		return
	}
	bug.NodeId = info.NodeId

	if info.Username == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.NodeNotBound, c))
		return
	}
	user, err := dao.GetUserByUsername(c.Request.Context(), info.Username)
	if err != nil {
		log.Errorf("BugReportHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.UserNotFound, c))
		return
	}
	bug.Username = info.Username
	bug.Email = user.UserEmail

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	n, err := dao.BugsCountByBuilderCtx(c.Request.Context(), squirrel.Select().Where("node_id=?", info.NodeId).
		Where("username = ?", info.Username).Where("created_at <= ?", endOfDay).Where("created_at >= ?", startOfDay))
	if err != nil {
		log.Errorf("BugReportHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if n >= 5 {
		c.JSON(http.StatusOK, respErrorCode(errors.ReportToManyBugs, c))
		return
	}

	bug.State = dao.BugStateWaiting

	if err := dao.BugsAddCtx(c.Request.Context(), &bug); err != nil {
		log.Errorf("BugReportHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func MyBugReportListHandler(c *gin.Context) {
	size, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))

	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	sb := squirrel.Select().From("bugs").Where("code = ?", code).OrderBy("updated_at DESC")
	state, _ := strconv.Atoi(c.Query("state"))
	if state > 0 {
		sb = sb.Where("state = ?", state)
	}

	list, n, err := dao.BugsListPageCtx(c.Request.Context(), page, size, sb, "id, username, email, node_id, telegram_id, description, feedback_type, feedback, pics, state, reward_type, reward, updated_at")
	if err != nil {
		log.Errorf("BugReportListHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"list":  list,
		"total": n,
	}))
}

func LocatorFromConfigHandler(c *gin.Context) {
	c.JSON(http.StatusOK, respJSON(JsonObject{
		"list": config.Cfg.Locators,
	}))
}

func GetEdgeConfigHandler(c *gin.Context) {
	node_id := c.Query("node_id")
	if node_id == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	config, err := dao.GetEdgeConfig(c.Request.Context(), node_id)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("GetEdgeConfigHandler: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, respErrorCode(errors.NotFound, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"config": config,
	}))
}

func SetEdgeConfigHandler(c *gin.Context) {
	var cfg *model.EdgeConfig

	if err := c.BindJSON(&cfg); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	cfg.CreatedAt = time.Now()
	cfg.UpdatedAt = time.Now()

	if err := dao.SetEdgeConfig(c.Request.Context(), cfg); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"msg": "success",
	}))
}

const (
	BatchUrlZsetKey = "titan:edge:bacth:set"
	BatchUrlHSetKey = "titan:edge:bacth:hset"
)

type BatchSetReq struct {
	LoggedIn bool       `json:"loggedIn"`
	Config   EdgeBatchX `json:"config"`
}
type EdgeBatchX struct {
	PhoneModel string
	OS         string
	Mac        string
	Time       time.Time
	Operation  string
	UrlConfig  map[string]BatchUrlConfig
}

type BatchUrlConfig struct {
	Like   bool
	Follow bool
}

func BatchReportHandler(c *gin.Context) {
	var req BatchSetReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}
	// if req.LoggedIn {
	// 	if _, err := dao.RedisCache.ZRem(c.Request.Context(), BatchUrlZsetKey, req.Config.Mac).Result(); err != nil {
	// 		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, respJSON(JsonObject{"msg": "success"}))
	// 	return
	// }

	req.Config.Time = time.Now()
	data, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	_, err = dao.RedisCache.ZRem(c.Request.Context(), BatchUrlZsetKey, req.Config.Mac).Result()
	if err != nil {
		log.Errorf("Failed to remove %s: %v", req.Config.Mac, err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if _, err := dao.RedisCache.ZAdd(c.Request.Context(), BatchUrlZsetKey, redis.Z{
		Score:  float64(req.Config.Time.Unix()),
		Member: req.Config.Mac,
	}).Result(); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if err := dao.RedisCache.HSet(c.Request.Context(), BatchUrlHSetKey, req.Config.Mac, data).Err(); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{"msg": "success"}))
}

func BatchGetHandler(c *gin.Context) {
	// 获取分页参数
	// size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	// page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	// if size <= 0 {
	// 	size = 10
	// }
	// if page <= 0 {
	// 	page = 1
	// }

	// start := int64((page - 1) * size)
	// stop := int64(page*size - 1)

	// 获取 ZSET 的总元素数量
	// total, err := dao.RedisCache.ZCard(c.Request.Context(), BatchUrlZsetKey).Result()
	// if err != nil {
	// 	c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
	// 	return
	// }

	// results, err := dao.RedisCache.ZRangeWithScores(c.Request.Context(), BatchUrlZsetKey, start, stop).Result()
	// if err != nil {
	// 	c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
	// 	return
	// }
	// 获取所有 zset 中的成员
	members, err := dao.RedisCache.ZRange(c.Request.Context(), BatchUrlZsetKey, 0, -1).Result()
	if err != nil {
		log.Errorf("ZRange error: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}
	// 存储反序列化后的数据
	var batchConfigs []BatchSetReq

	for _, mac := range members {

		data, err := dao.RedisCache.HGet(c.Request.Context(), BatchUrlHSetKey, mac).Result()
		if err != nil {
			log.Errorf("HGet error: %v", err)
			c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
			return
		}

		var config BatchSetReq
		// 反序列化数据
		err = json.Unmarshal([]byte(data), &config)
		if err != nil {
			c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
			return
		}
		batchConfigs = append(batchConfigs, config)
	}

	// 返回分页数据
	c.JSON(http.StatusOK, respJSON(JsonObject{
		"list": batchConfigs,
	}))
}

func BatchDelHandler(c *gin.Context) {
	mac := c.Query("mac")
	if mac == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return

	}

	_, err := dao.RedisCache.ZRem(c.Request.Context(), BatchUrlZsetKey, mac).Result()
	if err != nil {
		log.Errorf("Failed to remove %s: %v", mac, err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if _, err := dao.RedisCache.HDel(c.Request.Context(), BatchUrlHSetKey, mac).Result(); err != nil {
		log.Errorf("Failed to remove %s: %v", mac, err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{"msg": "success"}))

}

type BatchAddress struct {
	Name    string
	Url     string
	AddTime time.Time
	Enable  bool
}

const BatchAddressZsetKey = "titan:batch:address:set"

func BatchAddressSetHandler(c *gin.Context) {
	var address BatchAddress
	if err := c.BindJSON(&address); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	address.AddTime = time.Now()
	data, err := json.Marshal(address)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	score := float64(address.AddTime.Unix())

	_, err = dao.RedisCache.TxPipelined(c.Request.Context(), func(pipe redis.Pipeliner) error {
		members, err := dao.RedisCache.ZRangeByScore(c.Request.Context(), BatchAddressZsetKey, &redis.ZRangeBy{
			Min: "-inf", Max: "+inf",
		}).Result()
		if err != nil {
			return err
		}

		for _, member := range members {
			var existing BatchAddress
			if err := json.Unmarshal([]byte(member), &existing); err != nil {
				continue
			}
			if existing.Url == address.Url {
				_, err = pipe.ZRem(c.Request.Context(), BatchAddressZsetKey, member).Result()
				if err != nil {
					return err
				}
			}
		}

		// 新增新的 BatchAddress
		_, err = pipe.ZAdd(c.Request.Context(), BatchAddressZsetKey, redis.Z{
			Score:  score,
			Member: data,
		}).Result()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{"msg": "success"}))
}

func BatchAddressDelHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	// 开始事务删除
	_, err := dao.RedisCache.TxPipelined(c.Request.Context(), func(pipe redis.Pipeliner) error {
		// 查找所有 BatchAddress
		members, err := dao.RedisCache.ZRangeByScore(c.Request.Context(), BatchAddressZsetKey, &redis.ZRangeBy{
			Min: "-inf", Max: "+inf",
		}).Result()
		if err != nil {
			return err
		}

		// 遍历找到匹配的 URL 条目
		for _, member := range members {
			var existing BatchAddress
			if err := json.Unmarshal([]byte(member), &existing); err != nil {
				log.Errorf("Unmarshal BatchAddress fail: %v, err: %v", existing, err)
				continue
			}
			if existing.Url == url {
				// 删除匹配的条目
				_, err = pipe.ZRem(c.Request.Context(), BatchAddressZsetKey, member).Result()
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{"msg": "success"}))
}

func UserBatchAddressHandler(c *gin.Context) {
	var res = make([]BatchAddress, 0)
	members, err := dao.RedisCache.ZRangeByScore(c.Request.Context(), BatchAddressZsetKey, &redis.ZRangeBy{
		Min: "-inf", Max: "+inf",
	}).Result()
	if err != nil {
		log.Errorf("load redis %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	for _, member := range members {
		var existing BatchAddress
		if err := json.Unmarshal([]byte(member), &existing); err != nil {
			log.Errorf("parse redis %v: %v", member, err)
			continue
		}
		if existing.Enable {
			res = append(res, existing)
		}
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"list": res,
	}))
}

func BatchAddressListHandler(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if size <= 0 {
		size = 10
	}
	if page <= 0 {
		page = 1
	}

	start := int64((page - 1) * size)
	stop := int64(page*size - 1)

	rdb := dao.RedisCache

	total, err := rdb.ZCard(c.Request.Context(), BatchAddressZsetKey).Result()
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if total == 0 {
		c.JSON(http.StatusOK, respJSON(JsonObject{
			"data":  []BatchAddress{},
			"page":  page,
			"size":  size,
			"total": total,
		}))
		return
	}

	members, err := rdb.ZRange(c.Request.Context(), BatchAddressZsetKey, start, stop).Result()
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	var addresses []BatchAddress

	for _, member := range members {
		var address BatchAddress
		if err := json.Unmarshal([]byte(member), &address); err != nil {
			continue
		}
		addresses = append(addresses, address)
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"data":  addresses,
		"page":  page,
		"size":  size,
		"total": total,
	}))
}
