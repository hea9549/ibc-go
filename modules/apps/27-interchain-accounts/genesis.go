package interchain_accounts

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ibc-go/v2/modules/apps/27-interchain-accounts/keeper"
	"github.com/cosmos/ibc-go/v2/modules/apps/27-interchain-accounts/types"
	host "github.com/cosmos/ibc-go/v2/modules/core/24-host"
)

// InitGenesis initializes the interchain accounts application state from a provided genesis state
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, state types.GenesisState) {
	for _, portID := range state.Ports {
		if !keeper.IsBound(ctx, portID) {
			cap := keeper.BindPort(ctx, portID)
			if err := keeper.ClaimCapability(ctx, cap, host.PortPath(portID)); err != nil {
				panic(fmt.Sprintf("could not claim port capability: %v", err))
			}
		}
	}

	for _, ch := range state.ActiveChannels {
		keeper.SetActiveChannelID(ctx, ch.PortId, ch.ChannelId)
	}

	for _, acc := range state.InterchainAccounts {
		keeper.SetInterchainAccountAddress(ctx, acc.PortId, acc.AccountAddress)
	}
}

// ExportGenesis returns the interchain accounts exported genesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		keeper.GetAllPorts(ctx),
		keeper.GetAllActiveChannels(ctx),
		keeper.GetAllInterchainAccounts(ctx),
	)
}