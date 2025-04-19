from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.common_error_detail_element_value import CommonErrorDetailElementValue


T = TypeVar("T", bound="CommonErrorDetailElement")


@_attrs_define
class CommonErrorDetailElement:
    """
    Attributes:
        field (Union[Unset, str]):
        message (Union[Unset, str]):
        param (Union[Unset, str]):
        struct_field (Union[Unset, str]):
        tag (Union[Unset, str]):
        value (Union[Unset, CommonErrorDetailElementValue]):
    """

    field: Union[Unset, str] = UNSET
    message: Union[Unset, str] = UNSET
    param: Union[Unset, str] = UNSET
    struct_field: Union[Unset, str] = UNSET
    tag: Union[Unset, str] = UNSET
    value: Union[Unset, "CommonErrorDetailElementValue"] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        field = self.field

        message = self.message

        param = self.param

        struct_field = self.struct_field

        tag = self.tag

        value: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.value, Unset):
            value = self.value.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if field is not UNSET:
            field_dict["field"] = field
        if message is not UNSET:
            field_dict["message"] = message
        if param is not UNSET:
            field_dict["param"] = param
        if struct_field is not UNSET:
            field_dict["struct_field"] = struct_field
        if tag is not UNSET:
            field_dict["tag"] = tag
        if value is not UNSET:
            field_dict["value"] = value

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.common_error_detail_element_value import CommonErrorDetailElementValue

        d = dict(src_dict)
        field = d.pop("field", UNSET)

        message = d.pop("message", UNSET)

        param = d.pop("param", UNSET)

        struct_field = d.pop("struct_field", UNSET)

        tag = d.pop("tag", UNSET)

        _value = d.pop("value", UNSET)
        value: Union[Unset, CommonErrorDetailElementValue]
        if isinstance(_value, Unset):
            value = UNSET
        else:
            value = CommonErrorDetailElementValue.from_dict(_value)

        common_error_detail_element = cls(
            field=field,
            message=message,
            param=param,
            struct_field=struct_field,
            tag=tag,
            value=value,
        )

        common_error_detail_element.additional_properties = d
        return common_error_detail_element

    @property
    def additional_keys(self) -> list[str]:
        return list(self.additional_properties.keys())

    def __getitem__(self, key: str) -> Any:
        return self.additional_properties[key]

    def __setitem__(self, key: str, value: Any) -> None:
        self.additional_properties[key] = value

    def __delitem__(self, key: str) -> None:
        del self.additional_properties[key]

    def __contains__(self, key: str) -> bool:
        return key in self.additional_properties
