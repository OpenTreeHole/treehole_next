from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.models_floor import ModelsFloor


T = TypeVar("T", bound="ModelsHoleFloors")


@_attrs_define
class ModelsHoleFloors:
    """返回给前端的楼层列表，包括首楼、尾楼和预加载的前 n 个楼层

    Attributes:
        first_floor (Union[Unset, ModelsFloor]):
        last_floor (Union[Unset, ModelsFloor]):
        prefetch (Union[Unset, list['ModelsFloor']]): 预加载的楼层
    """

    first_floor: Union[Unset, "ModelsFloor"] = UNSET
    last_floor: Union[Unset, "ModelsFloor"] = UNSET
    prefetch: Union[Unset, list["ModelsFloor"]] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        first_floor: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.first_floor, Unset):
            first_floor = self.first_floor.to_dict()

        last_floor: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.last_floor, Unset):
            last_floor = self.last_floor.to_dict()

        prefetch: Union[Unset, list[dict[str, Any]]] = UNSET
        if not isinstance(self.prefetch, Unset):
            prefetch = []
            for prefetch_item_data in self.prefetch:
                prefetch_item = prefetch_item_data.to_dict()
                prefetch.append(prefetch_item)

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if first_floor is not UNSET:
            field_dict["first_floor"] = first_floor
        if last_floor is not UNSET:
            field_dict["last_floor"] = last_floor
        if prefetch is not UNSET:
            field_dict["prefetch"] = prefetch

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.models_floor import ModelsFloor

        d = dict(src_dict)
        _first_floor = d.pop("first_floor", UNSET)
        first_floor: Union[Unset, ModelsFloor]
        if isinstance(_first_floor, Unset):
            first_floor = UNSET
        else:
            first_floor = ModelsFloor.from_dict(_first_floor)

        _last_floor = d.pop("last_floor", UNSET)
        last_floor: Union[Unset, ModelsFloor]
        if isinstance(_last_floor, Unset):
            last_floor = UNSET
        else:
            last_floor = ModelsFloor.from_dict(_last_floor)

        prefetch = []
        _prefetch = d.pop("prefetch", UNSET)
        for prefetch_item_data in _prefetch or []:
            prefetch_item = ModelsFloor.from_dict(prefetch_item_data)

            prefetch.append(prefetch_item)

        models_hole_floors = cls(
            first_floor=first_floor,
            last_floor=last_floor,
            prefetch=prefetch,
        )

        models_hole_floors.additional_properties = d
        return models_hole_floors

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
