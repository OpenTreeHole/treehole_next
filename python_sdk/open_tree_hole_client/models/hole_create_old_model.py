from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.tag_create_model import TagCreateModel


T = TypeVar("T", bound="HoleCreateOldModel")


@_attrs_define
class HoleCreateOldModel:
    """
    Attributes:
        content (str):
        division_id (Union[Unset, int]):
        special_tag (Union[Unset, str]): Admin and Operator only
        tags (Union[Unset, list['TagCreateModel']]): All users
    """

    content: str
    division_id: Union[Unset, int] = UNSET
    special_tag: Union[Unset, str] = UNSET
    tags: Union[Unset, list["TagCreateModel"]] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        content = self.content

        division_id = self.division_id

        special_tag = self.special_tag

        tags: Union[Unset, list[dict[str, Any]]] = UNSET
        if not isinstance(self.tags, Unset):
            tags = []
            for tags_item_data in self.tags:
                tags_item = tags_item_data.to_dict()
                tags.append(tags_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "content": content,
            }
        )
        if division_id is not UNSET:
            field_dict["division_id"] = division_id
        if special_tag is not UNSET:
            field_dict["special_tag"] = special_tag
        if tags is not UNSET:
            field_dict["tags"] = tags

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.tag_create_model import TagCreateModel

        d = dict(src_dict)
        content = d.pop("content")

        division_id = d.pop("division_id", UNSET)

        special_tag = d.pop("special_tag", UNSET)

        tags = []
        _tags = d.pop("tags", UNSET)
        for tags_item_data in _tags or []:
            tags_item = TagCreateModel.from_dict(tags_item_data)

            tags.append(tags_item)

        hole_create_old_model = cls(
            content=content,
            division_id=division_id,
            special_tag=special_tag,
            tags=tags,
        )

        hole_create_old_model.additional_properties = d
        return hole_create_old_model

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
