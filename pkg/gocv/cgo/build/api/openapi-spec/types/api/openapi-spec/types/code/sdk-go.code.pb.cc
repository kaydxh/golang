// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: api/openapi-spec/types/code/sdk-go.code.proto

#include "api/openapi-spec/types/code/sdk-go.code.pb.h"

#include <algorithm>

#include <google/protobuf/stubs/common.h>
#include <google/protobuf/stubs/port.h>
#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/wire_format_lite_inl.h>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/generated_message_reflection.h>
#include <google/protobuf/reflection_ops.h>
#include <google/protobuf/wire_format.h>
// This is a temporary google only hack
#ifdef GOOGLE_PROTOBUF_ENFORCE_UNIQUENESS
#include "third_party/protobuf/version.h"
#endif
// @@protoc_insertion_point(includes)

namespace sdk {
namespace types {
namespace code {
class CgoErrorDefaultTypeInternal {
 public:
  ::google::protobuf::internal::ExplicitlyConstructed<CgoError>
      _instance;
} _CgoError_default_instance_;
}  // namespace code
}  // namespace types
}  // namespace sdk
namespace protobuf_api_2fopenapi_2dspec_2ftypes_2fcode_2fsdk_2dgo_2ecode_2eproto {
static void InitDefaultsCgoError() {
  GOOGLE_PROTOBUF_VERIFY_VERSION;

  {
    void* ptr = &::sdk::types::code::_CgoError_default_instance_;
    new (ptr) ::sdk::types::code::CgoError();
    ::google::protobuf::internal::OnShutdownDestroyMessage(ptr);
  }
  ::sdk::types::code::CgoError::InitAsDefaultInstance();
}

::google::protobuf::internal::SCCInfo<0> scc_info_CgoError =
    {{ATOMIC_VAR_INIT(::google::protobuf::internal::SCCInfoBase::kUninitialized), 0, InitDefaultsCgoError}, {}};

void InitDefaults() {
  ::google::protobuf::internal::InitSCC(&scc_info_CgoError.base);
}

::google::protobuf::Metadata file_level_metadata[1];
const ::google::protobuf::EnumDescriptor* file_level_enum_descriptors[1];

const ::google::protobuf::uint32 TableStruct::offsets[] GOOGLE_PROTOBUF_ATTRIBUTE_SECTION_VARIABLE(protodesc_cold) = {
  ~0u,  // no _has_bits_
  GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(::sdk::types::code::CgoError, _internal_metadata_),
  ~0u,  // no _extensions_
  ~0u,  // no _oneof_case_
  ~0u,  // no _weak_field_map_
  GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(::sdk::types::code::CgoError, error_code_),
  GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(::sdk::types::code::CgoError, error_message_),
  GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(::sdk::types::code::CgoError, sdk_error_code_),
  GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(::sdk::types::code::CgoError, sdk_error_message_),
};
static const ::google::protobuf::internal::MigrationSchema schemas[] GOOGLE_PROTOBUF_ATTRIBUTE_SECTION_VARIABLE(protodesc_cold) = {
  { 0, -1, sizeof(::sdk::types::code::CgoError)},
};

static ::google::protobuf::Message const * const file_default_instances[] = {
  reinterpret_cast<const ::google::protobuf::Message*>(&::sdk::types::code::_CgoError_default_instance_),
};

void protobuf_AssignDescriptors() {
  AddDescriptors();
  AssignDescriptors(
      "api/openapi-spec/types/code/sdk-go.code.proto", schemas, file_default_instances, TableStruct::offsets,
      file_level_metadata, file_level_enum_descriptors, NULL);
}

void protobuf_AssignDescriptorsOnce() {
  static ::google::protobuf::internal::once_flag once;
  ::google::protobuf::internal::call_once(once, protobuf_AssignDescriptors);
}

void protobuf_RegisterTypes(const ::std::string&) GOOGLE_PROTOBUF_ATTRIBUTE_COLD;
void protobuf_RegisterTypes(const ::std::string&) {
  protobuf_AssignDescriptorsOnce();
  ::google::protobuf::internal::RegisterAllTypes(file_level_metadata, 1);
}

void AddDescriptorsImpl() {
  InitDefaults();
  static const char descriptor[] GOOGLE_PROTOBUF_ATTRIBUTE_SECTION_VARIABLE(protodesc_cold) = {
      "\n-api/openapi-spec/types/code/sdk-go.cod"
      "e.proto\022\016sdk.types.code\"~\n\010CgoError\022(\n\ne"
      "rror_code\030\001 \001(\0162\024.sdk.types.code.Code\022\025\n"
      "\rerror_message\030\002 \001(\t\022\026\n\016sdk_error_code\030\003"
      " \001(\005\022\031\n\021sdk_error_message\030\004 \001(\t*\254\002\n\004Code"
      "\022\006\n\002OK\020\000\022\014\n\010Canceled\020\001\022\013\n\007Unknown\020\002\022\023\n\017I"
      "nvalidArgument\020\003\022\024\n\020DeadlineExceeded\020\004\022\014"
      "\n\010NotFound\020\005\022\021\n\rAlreadyExists\020\006\022\024\n\020Permi"
      "ssionDenied\020\007\022\025\n\021ResourceExhausted\020\010\022\026\n\022"
      "FailedPrecondition\020\t\022\013\n\007Aborted\020\n\022\016\n\nOut"
      "OfRange\020\013\022\021\n\rUnimplemented\020\014\022\014\n\010Internal"
      "\020\r\022\017\n\013Unavailable\020\016\022\014\n\010DataLoss\020\017\022\023\n\017Una"
      "uthenticated\020\020BCZAgithub.com/kaydxh/gola"
      "ng/pkg/cgo/api/openapi-spec/types/code;c"
      "odeb\006proto3"
  };
  ::google::protobuf::DescriptorPool::InternalAddGeneratedFile(
      descriptor, 571);
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedFile(
    "api/openapi-spec/types/code/sdk-go.code.proto", &protobuf_RegisterTypes);
}

void AddDescriptors() {
  static ::google::protobuf::internal::once_flag once;
  ::google::protobuf::internal::call_once(once, AddDescriptorsImpl);
}
// Force AddDescriptors() to be called at dynamic initialization time.
struct StaticDescriptorInitializer {
  StaticDescriptorInitializer() {
    AddDescriptors();
  }
} static_descriptor_initializer;
}  // namespace protobuf_api_2fopenapi_2dspec_2ftypes_2fcode_2fsdk_2dgo_2ecode_2eproto
namespace sdk {
namespace types {
namespace code {
const ::google::protobuf::EnumDescriptor* Code_descriptor() {
  protobuf_api_2fopenapi_2dspec_2ftypes_2fcode_2fsdk_2dgo_2ecode_2eproto::protobuf_AssignDescriptorsOnce();
  return protobuf_api_2fopenapi_2dspec_2ftypes_2fcode_2fsdk_2dgo_2ecode_2eproto::file_level_enum_descriptors[0];
}
bool Code_IsValid(int value) {
  switch (value) {
    case 0:
    case 1:
    case 2:
    case 3:
    case 4:
    case 5:
    case 6:
    case 7:
    case 8:
    case 9:
    case 10:
    case 11:
    case 12:
    case 13:
    case 14:
    case 15:
    case 16:
      return true;
    default:
      return false;
  }
}


// ===================================================================

void CgoError::InitAsDefaultInstance() {
}
#if !defined(_MSC_VER) || _MSC_VER >= 1900
const int CgoError::kErrorCodeFieldNumber;
const int CgoError::kErrorMessageFieldNumber;
const int CgoError::kSdkErrorCodeFieldNumber;
const int CgoError::kSdkErrorMessageFieldNumber;
#endif  // !defined(_MSC_VER) || _MSC_VER >= 1900

CgoError::CgoError()
  : ::google::protobuf::Message(), _internal_metadata_(NULL) {
  ::google::protobuf::internal::InitSCC(
      &protobuf_api_2fopenapi_2dspec_2ftypes_2fcode_2fsdk_2dgo_2ecode_2eproto::scc_info_CgoError.base);
  SharedCtor();
  // @@protoc_insertion_point(constructor:sdk.types.code.CgoError)
}
CgoError::CgoError(const CgoError& from)
  : ::google::protobuf::Message(),
      _internal_metadata_(NULL) {
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  error_message_.UnsafeSetDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  if (from.error_message().size() > 0) {
    error_message_.AssignWithDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), from.error_message_);
  }
  sdk_error_message_.UnsafeSetDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  if (from.sdk_error_message().size() > 0) {
    sdk_error_message_.AssignWithDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), from.sdk_error_message_);
  }
  ::memcpy(&error_code_, &from.error_code_,
    static_cast<size_t>(reinterpret_cast<char*>(&sdk_error_code_) -
    reinterpret_cast<char*>(&error_code_)) + sizeof(sdk_error_code_));
  // @@protoc_insertion_point(copy_constructor:sdk.types.code.CgoError)
}

void CgoError::SharedCtor() {
  error_message_.UnsafeSetDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  sdk_error_message_.UnsafeSetDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  ::memset(&error_code_, 0, static_cast<size_t>(
      reinterpret_cast<char*>(&sdk_error_code_) -
      reinterpret_cast<char*>(&error_code_)) + sizeof(sdk_error_code_));
}

CgoError::~CgoError() {
  // @@protoc_insertion_point(destructor:sdk.types.code.CgoError)
  SharedDtor();
}

void CgoError::SharedDtor() {
  error_message_.DestroyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  sdk_error_message_.DestroyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}

void CgoError::SetCachedSize(int size) const {
  _cached_size_.Set(size);
}
const ::google::protobuf::Descriptor* CgoError::descriptor() {
  ::protobuf_api_2fopenapi_2dspec_2ftypes_2fcode_2fsdk_2dgo_2ecode_2eproto::protobuf_AssignDescriptorsOnce();
  return ::protobuf_api_2fopenapi_2dspec_2ftypes_2fcode_2fsdk_2dgo_2ecode_2eproto::file_level_metadata[kIndexInFileMessages].descriptor;
}

const CgoError& CgoError::default_instance() {
  ::google::protobuf::internal::InitSCC(&protobuf_api_2fopenapi_2dspec_2ftypes_2fcode_2fsdk_2dgo_2ecode_2eproto::scc_info_CgoError.base);
  return *internal_default_instance();
}


void CgoError::Clear() {
// @@protoc_insertion_point(message_clear_start:sdk.types.code.CgoError)
  ::google::protobuf::uint32 cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  error_message_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  sdk_error_message_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  ::memset(&error_code_, 0, static_cast<size_t>(
      reinterpret_cast<char*>(&sdk_error_code_) -
      reinterpret_cast<char*>(&error_code_)) + sizeof(sdk_error_code_));
  _internal_metadata_.Clear();
}

bool CgoError::MergePartialFromCodedStream(
    ::google::protobuf::io::CodedInputStream* input) {
#define DO_(EXPRESSION) if (!GOOGLE_PREDICT_TRUE(EXPRESSION)) goto failure
  ::google::protobuf::uint32 tag;
  // @@protoc_insertion_point(parse_start:sdk.types.code.CgoError)
  for (;;) {
    ::std::pair<::google::protobuf::uint32, bool> p = input->ReadTagWithCutoffNoLastTag(127u);
    tag = p.first;
    if (!p.second) goto handle_unusual;
    switch (::google::protobuf::internal::WireFormatLite::GetTagFieldNumber(tag)) {
      // .sdk.types.code.Code error_code = 1;
      case 1: {
        if (static_cast< ::google::protobuf::uint8>(tag) ==
            static_cast< ::google::protobuf::uint8>(8u /* 8 & 0xFF */)) {
          int value;
          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   int, ::google::protobuf::internal::WireFormatLite::TYPE_ENUM>(
                 input, &value)));
          set_error_code(static_cast< ::sdk::types::code::Code >(value));
        } else {
          goto handle_unusual;
        }
        break;
      }

      // string error_message = 2;
      case 2: {
        if (static_cast< ::google::protobuf::uint8>(tag) ==
            static_cast< ::google::protobuf::uint8>(18u /* 18 & 0xFF */)) {
          DO_(::google::protobuf::internal::WireFormatLite::ReadString(
                input, this->mutable_error_message()));
          DO_(::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
            this->error_message().data(), static_cast<int>(this->error_message().length()),
            ::google::protobuf::internal::WireFormatLite::PARSE,
            "sdk.types.code.CgoError.error_message"));
        } else {
          goto handle_unusual;
        }
        break;
      }

      // int32 sdk_error_code = 3;
      case 3: {
        if (static_cast< ::google::protobuf::uint8>(tag) ==
            static_cast< ::google::protobuf::uint8>(24u /* 24 & 0xFF */)) {

          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::int32, ::google::protobuf::internal::WireFormatLite::TYPE_INT32>(
                 input, &sdk_error_code_)));
        } else {
          goto handle_unusual;
        }
        break;
      }

      // string sdk_error_message = 4;
      case 4: {
        if (static_cast< ::google::protobuf::uint8>(tag) ==
            static_cast< ::google::protobuf::uint8>(34u /* 34 & 0xFF */)) {
          DO_(::google::protobuf::internal::WireFormatLite::ReadString(
                input, this->mutable_sdk_error_message()));
          DO_(::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
            this->sdk_error_message().data(), static_cast<int>(this->sdk_error_message().length()),
            ::google::protobuf::internal::WireFormatLite::PARSE,
            "sdk.types.code.CgoError.sdk_error_message"));
        } else {
          goto handle_unusual;
        }
        break;
      }

      default: {
      handle_unusual:
        if (tag == 0) {
          goto success;
        }
        DO_(::google::protobuf::internal::WireFormat::SkipField(
              input, tag, _internal_metadata_.mutable_unknown_fields()));
        break;
      }
    }
  }
success:
  // @@protoc_insertion_point(parse_success:sdk.types.code.CgoError)
  return true;
failure:
  // @@protoc_insertion_point(parse_failure:sdk.types.code.CgoError)
  return false;
#undef DO_
}

void CgoError::SerializeWithCachedSizes(
    ::google::protobuf::io::CodedOutputStream* output) const {
  // @@protoc_insertion_point(serialize_start:sdk.types.code.CgoError)
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  // .sdk.types.code.Code error_code = 1;
  if (this->error_code() != 0) {
    ::google::protobuf::internal::WireFormatLite::WriteEnum(
      1, this->error_code(), output);
  }

  // string error_message = 2;
  if (this->error_message().size() > 0) {
    ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
      this->error_message().data(), static_cast<int>(this->error_message().length()),
      ::google::protobuf::internal::WireFormatLite::SERIALIZE,
      "sdk.types.code.CgoError.error_message");
    ::google::protobuf::internal::WireFormatLite::WriteStringMaybeAliased(
      2, this->error_message(), output);
  }

  // int32 sdk_error_code = 3;
  if (this->sdk_error_code() != 0) {
    ::google::protobuf::internal::WireFormatLite::WriteInt32(3, this->sdk_error_code(), output);
  }

  // string sdk_error_message = 4;
  if (this->sdk_error_message().size() > 0) {
    ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
      this->sdk_error_message().data(), static_cast<int>(this->sdk_error_message().length()),
      ::google::protobuf::internal::WireFormatLite::SERIALIZE,
      "sdk.types.code.CgoError.sdk_error_message");
    ::google::protobuf::internal::WireFormatLite::WriteStringMaybeAliased(
      4, this->sdk_error_message(), output);
  }

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    ::google::protobuf::internal::WireFormat::SerializeUnknownFields(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()), output);
  }
  // @@protoc_insertion_point(serialize_end:sdk.types.code.CgoError)
}

::google::protobuf::uint8* CgoError::InternalSerializeWithCachedSizesToArray(
    bool deterministic, ::google::protobuf::uint8* target) const {
  (void)deterministic; // Unused
  // @@protoc_insertion_point(serialize_to_array_start:sdk.types.code.CgoError)
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  // .sdk.types.code.Code error_code = 1;
  if (this->error_code() != 0) {
    target = ::google::protobuf::internal::WireFormatLite::WriteEnumToArray(
      1, this->error_code(), target);
  }

  // string error_message = 2;
  if (this->error_message().size() > 0) {
    ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
      this->error_message().data(), static_cast<int>(this->error_message().length()),
      ::google::protobuf::internal::WireFormatLite::SERIALIZE,
      "sdk.types.code.CgoError.error_message");
    target =
      ::google::protobuf::internal::WireFormatLite::WriteStringToArray(
        2, this->error_message(), target);
  }

  // int32 sdk_error_code = 3;
  if (this->sdk_error_code() != 0) {
    target = ::google::protobuf::internal::WireFormatLite::WriteInt32ToArray(3, this->sdk_error_code(), target);
  }

  // string sdk_error_message = 4;
  if (this->sdk_error_message().size() > 0) {
    ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
      this->sdk_error_message().data(), static_cast<int>(this->sdk_error_message().length()),
      ::google::protobuf::internal::WireFormatLite::SERIALIZE,
      "sdk.types.code.CgoError.sdk_error_message");
    target =
      ::google::protobuf::internal::WireFormatLite::WriteStringToArray(
        4, this->sdk_error_message(), target);
  }

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    target = ::google::protobuf::internal::WireFormat::SerializeUnknownFieldsToArray(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()), target);
  }
  // @@protoc_insertion_point(serialize_to_array_end:sdk.types.code.CgoError)
  return target;
}

size_t CgoError::ByteSizeLong() const {
// @@protoc_insertion_point(message_byte_size_start:sdk.types.code.CgoError)
  size_t total_size = 0;

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    total_size +=
      ::google::protobuf::internal::WireFormat::ComputeUnknownFieldsSize(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()));
  }
  // string error_message = 2;
  if (this->error_message().size() > 0) {
    total_size += 1 +
      ::google::protobuf::internal::WireFormatLite::StringSize(
        this->error_message());
  }

  // string sdk_error_message = 4;
  if (this->sdk_error_message().size() > 0) {
    total_size += 1 +
      ::google::protobuf::internal::WireFormatLite::StringSize(
        this->sdk_error_message());
  }

  // .sdk.types.code.Code error_code = 1;
  if (this->error_code() != 0) {
    total_size += 1 +
      ::google::protobuf::internal::WireFormatLite::EnumSize(this->error_code());
  }

  // int32 sdk_error_code = 3;
  if (this->sdk_error_code() != 0) {
    total_size += 1 +
      ::google::protobuf::internal::WireFormatLite::Int32Size(
        this->sdk_error_code());
  }

  int cached_size = ::google::protobuf::internal::ToCachedSize(total_size);
  SetCachedSize(cached_size);
  return total_size;
}

void CgoError::MergeFrom(const ::google::protobuf::Message& from) {
// @@protoc_insertion_point(generalized_merge_from_start:sdk.types.code.CgoError)
  GOOGLE_DCHECK_NE(&from, this);
  const CgoError* source =
      ::google::protobuf::internal::DynamicCastToGenerated<const CgoError>(
          &from);
  if (source == NULL) {
  // @@protoc_insertion_point(generalized_merge_from_cast_fail:sdk.types.code.CgoError)
    ::google::protobuf::internal::ReflectionOps::Merge(from, this);
  } else {
  // @@protoc_insertion_point(generalized_merge_from_cast_success:sdk.types.code.CgoError)
    MergeFrom(*source);
  }
}

void CgoError::MergeFrom(const CgoError& from) {
// @@protoc_insertion_point(class_specific_merge_from_start:sdk.types.code.CgoError)
  GOOGLE_DCHECK_NE(&from, this);
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  if (from.error_message().size() > 0) {

    error_message_.AssignWithDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), from.error_message_);
  }
  if (from.sdk_error_message().size() > 0) {

    sdk_error_message_.AssignWithDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), from.sdk_error_message_);
  }
  if (from.error_code() != 0) {
    set_error_code(from.error_code());
  }
  if (from.sdk_error_code() != 0) {
    set_sdk_error_code(from.sdk_error_code());
  }
}

void CgoError::CopyFrom(const ::google::protobuf::Message& from) {
// @@protoc_insertion_point(generalized_copy_from_start:sdk.types.code.CgoError)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

void CgoError::CopyFrom(const CgoError& from) {
// @@protoc_insertion_point(class_specific_copy_from_start:sdk.types.code.CgoError)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool CgoError::IsInitialized() const {
  return true;
}

void CgoError::Swap(CgoError* other) {
  if (other == this) return;
  InternalSwap(other);
}
void CgoError::InternalSwap(CgoError* other) {
  using std::swap;
  error_message_.Swap(&other->error_message_, &::google::protobuf::internal::GetEmptyStringAlreadyInited(),
    GetArenaNoVirtual());
  sdk_error_message_.Swap(&other->sdk_error_message_, &::google::protobuf::internal::GetEmptyStringAlreadyInited(),
    GetArenaNoVirtual());
  swap(error_code_, other->error_code_);
  swap(sdk_error_code_, other->sdk_error_code_);
  _internal_metadata_.Swap(&other->_internal_metadata_);
}

::google::protobuf::Metadata CgoError::GetMetadata() const {
  protobuf_api_2fopenapi_2dspec_2ftypes_2fcode_2fsdk_2dgo_2ecode_2eproto::protobuf_AssignDescriptorsOnce();
  return ::protobuf_api_2fopenapi_2dspec_2ftypes_2fcode_2fsdk_2dgo_2ecode_2eproto::file_level_metadata[kIndexInFileMessages];
}


// @@protoc_insertion_point(namespace_scope)
}  // namespace code
}  // namespace types
}  // namespace sdk
namespace google {
namespace protobuf {
template<> GOOGLE_PROTOBUF_ATTRIBUTE_NOINLINE ::sdk::types::code::CgoError* Arena::CreateMaybeMessage< ::sdk::types::code::CgoError >(Arena* arena) {
  return Arena::CreateInternal< ::sdk::types::code::CgoError >(arena);
}
}  // namespace protobuf
}  // namespace google

// @@protoc_insertion_point(global_scope)