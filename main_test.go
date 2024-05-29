package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Salvionied/apollo/serialization/PlutusData"
	"github.com/Salvionied/cbor/v2"
	"github.com/SundaeSwap-finance/kugo"
	"github.com/SundaeSwap-finance/ogmigo/v6"
	"github.com/wasmerio/wasmer-go/wasmer"
)

func TestLoadWasm(t *testing.T) {
	ogmigoClient := ogmigo.New(ogmigo.WithEndpoint("ws://localhost:1337"))
	kugoClient := kugo.New(kugo.WithEndpoint("http://localhost:1442"))
	// ogmios := OgmiosChainContext.NewOgmiosChainContext(*ogmigoClient, *kugoClient)
	//ogmios.Init()
	fmt.Println(ogmigoClient, kugoClient)
	// ogmios.Utxos()

	//addr, _ := Address.DecodeAddress("addr_test1wzamv8e88falzhvvc7gmrwjwwugqrf8t3pfpa32aun2wjdqh65j75")
	//utxos := ogmios.Utxos(addr)
	//fmt.Println(utxos)

	utxos := getAddressUtxos(*kugoClient, "addr_test1wzqahltmnzlcpjrn58lldzm0j96p24n0flm2ha9zcx8tjxqmqlqg9", true)
	decoded, err := hex.DecodeString(utxos[0].InlineDatum)
	if err != nil {
		log.Fatal(err)
	}

	var x PlutusData.PlutusData
	err = cbor.Unmarshal(decoded, &x)

	if err != nil {
		log.Fatal(err)
	}
	l := PlutusData.DatumOptionInline(&x)
	fmt.Println(l)

	// cborData, _ := l.MarshalCBOR()

	// fmt.Println(cborData)

	// // var datum HandlerDatumSchema
	// // err1 := cbor.Unmarshal(cborData, &datum)
	// // if err1 != nil {
	// // 		fmt.Println("Error unmarshaling CBOR:", err)
	// // 		return
	// // }
}

func TestDecodeCbor(t *testing.T) {
	policy, _ := hex.DecodeString("7f09a7d17522a33b8be76826dd080abdff348218fafd60fbf991de11")
	name, _ := hex.DecodeString("68616e646c6572")
	handler := HandlerDatumSchema{
		State: HandlerStateSchema{
			NextClientSequence: 1,
			NextConnectionSequence: 2,
			NextChannelSequence: 3,
			BoundPort: []int64{},
		},
		Token: AuthenTokenSchema{
			Policy: policy,
			Name: name,
		},
	}

	cborData, _ := cbor.Marshal(handler)
	fmt.Println(cborData)
	cborEncode, _ := hex.DecodeString(hex.EncodeToString(cborData))
	var decoded HandlerDatumSchema
	cbor.Unmarshal(cborEncode, &decoded)

	fmt.Println(decoded)
}

func TestLoadWasmerFile(t *testing.T) {
	wasmBytes, _ := os.ReadFile("cardano_multiplatform_lib_bg.wasm")

	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	// Compiles the module
	module, _ := wasmer.NewModule(store, wasmBytes)

	importObject := wasmer.NewImportObject()
	instance, _ := wasmer.NewInstance(module, importObject)
	fmt.Println(instance)

	//  Gets the `sum` exported function from the WebAssembly instance.
	PlutusData, _ := instance.Exports.GetTable("PlutusData")
	fmt.Println(PlutusData)
}
