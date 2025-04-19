from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="ModelsFavoriteGroup")


@_attrs_define
class ModelsFavoriteGroup:
    """
    Attributes:
        count (Union[Unset, int]):
        deleted (Union[Unset, bool]):
        favorite_group_id (Union[Unset, int]):
        name (Union[Unset, str]):  Default: '默认'.
        time_created (Union[Unset, str]):
        time_updated (Union[Unset, str]):
        user_id (Union[Unset, int]):
    """

    count: Union[Unset, int] = UNSET
    deleted: Union[Unset, bool] = UNSET
    favorite_group_id: Union[Unset, int] = UNSET
    name: Union[Unset, str] = "默认"
    time_created: Union[Unset, str] = UNSET
    time_updated: Union[Unset, str] = UNSET
    user_id: Union[Unset, int] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        count = self.count

        deleted = self.deleted

        favorite_group_id = self.favorite_group_id

        name = self.name

        time_created = self.time_created

        time_updated = self.time_updated

        user_id = self.user_id

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if count is not UNSET:
            field_dict["count"] = count
        if deleted is not UNSET:
            field_dict["deleted"] = deleted
        if favorite_group_id is not UNSET:
            field_dict["favorite_group_id"] = favorite_group_id
        if name is not UNSET:
            field_dict["name"] = name
        if time_created is not UNSET:
            field_dict["time_created"] = time_created
        if time_updated is not UNSET:
            field_dict["time_updated"] = time_updated
        if user_id is not UNSET:
            field_dict["user_id"] = user_id

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        count = d.pop("count", UNSET)

        deleted = d.pop("deleted", UNSET)

        favorite_group_id = d.pop("favorite_group_id", UNSET)

        name = d.pop("name", UNSET)

        time_created = d.pop("time_created", UNSET)

        time_updated = d.pop("time_updated", UNSET)

        user_id = d.pop("user_id", UNSET)

        models_favorite_group = cls(
            count=count,
            deleted=deleted,
            favorite_group_id=favorite_group_id,
            name=name,
            time_created=time_created,
            time_updated=time_updated,
            user_id=user_id,
        )

        models_favorite_group.additional_properties = d
        return models_favorite_group

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
