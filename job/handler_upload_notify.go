package job

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gnasnik/titan-explorer/core/dao"
	"github.com/gnasnik/titan-explorer/core/opasynq"
	"github.com/gnasnik/titan-explorer/core/storage"
	"github.com/hibiken/asynq"
	"github.com/jinzhu/copier"
)

type AssetUploadNotifyReq struct {
	ExtraID  string // 外部文件ID
	TenantID string // 租户ID
	UserID   string // 上传者ID

	AssetName   string
	AssetCID    string
	AssetType   string
	AssetSize   int64
	GroupID     int64
	CreatedTime time.Time
}

func assetUploadNotify(ctx context.Context, t *asynq.Task) error {

	var (
		payload opasynq.AssetUploadNotifyPayload
		err     error
	)

	err = json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		cronLog.Errorf("unable to parse message %+v", t.Payload())
		return err
	}

	defer func(err error) {
		if err != nil {

		}
	}(err)

	var body AssetUploadNotifyReq
	if err = copier.Copy(&body, &payload); err != nil {
		cronLog.Errorf("unable to copy asset %+v", payload)
		return err
	}

	tenantInfo, err := dao.GetTenantByBuilder(ctx, squirrel.Select("*").Where("tanant_id = ?", payload.TenantID))
	if err != nil {
		cronLog.Errorf("unable to find tenant info %+v", err)
		return err
	}

	pair, err := storage.LoadTenantKeyPairFromBlob([]byte(tenantInfo.ApiKey))
	if err != nil {
		cronLog.Errorf("unable to generate secret with pair %+v", err)
		return err
	}

	var (
		secret = pair.ApiSecret
		// key         = pair.ApiKey
		method      = "POST"
		url         = payload.NotifyUrl
		bodyData, _ = json.Marshal(body)
	)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyData))
	if err != nil {
		cronLog.Errorf("unable to generate req %+v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	ts := time.Now().Format(time.RFC3339)
	req.Header.Set("X-Timestamp", ts)

	c, err := rand.Int(rand.Reader, big.NewInt(1e16))
	if err != nil {
		cronLog.Errorf("unable to generate c %+v", err)
		return err
	}
	nonce := fmt.Sprintf("%d", c)
	req.Header.Set("X-Nonce", nonce)

	signature := genCallbackSignature(secret, method, url, string(bodyData), ts, nonce)
	req.Header.Set("X-Signature", signature)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		cronLog.Errorf("unable to send post to %s with signature %s with err %+v", url, signature, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		cronLog.Errorf("received non ok http code %d with body %s", resp.StatusCode, resp.Body)
		return fmt.Errorf("received non ok http code %d with body %s", resp.StatusCode, resp.Body)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		cronLog.Errorf("unable to read response body %+v", err)
		return err
	}

	if string(respData) != "success" {
		cronLog.Errorf("unexpected resp %s", respData)
		return fmt.Errorf("unexpected resp %s", respData)
	}

	cronLog.Infof("Notified client %s, status code %d, req %+v", url, resp.StatusCode, req)

	return nil
}

func genCallbackSignature(secret, method, path, body, timestamp, nonce string) string {
	data := method + path + body + timestamp + nonce
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
