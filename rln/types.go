package rln

import (
	"bytes"
	"encoding/binary"
	"time"
)

// Each node of the Merkle tee is a Poseidon hash which is a 32 byte value
type MerkleNode = [32]byte

type Nullifier = [32]byte

type RLNIdentifier = [32]byte

type ZKSNARK = [128]byte

type IDTrapdoor = [32]byte

type IDNullifier = [32]byte

// identity key as defined in https://hackmd.io/tMTLMYmTR5eynw2lwK9n1w?view#Membership
type IDSecretHash = [32]byte

// IDCommitment is hash of identity key as defined in https://hackmd.io/tMTLMYmTR5eynw2lwK9n1w?view#Membership
type IDCommitment = [32]byte

type IdentityCredential = struct {
	IDTrapdoor  IDTrapdoor  `json:"idTrapdoor"`
	IDNullifier IDNullifier `json:"idNullifier"`
	// user's identity key (a secret key) which is selected randomly
	// see details in https://hackmd.io/tMTLMYmTR5eynw2lwK9n1w?view#Membership
	IDSecretHash IDSecretHash `json:"idSecretHash"`
	// hash of user's identity key generated by
	// Poseidon hash function implemented in rln lib
	// more details in https://hackmd.io/tMTLMYmTR5eynw2lwK9n1w?view#Membership
	IDCommitment IDCommitment `json:"idCommitment"`
}

type RateLimitProof struct {
	// RateLimitProof holds the public inputs to rln circuit as
	// defined in https://hackmd.io/tMTLMYmTR5eynw2lwK9n1w?view#Public-Inputs
	// the `proof` field carries the actual zkSNARK proof
	Proof ZKSNARK `json:"proof"`
	// the root of Merkle tree used for the generation of the `proof`
	MerkleRoot MerkleNode `json:"root"`
	// the epoch used for the generation of the `proof`
	Epoch Epoch `json:"epoch"`
	// shareX and shareY are shares of user's identity key
	// these shares are created using Shamir secret sharing scheme
	// see details in https://hackmd.io/tMTLMYmTR5eynw2lwK9n1w?view#Linear-Equation-amp-SSS
	ShareX MerkleNode `json:"share_x"`
	ShareY MerkleNode `json:"share_y"`
	// nullifier enables linking two messages published during the same epoch
	// see details in https://hackmd.io/tMTLMYmTR5eynw2lwK9n1w?view#Nullifiers
	Nullifier Nullifier `json:"nullifier"`
	// Application specific RLN Identifier
	RLNIdentifier RLNIdentifier `json:"rlnIdentifier"`
}

type MembershipIndex = uint

type ProofMetadata struct {
	Nullifier         Nullifier
	ShareX            MerkleNode
	ShareY            MerkleNode
	ExternalNullifier Nullifier
}

func (p ProofMetadata) Equals(p2 ProofMetadata) bool {
	return bytes.Equal(p.Nullifier[:], p2.Nullifier[:]) && bytes.Equal(p.ShareX[:], p2.ShareX[:]) && bytes.Equal(p.ShareY[:], p2.ShareY[:]) && bytes.Equal(p.ExternalNullifier[:], p2.ExternalNullifier[:])
}

// the current implementation of the rln lib only supports a circuit for Merkle tree with depth 32
const MERKLE_TREE_DEPTH int = 20

// HASH_BIT_SIZE is the size of poseidon hash output in bits
const HASH_BIT_SIZE = 256

// HASH_HEX_SIZE is the size of poseidon hash output as the number of bytes
const HASH_HEX_SIZE = int(HASH_BIT_SIZE / 8)

// temporary variables to test waku-rln-relay performance in the static group mode

const STATIC_GROUP_SIZE = 100

// STATIC_GROUP_KEYS is a static list of 100 membership keys in the form of (identity key, identity commitment)
//
//	keys are created locally, using createMembershipList proc from waku_rln_relay_utils module, and the results are hardcoded in here
//	this list is temporary and is created to test the performance of waku-rln-relay for the static groups
//	in the later versions, this static hardcoded group will be replaced with a dynamic one
var STATIC_GROUP_KEYS [][]string

func init() {
	STATIC_GROUP_KEYS = [][]string{
		{"2de8ad2cd30e993ff870a31596df4161343d6f05cfe8501884b807eb22dec066", "11ddb7419b09e862027c0ff978d60fd8fece82b31ede62f5e5b931d9986c383d", "022eb3ca06ca2634c5e5a63bde6cc40455cbdda928236f6ca2495f749172ea53", "0d8bae46e9af5072c2042a4c6960fcfcb5ec945479a4683827aaf1ab8a63e07a"}, {"11ab22ee14e4f53ca76120588a33513a5f212dfa6e24175826cdea72482723a3", "263555637d443f73e6ff88ac5fd95cf1c525e52005ff080f2767a8e56fcb86f8", "19a76ef9a5a868bec4b62434e437d90680e36b9d0099302a9adce38a519c7040", "07061c488d0103f01b7334cd175dcdcbaea4bcb18a3847cd9186866536d246d0"}, {"0acc4beab8a4fe06b775ba52b36d7d8d7890fb62da415277de02bb8d7cbfaca2", "06649c699e795cd10588844aaa481481a718788ba7a138b413f7c735a0f76a3b", "00a0400ed331a6a74dbf3e361f7df612c318f1ebedae26f739e80827e36bb5df", "0cb30c063d0363d5547bfd01b62500b61d2f6d434d7728d96cd7b8219175d85c"}, {"04d56ee5bd17cbda45b2fd5b85c2a147c1d91ac675eaee991859abe314e8509c", "091b97505778841b632fbe6cfce2552e091ea1e71ee0266dfd542ad73e4da476", "2bd71275649625675d418d46347de784bf93952d3b29cb9e24a95962c570c3ff", "0a54271b9bfeeebd7467b4306aeef7b36ba71d4eaa6ef4b6ba311872d2cc1d4a"}, {"215833d966a2d3650c3b3d38829d2dca39b4b557c5e35a3783c41436d54415bd", "28c7c326b89a097da2a3254de930f98b05c3e1a1f6909c20afc4ef67d58948fd", "23d336beac7ea57e699537ae0270a3d7dd7476bb59835468be0b46b336513a00", "05f5d750eb243da48f07e3629451068fab8f9d7c551906af0caca00d2c94c50a"}, {"24520de883497eb171445202d7a82b98f6afa3827bb8064dd568e0546d08f49e", "157ec42e00e8d810ea94f834541a3d0bcc6029c4115269354d3baceb084bed9e", "0954402b0b26bdb7a2da1312cf65ef616dba71fa7c1459e421def2337897116a", "2bf5d1b7e98b51ae9cc300669d6e6d2a6c8d2de575519c6d53dda0c0d0bdde30"}, {"295fa4839e238b1410df304a44e2b59234bc9f87c1e8c8bc09075c64ae16e2b9", "2f90313f25959ed8f080c5a54df6991e460cad3d6312e3b18d029ad05f8e28a8", "0f0c817957603df4edbf7cc3a046e68f09ea7870a31192c48f1fc12e45032031", "1839a530aa56b1f50856a85926c1842e678754e9c3dc059e34eee0bc70f88dd4"}, {"1d0c99e03b6aa94e01f29decd6148230471a5b8e60753b659ed2a891c4e89698", "13d29e820f34e02477c7713a246e662705f33c92dc2e69695228d593e115260d", "006fe254e1c8fbe34555e51224ef97f0ea4eb005168cbd8885a5d99ee73da095", "119b0cebca96186bc219066d781258e04e386097f3c03080a494ae8d01e2bb31"}, {"217fd32fb7bd401ef0229aa2dd1c93c591772bc970249e52621f76a6d529fff3", "0113d56afd15a5ec23eff9d434f4547fdbb3c82029a622391585f5000633f583", "1f677a1bab221a4d5d1c281b635133b9376906dceab0fb36510cb0c2b08efca2", "2971c74b6c9f324a30a83b32d211850230250ade39740055576f0dc85fb98cf0"}, {"2d420d4da783d075faf6aceb5405766f49d08a06a57932f9e0fc58ffb41b4198", "272849724f0faf574008a29597fa21d8f86cdaf3973cfea19a50b59119621a3f", "18189e9e10e820cf3c667c8346b9527213ff3a1945f687138316bb0bb5e65ca4", "056ca699b1c87b5e8598665ef129318935dd3e96492b79113d2451fbae9f1cc6"}, {"1fe33b9a7284967bcd69b14f648a7e4667a59a59e030030abc12864dcf64a681", "1fa7f7a31f8d313d919765e6b4a0e681ca3e1e6f95cbb48d893fad389e1d8859", "2692bb4ca076936d3079735d5c9a6add5a1f02a968b00599a1478114a04891d9", "12b4157333b5b07b785db8a30c73b5c474c5f2649ef47a63d78f1b998cbf59e4"}, {"1d053812ffbe4d4da6e21cb3b8728a60bb6cb2420e29998e78c4ce965c35ec8d", "09d840b025919e2c2b381ecc5bf9fc2dd599e83ea1c2825d17299cb34ab9f6b2", "1824d455c9577d65db341e59ebc7b6cf940bf54af49ae800fb0dc5f32014d249", "2012a26ab7b783808fbcf8b2cc58158a2db725d277c9520b65e5d36e5dc6b052"}, {"132a653c9148ae4f992f223df26541e11d39953e273e47a863ba631c318090b5", "20010c72714dc65285a0201f7bb1cf464e030b0fbf2d2be3177667066a9d07da", "1b3dd5b82723da7f939d882125322f91f74279d0cd858c45e09011a51a4de962", "2601a3ed3ff4a7718f0019a3939c08fc3a9a5a8180a68daaa7bf34cc486bc757"}, {"2f656d8a0c07ac9b2de82be9d3cbf4535cdf169cbc2b355155d9e0ed256cc33d", "1f7c9c1aa1dca7354f135849012aa5968b37ae729b2847694f65355895d41984", "00b8302bd41850608f88ef1b0e6949b08e850fd8d1044b3d23b7224aca5ea7d5", "2c1d654b251d138a3a6be82f784e54b05ce5b6501ce2d59b89021d20abca7516"}, {"168c45a2df3f7f5dc9305277cd8b78d023f817ba5ad198a1ea6b1f23aa002478", "0f50c493b68eaf26cfab59c903ef583f03348d9e839806af0522759d64afac60", "2d1e68d056c1928fe7a44b96cb2333888ef3a169be54f4dd5438fc0de23554d7", "10e58fa7abc99dcafc5f4e08181c095c78aaac093e3a5c5e356e8e53b2049a99"}, {"0b543803e520e8ca645d8d324399c14e5cd1e6d5c2e371f2316beda99403471a", "0ae42ba2beaeecc65174318588f875ed9d4bffdc8a40cb4862dfbbe39254dce4", "2a5b829ec8871ee4260f83688978295a5b553cfb92a8e307bcab0a9d0e2d9320", "0d67d9161644b0a09e230f28e95b2f70558ae09260383b97b6f53bb791e72aa5"}, {"273961e60d8b20fea6b34446911f6922e95b5948f0c673f260f8b2e40727ced7", "0c0cd7f50d1472d3246ae86f605e3c3adc8865c8ad4799ad0f738ac710f71661", "2764d7522256fc102b5ec886e11938cfbae7b85fddc4b8ee2b40180178bdc118", "26d8a135115208aa2321b57befda6602fcc048cb988426d1d2f4888aa7b6a652"}, {"2d6d3b281a2d50cc5bafe3255c353247642966092a48ab570510c865c20f3691", "0300a6f79f78c1e0dbefe938800922046fc2d5ffc4e7d5dd9fc3958145ad7858", "0f7123bb490708bcdc2b16685ce3515fd1c62680e3cd628c15ad663bd1dd10f8", "1ca0d6fb94e3acb89299d8c5724f55f14eb82875d5c6d3e769ed2afcf7ec3b6a"}, {"205b991d7685b23c17abff32b668b09a1236742c0f3b35f4bbf92f61cbf2946d", "0cf3674c678049376d9e0c24e6fc8c1ff7e5cbd00971a23aab01efe1dca3d712", "27bd23a401958d102829f7cc4f2812ffb457fd5d4af84d603bb6de2cd81a2d0a", "26feb04420e100cb694a9e54cc4bc09848c99981e3d1d240a299a84fffec42bc"}, {"1230a623220b9c5246888bb15f98c646c85adba671bedf1facc7061f7d3a2cdd", "1e5da8d9c840a6fc6c46616d1b49d8c83328f3cde51ffb2a6d418da63b0827c2", "0ecc20de9b74105a5da2db218396997059e9ed61c2a2aaadd278acf4c6d715b8", "03fd0919740fb0d2d71827d1e904ddc23f386706edd61dfd88cea89c87df5620"}, {"0116df196a2fbcaeae5ec4a5ae817aec779881cf9aaeb315f8906d6b9fad453b", "133f0939bf05bfa523e51c6e9e4c96c125f10dc733f5166b0eca1c902a08a698", "1784c2ec594b121b2784ea4a9b44f259bef8a499939406d7859041226bc9b80f", "06f5089423022496743bafca22605c04e570887e349d109aaa3342cae557879f"}, {"12c80c6764d68e018713802998171001446c58d493b3add50643679bc506a894", "02c63445f57a532a0fa201eeefe1a5cc065f301a6e01f339c1ab7a2fc3c346b7", "1487669cc19664d349e43eecef1202b4706d79baa0fde4a44396b33d8232c6f1", "21c5219a7a414bed8e4a41d6586cd0791fa512166fc4615b643fb26d24cdc16e"}, {"2f4ee97ca44c2508af67ef15ad5bf317adef767f2ba80ced53cd3a6136cdb1f8", "2f903d2bfad210c4556c790e271af0f4979c63516b6f61cfb34eecedfcc641ec", "1a5c10c4cfac10834e5dbe29cc7d0e9d07c25f90b9f528fe87fcbea3b3738c57", "11587341d004eb538887cda6cbb2ffc33ac819cd4884eb06585a706c4efc4228"}, {"2146e9aac9c2eb448f738c0c67069a8eeed3e4de24e76aaaf4ae1bd94637a35d", "190fd021acd70b4c4465048dcbb2f42178e84a2586672a6fd29653aadcb7e580", "1fb396f784fcf88e50e6f005d1097916f62262677caf055e867a711f9486a732", "0489a13f0398f99063227918013c9da644565e62896d97ac520ef99c35fc0ca7"}, {"2dfed889b0b401205fc5e0b69728a732fe6aa3a3abcec3fe55edd671041344e5", "16886eded16061974adc5889a7a289283bc9c4103447fc971fef722e0627d8bb", "1ec4c39a2ec0a2671ec26045f307418b7d0bd92076b5a5e5997f13a08ff34f82", "0829b6f724b6d33060c3545c04708aff599cfc783aabb2cf50395a9de5505795"}, {"0dc00687f68023c9f79e4b01c94caca14a318aa24123f924eccadb5acd3ae674", "1a6a4b9d8804c109888b1e67c31f19e816d485077722c57377dc98462882a778", "0771ca89a28322686f809dff401b626a205f05d98086f194c9cb0ab697010004", "0a3d44533de845efbd99458b03d0fdad1c94c0e94ff98203141906c316943caa"}, {"05a6550fa6a01e2ee1eb13676f3e1d10afdaa571555a85e6c36749948cbe5fd6", "296c1a3f9a2dce1ba04e126883f026b96e3b24c7f20a2622fefa80fec566dd72", "0786dafe3c7aa093f2a6d474a62574cf96ac755c9b7a22138000a4d86d0b6e28", "2d82b28b0fb43d8f66b0ac548c099f5b03811f5fd3d64e79f48fc6a11fb368e8"}, {"0d3799c90302e2622222b85f6919062bb3cf20920ab42a515f047bcbeca8bfcc", "07945b9c2ceb366148736f997ea3567b2b88bb7b2f1bbd7152f7cac3a15356eb", "0cd858026c5a2a8fcecd92568fae21bbae89c0de60a966edcdd842dcf20b5bb8", "0080ead6d9c19c515e29072519d458ce942337e725dcacbbe2a0037a7d357c1e"}, {"1aadfd7495ddad900e6298f8e3b05f76b1b38790982f00dc1cc3c2a0d7b7a0b0", "2d4d02617233f6c596469a454225d5312209f83d9bf7ef7042986f4bc0d2f263", "17dd395569064dac45abac3f165dab42b4e2522d4bf3973d7f427a9895382598", "0a77dc303ee2cdac0b1f0afc6d033ba1fff4d3b4a5d0585288fc0da640de4b36"}, {"15b569e4a54a91a7c18007a46fc5d9bb2367b383a0e1d2327b667cead37600d6", "222e9a22f44cbe8e704ce52aa4edf54dc98bdc7240911d12c4f41b56dcf6edfb", "1122b5efe8686520a0e24ba534ccec09428d921d9bd290284910f615627e0dfd", "0d5bf9b92a7a3d942a68773b0fa3ddbb31f8abff8c1cf7c175b03a73bb5f04d3"}, {"09f5e4a584ebfd2bd07837020ebb47ba881c864da25a96c6c99ebc8046073651", "26b819e34bce5e8e20acc0a8370385d476307c2512a1655d550ef16230df6b39", "1ae05de2fb9296093f5dd47502627c35a50a1a0adfe59e2ca1fe09d0e9ab5192", "09c90c0a107d8798593de1764f90fbace04bc8ac2695eda0ed3011b7f6084a29"}, {"0f5503c833ae007dbbb53bc06c6abf4c7c68f52996abb2b115e870544868be01", "2760cf6b2471fc666674216dab356c6167e978d9b36fe871caf959e32254f6b7", "2def3b0bc44bc2d52a6f3e66e43b4a8aa048c66260ef7d8ff8b3ee9837a1b423", "0cdbb3cb827cd02819f8a08149bb9bb957a180aeafa18706df01e0edc41a89eb"}, {"065a3762cc9fd71b13ab5dbac9905d73b9a2adafaf5379e0f6d1a3b788735a83", "09a50a6e15d90468e6ab37b988804e44308530a2d9113a62607bb436cba7d21f", "08fe67b723b3c7d0a3eb1e041d1fed917a21750c966fff993c648b1fedc90083", "0808d4a78df10cd77889fabc546147acb93f94ee617019531719ea52a853bcb3"}, {"2adbce35c09adaea1f97a2fc639e1f2a1a975ccfee9509cced16760c7c9a75c2", "2cdeddb0628dc4b46d60ae9c6fb2a7c3f7fd04689851d91db654f9ba30d77522", "061ea76a71766cbdab24bb18e513dedfbaecf7dae6121d90e756360b5ca21a76", "0fea769786c8eb3e48489a92e1aed1965e61ee40be7f3fa3c15544fc6cdde7f0"}, {"138227c591792141525d9a128ab021f427c968b41927c13110dd3debcd5bb79b", "2111e1d4acdf7b4b0d46708deee5f5a4b93d7138798245075e591ef43af6c86a", "24e8c1dd992d18ff1646a911a4506ae64b59adf99a757191d8b15f294f23d39f", "056a158c1aac62d6e5587120a40f84716583bed992d31399b31a8c5ff48811f3"}, {"1cfca16420e617a17fe29ab96b468f5f6c63ad0a2fcaf1627ae53994bd9bf6e5", "037c3af3a9591d1cffde88b2383a27a3472e219c14545e8975a5b02bd7db8735", "1fdb4510b85aec27c60c977e1eb04318c560549a8237fa75788305e9376d1c1c", "244e22cc2f20498c300bb6d4b482e0fb0b70cc38b74f100d27b7f1a536fc3b3f"}, {"165c1c813e61d12e8203a378b681a6e56d38a152e3ff39ef35ab8a74162f6739", "18fab342e96343933108180ed6a96e018e9f46764379bd25ec53a02625deaeb1", "0aca6ad6b91aa0e4ad51b8e32eec62c770d2dddf69381010020b4c294834204e", "227bd6205ee4047bf5c16fbe8914e69e41d28b92c9d2fc901319bf6946aacdcb"}, {"0aa16c2106010e06debee4558bf893415c689ab661c9f1538e62350315e7343e", "0305bd43adc89116b436e004e0fc7b49975631339dd0e5e608229552edf6364d", "1e37cd06cce3965a566a2080b2087be3e85af2776f9f293c4bf2edb92fd94102", "037931a869e246df195a3f367df6221abcfba5f4a230a5d335de6b3f1d4af07d"}, {"0e7ba80ad7b271230438317de4ee9f2a3789a495d85ea75e2320f3dbdddcc3bb", "113d9847fb222673b086881eaf6e90a72ca9cba676d22d3e4f3644c8e8244b79", "1c10febb6e6bc82aed81c3a80ec34aa6bb7a14dc45eeb1272e7cb4d19c3eb1e5", "1d657c9ee4a8a2b31a02e6f5ad79da3b71f22685da66d48be63b8c96ca7dca01"}, {"2f7756ac206f83fce0e090a93f81a08563ae66594415db7a86668533ea19069c", "27d4c173a56e32bbfc470de7380488fb03c4783c1465a6f6f4ab9d9be7b81337", "2bbe945d4f20dd95b444343b86cc44a69b91a26e60c1f05754c76941f200c28f", "2e046ee0d56a92d75055b7c72766420cac23ee530076fd81917704141e4a3037"}, {"25a136bee3451a9bbf51a30e507eb281cbd6a56f160a51b849643d0b054e703f", "1839a68781e571545fdab3e8b04f9997006896fe580ce0e40d5c0f30a17b743c", "18c1217077424801cbeff7c71ea24327a3c7eaa01c43681d14f045982afe6869", "04cff8f2dd2248c6c34f312884af6b7ec6875570e795daf7cdf3b302a2f8d360"}, {"149975c4626c7c621e05af3c1059b90bab8fdbd8d11378343f8a13d36401e88d", "1978d1d3ff6483f9b26adfd2fc287d4a3f4972a9a00ae4a16490807c5238b6d2", "26b04d8f9d4071eafac8877bc305533f241cd03e857ca84bd689bdbca7f9f7f2", "2e2869ba9d0894a917ab3a8f5c03a6aa92737a2bf8739bf894c012168d730e33"}, {"05be359a5d595637d1f4e9948d17a611401a8ad583aa78d73aefe5cac3417e68", "265dacea532a6793395aa5e9cd3c2f19a9cdee9bd2229f9ec0ee0987f178439a", "1a28c25b66f03ac6749551d2a6056466f85b7b5e5db738110913734f6bafa5b6", "2d55b1ac28d0b2baf29bfbcbfde960f1bdf18bf28d52bf306d774a64c5a04faf"}, {"0acb23e044e518a88df2345f9cb6e331520f89af0f2c092a09b63716d83b029c", "15d2bc8ded95329ae9174bc01a4c12a48aa3399441dc1ff408b4fdaf181634d5", "2ccd822fd559347040e4fdfa2e7690d1af03f3d523cac219a7cf06f881f1859e", "0bb7770a6df49da10b3e2babe1f5881d2ef4f40e49d99a6a848ec1f2dca72e4c"}, {"1d9c9f55deea40422f37828e0816fe18b1acd665e96392d376bdec709d31862e", "02f9d59214040c035861d50b4e238e47a75460d03d58a217892f89c63a9685ea", "12cacfaa5924f4bba951f13cd3df068f2c11fe99166a4c59732ab1adeab418d0", "258e4273c1c463d211ca2a87866136192b0c0564cd91bf0a4fb037fea0a781e6"}, {"1611e0ef464c8903e6b6b23f9cdca49a09d66c1ccd4cc113f9a6bd6b20e9d1f1", "0744b00b48d373789adbbdd0be955e4887e2d68337ef477a5defa2b5ce6eea84", "162a1e3e6f8b7e91e8f3b44bf83ca18909190811bbd28b579be19b385c7bb524", "16ce06eaa8ef6e70e30564051e0f6d2b5752142cffff3937c6a33713dc87937d"}, {"115d5f4bc2dd803a71af2a43a230115dc2725a180beb7df1b88a2259e2053412", "10daf1f8f9e13d9ec374327c2b51b4cf2ba45e57ac2a45694aff64f6916aa889", "18433218813243ecef035e369a64b46a704a13fe6c80858c6ce90034f6357bb0", "04fdbb0aa628440f438a57f0f213dcda509ca8f1d0cea93556b30d4d01c5bc10"}, {"09de472ca52c06d076c3798403302eaa75803ed28e4d81f63f2bf488b69a88f6", "12759012df29f9c911a8b44c41d13974b7be3591eda2bb6f75ab41799e4712ad", "13d7626694617b0a492c90996660b665bf093f931891663bcaf7084c7cf59509", "15a646e677c3e76c021715ee48ea0430e1f279b739f8d8a129c46afabd26290c"}, {"1bc3c21386512fd3aa169484671be058fd471f9d4a899a6be4cc16a83b483c8a", "07e711c6993c9f370fd6100b468b7300ed1fbd3238712cf8e65052672171ce2c", "025e58dc2bc1483098410d9b3c4a17eab6102f89a6183aca299026b9f1615b16", "140229fe33aab22020f29e9b2432ba65065e4d6bd1494c359e6b086266f8aec6"}, {"1525b676c8981f9c5e544b44133e94f312baf24b42cbd459c41f92e24ea02cef", "123a0dc832ffcb70c94f1cf1bb8f1c7809f66a4773fccc3006d3110e60414f36", "2a645289284773a2efbb11823419b4ea22cd71ab2a8b814139e6375ba459acdd", "2d5105e8567086d71a2cce2046c449768bfea36645be00f17e4339fed306bfed"}, {"288857614a0f17a228c98105ae993e48bcdb25a876fff0a2b9e1db52aa7cfb12", "28d0e7dd79b3df14a6112bda58932a8a4ae65180b614919e7a8f32a3eeeeba39", "124d335735701f8f2e7af829f5b326c053c7fdc0bfbe9299f0d2cb4bf432f5ad", "13a602a334d937687a5a7710dd7085bb342dd5fc9f2261fccae20cd39ea7d45d"}, {"078e9560c3ac30232a3450e8f9a487c55e6e37728dd1abdf5ecebc56a2f1b691", "24d0478044359d16d5477892983d609a261f84b444214f2325346bfbdcd8452a", "1ec8e1362fe0816a633a02d4f1e80df5ad3dfbbc5fdc91ce3b4600e6e969c551", "0377f7c3c95b5799fa0bf03a76e7b85d436154f680a9eeb6b9b107a214787786"}, {"03da01bde649b5d23817740ef1fd2afe0d8fa3aa35add446e3c2165a5ee20c7f", "11c1a1006f436bb062d34cbfd4f65a37970c0cf4375b968e9afc39dedec5a1f5", "12ca4ee75899ddf27179b38880b4a0019b3d24cccd220a17ae2924b03fbd1594", "2c41274935645dadfdac964a5845201c4962560ee8cc4932969b374ad55ea4a4"}, {"1acc4748897f2319e9840230c603fee339eee1ae480317666de5bfe935b17a92", "1adcd69d17130e4e2d9349457735ebdf09d3279d81e577877ea969025647062b", "05b9178f7b8962377d922cf18942f88f994369079e7107d79f48b1c7ecdd0f75", "15a764bb5eb2cf17cbc6eddd7479e85d2de5a53d8961ad7e47e7481739580ae5"}, {"28045d48040ecea6e8ffa983f27422c12c9a32b6f93259c91d37b1b7282085e6", "146ac2fe6a594fbb4ad66cad23becdce68461fbfe9b1cf5e6d160f2e3b47299b", "1160d22d7b3298f58f65553d1332bf93669d26fd550b48317bfaaee8d5df09f3", "2ee983b7b21cde8ebdf78a4461abd510496f940a8524569e4df0203e449a415a"}, {"1a67133e1153ca967805fd36a138b113537f2dcd6bee07d89de2549541848fe8", "1a0214be9ac7f2ecb4e60953efa1d4da075e180b90de1320c48736b502272134", "27f3091c65ccf894f34a37a47c62368725969aaf0f83f11dd8693f5d3cee06ab", "0f09d248f2e5542e69a5d3f8cc4fdf88603fb563cb12a861f7511c98028a20bc"}, {"03a5d3830dcb353dc60546ba371901c2546bab426bcfa7a3aab986d67cd75a36", "01c05ed73a97503b50f58f265ea6809ba990c3fef2dc4c7ac364524c4682968d", "1ece48a1732f79a80629222c9183f5928c96d4cad35bc67ce066e2db5df1df18", "0c897ca901854557c05b01a58d15d4e39ea85f69840c29cf5a7077420a2014f0"}, {"30482bfb8cd460b37ad4b18d41a3ac4a19678bbf4eeb23df0a1f17ae293c2e59", "1be551bd7b10a217a68e982369e50e686642892bf52023a65fdb665a08b294d5", "14d3a66b20b691f7f5c676c4496bce6eac351a2fbdbb0551567e01682cbd563e", "245019483fc1811216ef3de643147c9a4dfe0c305165fa3dd5082c2828ae781d"}, {"0b9f0b77879dd07829838d350ea748dad8315df92eda8a09a06ae803b52433e8", "19d0ca56369e4f1bd8328de20ce700bda523f3061b9fa27c6e41c72707ef84f6", "2d731edfd266ae50b1ae70d22e3a9cea52bf052318240c6ff43ce25e0d26ac59", "2115a2054a4a502f18bd51c023220afae7f666c98bf27fdc7bbcb98a7c3220a7"}, {"053d4993e28cc2ccc565bf58cdd11707fa5a23bafacf48e0738851b7dc7cba86", "0b936c7981c946de2f6a212355f09abd15f17ee7cb1d839f68ce110484ad9b1d", "24d042c178799097ad2f73892cd2a667c54d9a27df62db765d783ecbc0e632e5", "28f150ec05f4aac631cf20bf530bc08df7dab9b6372df9a01210ea7ca68a90d9"}, {"0bb355cf5de29b2e15823def2571259dc74adc246e14bb815a31f29517a3adf4", "28ced5bad892909a16d7c1ee56eea5400d816a126e4e44fbbe8a48f1db8d014b", "1f64a07121a479444254e65b6c6cbcb94669fdde074f598b653654a62294f2ec", "27e1c7393ef111b8fb531cf65c72009a1b4ddbdaa8708d709a835652a2a7d49e"}, {"2f673f5f870d491c69c0435ee13b4c8e861e6eeae494416ee9429a78ead553e6", "13742c523d2fd99a3082fe6556c81ed61130e2965e5aa6a6ebe1b962ed1e9da5", "2eb3fb08e15d8a2f09895eaf0c599a7177e37b2eaf7e24b7a36a9e48c7d975ce", "138c90876961451a808920cfda3c2768d647e9b6659aa45888091d5b41c92c01"}, {"1b279e5fe04439739a57fef9ae89460ebfc9b25e9e14823bf0d0ea60854e3655", "0ff39bfc718daba8ffade09ee5a296dfd844a65fdb1ba041991581bcb4a9d1c0", "259de05683020849eb06d1416da4b4e5840af03316cf206ca01a30ada34582b6", "1354c7a215ee36091f67853a1300821c742d31a128542b1b4bfc0d849ae84369"}, {"0a42be19ff6acea717f247d0ddb3af24a22fa35b05d1f98bc99de3e1906885a7", "0b42f7a1b093dd2a6e193a48c8147b0042cd2a33b4947cae1df0da22061b4f69", "144c7d59971c8c57497f3dc42dbd3465305c181a5bef7849a0e5aa3fe765a2af", "2a8968087625cb77ac8ff378c126b18af6af4bc3853b42dfc7b7f13d9ed2b7c6"}, {"2cedf85a730cf1b685beda07c884488ee6af197eb48e484db7ff463405ea9451", "2b0f7c30e47a315a1621a962dbd7ebd9dace0060f5ad57d669eee0730d1b642c", "2875faea6e2d98e623589585409a16ffb363c30daba00e9f9c04bf4b418cf286", "1b0af6541739aa72b91b2da05a8efa65cdb26a1d8ecc0cc0c2b75a8bb5724196"}, {"2ae8bf17e5431c4e490ee932e16cb4d9e55d16b0d3ed890ba7a4c793158eef10", "19ac467898482e9dac2baba1562ea504a37247b72f93703f7bd7bd1926909ff7", "218201f9512cbfd8afec322da50972558b183a50a9fb8c7eb66ac8ef0ed5bc18", "1175ad4f6fff309ac6ee415975421776a910ce721b125bfd027dc70422d4a000"}, {"15d86af8963d33f8d74b7baacb24a7d05ca826ddfbd9e4122a4afd51d77a949e", "2eac91f3d0a76bcad0d776d4730cfd126ce6ac45a2a81a05eb65d857490108c3", "29233274802a5b2d8041635463b156d1cd779f0c5009007e4b11b17ce959cf20", "11702232063b55b92cc7491b3a1e28719be901a4a04efc54f844b322fa7ba85f"}, {"0e534fdbb94b67fab3591cb7fadfc9065c1f7105e43d152b843073f252b9d67d", "140894d438cd675e7b728f72c2a1caed5d59314343c3cf7da67a31c7933b8e68", "0838af0f022cf6eeee7f1d7219f42b7e4882be01d8719566a720a0879e612753", "2722e88784451b4909e8f9aaefe2d61d8446242f1e4d4a60a3e205bbef61d200"}, {"0d3858ba3f2abbea60584c7f3a1cca0fd5499e7f578a920c15689fae683e2c10", "0026cca08282ad624eebe4bd2040365e7733245b983b761eb6c3fe07a64a2e95", "07b05395164f77c9b3879c7b466d29445515a358947c9df9cf15f24dbd3aa771", "02ac50fd9d54446077163f0c44a7a071ce1950d97a1fa7323ef069e9eb6c7908"}, {"1d5b1cab3fcbcafd5df48d4ceb8b523ad69c15efa00fb13562f78a36fef1e025", "0d1bfa26a14e8b62348032f24ed9dcd1f1f26903b047bcdf8e55e46f31a8152e", "08b8d63746dc89688ba14b06f43f555c81fffa6bb61b53b10aea8a4a653a33aa", "00d635fc0e8f77ffa3551705b112949d6c4fbd8f4fb725be2a6fa5e80958f925"}, {"0f641c4ec00b3dc902e655e8b1adc205d510aee7e9f1c42dc8fa79d97d1e7c00", "151db8eb6d5afbc184eab5c8ccf78b5c480d96b03037a0c7bfc9447132b6a601", "18a12a1e7147b63046f266d2a0da18b4fe11833b935d24855bda1ef845c1c471", "28b0db4b66d7d70d9ee033e2327e3875455e74256769255725c6efec20755709"}, {"1a1c01ec0431a37fe70c552d7c0f61884f393f61f3f56e65a47e648f7fe21905", "2f5959c33370a833feb0cd0cafc172522f7cc557f1a087163c3e58888fd74d47", "2ed8528b478370dc1bda869a4d532940c09b91ea8d468f5743b20c76c9cfb962", "2a56e627a133a00a1ac09302d0999e93a7523aa2bd076ee7c7e8b126c0386a20"}, {"102c3b0452b4fe0941e7d48fa688a879350b5c853e2c784c7aebbd1f56dbd7a5", "063d1c4199129d0e9455635362c52c9a4e9cbc47b9a552511b4fd2f88c1ee5ae", "037b9e440035f9cd5ba1241bed7d231d0095fbbb9f9c4bcca2348e3d172e25ad", "2c1bfb2c73ba38f332b92270348b6c4e46f5b7db64df4b3286a1660aeec7a335"}, {"13eaab664792f60698f6d3a805b1e1d2b022db50c43ed36dcffe91adc6f6ea62", "1e09b7b25afb4d9e9e54f0c629cd91f9708ea68124e96644df0358e9de3f5a0c", "062a451e555686c3440435990c1d01761222244f37eec105a096b8c085c6ae25", "0f38234ac951490fffe7dafbbe51d8595e45b6cdbd35eac0064de517c70fa81b"}, {"0b43651f0b0428a35de9198a094dc89e24dd312f76f135380d6c1fbaf09ddb94", "0a299fcec2d58a5dbfbe0a6d3959f7eb917532dea24fb005a6af69578eecf2df", "253b83ec9924123a182ea50a78b22839f0f3da4e64c83a0cd2994c00a9ccd267", "02eb4bbb65921bb1c3490c4c8c2fb67f6df5723cdebd72ff6272ac1f167be36c"}, {"2079f45fb6956af1f2bafd60f1f5aaf88edfcaa457f3723586b0809b70f2d347", "223cfa7528ecbd533ef1377443604a16ae29954dbd1af91cb98f7e2985f22e8d", "1de1d663c28d45f511e5cfa2653a8b78ca615f9dad70e9308c78076414c11b43", "160ef23a48521ef3d0580309153f8794a6012c04c7822363e217a06ed9338d7c"}, {"09303702cf5b7ae9b8c344cfc4f5ec65d01122e625b6cc5d8af93bd62c406b7e", "0794bba36cbd98002023c13172a687b872d63939f6a9230b0e0045a3739adda0", "195a2521d450622ad65bbe3af7ab543cbd4f1ddcf996fac3b4167d1b2c111b4a", "22c7096816650e9422380bd97620b0b531b73b964513d1c04567a297c69d55c5"}, {"262e1fef656acd766d8a4b82fd239d0b5944ab1bb6b57816392dc9b0e9234462", "187ece88f5c70b2a65fcf1693dd81b580836f24e7c812b42e9e1e43653388935", "15bc5e9790848aa852a77fea27e7b3889928e6fc8939f632a0aea1aab8368d55", "027dce453827b37e79fef82b56aef1343485310fba4e767a521654db5da30e6d"}, {"147b474a0cbbfc629d30880b6b2d2d2397beec1124471aa3b9f693ed5168496e", "026ae7fe174d87aca57671f97c4dbc5e09a322103cc56d81d3f2d278474e72e7", "1c85cc18a5d56ebc22db1612948045f52f59db5620af270204b11d7aad2ae007", "1db4d3ef5a93bdddcaed9e4c97f7ee00b0cecd1623a3c456f81e44ea688dbaaa"}, {"12e3ea5692729aa114ef93abf94c1333a25943eb830a3620557f2827de98cb6e", "2935b4c6947ac07770473837a9f1e52ae5ee3fe73b1e7b229ab51be7a2d946b4", "1b5ed93c65d71c6ef23e3ff6c1e6110e9a86822f79e537f42a2217f726f3fa00", "03a4a5f4dc95358f9274767da666d6167528ed2d844bc68118de01a0d98c9299"}, {"1eaccc5c180dc5363f973c0a8b5f66dd1278c837769b48722ba5b058dc43be43", "2ea223ace807093bb610a287f523d193560073b9b2485975a76f065cd762a7e1", "224677f701bf32c3d07ba58efa1347cac50c2ddd85bf275dbcc6f8949afdaec1", "09c6c0800a5d66aa73e5332182a2e504c5de9acf064f56d8f68827c3a41b6b8c"}, {"06da02e21a1fa7023b4aa612515bc144e8e387e574f12d6220e7bbb35d432845", "0bc7c2f858f957690d4b26fbccb16665efa4421999d9b51999d19ed4f0fc83eb", "2755b25d79e05e3503c8f7d6aa55cf12cc2d7e56de705287cae328b89e0cdd77", "25275c6705b80bbd8d60f370de08154dd82527b25486bab21dae3ed632276d5b"}, {"0b3c0c13f2436948d0c8549350060c218af3b44f60ddb2cbcd891e1ef565ac80", "27231ed9c7230915aaf9597711abab241171e0e27a4b4a959a502f4e1c56640e", "0194fb9ac68b223d73022523023fe5769a018117591de6891cd152fe8847d55e", "150162d0063ea23267b36654b0dcd02bf312ba8521b4fda787db0a03071402a2"}, {"28a9cd34f4244ad6c93879c90241f0b9646ca951d7188d750e7f3d950ea90a2e", "00d3149066fe7f78404185db5aaa2d5db9ff734e282830395302b0fcded9823c", "2c782b5e544fa56ad3447f0dbfbf9e4dca3e60cc7830fdbdb21e5100f24db2fd", "00ee274ea6df380969371c0daf34fd2420dc55890ffa6a8332f19dbbe5b552fb"}, {"0ae796b57a05ab6977f73d0171078e32ebe92dfae94292a9bff12015dea15534", "001a8ff98d149891ae02c5f3e9ef14027085663d5fe4f655f630557c68f75d03", "28779cdfbe806a646915c115149dbd990d4aab7f4613f7d604593f90609ea56f", "15346345d59885de47eda8093b55b75fdee1e3b81e06cbe18816bee8f79b7371"}, {"0846615c138473cd30621416f2efac6b949bb334534c90cf149b9526a4c2ea46", "09671ae0e2a806e623b4d60ed15eea55e501ae41c781064289009aff152ad133", "1b8860f58938bd656cc7e3bd45e70ec1e5d88167e7f465f175e8940bc71ffd64", "1c73e79002a4d6f102921cf947c6b0a51885acfc0e7234a6ead5be034c13b97c"}, {"278f97261eacbe09ce5e35436bff601f56ce67a62c3f619924e1409072d013d5", "2650415f4edcced1bcfcc18211e84e2eda41a6d7d7afef402d0695730ce9eaa7", "2a55ef81807af7d514bd1790e32d2dea7d6624acddfd379fc1db525c21d2d246", "0aa3a82171b1ea6703afec3707fb98a6f3fe069d9a1a193e9fb21216af8adfb6"}, {"18da37f89c631baec3fe0dfc5fbc2b7abb8a0761c81cd3fd511a56469c5820fe", "3048d89b856b0730550d2a91104c4386db5c5e134dd951ed95293ed66dc46823", "17152a1e8c25d7fd372c352d2f607355eab904d50f3c94ce35f77c19ace8746f", "1deda369a6182a492f049ac9817f6ef3ebc158817b6b2de1926fbff8b874041c"}, {"150fb68164e011be67ffde47c166279ffff5d701296d9250d039a8e755fcdc7d", "1345557d87b9db9d6d5ad1e2cd3bd266bcb27f5c093414b34806384e03c062cb", "0b35ad4db700375aca9e444dda33a08ef942f21bbd6ddf30f7e268203782fc4a", "26c1bad66fcbe4fee85f1f9e73349d18371043a0dee0599055d617ec71996ee4"}, {"1a148f8a457b42830c34545a4accdfc1fde289f09d40ab38f11a967fc541f9fe", "25763da3d512f69548640fe115796be216f728649ba652697198861b688d05e0", "2576def12e58d34cccd48c12e963b93a929b91ed2906eb6e8b9848a471787e86", "14c2768d1441466bcd54bfe022f50799f32348790f85a61db590197c093da99a"}, {"18a8d241b5c76d61db23696d79c634425214c8f3dc13d414411ca61e8a72223c", "0963977a9f82b767a9e2aa7bc2a70b82f440ae2882bb78ff46fc7199f2317e2b", "0501378508f3f4c4f569eeda41626fd94d59af82c05b4fe127a7c3b9b20dad1c", "29c1ef75583cbd175eb34b2583b43e7e4ecd68008268e31c945127dc15737836"}, {"2f420a3b42cabc78c3804794232a9276a947c497e8e590640e4ac142f6bee057", "27457c22d9a344a4609c35473ba64046251a404ff60da967161366eebdaeaecb", "016faa994d3130859029a6e8df7ab268361a3a3ed3e85bf13ddb8b5edbafac85", "22f16cc6150e13d909d7b928ac7f02dfbdd2ee1be7c9d052411becb1f6134670"}, {"0fbb74f3c014aceec5adb459988a5a3f66649ed09788ffbf7e2b21691073f87b", "02bc677759f9e91d6b5d92703a350f0742c918cbd9f2be43399e2e373810a57d", "270211fbee7723a8c9b5d1640a843f29335953587cb641eef5bbd4d035c3e43a", "00908849226addddd34e7cf6c23f289b08dcf801a7fc1eaaa271143df3cf44ea"}, {"1dfd5b762b89ff0c9e63bcc16159449993f32835708bf6d22974da45448fc174", "2a615689d584dd478176c92e4b2156c429cb9f4ec29f32029807cabf1f6f5729", "2f204d7ebba1c7f34872c21a0c12f3ab66ef17394e720a465cea6932766779c6", "086398612dd2fe9c17f91c89568f91e3cbc722381cbc713bf98ed1c47d7d4e7f"}, {"0e2082a73fa233c9a81849ea7d8892ee3fb236c97076eb7c2a5fe4ccdf426cc5", "029576c3743409b9cd59c2117f6f0a8dc9d8684b72630928244a33c6288282f3", "13641d459a687b166340190a9c79d9f11cf0ea57c14d975310328a9e24b8da89", "2108d0dba482c072da3bc0d973b7c498ba8edcb3cf38537da0ee3a7e715e0fcf"}, {"2f4fb50043d4f39e59e2a86ab5df51a7463929fd9bc082c5bf9136182ea86e19", "14d5e53acd259b59fe35123380d0aef9fbe7c0429c456ffbcede7db0762c2a2c", "20f4c1f1fcf7a77160b02c7b3abe1d5851cb7be7cc3a6d62134fe449e86ca99b", "0c0311cd9978a19fe59f6f480059120c2759f74584e3981813f13d871ddb9814"}, {"124aee738cdabc3fb60e24f2dd845a7cdbd3939aff9e9f5cfc98ab6700c6c2ee", "19bb3a958205dfd40fc9c94e375016cff6a17b3b3075ebe8b6b34a09e7b1fafd", "2595be6d65d511b6e0fe790c5f4aad0a5213d40cd490febbf4100f4f44314db2", "1dce61aec04e8d0aa0aef5fc0492b1fbc292e99a1e485e68456c0985ae4095af"}, {"1ea0d1c7bcecdb91fdfc20c08d0f8fd4e27cc33e667c34e00bf9dceb6d5613bc", "1ea08f34053277d81440078e2f949030bed602c3b5f10e639f049813a6a66bad", "2b6ef02aa647456c23962109da0b43955383edcd2ca4531ce1fcaa84ef6c20e1", "2df64820d779c6066ca83e9121e94ce9e401185c05cb9422cc23be91ce654440"}, {"22e844870735033cf76801a3c7c35a37ed4b2ccbf21c745c745337d20ed0c26a", "0721d0db4965e7b114a91e3f39cb17038bfb0d4f24db1dd067def7353656d5eb", "1aa4cafdf3441b1209d8a46accda68eca48b24b6072280b688703cde443f98f1", "15b116f08da9d67ead6b0454dbd96b33c6ecfbed71751bfba7e69849f5e66bba"}, {"04820a29928ac4b2099d12bf7bfe733d39b76b5cc70cfd95150e37072847e635", "19f42966dcc7d571259a78a79e70fb9a7d0cbfc76e552c34581a3bbf0528d865", "05c40968820e343b3957125f484a59327bea56694bb9ee0f2b4d0281b7bdb92c", "143e3170fb81d62a3fd2f37eb452b8ac41fd64a71e9b066ef71b035df0411b53"},
	}
}

// STATIC_GROUP_MERKLE_ROOT is the root of the Merkle tree constructed from the STATIC_GROUP_KEYS above
// only identity commitments are used for the Merkle tree construction
// the root is created locally, using createMembershipList proc from waku_rln_relay_utils module, and the result is hardcoded in here
const STATIC_GROUP_MERKLE_ROOT = "25caa6e82a7476394b0ad5bfbca174a0a842479e70eaaeee14fa8096e49072ca"

const EPOCH_UNIT_SECONDS = uint64(10) // the rln-relay epoch length in seconds

type Epoch [32]byte

func BytesToEpoch(b []byte) Epoch {
	var result Epoch
	copy(result[:], b)
	return result
}

func ToEpoch(t uint64) Epoch {
	var result Epoch
	binary.LittleEndian.PutUint64(result[:], t)
	return result
}

func (e Epoch) Uint64() uint64 {
	return binary.LittleEndian.Uint64(e[:])
}

// CalcEpoch returns the corresponding rln `Epoch` value for a time.Time
func CalcEpoch(t time.Time) Epoch {
	return ToEpoch(uint64(t.Unix()) / EPOCH_UNIT_SECONDS)
}

// GetCurrentEpoch gets the current rln Epoch time
func GetCurrentEpoch() Epoch {
	return CalcEpoch(time.Now())
}

// Diff returns the difference between the two rln `Epoch`s `e1` and `e2`
func Diff(e1, e2 Epoch) int64 {
	epoch1 := e1.Uint64()
	epoch2 := e2.Uint64()
	return int64(epoch1) - int64(epoch2)
}

func (e Epoch) Time() time.Time {
	return time.Unix(int64(e.Uint64()*EPOCH_UNIT_SECONDS), 0)
}
