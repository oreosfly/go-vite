package contracts

import (
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/monitor"
	"github.com/vitelabs/go-vite/vm_context/vmctxt_interface"
	"math/big"
	"time"
)

var (
	AddressRegister, _       = types.BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
	AddressVote, _           = types.BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2})
	AddressPledge, _         = types.BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3})
	AddressConsensusGroup, _ = types.BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4})
	AddressMintage, _        = types.BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5})
)

type StorageDatabase interface {
	GetStorage(addr *types.Address, key []byte) []byte
	NewStorageIterator(prefix []byte) vmctxt_interface.StorageIterator
}

func GetTokenById(db StorageDatabase, tokenId types.TokenTypeId) *TokenInfo {
	data := db.GetStorage(&AddressMintage, GetMintageKey(tokenId))
	if len(data) > 0 {
		tokenInfo := new(TokenInfo)
		ABIMintage.UnpackVariable(tokenInfo, VariableNameMintage, data)
		return tokenInfo
	}
	return nil
}

func GetTokenMap(db StorageDatabase) map[types.TokenTypeId]*TokenInfo {
	defer monitor.LogTime("vm", "GetTokenMap", time.Now())
	iterator := db.NewStorageIterator(nil)
	tokenInfoMap := make(map[types.TokenTypeId]*TokenInfo)
	for {
		key, value, ok := iterator.Next()
		if !ok {
			break
		}
		tokenId := GetTokenIdFromMintageKey(key)
		tokenInfo := new(TokenInfo)
		ABIMintage.UnpackVariable(tokenInfo, VariableNameMintage, value)
		tokenInfoMap[tokenId] = tokenInfo
	}
	return tokenInfoMap
}

func GetRegisterList(db StorageDatabase, gid types.Gid) []*Registration {
	defer monitor.LogTime("vm", "GetRegisterList", time.Now())
	iterator := db.NewStorageIterator(gid.Bytes())
	registerList := make([]*Registration, 0)
	for {
		_, value, ok := iterator.Next()
		if !ok {
			break
		}
		registration := new(Registration)
		ABIRegister.UnpackVariable(registration, VariableNameRegistration, value)
		if registration.IsActive() {
			registerList = append(registerList, registration)
		}
	}
	return registerList
}

func GetVoteList(db StorageDatabase, gid types.Gid) []*VoteInfo {
	defer monitor.LogTime("vm", "GetVoteList", time.Now())
	if gid == types.DELEGATE_GID {
		gid = types.SNAPSHOT_GID
	}
	iterator := db.NewStorageIterator(gid.Bytes())
	voteInfoList := make([]*VoteInfo, 0)
	for {
		key, value, ok := iterator.Next()
		if !ok {
			break
		}
		voterAddr := GetAddrFromVoteKey(key)
		nodeName := new(string)
		ABIVote.UnpackVariable(nodeName, VariableNameVoteStatus, value)
		voteInfoList = append(voteInfoList, &VoteInfo{voterAddr, *nodeName})
	}
	return voteInfoList
}

func GetPledgeBeneficialAmount(db StorageDatabase, beneficial types.Address) *big.Int {
	key := GetPledgeBeneficialKey(beneficial)
	beneficialAmount := new(VariablePledgeBeneficial)
	err := ABIPledge.UnpackVariable(beneficialAmount, VariableNamePledgeBeneficial, db.GetStorage(&AddressPledge, key))
	if err == nil {
		return beneficialAmount.Amount
	}
	return big.NewInt(0)
}

func GetPledgeAmount(db StorageDatabase, addr types.Address) []*PledgeInfo {
	iterator := db.NewStorageIterator(addr.Bytes())
	pledgeInfoList := make([]*PledgeInfo, 0)
	for {
		key, value, ok := iterator.Next()
		if !ok {
			break
		}
		if IsPledgeKey(key) {
			pledgeInfo := new(PledgeInfo)
			ABIPledge.UnpackVariable(pledgeInfo, VariableNamePledgeInfo, value)
			pledgeInfo.BeneficialAddr = GetBeneficialFromPledgeKey(key)
			pledgeInfoList = append(pledgeInfoList, pledgeInfo)
		}
	}
	return pledgeInfoList
}

func GetActiveConsensusGroupList(db StorageDatabase) []*ConsensusGroupInfo {
	defer monitor.LogTime("vm", "GetActiveConsensusGroupList", time.Now())
	iterator := db.NewStorageIterator(nil)
	consensusGroupInfoList := make([]*ConsensusGroupInfo, 0)
	for {
		key, value, ok := iterator.Next()
		if !ok {
			break
		}
		consensusGroupInfo := new(ConsensusGroupInfo)
		ABIConsensusGroup.UnpackVariable(consensusGroupInfo, VariableNameConsensusGroupInfo, value)
		if consensusGroupInfo.IsActive() {
			consensusGroupInfo.Gid = GetGidFromConsensusGroupKey(key)
			consensusGroupInfoList = append(consensusGroupInfoList, consensusGroupInfo)
		}
	}
	return consensusGroupInfoList
}

func GetConsensusGroup(db StorageDatabase, gid types.Gid) *ConsensusGroupInfo {
	data := db.GetStorage(&AddressConsensusGroup, GetConsensusGroupKey(gid))
	if len(data) > 0 {
		consensusGroupInfo := new(ConsensusGroupInfo)
		ABIConsensusGroup.UnpackVariable(consensusGroupInfo, VariableNameConsensusGroupInfo, data)
		consensusGroupInfo.Gid = gid
		return consensusGroupInfo
	}
	return nil
}