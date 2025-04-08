from collections.abc import Mapping
from typing import Any, TypeVar, Union, cast

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="FavouriteMoveModel")


@_attrs_define
class FavouriteMoveModel:
    """
    Attributes:
        from_favorite_group_id (int):
        to_favorite_group_id (int):
        hole_ids (Union[Unset, list[int]]):
    """

    from_favorite_group_id: int
    to_favorite_group_id: int
    hole_ids: Union[Unset, list[int]] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        from_favorite_group_id = self.from_favorite_group_id

        to_favorite_group_id = self.to_favorite_group_id

        hole_ids: Union[Unset, list[int]] = UNSET
        if not isinstance(self.hole_ids, Unset):
            hole_ids = self.hole_ids

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "from_favorite_group_id": from_favorite_group_id,
                "to_favorite_group_id": to_favorite_group_id,
            }
        )
        if hole_ids is not UNSET:
            field_dict["hole_ids"] = hole_ids

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        from_favorite_group_id = d.pop("from_favorite_group_id")

        to_favorite_group_id = d.pop("to_favorite_group_id")

        hole_ids = cast(list[int], d.pop("hole_ids", UNSET))

        favourite_move_model = cls(
            from_favorite_group_id=from_favorite_group_id,
            to_favorite_group_id=to_favorite_group_id,
            hole_ids=hole_ids,
        )

        favourite_move_model.additional_properties = d
        return favourite_move_model

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
