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
		{"c12d11c3b8ee882559ef48f7a42633d81b1e13fc589e6caca028281a1098012c", "b3ec8a9138817be401b9ed2c683f6969d87f70ad87cbc514dee24146a542a71c"},
		{"1a19a8d1709ffa10882673962fa4b6fcecdd2ae3c95f9cc2db294633fd775109", "eb599e4681f8fd50730d22d2b0d0e9e4efcc659d2d310bd6cb3ff5600a81300e"},
		{"6317e72b74ad1395ce77777bcde06b7b5dd02ca69ad6fdc7affeb3bf4d6d1c13", "151a533bc8927e97be9ba2173644282e3aeada4f9a65c4cc72eda27f3ba10c11"},
		{"d5d375f61fc9f9b4702aee9473ce304abca838d224dbc0dcfabffa0e3d04452f", "f538609579c413bca2d395a772b026498eba0c00a1cdd5f7127d526bd96da723"},
		{"f47b36d2df712cb16eb21797c2d0672e12a60e0a7807edcfb15e7213f072a20c", "d4d12cce11c6d8311c15353d35926f298b886ac66cf6420c47ad178231632608"},
		{"3a8360de57f4ab8cad701459a73019c71b0d84927dcec0384bafb2586356080d", "08a8b713036264b878bdb8051cd6a6ccc9acf9e094daff690d167e699a90c628"},
		{"661b5eb69735c9e8a181bcbf34563b96aff763d4996d60ef88c439e82549b622", "c1fffac2bc6d8b78eb24ba052d72bb88fe5cdb40e86eb3ebadbd57aff47b1e2c"},
		{"79fe1ac6ee536412d737091c7e53f003cfd4e4d1a96b1ebc1d27faff4527101c", "ae1b33ed18cb164c4b3227d8acaf75eb480a58c07d04792361d3e7688437572e"},
		{"88f35905877c4417d418a75994eb2ccf7da052b1032bae782c935b107bcaf12e", "0e3cf8fee2c863470dfefabea0162e98a3cb0b5aed1387d9aa2990710ae6b705"},
		{"771b20194afc2d043d133213a39c99f88a50f6459eb7bb7a8b19cd468fdace1a", "408f8624260e2a85ea354959cff0e7550d89ba666e4a5d646fe1ab8a9b253d2f"},
		{"1e40c14bbe937cc3baac06f4b9c1e9d15374064a12462a505ee3c85b59e12d20", "c5ea8270d54c3919ffc5659bef2b4f00f55a9b7e8a655875ca864837cf66561f"},
		{"23dae2b032cf45c17378778786b91fe53c7aad7928391c5c4613b5683ef22c16", "9f6e20cc9f6df9e64cabcbb765c92faabe25257d92a00f746d13cf1e5f113f1f"},
		{"f3213ee18ebd73ee5813ef6267d54e0af7bff2a62ebd329adfe41a42e1d7451c", "839c3c58a2ecc4c181b8f298ce6aaeeace4e0777e8774cc1fbc4279c3e001724"},
		{"f96996f6df01ebe12b75f96ec634458e5e96d8bc9b2983b211d943d8093d0b0b", "8ce6a29f023aa78d5a4de75b3e12c3cc673bba9dcafd5a2968f4a1b9707c5b11"},
		{"c9d5403ed3ee20b29cb49c16312464a52109c553cb7c3e2e760488bc1f4f1b1b", "96a4069fa8485b11c9fcd9502c384caa6a1db2a90d45559a2aba51d5ea21782b"},
		{"3de4f01667d317d52e2718926646dcf559f2fb7266b6977a47401b976b457b06", "6eff6bbee3d3bd885c7f900b362dbf245e97a523afdbb62d25eefb9afeba4c0f"},
		{"89453ee367a782654502062cba7da961d638b4160121ba0ac88abe7a2f82ee18", "551f19081516e49a4e94dd97af53da2b35b603c321ad02d29f85b9d5ff12f802"},
		{"32c5a512efe5d41504bea17a4e3bb865c8d54e9c6732a3bce77962e52e1fc414", "3b8dab877f10640659e99bc6dba2367664aca2f4e3f87b38ffd532f326332f12"},
		{"a74c13d3813e7452fd8680074c31dcc8d6e79e95086175205129f77f37feb129", "c55985f3a99a3fd1a2619a54b87e433e4114f59d1a59c3a911e957790bec471e"},
		{"9b1c420aa3dc252e9f290bba21136c97d4bea618ee298a4167eeb445d3b6d517", "66f0efae6899a6f851a2dab31df7e936238b8aba4e961c9b65c6b1d113f13e23"},
		{"d4821c14ed5e7b8f4febc8cecc26e1d0ae6fd97a7324566c06920300110de414", "9653e0cbb946b534f0468e42bf124e5806adb4bdb93e91665610d69037b28b1d"},
		{"26c453348bae0b3398691e39d5032021c15f7fcf8efb5666f2ffb3bf0c609804", "a366c39b7cfd462063aa394c31642d36cd3ea5fe89ed0d7db423f2791cf74429"},
		{"4a50a8c66d78ea1850ef9e6ff3d082fb6aea7a380d267b89e449826c8e7c1319", "e353c38ea40b59565dbbdcbc03c0c9f99a017551ae912afca8d5b6cb028bbd2d"},
		{"d04788fe9a750986d14596580b48edb0415b18466872d4b431ebbe80c0276f17", "234ecb8bd50f758f80a48484b6997b4d0a64a1f67941aca1bb31661797d63010"},
		{"c5a3e4885e1b16a26ba3d938659e6a4d37e3f66d3f4a8d82ff6e22145ea58303", "2921633b6bc30dacadba9ee3956ec7e3024971d9db600ef99c4ffeb1ccb8e425"},
		{"485887237db885ac07ce8eaf4e881fdc3fbcd4454cf0c56f0bee6b7213de570c", "061186e353aa3a59c4bc1d98e926c17b3450081dbcae63d2d0841a3fb3cef422"},
		{"ca0f9ff876c78957564b303c9e99598036293efe635ec29e0e4bbc59ec59d106", "48e04ce11bf78ef28261067eec8e5a47ab8632b2d35dcf2e28d229e1e2894714"},
		{"9d7965f433303388ced9097e0563c2871c7ce0b286f108bb53e7a68f77102b24", "b6afb6e2de8fd30417e4b8d1fe4559ec73aa9e96726d0448eef104a0f099eb2f"},
		{"db1ef92e473d8bdad5654525d9a9fd9fc0febfe7101eed67c8031d697fff5913", "34d5b8bb8893c4f4fcf0aa4cb6bc13187bd4867bf0b4b32b57387bd371406f01"},
		{"d43e059b5a5a2cb6b4200ac3832fd4ae6a33c69bcd784eaa3e662007a43c2614", "560683915ff850883b2344e9c64543cd40b2a544c099edb1e37932a7c21a1d12"},
		{"a1cf07a46e8696f4a6f6838d246c4e9fbfe6db33149c99fa563f233b16317e01", "3904003e9ec020a567d23301a8f381a7395d129020ad320fb2b11f57680de027"},
		{"178c9c8612a61f62506da40443cbf6d6fccbc9406303b6f88d9536b42c506826", "2c81906219408328fa05a005247c9baf796c459ecc3ab0e1a70195c180e47705"},
		{"f84b9362f81ec147c40f43cde64f3ce883bd80b40230c435978794b54431be1c", "ca524f39724400999116252fdd67316cc0caf586c3ee0bd98c132ab2fdb7f30f"},
		{"a2fbcc2ebb6f728e42c2967bde68461af69c2b10c5305fd40053eb01d1db1e22", "4ef48e82ffc90c273c6a1627eed225a1ecf5d34bfa33026758306601a08ee71e"},
		{"91a7de9363388d15501cf72449b053a036ec5fa16faddb0bfdb6aca0a0c1f409", "fa5bc2eb977165e92a45d92d5da48e0b1e95e2d13e2d8d42dcf9e99f8761f20a"},
		{"6e2598bf6a6975a578abc5615e0791c678ff1776176a771f025c17a67777791c", "22afc07a5715a0d1a47ba27403e83660837d2c7b9a5902c22c0fed861ff5ac14"},
		{"e788d7b78798f2edc1d5575e35dfa3c17b6c15b6642df72ea6ee28297422b011", "ba9a4176a20d61efabf8b3a6e2197b8dcd26b0337c26b567c2fc4b3ccf67aa15"},
		{"bd13c15935c3a49b2f19058e784d3bf700f4c06c0641fa771822194e543a3200", "1535c97c68abc851042f117cf98be4130a25a49acf5f9c910babef342db1fa1a"},
		{"7718d0013fe1be1715041b7df3372f21185821111966fc40c5c29b948fecf60f", "e476d8441b12a235c48c24cf1a4edd1b9384c2531d70dbaeaab891aea4c39a09"},
		{"857adf44efeec3ee71001be5172f0796a56021cbc94273ae4c8a58356a0d2003", "358eae8e81fd089c3807354d20cf1f878d39b1ce757126e787d4487af65d7821"},
		{"8dd2491ce49ef575e8e0ebfd675b6b831e8d19c90d6110ebe57a60d3a9fff622", "88ef9b9cadb4395c03d57ca9c0a84fc76988b1285d716d4ed3a6340aa7f85a28"},
		{"95a421fd9f866bf28eae38fff084ed0d300ac08c3c020d73e6c0a432e5731313", "d86722ca41b4dacfcf1bbcce9a232979722e228e15fb3e2048b8cc271b021726"},
		{"8509921c8c87eeddba208836e3a70d570b39d14d8fc89a0cd988ace585a3ea2d", "1970a24152128fd6c74ed49315ff705d5af4a58b4dac87d8c82f9be6a6d77507"},
		{"f0591ba2f822317b6d5d8b771474ae9518e4d36518469965d83d84d5795ea513", "1d78b5d07a822537a1bd8e8a2fe2fa9acd4d858aae251f5e33e57d1f7c462300"},
		{"a302906a3fbf5dd8753edad674bc00b9397d1a5bc3dd1d229359044ffa346b0b", "de7690a0fceb4c071f52a09a1fe3e872a74a33c698792a0c30e26fbc8d8b4d20"},
		{"6b425f3cfd5f66616556d9e16698fa1d2cb2e6ea6149b75089c0c403d52bbf07", "3254d4f64d9fd0ab8269bff02865dc115841f1717ca4408c8fd21830deba4900"},
		{"570a3f9bb4a293fde27fd13f1407a0aef5c1e1025e2417af400d5c40a043222a", "f3d481d495572a89216be3bf4d3ba719d2c81f59f67ff825f2ac0bed67ab2a11"},
		{"d2a4336cfe79faa8695f88d74b7786ef418bac6021a9c4ba1c3db8e433fda122", "618888220de5b3f2eb1470ea0ab8188d5385b21e1eef64a691b2f31d066be12d"},
		{"5e9db678cd1dfd7e0c598236d25f27b34139e26e5b15b032a68de05b0e394e28", "4f2379dc6a1212d0b7029dc3248d0546d003edc23329c848ea62442e3b2a280a"},
		{"042dbc17ec31dbd098c87c98fa9cd5d8ce7716045ef9d93aab3c9d6bf6f86e21", "d30841c4768e3b902d9def72131244717d2a0341540e71b51321aabe81cbcb08"},
		{"c5cda9e62ddff24a2f14c8ec8ffd7746e230b3023bc2f87353a6eba7d1e55f1f", "d6252a48c7baa1b9194d0d12a8a07b97f2b624234b48f5eace2d1adc958a8118"},
		{"b751b8e0c753c8dd5a07293c0dcc51448a49be3cfad6c8d3fcb8e15703a1f402", "297af6aed5d949eb9ca3ce7f0f16ab270fd509ca350376cecb844fc55606f523"},
		{"2c72a6ac20aa6c8ad2500bab50c90fa8c5b2150a17d3f1d249faf29dc48ee81a", "c05528b87b7d9b7f1c96937116cb5b6c1d66fdd7678332e257d95601e98bf108"},
		{"843a2f33499e417fb3370d2b35170dfd89ae3d7296bc2552611a1f04542f2b15", "85f5166a1b5c384f6bc9f59e779c9f866c4a4d00443372cd433b5096a7a77e08"},
		{"329f698e99433a9acfe5bde3662d8e2c05b5b68024d29af1a59eb63d3722e40c", "910b67959ff965ae27ae8679e07bc2dfd3b6f567bdf74f07b7dc3b055d883430"},
		{"2e00f33354bcace1c798690fdae14a40b8b0d5d922c5e7d9b8a7bb17ec72a40e", "4b50726e2c50f4e404bbc39eea2a8fd711a6cbd194489c4bedce99f32cebb81e"},
		{"73c09da2c4cd22b3890ada1d6045a6877d558ea5c3a7088fcdd3b77b229b7620", "f500793aaae728efa2029825185175fffc286159319347d10586b8a1de01b613"},
		{"25c8efe9ff791b4a0f4478a6dda0867d8df396aa51044c6d6b1ed9427d117c20", "eb57c5d562ee43c72d8972ae0e8c170b3a7f0e4c89ba67e82186229adb904706"},
		{"863f44e00121079c54d36d7cccc1da51ff5900610386fdb8bc36b3a47483d72e", "c30fd9b1b05ac1a347f432d65b68c82476b4ec0994fa00cfd90f1f7db1571d2c"},
		{"6a7311e3f18945a8709eb5e90021a8139375b5b68af6c9cad121615a80ee3f07", "1a3d8faa7c7d38d5acb627def5b070d8f5719189f7a25e3861c0a9a879cc611c"},
		{"97d27ce44b476664863f34a2073278dd5ef1c8623771a9813fedc3a1455ce92f", "9fc429eafee88fad27dd8a0b05087a9282c926353152c8174e774f34128a7d13"},
		{"a88ab45b5ea8cd975399fa39d3ea5b04b12adc705732b54ba6e5af494863c310", "2429cf8b01347e32d2774cc4070928d7ff96ff585e6f39e0a2e06fabce53c81c"},
		{"51eba466f4662972616dfc4fe846425b245ca1405730b6809882f51f413b8526", "f9112ddb4c80bb385a3938959a750e091c3bb9b6e16d717db46c28efbe273a1d"},
		{"45e9ff284aa8b4c825ebe16165953b186bbc0b62f209f84dac2eee3382a94e2b", "d932afbbe10120b68c573e1844a4f8f87bc93ff9d359d7c15621952e4ef9821c"},
		{"d1807c403b8ed2e8022db73486ff6dd2471872404accb8208cda3d757079041c", "7aed51eb6e3f042a32e44f7add13f9d8cc675839232323094692fa9ec0385e19"},
		{"b9c93861237f423f8cb2e96e3a92ba986f290f3852475d9b62cb21a445cdc201", "e25ae2bb31b01d5d80186f906af11d4c7a6ed172a5aefbabe3b3eeece6750816"},
		{"558ad70ccba7882b6f20cd8098f52b8288afdee8b346bf4db33b5deb8153c71a", "c651377b6f9deb188dfc868df0157ee50dd5f9f7d92ca0e69e82f03355af9821"},
		{"b6e4ff38fc18fcb2ca63486314db80183b35f1dc8082e8dffae0726a1c284c25", "2eedb645aa09985bf178bbc4c5417f8c1a9907440066096111292f2e72e9a01b"},
		{"c264ab7d9008339abbc1be91bb96eed30cc5d051d8833a3f5cc94674fccd8627", "8fd732c230f79e11d56d8f7cacd5f7095e4ad1a80a3c79b1cf42d9733001fe2b"},
		{"8b29b2811047827f356a57f7166f8b3dd4a3aac23b02522daf007c677295801a", "e4b4d00d5d3eeb087c2edfdede5eb92ad39974c359172913abc78e5a5c78ff13"},
		{"30be5db463aef5665c8699f2e5fc69ea2ca209290771e2aaac3b60caee6cf22e", "3c5e974d664c03b13adbe5ebcf9b03491ed0e4c50095297d7b3115804274d70e"},
		{"db16b337102ce1b0932fe6e841fc1e7c01473ed4f3765934f2275b821d5b5d2d", "dbcb04a56099034b4eddf402c08810f5842a74d5312cd5fd86d9378a2da54323"},
		{"6f75a23af554d0b3f5ee5a48b5ec1ad8fe9a6f7c2c64b9e44bd9deb644212e17", "68cebf8d52280b6484bd14f9b6bdbf89a485fc6f6129f49494bd7c1b40c90624"},
		{"e3d00baa245cb4f99dcc282cb33121dfa42c3ae1524139c5be17a043cdf65a28", "5710f34c928c76f21871bcb63731f3417f1056437397b083d095e7fd3f49790f"},
		{"f382f322140415ed6692583c594e8d8fc5bea0f027a159ad01df4a3942771100", "4bfb6da22da207b0935868b7ac4574bea7f3358f4a281837e56b1fa3147cb40f"},
		{"33e74fde6f16209c57b24d496fc87ce2270dd2f3b04a9a5a701ec743ed9e1d04", "e1e847e1ba408253c0539af6a7ff0a8700802ac26f8f7aa68906471613f8371f"},
		{"3e4b5dc67f25293d3c432cfe6e37ac7905ae19e62c7836c8e1a05b5822ad432c", "c978a79b21c177d102af936de352d5fb2862396157628c8c53b259eaadd60303"},
		{"1468956b2009da0abc540721681516d2d836fbb19692276d07345b6706a53129", "eef70d99244f8e5de8c938b56d3079990652e399edf4996c7ba3090bd20e652a"},
		{"c641b2667bc124b26572f9fbbec9ecf839db74c9edec9a75168579b71cbb9901", "83b957a57b5ecacc1a4b0231795be7013e488ab0ba2e7cb4122152aa2a14ba18"},
		{"5bc5c3903a9a19dd230310422c11dc42c590c949580f37dedc6bfb528be5c62f", "5e2dcafd8d018dd8d3ad1e5a7adf58605cd8628dadd96ee48f32bc0f8c4be41b"},
		{"9766e135d8b9aa253c90202454fa824b03b9d2d25e0b6c18cd99d87cb328590e", "d3f8885f3dfa8a0416937ef89f89d2ee7df9e71852f09f812ac6d7934fccb60d"},
		{"feb64610db2ac2f01869a198f5a3fc524d6cf0bd171f118bd291c50db1d54a1c", "3e27cd9b28b288fb3d1953a7355c986c88428a0a95b56ee39f7e5aeb0bbfdb0f"},
		{"571bf13dc817ae45281208cf712cee1917900e203be6d617984bb493e0c24c25", "132f0e795cb6f5127e9fbaff53b28e4baf05df08b92dcbdc05d8ea2638d9e70c"},
		{"fb649a934864788acacdfa654bc262cb71af2842f0e0b65054f37e8bc5332d08", "67a5a6e195a43d1492461d65ac8dd2b254b1467bfe85342cec8fe6ae9892ca1e"},
		{"bd6d4da3fcb81b710dbba70051fe8d565c47419517fea3d20639667edb415413", "07fc16e5a1523005ca08be860c4dc58413b773c1c15f9079c6d373f9e2f93228"},
		{"e49799e28327a6e8f4a1d1ad7290345a37a263284517b094d325e7e17593971d", "641a4439c23e414ef21dd7a563cd75f533a3e26b1e11f5d207d29629d5c4d416"},
		{"c0ba522a52198ab0a79d935b17eb57611d141f0ac3864e2a37439e4996591e24", "78fe086005cd3baaa5315cd138530ef4d7f6febc6e427cd71625329c56419312"},
		{"871ee440ba18913aeec0d7fac20c9671e4ceba8e1cde2dc74e2636ca57de6922", "9d2e7c22c6b8b723b06e5960c92b1a7e6cc4cd11619ec7f21b7c1543103aae2a"},
		{"ab5a1272cd4e16be953511a5c5ef9ea24f0072f8bd976314d260757ed0b30c12", "ebfc21341bfde18f6f7fe1b883d83b43278a635b5d699525aeaa2eec2aba211e"},
		{"2c9e91994096c903e90144689053f6f3d9645bc6e11ce48e82facfb03551c41b", "f9cb618cc78c0e630f3035da914c8606ac1b6629657210308509cf6724748300"},
		{"ba41a3b75d7fdd1962feddb3ebabfb1ed01480334dd3bcae3e45f80db0353123", "39eecd2d4a751206f4aeed3dba6d9acc0aefddba1897eb91731f767ed94cfd07"},
		{"ccf1feef0bca203265ffda1e22d88c7d23db4244658f8b3629cc1c7bec17fe02", "093c3209e63e409899050e2b2e17b6397a9e6c9f267056b1300814d9bfadf80c"},
		{"7c2f59be680d820f1fdd4b95982b31931cf3d218088e36f1400d07089f1c2211", "67eb216710fae6f8cdc776e8edbe6adedb670d2ae92a399e80d35ed1dd82de16"},
		{"94063e3fa709f74b22761cbc400d3b7971b0e32d75de9618c11caa06c6d0c012", "c2eefc502f09e9098c554d7db21cc4ebe3432baed062fb7f1a70d3ea76044d18"},
		{"bd8a78715e32d4d7b263b2f358509157a8f1488a48860cb4dae04501e5040926", "f9dac2ab11885c3478469582ec619714623485572a65839aa6a6254c7fbfc914"},
		{"c742d8f410f594be95b9c70f30ab2b3c752388f5d5c139653e3e1f46a3ea1c2a", "39a28b57b0341b76c9a6d8d4502702aa79f03b6b4c71b4a8b16ee73289f9a405"},
		{"b74ce76b34b4e0bea87c576b4185f6e0e2fe60a1ee29a7be6685ab06f84b340d", "fabe6f436f34a98de98776d7170a537afdc4e697933ef83f3ee083619eb6550b"},
		{"89f3b3d0a0563fdb52d340d60bd4a94acb8e9fcb1a078b3784f5d5dd0a76bc2f", "54e5e2dc8bee937a903dbb41fda7d26855d1a852c10f86e60fadae5284a2d82c"},
		{"6a2b21264c42a6fe6968eb4d9539f7d3bd02b0598c58c2a4e709249016720b0d", "a02038d629f056214390c7c3d07b29d9fd2187e671bf68edfc4c4e6d215e2a1d"},
		{"a7f518c047cb8af54cbab674f684d2114517a5ece15b38511333fe60fa75b10a", "e092a4f17f93aabe3b062cd0a41321a3cef624c1b6cfd943d3a5f1834cb2ae03"},
	}
}

// STATIC_GROUP_MERKLE_ROOT is the root of the Merkle tree constructed from the STATIC_GROUP_KEYS above
// only identity commitments are used for the Merkle tree construction
// the root is created locally, using createMembershipList proc from waku_rln_relay_utils module, and the result is hardcoded in here
const STATIC_GROUP_MERKLE_ROOT = "805be2ac92bc8b21bf093440f5a8055a8a4ec7bf5c5af5e22680d9123a4a5c2b"

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
