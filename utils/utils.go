package utils

import (
	"reflect"
	"sync"
	"time"

	"github.com/anvie/port-scanner"
	"github.com/prazd/nodes_mon_bot/config"
	"github.com/prazd/nodes_mon_bot/state"
	"github.com/prazd/nodes_mon_bot/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

type NodesInfo struct {
	State *state.SingleState
	Port int
	Addresses []string
}

func Worker(wg *sync.WaitGroup, addr string, port int, r *state.SingleState) {
	defer wg.Done()
	ps := portscanner.NewPortScanner(addr, 1*time.Second, 1)
	isAlive := ps.IsOpen(port)
	r.Set(addr, isAlive)
}

func GetHostInfo(curr string, configData config.Config) ([]string, int) {

	var addresses []string
	var port int

	switch curr {
	case "eth":
		addresses = configData.EthNodes.Addresses
		port = configData.EthNodes.Port
	case "etc":
		addresses = configData.EtcNodes.Addresses
		port = configData.EtcNodes.Port
	case "btc":
		addresses = configData.BtcNodes.Addresses
		port = configData.BtcNodes.Port
	case "bch":
		addresses = configData.BchNodes.Addresses
		port = configData.BchNodes.Port
	case "ltc":
		addresses = configData.LtcNodes.Addresses
		port = configData.LtcNodes.Port
	case "xlm":
		addresses = configData.XlmNodes.Addresses
		port = configData.XlmNodes.Port
	}

	return addresses, port
}

func GetMessage(result map[string]bool) string {
	var message string
	for address, status := range result {
		message += address
		switch status {
		case true:
			message += ": ✔"
		case false:
			message += ": ✖"
		}
		message += "\n"
	}

	return message
}

func NodesStatus(curr string, configData config.Config) string {

	nodesState := state.NewSingleState()

	addresses, port := GetHostInfo(curr, configData)

	RunWorkers(addresses, port, nodesState)

	message := GetMessage(nodesState.Result)

	return message
}

func RunWorkers(addresses []string, port int, state *state.SingleState){
	var wg sync.WaitGroup

	for i := 0; i < len(addresses); i++ {
		wg.Add(1)
		go Worker(&wg, addresses[i], port, state)
	}
	wg.Wait()
}

func isAllNodesUp(addresses []string, port int, state *state.SingleState) bool{
	RunWorkers(addresses, port, state)
	for _,j := range state.Result{
		if j == false{
			return false
		}
	}
	return true
}

func GetAllNodes(configData config.Config) map[string]NodesInfo {
	return map[string]NodesInfo{
		"eth": NodesInfo{
			State:     state.NewSingleState(),
			Port:      configData.EthNodes.Port,
			Addresses: configData.EthNodes.Addresses,
		},
		"etc": NodesInfo{
			State:     state.NewSingleState(),
			Port:      configData.EtcNodes.Port,
			Addresses: configData.EtcNodes.Addresses,
		},
		"btc": NodesInfo{
			State:     state.NewSingleState(),
			Port:      configData.BtcNodes.Port,
			Addresses: configData.BtcNodes.Addresses,
		},
		"ltc": NodesInfo{
			State:     state.NewSingleState(),
			Port:      configData.LtcNodes.Port,
			Addresses: configData.LtcNodes.Addresses,
		},
		"bch": NodesInfo{
			State:     state.NewSingleState(),
			Port:      configData.BchNodes.Port,
			Addresses: configData.BchNodes.Addresses,
		},
		"xlm": NodesInfo{
			State:     state.NewSingleState(),
			Port:      configData.XlmNodes.Port,
			Addresses: configData.XlmNodes.Addresses,
		},
	}
}

func FullCheckOfNode(configData config.Config, bot *tb.Bot){

	allNodes := GetAllNodes(configData)

	for {
		for currency, nodesInfo := range allNodes {
				up := isAllNodesUp(nodesInfo.Addresses, nodesInfo.Port, nodesInfo.State)
				if !up {
					ids := db.GetAllSubscribers()
					if ids == nil{
						continue
					}
					message := GetMessage(nodesInfo.State.Result)
					for i:=0; i<len(ids); i++ {
						bot.Send(&tb.User{ID:ids[i]},"Subscribe message:\nCurrency: " + currency +"\n"+  message)
					}
				}
		}
		time.Sleep(time.Second * 60)
	}
}

func Contains(params ...interface{}) bool {
	v := reflect.ValueOf(params[0])
	arr := reflect.ValueOf(params[1])

	var t = reflect.TypeOf(params[1]).Kind()

	if t != reflect.Slice && t != reflect.Array {
		panic("Type Error! Second argument must be an array or a slice.")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == v.Interface() {
			return true
		}
	}
	return false
}

func Start(id int)(error){
	inDb, err := db.IsInDb(id)
	if err != nil{
		return err
	}

	if !inDb{
		err = db.CreateUser(id)
		if err != nil{
			return err
		}
	}
	return nil
}
