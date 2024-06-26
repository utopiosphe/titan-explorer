package oss

import (
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gnasnik/titan-explorer/config"
)

var (
	tokenCacheKey     = "cache:sts:token"
	ossClientCacheKey = "cache:oss:token"
)

var (
	OssInstance OssAPI
)

type OssAPI interface {
	SignUrl(bucket, objectKey string, expire int64) (string, error)
	Upload(bucket, obj string, buf io.Reader) error
}

type ossAPI struct {
	endpoint     string
	accessId     string
	accessSecret string
	client       *oss.Client
}

func InitFromCfg(cfg config.OssConfig) error {
	client, err := oss.New(cfg.EndPoint, cfg.AccessId, cfg.AccessKey)
	if err != nil {
		return err
	}
	OssInstance = &ossAPI{
		endpoint:     cfg.EndPoint,
		accessId:     cfg.AccessId,
		accessSecret: cfg.AccessKey,
		client:       client,
	}
	return nil
}

type Option func(*ossAPI)

func NewMustOssAPI(endpint, id, secret string) OssAPI {
	client, err := oss.New(endpint, id, secret)
	if err != nil {
		panic(err)
	}
	return &ossAPI{
		endpoint:     endpint,
		accessId:     id,
		accessSecret: secret,
		client:       client,
	}
}

func (o *ossAPI) Upload(bucket, obj string, buf io.Reader) error {
	bk, err := o.client.Bucket(bucket)
	if err != nil {
		return err
	}

	return bk.PutObject(obj, buf, oss.ObjectACL(oss.ACLPublicRead))
}

func (o *ossAPI) SignUrl(bucket, objectKey string, expire int64) (string, error) {
	bk, err := o.client.Bucket(bucket)
	if err != nil {
		return "", err
	}

	srcUrl, err := bk.SignURL(objectKey, oss.HTTPGet, expire)
	if err != nil {
		return "", nil
	}
	return srcUrl, nil
}
