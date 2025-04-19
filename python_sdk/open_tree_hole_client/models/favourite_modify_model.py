from collections.abc import Mapping
from typing import Any, TypeVar, Union, cast

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="FavouriteModifyModel")


@_attrs_define
class FavouriteModifyModel:
    """
    Attributes:
        favorite_group_id (Union[Unset, int]):
        hole_ids (Union[Unset, list[int]]):
    """

    favorite_group_id: Union[Unset, int] = UNSET
    hole_ids: Union[Unset, list[int]] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        favorite_group_id = self.favorite_group_id

        hole_ids: Union[Unset, list[int]] = UNSET
        if not isinstance(self.hole_ids, Unset):
            hole_ids = self.hole_ids

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if favorite_group_id is not UNSET:
            field_dict["favorite_group_id"] = favorite_group_id
        if hole_ids is not UNSET:
            field_dict["hole_ids"] = hole_ids

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        favorite_group_id = d.pop("favorite_group_id", UNSET)

        hole_ids = cast(list[int], d.pop("hole_ids", UNSET))

        favourite_modify_model = cls(
            favorite_group_id=favorite_group_id,
            hole_ids=hole_ids,
        )

        favourite_modify_model.additional_properties = d
        return favourite_modify_model

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
