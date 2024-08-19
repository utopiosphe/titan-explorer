package dao

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

const (
	tableAreaMap = "area_map"
)

// GetAreaCnByAreaEn 通过英文的区域名称获取中文的区域名称
func GetAreaCnByAreaEn(ctx context.Context, areaEn []string) ([]string, error) {
	var areaCn []string

	query, args, err := squirrel.Select("area_cn").From(tableAreaMap).Where(squirrel.Eq{
		"area_en": areaEn,
	}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("generate get areas error:%w", err)
	}

	err = DB.SelectContext(ctx, &areaCn, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get areas error:%w", err)
	}

	return areaCn, nil
}

// GetAreaEnByAreaCn 通过中文的区域名称获取英文的区域名称
func GetAreaEnByAreaCn(ctx context.Context, areaCn []string) ([]string, error) {
	var areaEn []string

	query, args, err := squirrel.Select("area_en").From(tableAreaMap).Where(squirrel.Eq{
		"area_cn": areaCn,
	}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("generate get areas error:%w", err)
	}

	err = DB.SelectContext(ctx, &areaEn, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get areas error:%w", err)
	}

	return areaEn, nil
}
