package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"

	"github.com/Salvionied/apollo/serialization/PlutusData"
	"github.com/Salvionied/apollo/txBuilding/Backend/Base"
	"github.com/Salvionied/cbor/v2"
	"github.com/SundaeSwap-finance/kugo"
	chainsyncv5 "github.com/SundaeSwap-finance/ogmigo/ouroboros/chainsync"
	"github.com/SundaeSwap-finance/ogmigo/v6"
)

type HandlerStateSchema struct {
	NextClientSequence int64
	NextConnectionSequence int64
	NextChannelSequence int64
	BoundPort []int64
}
type AuthenTokenSchema struct {
	Policy []byte
	Name []byte
}

type HandlerDatumSchema struct {
	State HandlerStateSchema
	Token AuthenTokenSchema
}

func main() {
	ogmigoClient := ogmigo.New(ogmigo.WithEndpoint("ws://localhost:1337"))
	kugoClient := kugo.New(kugo.WithEndpoint("http://localhost:1442"))
	// ogmios := OgmiosChainContext.NewOgmiosChainContext(*ogmigoClient, *kugoClient)
	//ogmios.Init()
	fmt.Println(ogmigoClient, kugoClient)
	// ogmios.Utxos()

	//addr, _ := Address.DecodeAddress("addr_test1wzamv8e88falzhvvc7gmrwjwwugqrf8t3pfpa32aun2wjdqh65j75")
	//utxos := ogmios.Utxos(addr)
	//fmt.Println(utxos)

	utxos := getAddressUtxos(*kugoClient, "addr_test1wzamv8e88falzhvvc7gmrwjwwugqrf8t3pfpa32aun2wjdqh65j75", true)
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
	cborBytes, _ := l.MarshalCBOR()

	fmt.Println(cborBytes)
}

func getAddressUtxos(kugoClient kugo.Client, address string, gather bool) []Base.AddressUTXO {
	ctx := context.Background()
	addressUtxos := make([]Base.AddressUTXO, 0)
	matches, err := kugoClient.Matches(ctx, kugo.OnlyUnspent(), kugo.Address(address))
	if err != nil {
		log.Fatal(err, "OgmiosChainContext: AddressUtxos: kupo request failed")
	}
	for _, match := range matches {
		datum := ""
		if match.DatumType == "inline" {
			datum, err = kugoClient.Datum(ctx, match.DatumHash)
			if err != nil {
				log.Fatal(err, "OgmiosChainContext: AddressUtxos: kupo datum request failed")
			}
		}
		addressUtxos = append(addressUtxos, Base.AddressUTXO{
			TxHash:      match.TransactionID,
			OutputIndex: match.OutputIndex,
			Amount:      chainsyncValue_toAddressAmount(chainsyncv5.Value{}),
			// We probably don't need this info and kupo doesn't provide it in this query
			Block:       "",
			DataHash:    match.DatumHash,
			InlineDatum: datum,
		})
	}
	return addressUtxos
}

func chainsyncValue_toAddressAmount(v chainsyncv5.Value) []Base.AddressAmount {
	amts := make([]Base.AddressAmount, 0)
	amts = append(amts, Base.AddressAmount{
		Unit:     "lovelace",
		Quantity: strconv.FormatInt(v.Coins.Int64(), 10),
	})
	for assetId, quantity := range v.Assets {
		a := string(assetId)
		policy := a[:56] // always 28 bytes
		token := a[57:]  // skip the '.'
		amts = append(amts, Base.AddressAmount{
			Unit:     policy + token,
			Quantity: strconv.FormatInt(quantity.Int64(), 10),
		})
	}
	return amts
}
