package constants

const (
	PriceStatsUrl = `https://api.binance.com/api/v3/ticker/24hr`
	BTC = "BTC"
	ETH = "ETH"
	USDT = "USDT"
	BNB = "BNB"
	CRON_EVERY_FIVE_MINUTES = "*/5 * * * *"
	CRON_EVERY_TEN_MINUTES = "*/10 * * * *"
	CRON_EVERY_FIFTEEN_MINUTES = "*/15 * * * *"
	CRON_EVERY_THIRTY_MINUTES = "*/30 * * * *"
	CRON_EVERY_HOUR = "0 * * * *"
	ONE_SECOND = 1
	ONE_MINUTE = ONE_SECOND * 60
	FIVE_MINUTES = ONE_MINUTE * 5
	FIFTEEN_MINUTES = ONE_MINUTE * 15
	ONE_HOUR = ONE_MINUTE * 60
	FOUR_HOURS = ONE_HOUR * 4
	ONE_DAY = ONE_HOUR * 24
	ONE_MONTH = ONE_DAY * 30
)
