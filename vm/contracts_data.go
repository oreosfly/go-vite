package vm

import (
	"github.com/vitelabs/go-vite/common/types"
	"math/big"
)

func ToRegisterData(gid Gid) []byte {
	return joinBytes(DataRegister, LeftPadBytes(gid.Bytes(), 32))
}
func ToCancelRegisterData(gid Gid) []byte {
	return joinBytes(DataCancelRegister, LeftPadBytes(gid.Bytes(), 32))
}
func ToRewardData(gid Gid, rewardHeightEnd *big.Int) []byte {
	if rewardHeightEnd == nil {
		return joinBytes(DataReward, LeftPadBytes(gid.Bytes(), 32))
	} else {
		return joinBytes(DataReward, LeftPadBytes(gid.Bytes(), 32), LeftPadBytes(rewardHeightEnd.Bytes(), 32))
	}
}
func ToVoteData(gid Gid, addr types.Address) []byte {
	return joinBytes(DataVote, LeftPadBytes(gid.Bytes(), 32), LeftPadBytes(addr.Bytes(), 32))
}
func ToCancelVoteData(gid Gid) []byte {
	return joinBytes(DataCancelVote, LeftPadBytes(gid.Bytes(), 32))
}
func ToMortgageData(beneficial types.Address, withdrawTime int64) []byte {
	return joinBytes(DataMortgage, LeftPadBytes(beneficial.Bytes(), 32), LeftPadBytes(new(big.Int).SetInt64(withdrawTime).Bytes(), 32))
}
func ToCancelMortgageData(beneficial types.Address, amount *big.Int) []byte {
	return joinBytes(DataCancelMortgage, LeftPadBytes(beneficial.Bytes(), 32), LeftPadBytes(amount.Bytes(), 32))
}
func ToCreateConsensusGroupData(gid Gid, group ConsensusGroup) []byte {
	tmp := new(big.Int)
	countingRuleParamSize := ((len(group.CountingRuleParam) + 31) / 32) * 32
	registerConditionParamSize := ((len(group.RegisterConditionParam) + 31) / 32) * 32
	VoteConditionParamSize := ((len(group.VoteConditionParam) + 31) / 32) * 32
	return joinBytes(DataCreateConsensusGroup,
		LeftPadBytes(gid.Bytes(), 32),
		LeftPadBytes(tmp.SetUint64(uint64(group.NodeCount)).Bytes(), 32),
		LeftPadBytes(tmp.SetInt64(group.Interval).Bytes(), 32),
		LeftPadBytes(tmp.SetUint64(uint64(group.CountingRuleId)).Bytes(), 32),
		LeftPadBytes(tmp.SetUint64(288).Bytes(), 32),
		LeftPadBytes(tmp.SetUint64(uint64(group.RegisterConditionId)).Bytes(), 32),
		LeftPadBytes(tmp.SetUint64(288+32+uint64(countingRuleParamSize)).Bytes(), 32),
		LeftPadBytes(tmp.SetUint64(uint64(group.VoteConditionId)).Bytes(), 32),
		LeftPadBytes(tmp.SetUint64(288+32+uint64(countingRuleParamSize)+32+uint64(registerConditionParamSize)).Bytes(), 32),
		LeftPadBytes(tmp.SetUint64(uint64(len(group.CountingRuleParam))).Bytes(), 32),
		RightPadBytes(group.CountingRuleParam, countingRuleParamSize),
		LeftPadBytes(tmp.SetUint64(uint64(len(group.RegisterConditionParam))).Bytes(), 32),
		RightPadBytes(group.RegisterConditionParam, registerConditionParamSize),
		LeftPadBytes(tmp.SetUint64(uint64(len(group.VoteConditionParam))).Bytes(), 32),
		RightPadBytes(group.VoteConditionParam, VoteConditionParamSize),
	)
}

type ConsensusGroup struct {
	NodeCount              uint8
	Interval               int64
	CountingRuleId         uint8
	CountingRuleParam      []byte
	RegisterConditionId    uint8
	RegisterConditionParam []byte
	VoteConditionId        uint8
	VoteConditionParam     []byte
}

type CountingRuleCode string

const (
	CountingRuleOfBalance       CountingRuleCode = "counting1"
	RegisterConditionOfSnapshot CountingRuleCode = "register1"
	VoteConditionOfDefault      CountingRuleCode = "vote1"
	VoteConditionOfBalance      CountingRuleCode = "vote2"
)

type createConsensusGroupCondition interface {
	checkParam(param []byte, db VmDatabase) bool
}

var SimpleCountingRuleList = map[CountingRuleCode]createConsensusGroupCondition{
	CountingRuleOfBalance:       &countingRuleOfBalance{},
	RegisterConditionOfSnapshot: &registerConditionOfSnapshot{},
	VoteConditionOfDefault:      &voteConditionOfDefault{},
	VoteConditionOfBalance:      &voteConditionOfBalance{},
}

type countingRuleOfBalance struct{}

func (c countingRuleOfBalance) checkParam(param []byte, db VmDatabase) bool {
	if len(param) != 32 {
		return false
	}
	if tokenId, err := types.BytesToTokenTypeId(LeftPadBytes(new(big.Int).SetBytes(param).Bytes(), 10)); err != nil || !db.IsExistToken(tokenId) {
		return false
	}
	return true
}

type registerConditionOfSnapshot struct{}

func (c registerConditionOfSnapshot) checkParam(param []byte, db VmDatabase) bool {
	if len(param) != 96 {
		return false
	}
	if tokenId, err := types.BytesToTokenTypeId(LeftPadBytes(new(big.Int).SetBytes(param[32:64]).Bytes(), 10)); err != nil || !db.IsExistToken(tokenId) {
		return false
	}
	return true
}

type voteConditionOfDefault struct{}

func (c voteConditionOfDefault) checkParam(param []byte, db VmDatabase) bool {
	if len(param) != 0 {
		return false
	}
	return true
}

type voteConditionOfBalance struct{}

func (c voteConditionOfBalance) checkParam(param []byte, db VmDatabase) bool {
	if len(param) != 64 {
		return false
	}
	if tokenId, err := types.BytesToTokenTypeId(LeftPadBytes(new(big.Int).SetBytes(param[32:64]).Bytes(), 10)); err != nil || !db.IsExistToken(tokenId) {
		return false
	}
	return true
}
