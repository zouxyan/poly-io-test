package neo

import (
	"encoding/hex"
	"fmt"
	"github.com/joeqian10/neo-gogogo/helper"
	"github.com/polynetwork/poly-io-test/config"
	"github.com/polynetwork/poly-io-test/log"
	"testing"
)

func TestNewNeoInvoker(t *testing.T) {
	config.DefConfig.Init("./config-test-cp.json")
	invoker, err := NewNeoInvoker()
	if err != nil {
		t.Fatal(err)
	}
	res := invoker.Cli.GetBlockCount()
	if res.HasError() {
		t.Fatal(res.Error.Message)
	}
	fmt.Println(res.Result)

	res1 := invoker.Cli.GetContractState("0xa837ba329255884b40581ba8a3d29820acf44316")
	fmt.Println(res1)
}

func Test_GetOperator(t *testing.T) {
	config.DefConfig.Init("../../config.json")
	invoker, err := NewNeoInvoker()
	if err != nil {
		t.Fatal(err)
	}
	ccmc, _ := helper.UInt160FromString(config.DefConfig.NeoCCMC)
	neoLock, _ := helper.UInt160FromString(config.DefConfig.NeoLockProxy)

	// check if initialized
	{
		res, err := invoker.GetStorage(ccmc.Bytes(), hex.EncodeToString([]byte("IsInitGenesisBlock")))
		if err != nil {
			log.Errorf("Initialized height err: %v", err)
			return
		}
		log.Infof("ccmc initialized poly height on neo ccmc is: %s", res)
	}

	{
		res, err := invoker.GetProxyOperator(neoLock.Bytes())
		if err != nil {
			log.Errorf("GetOperator err: %v", err)
			return
		}
		log.Infof("neoLockProxy operator: %s", res)
	}

	var toChainId uint64 = 101
	{
		res, err := invoker.GetProxyHash(neoLock.Bytes(), toChainId)
		if err != nil {
			log.Errorf("GetProxyHash err: %v", err)
			return
		}
		log.Infof("neoLockProxy: %x, toChainId: %d, toChainProxyHash (little) : %s", neoLock.Bytes(), toChainId, res)
	}

	fromAssetHashs := []string{
		"0xb4dea34fa0b3dd6b421733d7c6c3c083a7fb6633",
		"0x2425c293d7c2cffd9a7944226ab3b702d3985cc6",
	}
	{
		froms := make([][]byte, 0)
		for _, from := range fromAssetHashs {
			f, _ := ParseNeoAddr(from)
			froms = append(froms, f)
		}
		res, err := invoker.GetAssetHashs(neoLock.Bytes(), toChainId, froms)
		if err != nil {
			log.Errorf("GetProxyHash err: %v", err)
			return
		}
		bals, err := invoker.GetAssetBalances(neoLock.Bytes(), froms)
		if err != nil {
			log.Errorf("GetAssetBalances err: %v", err)
			return
		}

		for i, _ := range fromAssetHashs {
			log.Infof("neoLockProxy GetAssetHashs, toChainId: %d, from: %s, proxyBalance: %s, to(little): %s", toChainId, fromAssetHashs[i], bals[i].String(), res[i])
		}
	}

}
