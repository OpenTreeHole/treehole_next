from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.tag_create_model import TagCreateModel


T = TypeVar("T", bound="HoleModifyModel")


@_attrs_define
class HoleModifyModel:
    """
    Attributes:
        division_id (Union[Unset, int]): Admin and owner only
        hidden (Union[Unset, bool]): Admin only
        lock (Union[Unset, bool]): admin only
        tags (Union[Unset, list['TagCreateModel']]): All users
        unhidden (Union[Unset, bool]): admin only
    """

    division_id: Union[Unset, int] = UNSET
    hidden: Union[Unset, bool] = UNSET
    lock: Union[Unset, bool] = UNSET
    tags: Union[Unset, list["TagCreateModel"]] = UNSET
    unhidden: Union[Unset, bool] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        division_id = self.division_id

        hidden = self.hidden

        lock = self.lock

        tags: Union[Unset, list[dict[str, Any]]] = UNSET
        if not isinstance(self.tags, Unset):
            tags = []
            for tags_item_data in self.tags:
                tags_item = tags_item_data.to_dict()
                tags.append(tags_item)

        unhidden = self.unhidden

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if division_id is not UNSET:
            field_dict["division_id"] = division_id
        if hidden is not UNSET:
            field_dict["hidden"] = hidden
        if lock is not UNSET:
            field_dict["lock"] = lock
        if tags is not UNSET:
            field_dict["tags"] = tags
        if unhidden is not UNSET:
            field_dict["unhidden"] = unhidden

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.tag_create_model import TagCreateModel

        d = dict(src_dict)
        division_id = d.pop("division_id", UNSET)

        hidden = d.pop("hidden", UNSET)

        lock = d.pop("lock", UNSET)

        tags = []
        _tags = d.pop("tags", UNSET)
        for tags_item_data in _tags or []:
            tags_item = TagCreateModel.from_dict(tags_item_data)

            tags.append(tags_item)

        unhidden = d.pop("unhidden", UNSET)

        hole_modify_model = cls(
            division_id=division_id,
            hidden=hidden,
            lock=lock,
            tags=tags,
            unhidden=unhidden,
        )

        hole_modify_model.additional_properties = d
        return hole_modify_model

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
