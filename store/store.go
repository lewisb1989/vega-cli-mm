package store

import (
	apipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	vegapb "code.vegaprotocol.io/vega/protos/vega"
	"fmt"
	"github.com/sasha-s/go-deadlock"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/maps"
)

type PriceSource string

const (
	Binance   PriceSource = "Binance"
	Pyth      PriceSource = "Pyth"
	Chainlink PriceSource = "Chainlink"
)

type KeyPair struct {
	PrivateKey string
	PublicKey  string
}

func NewKeyPair(privateKey string, publicKey string) *KeyPair {
	return &KeyPair{PrivateKey: privateKey, PublicKey: publicKey}
}

type MarketConfig struct {
	VegaId        string          `json:"vegaId"`
	ExternalId    string          `json:"externalId"`
	PriceSource   PriceSource     `json:"priceSource"`
	Spread        float64         `json:"spread"`
	ExposureLimit float64         `json:"exposureLimit"`
	LpRatio       float64         `json:"lpRatio"`
	KeyPair       *KeyPair        `json:"-"`
	BidPrice      decimal.Decimal `json:"-"`
	AskPrice      decimal.Decimal `json:"-"`
}

type Store struct {
	marketConfig            map[string]*MarketConfig
	accounts                map[string]*apipb.AccountBalance
	assets                  map[string]*vegapb.Asset
	positions               map[string]*vegapb.Position
	orders                  map[string]*vegapb.Order
	marketData              map[string]*vegapb.MarketData
	markets                 map[string]*vegapb.Market
	liquidityProvisions     map[string]*vegapb.LiquidityProvision
	networkParameters       map[string]*vegapb.NetworkParameter
	accountsLock            deadlock.RWMutex
	marketConfigLock        deadlock.RWMutex
	assetsLock              deadlock.RWMutex
	positionsLock           deadlock.RWMutex
	ordersLock              deadlock.RWMutex
	marketDataLock          deadlock.RWMutex
	marketsLock             deadlock.RWMutex
	liquidityProvisionsLock deadlock.RWMutex
	networkParametersLock   deadlock.RWMutex
}

func NewStore() *Store {
	return &Store{
		marketConfig:        map[string]*MarketConfig{},
		accounts:            map[string]*apipb.AccountBalance{},
		assets:              map[string]*vegapb.Asset{},
		positions:           map[string]*vegapb.Position{},
		orders:              map[string]*vegapb.Order{},
		marketData:          map[string]*vegapb.MarketData{},
		markets:             map[string]*vegapb.Market{},
		liquidityProvisions: map[string]*vegapb.LiquidityProvision{},
		networkParameters:   map[string]*vegapb.NetworkParameter{},
	}
}

func (s *Store) SaveMarketConfig(market *MarketConfig) {
	s.marketConfigLock.Lock()
	defer s.marketConfigLock.Unlock()
	s.marketConfig[market.VegaId] = market
}

func (s *Store) SaveMarketData(marketData *vegapb.MarketData) {
	s.marketsLock.Lock()
	defer s.marketDataLock.Unlock()
	s.marketData[marketData.Market] = marketData
}

func (s *Store) SaveMarket(market *vegapb.Market) {
	s.marketsLock.Lock()
	defer s.marketsLock.Unlock()
	s.markets[market.Id] = market
}

func (s *Store) SaveAsset(asset *vegapb.Asset) {
	s.assetsLock.Lock()
	defer s.assetsLock.Unlock()
	s.assets[asset.Id] = asset
}

func (s *Store) SaveAccount(account *apipb.AccountBalance) {
	s.accountsLock.Lock()
	defer s.accountsLock.Unlock()
	// TODO: Check if this is correct with Jeremy
	id := fmt.Sprintf("%s%s%s%s", account.Asset, account.Owner, account.Type, account.MarketId)
	s.accounts[id] = account
}

func (s *Store) SaveOrder(order *vegapb.Order) {
	s.ordersLock.Lock()
	defer s.ordersLock.Unlock()
	s.orders[order.Id] = order
}

func (s *Store) SavePosition(position *vegapb.Position) {
	s.positionsLock.Lock()
	defer s.positionsLock.Unlock()
	// TODO: Check if this is correct with Jeremy
	id := fmt.Sprintf("%s%s", position.PartyId, position.MarketId)
	s.positions[id] = position
}

func (s *Store) SaveLiquidityProvision(liquidityProvision *vegapb.LiquidityProvision) {
	s.liquidityProvisionsLock.Lock()
	defer s.liquidityProvisionsLock.Unlock()
	s.liquidityProvisions[liquidityProvision.Id] = liquidityProvision
}

func (s *Store) SaveNetworkParameter(networkParameter *vegapb.NetworkParameter) {
	s.networkParametersLock.Lock()
	defer s.networkParametersLock.Unlock()
	s.networkParameters[networkParameter.Key] = networkParameter
}

func (s *Store) GetMarketConfig() []*MarketConfig {
	s.marketConfigLock.RLock()
	defer s.marketConfigLock.RUnlock()
	return maps.Values(s.marketConfig)
}

func (s *Store) GetNetworkParameter(key string) *vegapb.NetworkParameter {
	s.networkParametersLock.RLock()
	defer s.networkParametersLock.RUnlock()
	if s.networkParameters[key] == nil {
		return nil
	}
	return &vegapb.NetworkParameter{
		Key:   s.networkParameters[key].Key,
		Value: s.networkParameters[key].Value,
	}
}
