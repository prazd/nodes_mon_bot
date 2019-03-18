package keyboard

import tb "gopkg.in/tucnak/telebot.v2"

var (
	EthButton = tb.InlineButton{
		Unique: "ETH",
		Text:   "ETH",
	}

	EtcButton = tb.InlineButton{
		Unique: "ETC",
		Text:   "ETC",
	}

	BtcButton = tb.InlineButton{
		Unique: "BTC",
		Text:   "BTC",
	}

	BchButton = tb.InlineButton{
		Unique: "BCH",
		Text:   "BCH",
	}

	LtcButton = tb.InlineButton{
		Unique: "LTC",
		Text:   "LTC",
	}

	MainMenu = [][]tb.InlineButton{
		[]tb.InlineButton{EthButton, EtcButton},
		[]tb.InlineButton{BtcButton, BchButton, LtcButton},
	}

	// TODO: Balance check
)