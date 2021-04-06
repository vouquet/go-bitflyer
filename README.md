go-bitflyer
===

* the library for easy use of [https://lightning.bitflyer.com/docs/api?lang=ja](https://lightning.bitflyer.com/docs/api?lang=ja)
	* cannot order yet.

## sample

* easy
```

import (
	"fmt"
	"context"
)

import "github.com/vouquet/go-bitflyer/bitflyer"

func main() {
	API_KEY = "your api key"
	SECRET_KEY = "your secret key"

	shop, err := bitflyer.NewBitflyer.(API_KEY, SECRET_KEY, context.Background())
	if err != nil {
		panic(err)
	}
	defer shop.Close()

	target := []string{
		bitflyer.PRODUCTCODE_BTC_JPY,
		bitflyer.PRODUCTCODE_ETH_JPY,
	}

	rates, err = shop.GetRates()
	if err != nil {
		panic(err)
	}

	for symbol, rate := range rates {
		fmt.Println(symbol, rate)
	}
}
```

