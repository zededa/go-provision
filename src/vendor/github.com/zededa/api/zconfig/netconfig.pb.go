// Code generated by protoc-gen-go. DO NOT EDIT.
// source: netconfig.proto

package zconfig

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type NetworkConfig struct {
	Id   string      `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Type NetworkType `protobuf:"varint,5,opt,name=type,enum=NetworkType" json:"type,omitempty"`
	// network ip specification
	Ip  *Ipspec               `protobuf:"bytes,6,opt,name=ip" json:"ip,omitempty"`
	Dns []*ZnetStaticDNSEntry `protobuf:"bytes,7,rep,name=dns" json:"dns,omitempty"`
	// enterprise proxy
	EntProxy *ProxyConfig `protobuf:"bytes,8,opt,name=entProxy" json:"entProxy,omitempty"`
}

func (m *NetworkConfig) Reset()                    { *m = NetworkConfig{} }
func (m *NetworkConfig) String() string            { return proto.CompactTextString(m) }
func (*NetworkConfig) ProtoMessage()               {}
func (*NetworkConfig) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{0} }

func (m *NetworkConfig) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *NetworkConfig) GetType() NetworkType {
	if m != nil {
		return m.Type
	}
	return NetworkType_NETWORKTYPENOOP
}

func (m *NetworkConfig) GetIp() *Ipspec {
	if m != nil {
		return m.Ip
	}
	return nil
}

func (m *NetworkConfig) GetDns() []*ZnetStaticDNSEntry {
	if m != nil {
		return m.Dns
	}
	return nil
}

func (m *NetworkConfig) GetEntProxy() *ProxyConfig {
	if m != nil {
		return m.EntProxy
	}
	return nil
}

type NetworkAdapter struct {
	Name      string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Id        string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	NetworkId string `protobuf:"bytes,3,opt,name=networkId" json:"networkId,omitempty"`
	Addr      string `protobuf:"bytes,4,opt,name=addr" json:"addr,omitempty"`
	Hostname  string `protobuf:"bytes,5,opt,name=hostname" json:"hostname,omitempty"`
	// more configuration for getting addr/EID
	CryptoEid     string `protobuf:"bytes,10,opt,name=cryptoEid" json:"cryptoEid,omitempty"`
	Lispsignature string `protobuf:"bytes,6,opt,name=lispsignature" json:"lispsignature,omitempty"`
	Pemcert       []byte `protobuf:"bytes,7,opt,name=pemcert,proto3" json:"pemcert,omitempty"`
	Pemprivatekey []byte `protobuf:"bytes,8,opt,name=pemprivatekey,proto3" json:"pemprivatekey,omitempty"`
	// used in case of P2V, where we want to specify a macAddress
	// to vif, that is simulated towards app
	MacAddress string `protobuf:"bytes,9,opt,name=macAddress" json:"macAddress,omitempty"`
	// firewall
	Acls []*ACE `protobuf:"bytes,40,rep,name=acls" json:"acls,omitempty"`
}

func (m *NetworkAdapter) Reset()                    { *m = NetworkAdapter{} }
func (m *NetworkAdapter) String() string            { return proto.CompactTextString(m) }
func (*NetworkAdapter) ProtoMessage()               {}
func (*NetworkAdapter) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{1} }

func (m *NetworkAdapter) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *NetworkAdapter) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *NetworkAdapter) GetNetworkId() string {
	if m != nil {
		return m.NetworkId
	}
	return ""
}

func (m *NetworkAdapter) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *NetworkAdapter) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *NetworkAdapter) GetCryptoEid() string {
	if m != nil {
		return m.CryptoEid
	}
	return ""
}

func (m *NetworkAdapter) GetLispsignature() string {
	if m != nil {
		return m.Lispsignature
	}
	return ""
}

func (m *NetworkAdapter) GetPemcert() []byte {
	if m != nil {
		return m.Pemcert
	}
	return nil
}

func (m *NetworkAdapter) GetPemprivatekey() []byte {
	if m != nil {
		return m.Pemprivatekey
	}
	return nil
}

func (m *NetworkAdapter) GetMacAddress() string {
	if m != nil {
		return m.MacAddress
	}
	return ""
}

func (m *NetworkAdapter) GetAcls() []*ACE {
	if m != nil {
		return m.Acls
	}
	return nil
}

func init() {
	proto.RegisterType((*NetworkConfig)(nil), "NetworkConfig")
	proto.RegisterType((*NetworkAdapter)(nil), "NetworkAdapter")
}

func init() { proto.RegisterFile("netconfig.proto", fileDescriptor9) }

var fileDescriptor9 = []byte{
	// 405 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0xcf, 0x6e, 0x1b, 0x21,
	0x10, 0xc6, 0xb5, 0xeb, 0x4d, 0x6c, 0x4f, 0x36, 0xae, 0x44, 0x0f, 0x45, 0x51, 0xff, 0xac, 0xa2,
	0x56, 0xda, 0x13, 0x96, 0xdc, 0x17, 0xa8, 0x9b, 0xfa, 0xd0, 0x4b, 0x54, 0x91, 0x9e, 0x72, 0x23,
	0x30, 0x71, 0x50, 0xbc, 0x80, 0x00, 0x27, 0xdd, 0xbc, 0x52, 0x5f, 0xa2, 0x8f, 0x56, 0x2d, 0x6c,
	0xdd, 0xe4, 0x36, 0xf3, 0xfb, 0x3e, 0x3e, 0x60, 0x34, 0xf0, 0xca, 0x60, 0x94, 0xd6, 0xdc, 0xea,
	0x2d, 0x73, 0xde, 0x46, 0x7b, 0x36, 0xbb, 0x7d, 0x1c, 0xab, 0x7a, 0x90, 0x3a, 0x93, 0xbb, 0xf3,
	0xdf, 0x05, 0x9c, 0x5e, 0x62, 0x7c, 0xb4, 0xfe, 0xfe, 0x22, 0xf9, 0xc9, 0x02, 0x4a, 0xad, 0x68,
	0xd1, 0x14, 0xed, 0x9c, 0x97, 0x5a, 0x91, 0x06, 0xaa, 0xd8, 0x3b, 0xa4, 0x47, 0x4d, 0xd1, 0x2e,
	0x56, 0x35, 0x1b, 0xdd, 0x3f, 0x7b, 0x87, 0x3c, 0x29, 0xe4, 0x0d, 0x94, 0xda, 0xd1, 0xe3, 0xa6,
	0x68, 0x4f, 0x56, 0x53, 0xa6, 0x5d, 0x70, 0x28, 0x79, 0xa9, 0x1d, 0xf9, 0x04, 0x13, 0x65, 0x02,
	0x9d, 0x36, 0x93, 0xf6, 0x64, 0xf5, 0x9a, 0x5d, 0x1b, 0x8c, 0x57, 0x51, 0x44, 0x2d, 0xbf, 0x5d,
	0x5e, 0x6d, 0x4c, 0xf4, 0x3d, 0x1f, 0x74, 0xd2, 0xc2, 0x0c, 0x4d, 0xfc, 0xe1, 0xed, 0xaf, 0x9e,
	0xce, 0x52, 0x4a, 0xcd, 0x52, 0x97, 0x5f, 0xc4, 0x0f, 0xea, 0xf9, 0x9f, 0x12, 0x16, 0xe3, 0xfd,
	0x6b, 0x25, 0x5c, 0x44, 0x4f, 0x08, 0x54, 0x46, 0x74, 0x38, 0x3e, 0x38, 0xd5, 0xe3, 0x17, 0xca,
	0xc3, 0x17, 0xde, 0xc2, 0xdc, 0xe4, 0x53, 0xdf, 0x15, 0x9d, 0x24, 0xfc, 0x1f, 0x0c, 0x09, 0x42,
	0x29, 0x4f, 0xab, 0x9c, 0x30, 0xd4, 0xe4, 0x0c, 0x66, 0x77, 0x36, 0xc4, 0x94, 0x7c, 0x94, 0xf8,
	0xa1, 0x1f, 0xd2, 0xa4, 0xef, 0x5d, 0xb4, 0x1b, 0xad, 0x28, 0xe4, 0xb4, 0x03, 0x20, 0x1f, 0xe1,
	0x74, 0xa7, 0x83, 0x0b, 0x7a, 0x6b, 0x44, 0xdc, 0x7b, 0x4c, 0x73, 0x99, 0xf3, 0x97, 0x90, 0x50,
	0x98, 0x3a, 0xec, 0x24, 0xfa, 0x48, 0xa7, 0x4d, 0xd1, 0xd6, 0xfc, 0x5f, 0x3b, 0x9c, 0x77, 0xd8,
	0x39, 0xaf, 0x1f, 0x44, 0xc4, 0x7b, 0xcc, 0x13, 0xa9, 0xf9, 0x4b, 0x48, 0xde, 0x03, 0x74, 0x42,
	0xae, 0x95, 0xf2, 0x18, 0x02, 0x9d, 0xa7, 0x2b, 0x9e, 0x11, 0x42, 0xa1, 0x12, 0x72, 0x17, 0x68,
	0x9b, 0x46, 0x5f, 0xb1, 0xf5, 0xc5, 0x86, 0x27, 0xf2, 0xf5, 0x0b, 0x7c, 0x90, 0xb6, 0x63, 0x4f,
	0xa8, 0x50, 0x09, 0x26, 0x77, 0x76, 0xaf, 0xd8, 0x3e, 0xa0, 0x7f, 0xd0, 0x12, 0xf3, 0x4e, 0x5c,
	0xbf, 0xdb, 0xea, 0x78, 0xb7, 0xbf, 0x61, 0xd2, 0x76, 0xcb, 0xec, 0x5b, 0x0a, 0xa7, 0x97, 0x4f,
	0x79, 0xa1, 0x6e, 0x8e, 0x93, 0xeb, 0xf3, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x0b, 0x82, 0x01,
	0x66, 0x64, 0x02, 0x00, 0x00,
}
