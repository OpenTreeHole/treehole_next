from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="ModelsTag")


@_attrs_define
class ModelsTag:
    """
    Attributes:
        id (Union[Unset, int]): / saved fields
        name (Union[Unset, str]): / base info
        nsfw (Union[Unset, bool]):
        tag_id (Union[Unset, int]): / generated field
        temperature (Union[Unset, int]):
    """

    id: Union[Unset, int] = UNSET
    name: Union[Unset, str] = UNSET
    nsfw: Union[Unset, bool] = UNSET
    tag_id: Union[Unset, int] = UNSET
    temperature: Union[Unset, int] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        id = self.id

        name = self.name

        nsfw = self.nsfw

        tag_id = self.tag_id

        temperature = self.temperature

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if id is not UNSET:
            field_dict["id"] = id
        if name is not UNSET:
            field_dict["name"] = name
        if nsfw is not UNSET:
            field_dict["nsfw"] = nsfw
        if tag_id is not UNSET:
            field_dict["tag_id"] = tag_id
        if temperature is not UNSET:
            field_dict["temperature"] = temperature

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        id = d.pop("id", UNSET)

        name = d.pop("name", UNSET)

        nsfw = d.pop("nsfw", UNSET)

        tag_id = d.pop("tag_id", UNSET)

        temperature = d.pop("temperature", UNSET)

        models_tag = cls(
            id=id,
            name=name,
            nsfw=nsfw,
            tag_id=tag_id,
            temperature=temperature,
        )

        models_tag.additional_properties = d
        return models_tag

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
