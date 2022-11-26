package rln

import (
	"bytes"
	"encoding/hex"
	"math"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestRLNSuite(t *testing.T) {
	suite.Run(t, new(RLNSuite))
}

type RLNSuite struct {
	suite.Suite
}

func (s *RLNSuite) TestMembershipKeyGen() {
	rln, err := NewRLN()
	s.NoError(err)

	key, err := rln.MembershipKeyGen()
	s.NoError(err)
	s.Len(key.IDKey, 32)
	s.Len(key.IDCommitment, 32)
	s.NotEmpty(key.IDKey)
	s.NotEmpty(key.IDCommitment)
	s.False(bytes.Equal(key.IDCommitment[:], make([]byte, 32)))
	s.False(bytes.Equal(key.IDKey[:], make([]byte, 32)))
}

func (s *RLNSuite) TestGetMerkleRoot() {
	rln, err := NewRLN()
	s.NoError(err)

	root1, err := rln.GetMerkleRoot()
	s.NoError(err)
	s.Len(root1, 32)

	root2, err := rln.GetMerkleRoot()
	s.NoError(err)
	s.Len(root2, 32)

	s.Equal(root1, root2)
}

func (s *RLNSuite) TestInsertMember() {
	rln, err := NewRLN()
	s.NoError(err)

	keypair, err := rln.MembershipKeyGen()
	s.NoError(err)

	err = rln.InsertMember(keypair.IDCommitment)
	s.NoError(err)
}

func (s *RLNSuite) TestInsertMembers() {
	rln, err := NewRLN()
	s.NoError(err)

	keypair, err := rln.MembershipKeyGen()
	s.NoError(err)

	err = rln.InsertMembers(0, []IDCommitment{keypair.IDCommitment})
	s.NoError(err)
}

func (s *RLNSuite) TestRemoveMember() {
	rln, err := NewRLN()
	s.NoError(err)

	err = rln.DeleteMember(MembershipIndex(0))
	s.NoError(err)
}

func (s *RLNSuite) TestMerkleTreeConsistenceBetweenDeletionAndInsertion() {
	rln, err := NewRLN()
	s.NoError(err)

	root1, err := rln.GetMerkleRoot()
	s.NoError(err)
	s.Len(root1, 32)

	keypair, err := rln.MembershipKeyGen()
	s.NoError(err)

	err = rln.InsertMember(keypair.IDCommitment)
	s.NoError(err)

	// read the Merkle Tree root after insertion
	root2, err := rln.GetMerkleRoot()
	s.NoError(err)
	s.Len(root2, 32)

	// delete the first member
	deleted_member_index := MembershipIndex(0)
	err = rln.DeleteMember(deleted_member_index)
	s.NoError(err)

	// read the Merkle Tree root after the deletion
	root3, err := rln.GetMerkleRoot()
	s.NoError(err)
	s.Len(root3, 32)

	// the root must change after the insertion
	s.NotEqual(root1, root2)

	// The initial root of the tree (empty tree) must be identical to
	// the root of the tree after one insertion followed by a deletion
	s.Equal(root1, root3)
}

func (s *RLNSuite) TestHash() {
	rln, err := NewRLN()
	s.NoError(err)

	// prepare the input
	msg := []byte("Hello")
	hash, err := rln.Hash(msg)
	s.NoError(err)

	expectedHash, _ := hex.DecodeString("4c6ea217404bd5f10e243bac29dc4f1ec36bf4a41caba7b4c8075c54abb3321e")
	s.Equal(expectedHash, hash[:])
}

func (s *RLNSuite) TestCreateListMembershipKeysAndCreateMerkleTreeFromList() {
	groupSize := 100
	list, root, err := CreateMembershipList(groupSize)
	s.NoError(err)
	s.Len(list, groupSize)
	s.Len(root, HASH_HEX_SIZE) // check the size of the calculated tree root
}

func (s *RLNSuite) TestCheckCorrectness() {
	groupKeys := STATIC_GROUP_KEYS

	// create a set of MembershipKeyPair objects from groupKeys
	groupKeyPairs, err := toMembershipKeyPairs(groupKeys)
	s.NoError(err)

	// extract the id commitments
	var groupIDCommitments []IDCommitment
	for _, c := range groupKeyPairs {
		groupIDCommitments = append(groupIDCommitments, c.IDCommitment)
	}

	// calculate the Merkle tree root out of the extracted id commitments
	root, err := CalcMerkleRoot(groupIDCommitments)
	s.NoError(err)

	expectedRoot, _ := hex.DecodeString(STATIC_GROUP_MERKLE_ROOT)

	s.Len(groupKeyPairs, STATIC_GROUP_SIZE)
	s.Equal(expectedRoot, root[:])
}

func (s *RLNSuite) TestValidProof() {
	rln, err := NewRLN()
	s.NoError(err)

	memKeys, err := rln.MembershipKeyGen()
	s.NoError(err)

	//peer's index in the Merkle Tree
	index := 5

	// Create a Merkle tree with random members
	for i := 0; i < 10; i++ {
		if i == index {
			// insert the current peer's pk
			err = rln.InsertMember(memKeys.IDCommitment)
			s.NoError(err)
		} else {
			// create a new key pair
			memberKeys, err := rln.MembershipKeyGen()
			s.NoError(err)

			err = rln.InsertMember(memberKeys.IDCommitment)
			s.NoError(err)
		}
	}

	// prepare the message
	msg := []byte("Hello")

	// prepare the epoch
	var epoch Epoch

	// generate proof
	proofRes, err := rln.GenerateProof(msg, *memKeys, MembershipIndex(index), epoch)
	s.NoError(err)

	// verify the proof
	verified, err := rln.Verify(msg, *proofRes)
	s.NoError(err)
	s.True(verified)

	// verify with roots
	root, err := rln.GetMerkleRoot()
	s.NoError(err)

	verified, err = rln.VerifyWithRoots(msg, *proofRes, [][32]byte{root})
	s.NoError(err)
	s.True(verified)
}

func (s *RLNSuite) TestInvalidProof() {
	rln, err := NewRLN()
	s.NoError(err)

	memKeys, err := rln.MembershipKeyGen()
	s.NoError(err)

	//peer's index in the Merkle Tree
	index := 5

	// Create a Merkle tree with random members
	for i := 0; i < 10; i++ {
		if i == index {
			// insert the current peer's pk
			err := rln.InsertMember(memKeys.IDCommitment)
			s.NoError(err)
		} else {
			// create a new key pair
			memberKeys, err := rln.MembershipKeyGen()
			s.NoError(err)

			err = rln.InsertMember(memberKeys.IDCommitment)
			s.NoError(err)
		}
	}

	// prepare the message
	msg := []byte("Hello")

	// prepare the epoch
	var epoch Epoch

	badIndex := 4

	// generate proof
	proofRes, err := rln.GenerateProof(msg, *memKeys, MembershipIndex(badIndex), epoch)
	s.NoError(err)

	// verify the proof (should not be verified)
	verified, err := rln.Verify(msg, *proofRes)
	s.NoError(err)
	s.False(verified)
}

func (s *RLNSuite) TestEpochConsistency() {
	// check edge cases
	var epoch uint64 = math.MaxUint64
	epochBytes := ToEpoch(epoch)
	decodedEpoch := epochBytes.Uint64()

	s.Equal(epoch, decodedEpoch)
}

func (s *RLNSuite) TestEpochComparison() {
	// check edge cases
	var time1 uint64 = math.MaxUint64
	var time2 uint64 = math.MaxUint64 - 1

	epoch1 := ToEpoch(time1)
	epoch2 := ToEpoch(time2)

	s.Equal(int64(1), Diff(epoch1, epoch2))
	s.Equal(int64(-1), Diff(epoch2, epoch1))
}
