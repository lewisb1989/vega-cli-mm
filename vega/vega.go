package vega

import (
	"code.vegaprotocol.io/vega/libs/ptr"
	apipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	vegapb "code.vegaprotocol.io/vega/protos/vega"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"vega-mm/auth"
	"vega-mm/logging"
	"vega-mm/store"
)

type Vega struct {
	authenticator                *auth.Authenticator
	store                        *store.Store
	coreNode                     string
	accountsConnected            map[string]bool
	marketDataConnected          bool
	ordersConnected              bool
	positionsConnected           map[string]bool
	liquidityProvisionsConnected map[string]bool
}

func NewVega(
	store *store.Store,
	coreNode string,
) *Vega {
	return &Vega{
		store:                        store,
		coreNode:                     coreNode,
		accountsConnected:            map[string]bool{},
		liquidityProvisionsConnected: map[string]bool{},
		positionsConnected:           map[string]bool{},
	}
}

func (v *Vega) GetAuthenticator() *auth.Authenticator {
	return v.authenticator
}

func (v *Vega) SetAuthenticator(authenticator *auth.Authenticator) {
	v.authenticator = authenticator
}

func (v *Vega) GetCoreNode() string {
	return v.coreNode
}

func (v *Vega) IsAccountsConnected(partyId string) bool {
	return v.accountsConnected[partyId]
}

func (v *Vega) IsMarketDataConnected() bool {
	return v.marketDataConnected
}

func (v *Vega) IsOrdersConnected() bool {
	return v.ordersConnected
}

func (v *Vega) IsPositionsConnected(partyId string) bool {
	return v.positionsConnected[partyId]
}

func (v *Vega) IsLiquidityProvisionsConnected(partyId string) bool {
	return v.liquidityProvisionsConnected[partyId]
}

func (v *Vega) GetAssets() []*vegapb.Asset {
	assets := make([]*vegapb.Asset, 0)
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not list assets: %v", err)
		return assets
	}
	req := &apipb.ListAssetsRequest{}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	resp, err := tradingDataService.ListAssets(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not list assets: %v", err)
		return assets
	}
	for _, edge := range resp.Assets.Edges {
		assets = append(assets, edge.Node)
	}
	return assets
}

func (v *Vega) GetMarkets() []*vegapb.Market {
	markets := make([]*vegapb.Market, 0)
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not list markets: %v", err)
		return markets
	}
	req := &apipb.ListMarketsRequest{}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	resp, err := tradingDataService.ListMarkets(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not list markets: %v", err)
		return markets
	}
	for _, edge := range resp.Markets.Edges {
		markets = append(markets, edge.Node)
	}
	return markets
}

func (v *Vega) GetMarketData() []*vegapb.MarketData {
	marketData := make([]*vegapb.MarketData, 0)
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not list market data: %v", err)
		return marketData
	}
	req := &apipb.ListLatestMarketDataRequest{}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	resp, err := tradingDataService.ListLatestMarketData(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not list market data: %v", err)
		return marketData
	}
	for _, data := range resp.MarketsData {
		marketData = append(marketData, data)
	}
	return marketData
}

func (v *Vega) GetAccounts(
	partyIds []string,
) []*apipb.AccountBalance {
	accounts := make([]*apipb.AccountBalance, 0)
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not list accounts: %v", err)
		return accounts
	}
	req := &apipb.ListAccountsRequest{Filter: &apipb.AccountFilter{PartyIds: partyIds}}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	resp, err := tradingDataService.ListAccounts(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not list accounts: %v", err)
		return accounts
	}
	for _, edge := range resp.Accounts.Edges {
		accounts = append(accounts, edge.Node)
	}
	return accounts
}

func (v *Vega) GetOrders(
	partyIds []string,
) []*vegapb.Order {
	orders := make([]*vegapb.Order, 0)
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not list orders: %v", err)
		return orders
	}
	req := &apipb.ListOrdersRequest{Filter: &apipb.OrderFilter{PartyIds: partyIds, LiveOnly: ptr.From(true)}}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	resp, err := tradingDataService.ListOrders(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not list orders: %v", err)
		return orders
	}
	for _, edge := range resp.Orders.Edges {
		orders = append(orders, edge.Node)
	}
	return orders
}

func (v *Vega) GetPositions(
	partyIds []string,
) []*vegapb.Position {
	positions := make([]*vegapb.Position, 0)
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not list positions: %v", err)
		return positions
	}
	req := &apipb.ListAllPositionsRequest{Filter: &apipb.PositionsFilter{PartyIds: partyIds}}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	resp, err := tradingDataService.ListAllPositions(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not list positions: %v", err)
		return positions
	}
	for _, edge := range resp.Positions.Edges {
		positions = append(positions, edge.Node)
	}
	return positions
}

func (v *Vega) GetNetworkParameters() []*vegapb.NetworkParameter {
	networkParameters := make([]*vegapb.NetworkParameter, 0)
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not list network parameters: %v", err)
		return networkParameters
	}
	req := &apipb.ListNetworkParametersRequest{}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	resp, err := tradingDataService.ListNetworkParameters(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not list network parameters: %v", err)
		return networkParameters
	}
	for _, edge := range resp.NetworkParameters.Edges {
		networkParameters = append(networkParameters, edge.Node)
	}
	return networkParameters
}

func (v *Vega) GetLiquidityProvisions(
	partyId string,
) []*vegapb.LiquidityProvision {
	liquidityProvisions := make([]*vegapb.LiquidityProvision, 0)
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not list liquidity provisions: %v", err)
		return liquidityProvisions
	}
	req := &apipb.ListLiquidityProvisionsRequest{PartyId: ptr.From(partyId)}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	resp, err := tradingDataService.ListLiquidityProvisions(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not list liquidity provisions: %v", err)
		return liquidityProvisions
	}
	for _, edge := range resp.LiquidityProvisions.Edges {
		liquidityProvisions = append(liquidityProvisions, edge.Node)
	}
	return liquidityProvisions
}

func (v *Vega) StreamMarketData(
	marketIds []string,
	callback func(marketData []*vegapb.MarketData),
) {
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not start market data stream: %v", err)
		return
	}
	req := &apipb.ObserveMarketsDataRequest{MarketIds: marketIds}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	stream, err := tradingDataService.ObserveMarketsData(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not start market data stream: %v", err)
		return
	}
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				logging.GetLogger().Warnf("could not recv market data: %v", err)
				v.marketDataConnected = false
				break
			} else {
				v.marketDataConnected = true
				callback(resp.MarketData)
			}
		}
	}()
}

func (v *Vega) StreamOrders(
	partyIds []string,
	callback func(orders []*vegapb.Order),
) {
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not start orders stream: %v", err)
		return
	}
	req := &apipb.ObserveOrdersRequest{PartyIds: partyIds}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	stream, err := tradingDataService.ObserveOrders(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not start orders stream: %v", err)
		return
	}
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				logging.GetLogger().Warnf("could not recv orders: %v", err)
				v.ordersConnected = false
				break
			} else {
				v.ordersConnected = true
				switch r := resp.Response.(type) {
				case *apipb.ObserveOrdersResponse_Snapshot:
					callback(r.Snapshot.Orders)
				case *apipb.ObserveOrdersResponse_Updates:
					callback(r.Updates.Orders)
				}
			}
		}
	}()
}

func (v *Vega) StreamPositions(
	partyId string,
	callback func(positions []*vegapb.Position),
) {
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not start positions stream: %v", err)
		return
	}
	req := &apipb.ObservePositionsRequest{PartyId: ptr.From(partyId)}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	stream, err := tradingDataService.ObservePositions(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not start positions stream: %v", err)
		return
	}
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				logging.GetLogger().Warnf("could not recv positions: %v", err)
				v.positionsConnected[partyId] = false
				break
			} else {
				v.positionsConnected[partyId] = true
				switch r := resp.Response.(type) {
				case *apipb.ObservePositionsResponse_Snapshot:
					callback(r.Snapshot.Positions)
				case *apipb.ObservePositionsResponse_Updates:
					callback(r.Updates.Positions)
				}
			}
		}
	}()
}

func (v *Vega) StreamLiquidityProvisions(
	partyId string,
	callback func(liquidityProvisions []*vegapb.LiquidityProvision),
) {
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not start liquidity provisions stream: %v", err)
		return
	}
	req := &apipb.ObserveLiquidityProvisionsRequest{PartyId: ptr.From(partyId)}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	stream, err := tradingDataService.ObserveLiquidityProvisions(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not start liquidity provisions stream: %v", err)
		return
	}
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				logging.GetLogger().Warnf("could not recv liquidity provisions: %v", err)
				v.liquidityProvisionsConnected[partyId] = false
				break
			} else {
				v.liquidityProvisionsConnected[partyId] = true
				callback(resp.LiquidityProvisions)
			}
		}
	}()
}

func (v *Vega) StreamAccounts(
	partyId string,
	callback func(accounts []*apipb.AccountBalance),
) {
	node, err := grpc.Dial(v.coreNode, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger().Warnf("could not start accounts stream: %v", err)
		return
	}
	req := &apipb.ObserveAccountsRequest{PartyId: partyId}
	tradingDataService := apipb.NewTradingDataServiceClient(node)
	stream, err := tradingDataService.ObserveAccounts(context.Background(), req)
	if err != nil {
		logging.GetLogger().Warnf("could not start accounts stream: %v", err)
		return
	}
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				logging.GetLogger().Warnf("could not recv accounts: %v", err)
				v.accountsConnected[partyId] = false
				break
			} else {
				v.accountsConnected[partyId] = true
				switch r := resp.Response.(type) {
				case *apipb.ObserveAccountsResponse_Snapshot:
					callback(r.Snapshot.Accounts)
				case *apipb.ObserveAccountsResponse_Updates:
					callback(r.Updates.Accounts)
				}
			}
		}
	}()
}

func (v *Vega) SubmitBatchMarketInstruction() {
	// TODO - submit batch market instruction
}
