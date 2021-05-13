package binanceUtil

import (
	"io/ioutil"
	"net/http"

	"github.com/kingokeke/go-example-1/constants"
	"github.com/kingokeke/go-example-1/utils"
)

func GetPriceInfoAsByteArray() []byte {
	response, e := http.Get(constants.PriceStatsUrl)
	utils.CheckError(e)
	utils.LogToGeneral("Successfully retrieved price data from exchange...")

	defer response.Body.Close()

	bodyBytes, e := ioutil.ReadAll(response.Body)
	utils.CheckError(e)

	return bodyBytes
}
