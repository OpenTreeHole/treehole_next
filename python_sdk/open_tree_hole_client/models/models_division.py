from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.models_hole import ModelsHole


T = TypeVar("T", bound="ModelsDivision")


@_attrs_define
class ModelsDivision:
    """
    Attributes:
        description (Union[Unset, str]):
        division_id (Union[Unset, int]): / generated field
        hidden (Union[Unset, bool]):
        id (Union[Unset, int]): / saved fields
        name (Union[Unset, str]): / base info
        pinned (Union[Unset, list['ModelsHole']]): return pinned hole to frontend
        time_created (Union[Unset, str]):
        time_updated (Union[Unset, str]):
    """

    description: Union[Unset, str] = UNSET
    division_id: Union[Unset, int] = UNSET
    hidden: Union[Unset, bool] = UNSET
    id: Union[Unset, int] = UNSET
    name: Union[Unset, str] = UNSET
    pinned: Union[Unset, list["ModelsHole"]] = UNSET
    time_created: Union[Unset, str] = UNSET
    time_updated: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        description = self.description

        division_id = self.division_id

        hidden = self.hidden

        id = self.id

        name = self.name

        pinned: Union[Unset, list[dict[str, Any]]] = UNSET
        if not isinstance(self.pinned, Unset):
            pinned = []
            for pinned_item_data in self.pinned:
                pinned_item = pinned_item_data.to_dict()
                pinned.append(pinned_item)

        time_created = self.time_created

        time_updated = self.time_updated

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if description is not UNSET:
            field_dict["description"] = description
        if division_id is not UNSET:
            field_dict["division_id"] = division_id
        if hidden is not UNSET:
            field_dict["hidden"] = hidden
        if id is not UNSET:
            field_dict["id"] = id
        if name is not UNSET:
            field_dict["name"] = name
        if pinned is not UNSET:
            field_dict["pinned"] = pinned
        if time_created is not UNSET:
            field_dict["time_created"] = time_created
        if time_updated is not UNSET:
            field_dict["time_updated"] = time_updated

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.models_hole import ModelsHole

        d = dict(src_dict)
        description = d.pop("description", UNSET)

        division_id = d.pop("division_id", UNSET)

        hidden = d.pop("hidden", UNSET)

        id = d.pop("id", UNSET)

        name = d.pop("name", UNSET)

        pinned = []
        _pinned = d.pop("pinned", UNSET)
        for pinned_item_data in _pinned or []:
            pinned_item = ModelsHole.from_dict(pinned_item_data)

            pinned.append(pinned_item)

        time_created = d.pop("time_created", UNSET)

        time_updated = d.pop("time_updated", UNSET)

        models_division = cls(
            description=description,
            division_id=division_id,
            hidden=hidden,
            id=id,
            name=name,
            pinned=pinned,
            time_created=time_created,
            time_updated=time_updated,
        )

        models_division.additional_properties = d
        return models_division

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
