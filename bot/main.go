package main

import (
	"log"
	"os"
	"time"

	"encoding/json"
	"github.com/prazd/nodes_mon_bot/config"
	"github.com/prazd/nodes_mon_bot/keyboard"
	"github.com/prazd/nodes_mon_bot/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"path/filepath"
)

func ReadConfig() (*config.Config, error) {

	defaultConfigPath, _ := filepath.Abs("../config/config.json")
	configFile, err := os.Open(defaultConfigPath)
	if err != nil {
		return nil, err
	}

	defer configFile.Close()

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil{
		return nil, err
	}

	var conf config.Config

	if err = json.Unmarshal(byteValue, &conf); err != nil{
		return nil, err
	}

	return &conf, nil
}

func main() {

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("token"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	configData,err := ReadConfig()
	if err !=nil{
		log.Println(err)
		os.Exit(1)
	}

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "Hi!I can help you with nodes monitoring!", &tb.SendOptions{ParseMode: "Markdown"},
			&tb.ReplyMarkup{ResizeReplyKeyboard: true, ReplyKeyboard: keyboard.MainMenu})
	})

	// Handlers
	b.Handle(&keyboard.EthButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("eth", *configData))
	})

	b.Handle(&keyboard.EtcButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("etc", *configData))
	})

	b.Handle(&keyboard.BtcButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("btc", *configData))
	})

	b.Handle(&keyboard.BchButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("bch", *configData))
	})

	b.Handle(&keyboard.LtcButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("ltc", *configData))
	})

	b.Handle(&keyboard.XlmButton, func(m *tb.Message) {
		b.Send(m.Sender, utils.IsAlive("xlm", *configData))
	})

	b.Start()
}
