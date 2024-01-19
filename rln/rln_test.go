package rln

import (
	"bytes"
	"encoding/hex"
	"fmt"
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

func (s *RLNSuite) TestNew() {
	rln, err := NewRLN()
	s.NoError(err)

	root1, err := rln.GetMerkleRoot()
	s.NoError(err)
	s.Len(root1, 32)

	rln2, err := NewWithConfig(DefaultTreeDepth, nil)
	s.NoError(err)

	root2, err := rln2.GetMerkleRoot()
	s.NoError(err)
	s.Len(root2, 32)
	s.Equal(root1, root2)
}

func (s *RLNSuite) TestMembershipKeyGen() {
	rln, err := NewRLN()
	s.NoError(err)

	key, err := rln.MembershipKeyGen()
	s.NoError(err)
	s.Len(key.IDSecretHash, 32)
	s.Len(key.IDCommitment, 32)
	s.Len(key.IDTrapdoor, 32)
	s.Len(key.IDNullifier, 32)
	s.NotEmpty(key.IDSecretHash)
	s.NotEmpty(key.IDCommitment)
	s.NotEmpty(key.IDTrapdoor)
	s.NotEmpty(key.IDNullifier)
	s.False(bytes.Equal(key.IDCommitment[:], make([]byte, 32)))
	s.False(bytes.Equal(key.IDSecretHash[:], make([]byte, 32)))
	s.False(bytes.Equal(key.IDTrapdoor[:], make([]byte, 32)))
	s.False(bytes.Equal(key.IDNullifier[:], make([]byte, 32)))
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

	var commitments []IDCommitment
	for i := 0; i < 10; i++ {
		keypair, err := rln.MembershipKeyGen()
		s.NoError(err)
		commitments = append(commitments, keypair.IDCommitment)
	}

	err = rln.InsertMembers(0, commitments)
	s.NoError(err)

	numLeaves := rln.LeavesSet()
	s.Equal(uint(10), numLeaves)
}

func (s *RLNSuite) TestRemoveMember() {
	rln, err := NewRLN()
	s.NoError(err)

	keypair, err := rln.MembershipKeyGen()
	s.NoError(err)

	err = rln.InsertMember(keypair.IDCommitment)
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
	hash, err := rln.Sha256(msg)
	s.NoError(err)

	expectedHash, _ := hex.DecodeString("4c6ea217404bd5f10e243bac29dc4f1ec36bf4a41caba7b4c8075c54abb3321e")
	s.Equal(expectedHash, hash[:])
}

func (s *RLNSuite) TestPoseidon() {
	rln, err := NewRLN()
	s.NoError(err)

	// prepare the input
	msg1, _ := hex.DecodeString("126f4c026cd731979365f79bd345a46d673c5a3f6f588bdc718e6356d02b6fdc")
	msg2, _ := hex.DecodeString("1f0e5db2b69d599166ab16219a97b82b662085c93220382b39f9f911d3b943b1")
	hash, err := rln.Poseidon(msg1, msg2)
	s.NoError(err)

	expectedHash, _ := hex.DecodeString("83e4a6b2dea68aad26f04f32f37ac1e018188a0056b158b2aa026d34266d1f30")
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
	groupKeyPairs, err := ToIdentityCredentials(groupKeys)
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
	index := uint(5)

	// Create a Merkle tree with random members
	for i := uint(0); i < 10; i++ {
		if i == index {
			// insert the current peer's pk
			err = rln.InsertMember(memKeys.IDCommitment)
			s.NoError(err)

			fifthIndexLeaf, err := rln.GetLeaf(index)
			s.NoError(err)
			s.Equal(memKeys.IDCommitment, fifthIndexLeaf)
		} else {
			// create a new key pair
			memberKeys, err := rln.MembershipKeyGen()
			s.NoError(err)

			err = rln.InsertMember(memberKeys.IDCommitment)
			s.NoError(err)

			leaf, err := rln.GetLeaf(i)
			s.NoError(err)
			s.Equal(memberKeys.IDCommitment, leaf)
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

	verified, err = rln.Verify(msg, *proofRes, root)
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

	root, err := rln.GetMerkleRoot()
	s.NoError(err)

	// prepare the message
	msg := []byte("Hello")

	// prepare the epoch
	var epoch Epoch

	badIndex := 4

	// generate proof
	proofRes, err := rln.GenerateProof(msg, *memKeys, MembershipIndex(badIndex), epoch)
	s.NoError(err)

	// verify the proof (should not be verified)
	verified, err := rln.Verify(msg, *proofRes, root)
	s.NoError(err)
	s.False(verified)
}

func (s *RLNSuite) TestGetMerkleProof() {
	for _, treeDepth := range []TreeDepth{TreeDepth15, TreeDepth19, TreeDepth20} {
		treeDepthInt := int(treeDepth)

		rln, err := NewWithConfig(treeDepth, nil)
		s.NoError(err)

		leaf0 := [32]byte{0x00}
		leaf1 := [32]byte{0x01}
		leaf5 := [32]byte{0x05}

		rln.InsertMemberAt(0, leaf0)
		rln.InsertMemberAt(1, leaf1)
		rln.InsertMemberAt(5, leaf5)

		b1, err := rln.GetMerkleProof(0)
		s.NoError(err)
		s.Equal(treeDepthInt, len(b1.PathElements))
		s.Equal(treeDepthInt, len(b1.PathIndexes))
		// First path is right leaf [0, 1]
		s.EqualValues(leaf1, b1.PathElements[0])

		b2, err := rln.GetMerkleProof(4)
		s.NoError(err)
		s.Equal(treeDepthInt, len(b2.PathElements))
		s.Equal(treeDepthInt, len(b2.PathIndexes))
		// First path is right leaf [4, 5]
		s.EqualValues(leaf5, b2.PathElements[0])

		b3, err := rln.GetMerkleProof(10)
		s.NoError(err)
		s.Equal(treeDepthInt, len(b3.PathElements))
		s.Equal(treeDepthInt, len(b3.PathIndexes))
		// First path is right leaf. But its empty
		s.EqualValues([32]byte{0x00}, b3.PathElements[0])
	}
}

func (s *RLNSuite) TestGenerateRLNProofWithWitness() {
	rln, err := NewRLN()
	s.NoError(err)

	// Leaf we generate the proof for
	memberIndex := uint(4)
	memKeys, err := rln.MembershipKeyGen()
	s.NoError(err)

	// Create a Merkle tree with random members
	for i := 0; i < 16; i++ {
		if i == int(memberIndex) {
			err := rln.InsertMember(memKeys.IDCommitment)
			s.NoError(err)
		} else {
			memberKeys, err := rln.MembershipKeyGen()
			s.NoError(err)

			err = rln.InsertMember(memberKeys.IDCommitment)
			s.NoError(err)
		}
	}

	root, err := rln.GetMerkleRoot()
	s.NoError(err)

	fmt.Println("root from zerokit: ", root)

	// Inputs for proof generation
	msg := [32]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0}
	//msg := [32]byte{69, 7, 140, 46, 26, 131, 147, 30, 161, 68, 2, 5, 234, 195, 227, 223, 119, 187, 116, 97, 153, 70, 71, 254, 60, 149, 54, 109, 77, 79, 105, 20}
	var epoch = Epoch([32]byte{0x00, 0x00, 0x00, 0x00, 0x01})

	// We provide out custom witness
	merkleProof, err := rln.GetMerkleProof(memberIndex)
	s.NoError(err)

	//hashMsg, err := rln.Poseidon(msg[:])
	//s.NoError(err)
	hashMsg := [32]byte{69, 7, 140, 46, 26, 131, 147, 30, 161, 68, 2, 5, 234, 195, 227, 223, 119, 187, 116, 97, 153, 70, 71, 254, 60, 149, 54, 109, 77, 79, 105, 20}

	rlnWitness := RLNWitnessInput{
		// memberIndex key
		IdentityCredential: *memKeys,
		MerkleProof:        merkleProof,
		Data:               hashMsg[:],
		Epoch:              epoch,
		RlnIdentifier:      [32]byte{166, 140, 43, 8, 8, 22, 206, 113, 151, 128, 118, 40, 119, 197, 218, 174, 11, 117, 84, 228, 96, 211, 212, 140, 145, 104, 146, 99, 24, 192, 217, 4}, // TODO
	}

	fmt.Println("Witness secrethash ", rlnWitness.IdentityCredential.IDSecretHash)
	fmt.Println("Witness merkle path", rlnWitness.MerkleProof.PathElements)
	fmt.Println("Witness merkle indexes", rlnWitness.MerkleProof.PathIndexes)
	fmt.Println("Witness data", rlnWitness.Data)
	fmt.Println("Witness epoch", rlnWitness.Epoch)
	fmt.Println("Witness rln identifier", rlnWitness.RlnIdentifier)

	// generate proof
	proofRes, err := rln.GenerateRLNProofWithWitness(rlnWitness)
	s.NoError(err)

	//proofRes.ShareX = dataToReplace

	fmt.Println("Proof Epoch: ", proofRes)

	// TODO: for testing. proof without witness (are proofs deterministic? maybe not)
	proofRes2, err := rln.GenerateProof(msg[:], *memKeys, MembershipIndex(memberIndex), epoch)
	s.NoError(err)

	fmt.Println("Proof1 Epoch: ", proofRes2.Epoch)
	fmt.Println("Proof1 Nullifier: ", proofRes2.Nullifier)
	fmt.Println("Proof1 ShareX: ", proofRes2.ShareX)
	fmt.Println("Proof1 ShareY: ", proofRes2.ShareY)
	fmt.Println("Proof1 MerkleRoot: ", proofRes2.MerkleRoot)
	fmt.Println("Proof1 RlnIdentifier: ", proofRes2.RLNIdentifier)

	fmt.Println("Proof Epoch: ", proofRes.Epoch)
	fmt.Println("Proof Nullifier: ", proofRes.Nullifier)
	fmt.Println("Proof ShareX: ", proofRes.ShareX)
	fmt.Println("Proof ShareY: ", proofRes.ShareY)
	fmt.Println("Proof MerkleRoot: ", proofRes.MerkleRoot)
	fmt.Println("Proof RlnIdentifier: ", proofRes.RLNIdentifier)

	// Verifty old proofs
	verified1, err := rln.Verify(msg[:], *proofRes2, root)
	s.NoError(err)
	s.True(verified1)

	// verify the proof with the witness
	//msg := [32]byte{0x00, 0x00, 0x01}
	//verified, err := rln.Verify([]byte{0x00, 0x00, 0x01}, *proofRes, root)
	verified, err := rln.Verify(msg[:], *proofRes, root)
	s.NoError(err)

	_ = verified
	// TODO: Not working
	s.True(verified)

	// TODO: test a proof that shall not be verified
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
