package bitflyer

//check https://api.bitflyer.com/v1/getmarkets

const (
	PRODUCTCODE_BTC_JPY  string = "BTC_JPY"
	PRODUCTCODE_XRP_JPY  string = "XRP_JPY"
	PRODUCTCODE_ETH_JPY  string = "ETH_JPY"
	PRODUCTCODE_XLM_JPY  string = "XLM_JPY"
	PRODUCTCODE_MONA_JPY string = "MONA_JPY"
	PRODUCTCODE_ETH_BTC  string = "ETH_BTC"
	PRODUCTCODE_BCH_BTC  string = "BCH_BTC"

	CURRENCYCODE_BTC     string = "BTC"
	CURRENCYCODE_XRP     string = "XRP"
	CURRENCYCODE_ETH     string = "ETH"
	CURRENCYCODE_XLM     string = "XLM"
	CURRENCYCODE_MONA    string = "MONA"

	FEE_TRADE_RATE       float64 = 0.0015


	STATUS_OPEN string = "NORMAL"

	MODE_LIMIT  string = "LIMIT"
	MODE_MARKET string = "MARKET"

	ORDER_STATE_OPEN  string = "ACTIVE"
	ORDER_STATE_FIXED string = "COMPLETED"
)
