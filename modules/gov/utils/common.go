package utils

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/desmos-labs/juno/client"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	stakingutils "github.com/forbole/bdjuno/modules/staking/utils"

	"google.golang.org/grpc/codes"

	"github.com/forbole/bdjuno/database"
	authutils "github.com/forbole/bdjuno/modules/auth/utils"
	bankutils "github.com/forbole/bdjuno/modules/bank/utils"
	"github.com/forbole/bdjuno/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ErrProposalNotFound = "rpc error: code = %s desc = rpc error: code = %s desc = proposal %d doesn't exist: key not found"
)

func UpdateProposal(
	height int64, blockVals *tmctypes.ResultValidators, id uint64,
	govClient govtypes.QueryClient, bankClient banktypes.QueryClient, stakingClient stakingtypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	// Get the proposal
	res, err := govClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: id})
	if err != nil {
		// Get the error code
		var code string
		_, err := fmt.Sscanf(err.Error(), ErrProposalNotFound, &code, &code, &id)
		if err != nil {
			return err
		}

		if code == codes.NotFound.String() {
			// Handle case when a proposal is deleted from the chain (did not pass deposit period)
			return updateDeletedProposalStatus(id, db)
		}

		return fmt.Errorf("error while getting proposal: %s", err)
	}

	err = updateProposalStatus(res.Proposal, db)
	if err != nil {
		return fmt.Errorf("error while updating proposal status: %s", err)
	}

	err = updateProposalTallyResult(res.Proposal, govClient, db)
	if err != nil {
		return fmt.Errorf("error while updating proposal tally result: %s", err)
	}

	err = updateAccounts(res.Proposal, bankClient, db)
	if err != nil {
		return fmt.Errorf("error while updating account: %s", err)
	}

	err = updateProposalStakingPoolSnapshot(height, id, stakingClient, db)
	if err != nil {
		return fmt.Errorf("error while updating proposal staking pool snapshot: %s", err)
	}

	err = updateProposalValidatorStatusesSnapshot(height, id, blockVals, stakingClient, cdc, db)
	if err != nil {
		return fmt.Errorf("error while updating proposal validator statuses snapshot: %s", err)
	}

	return nil
}

// updateDeletedProposalStatus updates the proposal having the given id by setting its status
// to the one that represents a deleted proposal
func updateDeletedProposalStatus(id uint64, db *database.Db) error {
	stored, err := db.GetProposal(id)
	if err != nil {
		return err
	}

	return db.UpdateProposal(
		types.NewProposalUpdate(
			stored.ProposalID,
			types.ProposalStatusInvalid,
			stored.VotingStartTime,
			stored.VotingEndTime,
		),
	)
}

// updateProposalStatus updates the given proposal status
func updateProposalStatus(proposal govtypes.Proposal, db *database.Db) error {
	return db.UpdateProposal(
		types.NewProposalUpdate(
			proposal.ProposalId,
			proposal.Status.String(),
			proposal.VotingStartTime,
			proposal.VotingEndTime,
		),
	)
}

// updateProposalTallyResult updates the tally result associated with the given proposal
func updateProposalTallyResult(proposal govtypes.Proposal, govClient govtypes.QueryClient, db *database.Db) error {
	height, err := db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	header := client.GetHeightRequestHeader(height)
	res, err := govClient.TallyResult(
		context.Background(),
		&govtypes.QueryTallyResultRequest{ProposalId: proposal.ProposalId},
		header,
	)
	if err != nil {
		return fmt.Errorf("error while getting tally result: %s", err)
	}

	return db.SaveTallyResults([]types.TallyResult{
		types.NewTallyResult(
			proposal.ProposalId,
			res.Tally.Yes.Int64(),
			res.Tally.Abstain.Int64(),
			res.Tally.No.Int64(),
			res.Tally.NoWithVeto.Int64(),
			height,
		),
	})
}

// updateAccounts updates any account that might be involved in the proposal (eg. fund community recipient)
func updateAccounts(proposal govtypes.Proposal, bankClient banktypes.QueryClient, db *database.Db) error {
	content, ok := proposal.Content.GetCachedValue().(*distrtypes.CommunityPoolSpendProposal)
	if ok {
		height, err := db.GetLastBlockHeight()
		if err != nil {
			return fmt.Errorf("error while getting last block height: %s", err)
		}

		addresses := []string{content.Recipient}

		err = authutils.UpdateAccounts(addresses, db)
		if err != nil {
			return err
		}

		return bankutils.RefreshBalances(height, addresses, bankClient, db)
	}
	return nil
}

// updateProposalStakingPoolSnapshot updates the staking pool snapshot associated with the gov
// proposal having the provided id
func updateProposalStakingPoolSnapshot(
	height int64, proposalID uint64, stakingClient stakingtypes.QueryClient, db *database.Db,
) error {
	pool, err := stakingutils.GetStakingPool(height, stakingClient)
	if err != nil {
		return fmt.Errorf("error while getting staking pool: %s", err)
	}

	return db.SaveProposalStakingPoolSnapshot(
		types.NewProposalStakingPoolSnapshot(proposalID, pool),
	)
}

// updateProposalValidatorStatusesSnapshot updates the snapshots of the various validators for
// the proposal having the given id
func updateProposalValidatorStatusesSnapshot(
	height int64, proposalID uint64,
	blockVals *tmctypes.ResultValidators, stakingClient stakingtypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	validators, _, err := stakingutils.GetValidatorsWithStatus(height, stakingtypes.Bonded.String(), stakingClient, cdc)
	if err != nil {
		return fmt.Errorf("error while getting validators with bonded status: %s", err)
	}

	votingPowers := stakingutils.GetValidatorsVotingPowers(height, blockVals, db)

	statuses, err := stakingutils.GetValidatorsStatuses(height, validators, cdc)
	if err != nil {
		return fmt.Errorf("error while getting validator statuses: %s", err)
	}

	var snapshots = make([]types.ProposalValidatorStatusSnapshot, len(validators))
	for index, validator := range validators {
		consAddr, err := validator.GetConsAddr()
		if err != nil {
			return err
		}

		status, err := findStatus(consAddr.String(), statuses)
		if err != nil {
			return fmt.Errorf("error while searching for status: %s", err)
		}

		votingPower, err := findVotingPower(consAddr.String(), votingPowers)
		if err != nil {
			return fmt.Errorf("error while searching for voting power: %s", err)
		}

		snapshots[index] = types.NewProposalValidatorStatusSnapshot(
			proposalID,
			consAddr.String(),
			votingPower.VotingPower,
			status.Status,
			status.Jailed,
			height,
		)
	}

	return db.SaveProposalValidatorsStatusesSnapshots(snapshots)
}

func findVotingPower(consAddr string, powers []types.ValidatorVotingPower) (types.ValidatorVotingPower, error) {
	for _, votingPower := range powers {
		if votingPower.ConsensusAddress == consAddr {
			return votingPower, nil
		}
	}
	return types.ValidatorVotingPower{}, fmt.Errorf("voting power not found for validator with consensus address %s", consAddr)
}

func findStatus(consAddr string, statuses []types.ValidatorStatus) (types.ValidatorStatus, error) {
	for _, status := range statuses {
		if status.ConsensusAddress == consAddr {
			return status, nil
		}
	}
	return types.ValidatorStatus{}, fmt.Errorf("cannot find status for validator with consensus address %s", consAddr)
}
