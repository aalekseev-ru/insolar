/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package packets

import (
	"github.com/insolar/insolar/core"
)

type ClaimType uint8

const (
	TypeNodeJoinClaim = ClaimType(iota + 1)
	TypeCapabilityPollingAndActivation
	TypeNodeViolationBlame
	TypeNodeBroadcast
	TypeNodeLeaveClaim
)

type ReferendumClaim interface {
	Serializer
	Type() ClaimType
	Length() uint16
}

// NodeBroadcast is a broadcast of info. Must be brief and only one entry per node.
// Type 4.
type NodeBroadcast struct {
	EmergencyLevel uint8
	length         uint16
}

func (nb *NodeBroadcast) Type() ClaimType {
	return TypeNodeBroadcast
}

func (nb *NodeBroadcast) Length() uint16 {
	return nb.length
}

// CapabilityPoolingAndActivation is a type 3.
type CapabilityPoolingAndActivation struct {
	PollingFlags   uint16
	CapabilityType uint16
	CapabilityRef  [ReferenceLength]byte
	length         uint16
}

func (cpa *CapabilityPoolingAndActivation) Type() ClaimType {
	return TypeCapabilityPollingAndActivation
}

func (cpa *CapabilityPoolingAndActivation) Length() uint16 {
	return cpa.length
}

// NodeViolationBlame is a type 2.
type NodeViolationBlame struct {
	BlameNodeID   uint32
	TypeViolation uint8
	claimType     ClaimType
	length        uint16
}

func (nvb *NodeViolationBlame) Type() ClaimType {
	return TypeNodeViolationBlame
}

func (nvb *NodeViolationBlame) Length() uint16 {
	return nvb.length
}

// NodeJoinClaim is a type 1, len == 272.
type NodeJoinClaim struct {
	NodeID                  uint32
	RelayNodeID             uint32
	ProtocolVersionAndFlags uint32
	JoinsAfter              uint32
	NodeRoleRecID           uint32
	NodeRef                 core.RecordRef
	NodePK                  [64]byte
	//length uint16
}

func (njc *NodeJoinClaim) Type() ClaimType {
	return TypeNodeJoinClaim
}

func (njc *NodeJoinClaim) Length() uint16 {
	return 0
}

// NodeLeaveClaim can be the only be issued by the node itself and must be the only claim record.
// Should be executed with the next pulse. Type 1, len == 0.
type NodeLeaveClaim struct {
	length uint16
}

func (nlc *NodeLeaveClaim) Type() ClaimType {
	return TypeNodeLeaveClaim
}

func (nlc *NodeLeaveClaim) Length() uint16 {
	return nlc.length
}

func NewNodeJoinClaim() *NodeJoinClaim {
	return &NodeJoinClaim{
		//length: 272,
	}
}

func NewNodViolationBlame() *NodeViolationBlame {
	return &NodeViolationBlame{
		claimType: TypeNodeViolationBlame,
	}
}