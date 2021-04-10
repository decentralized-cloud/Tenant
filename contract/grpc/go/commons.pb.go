// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.6
// source: commons.proto

package __

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//*
// The different error types
type Error int32

const (
	// Indicates the operation was successful
	Error_NO_ERROR Error = 0
	// Indicates the operation fail with unknown error
	Error_UNKNOWN Error = 1
	// Indicates the project already exists
	Error_PROJECT_ALREADY_EXISTS Error = 2
	// Indicates the project does not exist
	Error_PROJECT_NOT_FOUND Error = 3
	// Indicates the provided values for he operation were invalid
	Error_BAD_REQUEST Error = 4
)

// Enum value maps for Error.
var (
	Error_name = map[int32]string{
		0: "NO_ERROR",
		1: "UNKNOWN",
		2: "PROJECT_ALREADY_EXISTS",
		3: "PROJECT_NOT_FOUND",
		4: "BAD_REQUEST",
	}
	Error_value = map[string]int32{
		"NO_ERROR":               0,
		"UNKNOWN":                1,
		"PROJECT_ALREADY_EXISTS": 2,
		"PROJECT_NOT_FOUND":      3,
		"BAD_REQUEST":            4,
	}
)

func (x Error) Enum() *Error {
	p := new(Error)
	*p = x
	return p
}

func (x Error) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Error) Descriptor() protoreflect.EnumDescriptor {
	return file_commons_proto_enumTypes[0].Descriptor()
}

func (Error) Type() protoreflect.EnumType {
	return &file_commons_proto_enumTypes[0]
}

func (x Error) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Error.Descriptor instead.
func (Error) EnumDescriptor() ([]byte, []int) {
	return file_commons_proto_rawDescGZIP(), []int{0}
}

var File_commons_proto protoreflect.FileDescriptor

var file_commons_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2a, 0x66, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x4f, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x00, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16,
	0x50, 0x52, 0x4f, 0x4a, 0x45, 0x43, 0x54, 0x5f, 0x41, 0x4c, 0x52, 0x45, 0x41, 0x44, 0x59, 0x5f,
	0x45, 0x58, 0x49, 0x53, 0x54, 0x53, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11, 0x50, 0x52, 0x4f, 0x4a,
	0x45, 0x43, 0x54, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x03, 0x12,
	0x0f, 0x0a, 0x0b, 0x42, 0x41, 0x44, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x10, 0x04,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_commons_proto_rawDescOnce sync.Once
	file_commons_proto_rawDescData = file_commons_proto_rawDesc
)

func file_commons_proto_rawDescGZIP() []byte {
	file_commons_proto_rawDescOnce.Do(func() {
		file_commons_proto_rawDescData = protoimpl.X.CompressGZIP(file_commons_proto_rawDescData)
	})
	return file_commons_proto_rawDescData
}

var file_commons_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_commons_proto_goTypes = []interface{}{
	(Error)(0), // 0: project.Error
}
var file_commons_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_commons_proto_init() }
func file_commons_proto_init() {
	if File_commons_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_commons_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_commons_proto_goTypes,
		DependencyIndexes: file_commons_proto_depIdxs,
		EnumInfos:         file_commons_proto_enumTypes,
	}.Build()
	File_commons_proto = out.File
	file_commons_proto_rawDesc = nil
	file_commons_proto_goTypes = nil
	file_commons_proto_depIdxs = nil
}
