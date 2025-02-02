// +build unittest

package quebecoin

import (
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/btc"
	"encoding/hex"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/martinboehm/btcutil/chaincfg"
)

func TestMain(m *testing.M) {
	c := m.Run()
	chaincfg.ResetParams()
	os.Exit(c)
}

func Test_GetAddrDescFromAddress_Mainnet(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{address: "QhZsGeweJd7zAM9LxbGfHK6npmhxStM2uE"},
			want:    "76a914e5f419d3b464c67152fb9d3ecc36932d5280673f88ac",
			wantErr: false,
		},
		{
			name:    "P2SH1",
			args:    args{address: "37P2coFRto5CRZcduV4N3mdnRb1EdZWtS7"},
			want:    "a9143e69d8c4772eb34d77c96aae58c041e887b404f387",
			wantErr: false,
		},
		{
			name:    "witness_v0_keyhash",
			args:    args{address: "my1qr9y3pd7wy7jjpqf87qsmp08ecppc0p2jxhfcfc"},
			want:    "0014194910b7ce27a5208127f021b0bcf9c043878552",
			wantErr: false,
		},
	}
	parser := NewQuebecoinParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_GetAddressesFromAddrDesc(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want2   bool
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{script: "76a914e5f419d3b464c67152fb9d3ecc36932d5280673f88ac"},
			want:    []string{"QhZsGeweJd7zAM9LxbGfHK6npmhxStM2uE"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2SH1",
			args:    args{script: "a9143e69d8c4772eb34d77c96aae58c041e887b404f387"},
			want:    []string{"37P2coFRto5CRZcduV4N3mdnRb1EdZWtS7"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "witness_v0_keyhash",
			args:    args{script: "0014194910b7ce27a5208127f021b0bcf9c043878552"},
			want:    []string{"my1qr9y3pd7wy7jjpqf87qsmp08ecppc0p2jxhfcfc"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "pubkey",
			args:    args{script: "2102c5c7165eb66f35a120f2f9d97fa61b1be6c621f9b868454b35a284fa7ecc831eac"},
			want:    []string{"QSHRwNfPRQYWtQknx8BYJ3Qiw8FEbXamWr"},
			want2:   false,
			wantErr: false,
		},
	}

	parser := NewQuebecoinParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.script)
			got, got2, err := parser.GetAddressesFromAddrDesc(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressesFromAddrDesc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

var (
	testTx1       bchain.Tx
	testTxPacked1 = "00004e208ab194a1180100000001163465df9bb21d89e90056f11887a398d5a313aef71e3974306459661a91588c000000006b4830450220129c9e9a27406796f3f7d7edcc446037b38ddb3ef94745cec8e7cde618a811140221008eb3b893cdd3725e99b74c020867821e1f74199065260586f5ef3c22b133dd2a012103e2e23d38dc8fa493cde4077f650ab9f22eacafd14a10b123994f38c9f35dfee9ffffffff025e90ec28050000001976a9141cba92fe1510b8c73550fd4d3e0b44acdffcd12d88ac79c268ba0a0000001976a9142f86cdfa98cac89143cf9e3d309cc072caccdf6f88ac00000000"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:       "0100000001163465df9bb21d89e90056f11887a398d5a313aef71e3974306459661a91588c000000006b4830450220129c9e9a27406796f3f7d7edcc446037b38ddb3ef94745cec8e7cde618a811140221008eb3b893cdd3725e99b74c020867821e1f74199065260586f5ef3c22b133dd2a012103e2e23d38dc8fa493cde4077f650ab9f22eacafd14a10b123994f38c9f35dfee9ffffffff025e90ec28050000001976a9141cba92fe1510b8c73550fd4d3e0b44acdffcd12d88ac79c268ba0a0000001976a9142f86cdfa98cac89143cf9e3d309cc072caccdf6f88ac00000000",
		Blocktime: 1393723468,
		Txid:      "b01e2eb866ed101ed117b4ad18b753929e85c42e3d8add76bdd16e5c00519dcc",
		LockTime:  0,
		Version:   1,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "4830450220129c9e9a27406796f3f7d7edcc446037b38ddb3ef94745cec8e7cde618a811140221008eb3b893cdd3725e99b74c020867821e1f74199065260586f5ef3c22b133dd2a012103e2e23d38dc8fa493cde4077f650ab9f22eacafd14a10b123994f38c9f35dfee9",
				},
				Txid:     "8c58911a6659643074391ef7ae13a3d598a38718f15600e9891db29bdf653416",
				Vout:     0,
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(22161428574),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a9141cba92fe1510b8c73550fd4d3e0b44acdffcd12d88ac",
					Addresses: []string{
						"QPDtY58kzPHSY4GkALyN2Dv9PyrKzXokSk",
					},
				},
			},
			{
				ValueSat: *big.NewInt(46077100665),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a9142f86cdfa98cac89143cf9e3d309cc072caccdf6f88ac",
					Addresses: []string{
						"QQwHMhjR2ioPG81K4kCbp33nWN3Zn7khgk",
					},
				},
			},
		},
	}
}

func Test_PackTx(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *quebecoinParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "quebecoin-1",
			args: args{
				tx:        testTx1,
				height:    20000,
				blockTime: 1393723468,
				parser:    NewQuebecoinParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_UnpackTx(t *testing.T) {
	type args struct {
		packedTx string
		parser   *quebecoinParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "quebecoin-1",
			args: args{
				packedTx: testTxPacked1,
				parser:   NewQuebecoinParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx1,
			want1:   20000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
