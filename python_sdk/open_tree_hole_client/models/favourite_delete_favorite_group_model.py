from collections.abc import Mapping
from typing import Any, TypeVar

from attrs import define as _attrs_define
from attrs import field as _attrs_field

T = TypeVar("T", bound="FavouriteDeleteFavoriteGroupModel")


@_attrs_define
class FavouriteDeleteFavoriteGroupModel:
    """
    Attributes:
        favorite_group_id (int):
    """

    favorite_group_id: int
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        favorite_group_id = self.favorite_group_id

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "favorite_group_id": favorite_group_id,
            }
        )

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        favorite_group_id = d.pop("favorite_group_id")

        favourite_delete_favorite_group_model = cls(
            favorite_group_id=favorite_group_id,
        )

        favourite_delete_favorite_group_model.additional_properties = d
        return favourite_delete_favorite_group_model

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
