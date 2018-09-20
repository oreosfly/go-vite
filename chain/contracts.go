package chain

import (
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/contracts"
	"github.com/vitelabs/go-vite/vm_context"
	"math/big"
)

func (c *Chain) GetContractGid(addr *types.Address) (*types.Gid, error) {
	gid, err := c.chainDb.Ac.GetContractGid(addr)
	if err != nil {
		c.log.Error("GetContractGid failed, error is "+err.Error(), "method", "GetContractGid")
		return nil, err
	}

	return gid, nil
}

func (c *Chain) GetRegisterList(snapshotHash types.Hash, gid types.Gid) []*contracts.Registration {
	vmContext, err := vm_context.NewVmContext(c, &snapshotHash, nil, &contracts.AddressRegister)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetRegisterList")
		return nil
	}
	return contracts.GetRegisterList(vmContext, gid)
}

func (c *Chain) GetVoteMap(snapshotHash types.Hash, gid types.Gid) []*contracts.VoteInfo {
	vmContext, err := vm_context.NewVmContext(c, &snapshotHash, nil, &contracts.AddressRegister)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetVoteList")
		return nil
	}
	return contracts.GetVoteList(vmContext, gid)
}

func (c *Chain) GetPledgeAmount(snapshotHash types.Hash, beneficial types.Address) *big.Int {
	vmContext, err := vm_context.NewVmContext(c, &snapshotHash, nil, &contracts.AddressRegister)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetPledgeAmount")
		return nil
	}
	return contracts.GetPledgeAmount(vmContext, beneficial)
}

func (c *Chain) GetConsensusGroupList(snapshotHash types.Hash) []*contracts.ConsensusGroupInfo {
	vmContext, err := vm_context.NewVmContext(c, &snapshotHash, nil, &contracts.AddressRegister)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetConsensusGroupList")
		return nil
	}
	return contracts.GetConsensusGroupList(vmContext)
}

// TODO
func (c *Chain) GetBalanceList(snapshotHash types.Hash, tokenTypeId types.TokenTypeId, addressList []types.Address) map[types.Address]*big.Int {
	return nil
}

// TODO
func (c *Chain) GetTokenInfoById(tokenId *types.TokenTypeId) (*contracts.TokenInfo, error) {
	return nil, nil
}
