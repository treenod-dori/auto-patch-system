package utils

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestSplitSQLQueries(t *testing.T) {
	query := "# 스크롤 패키지 43회차 보상\nDELETE FROM product_info WHERE mainGroup = 19;\nINSERT INTO product_info (productid_aos, productid_ios, dprice_ios, dprice_aos, yprice, subGroup, mainGroup) VALUES\n('a_lgpkpk_74_2', 'i_lgpkpk_74_2', '5.99 ', '5.99 ', '1000', 1, 19),\n('a_lgpkpk_74_3', 'i_lgpkpk_74_3', '5.99 ', '5.99 ', '1000', 2, 19),\n('a_lgpkpk_74_4', 'i_lgpkpk_74_4', '5.99 ', '5.99 ', '1000', 3, 19),\n('a_lgpkpk_74_5', 'i_lgpkpk_74_5', '5.99 ', '5.99 ', '1000', 4, 19),\n('a_lgpkpk_74_6', 'i_lgpkpk_74_6', '5.99 ', '5.99 ', '1000', 5, 19);\n\nDELETE FROM ref_scrollPackage WHERE packageVer = 43;\nINSERT INTO ref_scrollPackage (packageInfo, packageCount, msgNo, packageVer, finalRewardType, finalRewardValue, priceType) VALUES\n('[{\"id\":1,\"gCode\":0,\"price\":0,\"isSpecial\":true,\"reward\":[]},{\"id\":2,\"gCode\":0,\"price\":0,\"isSpecial\":true,\"reward\":[{\"type\":7,\"value\":1}]},{\"id\":3,\"gCode\":0,\"price\":0,\"isSpecial\":true,\"reward\":[{\"type\":6,\"value\":2}]},{\"id\":4,\"gCode\":0,\"price\":0,\"isSpecial\":false,\"reward\":[{\"type\":5,\"value\":5}]},{\"id\":5,\"gCode\":0,\"price\":0,\"isSpecial\":false,\"reward\":[{\"type\":26,\"value\":2}]}]', 5, 0, 43, 101, 193, 0);\n\nDELETE FROM ref_scrollPackageExtraReward WHERE version = 43;\nINSERT INTO ref_scrollPackageExtraReward (rewardType, rewardValue, version) VALUES\n(18, 15, 43);"
	rawData := []byte(query)

	decodedData, _ := base64.StdEncoding.DecodeString(base64.StdEncoding.EncodeToString(rawData))
	fmt.Println(string(decodedData)) // 원본 SQL 출력됨

}
