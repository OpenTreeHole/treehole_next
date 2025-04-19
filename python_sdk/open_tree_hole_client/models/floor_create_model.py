from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="FloorCreateModel")


@_attrs_define
class FloorCreateModel:
    """
    Attributes:
        content (str):
        reply_to (Union[Unset, int]): id of the floor to which replied
        special_tag (Union[Unset, str]): Admin and Operator only
    """

    content: str
    reply_to: Union[Unset, int] = UNSET
    special_tag: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        content = self.content

        reply_to = self.reply_to

        special_tag = self.special_tag

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "content": content,
            }
        )
        if reply_to is not UNSET:
            field_dict["reply_to"] = reply_to
        if special_tag is not UNSET:
            field_dict["special_tag"] = special_tag

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        content = d.pop("content")

        reply_to = d.pop("reply_to", UNSET)

        special_tag = d.pop("special_tag", UNSET)

        floor_create_model = cls(
            content=content,
            reply_to=reply_to,
            special_tag=special_tag,
        )

        floor_create_model.additional_properties = d
        return floor_create_model

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
