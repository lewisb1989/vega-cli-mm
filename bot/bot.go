package bot

import (
	apipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	vegapb "code.vegaprotocol.io/vega/protos/vega"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
	"vega-cli-mm/auth"
	"vega-cli-mm/logging"
	"vega-cli-mm/store"
	"vega-cli-mm/vega"
)

type Bot struct {
	store *store.Store
	vega  *vega.Vega
}

func NewBot(
	store *store.Store,
	vega *vega.Vega,
) *Bot {
	return &Bot{
		store: store,
		vega:  vega,
	}
}

func (b *Bot) initWallet() {
	secretFile, err := os.Open(".secret")
	if err != nil {
		logging.Panic(fmt.Sprintf("error loading .secret: %v", err))
	}
	mnemonicBytes, err := io.ReadAll(secretFile)
	if err != nil {
		logging.Panic(fmt.Sprintf("error loading .secret: %v", err))
	}
	wallet := auth.NewWallet(string(mnemonicBytes))
	authenticator := auth.NewAuthenticator(b.vega.GetCoreNode(), wallet, b.store)
	b.vega.SetAuthenticator(authenticator)
}

func (b *Bot) loadMarkets() {
	marketsFile, err := os.Open("markets.json")
	if err != nil {
		logging.Panic(fmt.Sprintf("error loading markets.json: %v", err))
	}
	var markets []*store.MarketConfig
	marketsJson, err := io.ReadAll(marketsFile)
	if err != nil {
		logging.Panic(fmt.Sprintf("error loading markets.json: %v", err))
	}
	err = json.Unmarshal(marketsJson, &markets)
	if err != nil {
		logging.Panic(fmt.Sprintf("error loading markets.json: %v", err))
	}
	err = marketsFile.Close()
	if err != nil {
		logging.Panic(fmt.Sprintf("error closing markets.json: %v", err))
	}
	for i, market := range markets {
		market.KeyPair = b.vega.GetAuthenticator().GetWallet().Get(uint(i))
		b.store.SaveMarketConfig(market)
	}
}

func (b *Bot) updateReferencePrices() {
	go func() {
		for range time.NewTicker(time.Second).C {
			for _, config := range b.store.GetMarketConfig() {
				// TODO - update reference price for market
				/**
				* 1) Check the price source
				* 2) Use the relevant code to grab the price (for Binance it comes via WS, the ETH stuff will be sync)
				* 3) Update config and save it in the store
				 */
				b.store.SaveMarketConfig(config)
			}
		}
	}()
}

func (b *Bot) connectToVegaStreams() {
	go func() {
		for range time.NewTicker(time.Second).C {
			var marketIds []string
			var partyIds []string
			for _, config := range b.store.GetMarketConfig() {
				partyId := config.KeyPair.PublicKey
				marketIds = append(marketIds, config.VegaId)
				partyIds = append(partyIds, partyId)
			}
			if !b.vega.IsMarketDataConnected() {
				b.vega.StreamMarketData(marketIds, func(marketData []*vegapb.MarketData) {
					for _, data := range marketData {
						b.store.SaveMarketData(data)
					}
				})
			}
			if !b.vega.IsOrdersConnected() {
				b.vega.StreamOrders(partyIds, func(orders []*vegapb.Order) {
					for _, order := range orders {
						b.store.SaveOrder(order)
					}
				})
			}
			for _, partyId := range partyIds {
				if !b.vega.IsAccountsConnected(partyId) {
					b.vega.StreamAccounts(partyId, func(accounts []*apipb.AccountBalance) {
						for _, account := range accounts {
							b.store.SaveAccount(account)
						}
					})
				}
				if !b.vega.IsLiquidityProvisionsConnected(partyId) {
					b.vega.StreamLiquidityProvisions(partyId, func(liquidityProvisions []*vegapb.LiquidityProvision) {
						for _, lp := range liquidityProvisions {
							b.store.SaveLiquidityProvision(lp)
						}
					})
				}
				if !b.vega.IsPositionsConnected(partyId) {
					b.vega.StreamPositions(partyId, func(positions []*vegapb.Position) {
						for _, position := range positions {
							b.store.SavePosition(position)
						}
					})
				}
			}
		}
	}()
}

func (b *Bot) syncVegaData() {
	go func() {
		for range time.NewTicker(time.Second * 15).C {
			var marketIds []string
			var partyIds []string
			assets := b.vega.GetAssets()
			for _, asset := range assets {
				b.store.SaveAsset(asset)
			}
			markets := b.vega.GetMarkets()
			for _, market := range markets {
				b.store.SaveMarket(market)
			}
			marketData := b.vega.GetMarketData()
			for _, data := range marketData {
				b.store.SaveMarketData(data)
			}
			networkParameters := b.vega.GetNetworkParameters()
			for _, param := range networkParameters {
				b.store.SaveNetworkParameter(param)
			}
			for _, config := range b.store.GetMarketConfig() {
				partyId := config.KeyPair.PublicKey
				marketIds = append(marketIds, config.VegaId)
				partyIds = append(partyIds, partyId)
				liquidityProvisions := b.vega.GetLiquidityProvisions(partyId)
				for _, lp := range liquidityProvisions {
					b.store.SaveLiquidityProvision(lp)
				}
			}
			orders := b.vega.GetOrders(partyIds)
			for _, order := range orders {
				b.store.SaveOrder(order)
			}
			positions := b.vega.GetPositions(partyIds)
			for _, position := range positions {
				b.store.SavePosition(position)
			}
			accounts := b.vega.GetAccounts(partyIds)
			for _, account := range accounts {
				b.store.SaveAccount(account)
			}
		}
	}()
}

func (b *Bot) updateLiquidityCommitment() {
	go func() {
		for range time.NewTicker(time.Second).C {
			for _, config := range b.store.GetMarketConfig() {
				// TODO - connect liquidity commitment on Vega
				/**
				* 1) Get Vega balance to determine commitment size
				* 2) Get existing commitment to ensure we don't amend size down unless threshold exceeded
				* 3) Amend commitment amount, or create it if it doesn't exist
				 */
				logging.GetLogger().Infof("update vega liquidity commitment for market: %s", config.VegaId)
			}
		}
	}()
}

func (b *Bot) updateQuotes() {
	go func() {
		for range time.NewTicker(time.Second).C {
			for _, config := range b.store.GetMarketConfig() {
				// TODO - update quotes on Vega
				/**
				* 1) Get reference price for market
				* 2) Get current open volume on Vega
				* 3) Build distribution of orders with relevant skew
				* 4) Submit to Vega using batch market instruction
				 */
				logging.GetLogger().Infof("update vega quotes for market: %s", config.VegaId)
			}
		}
	}()
}

func (b *Bot) Start() {
	b.initWallet()
	b.loadMarkets()
	b.syncVegaData()
	b.connectToVegaStreams()
	b.updateReferencePrices()
	b.updateLiquidityCommitment()
	b.updateQuotes()
}
