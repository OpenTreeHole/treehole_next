from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="FloorModifySensitiveFloorRequest")


@_attrs_define
class FloorModifySensitiveFloorRequest:
    """
    Attributes:
        is_actual_sensitive (Union[Unset, bool]):
    """

    is_actual_sensitive: Union[Unset, bool] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        is_actual_sensitive = self.is_actual_sensitive

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if is_actual_sensitive is not UNSET:
            field_dict["is_actual_sensitive"] = is_actual_sensitive

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        is_actual_sensitive = d.pop("is_actual_sensitive", UNSET)

        floor_modify_sensitive_floor_request = cls(
            is_actual_sensitive=is_actual_sensitive,
        )

        floor_modify_sensitive_floor_request.additional_properties = d
        return floor_modify_sensitive_floor_request

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
