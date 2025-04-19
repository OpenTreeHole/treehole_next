from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.models_division import ModelsDivision
    from ..models.models_floor import ModelsFloor


T = TypeVar("T", bound="ModelsPunishment")


@_attrs_define
class ModelsPunishment:
    """
    Attributes:
        created_at (Union[Unset, str]): time when this punishment creates
        day (Union[Unset, int]):
        division (Union[Unset, ModelsDivision]):
        division_id (Union[Unset, int]):
        duration (Union[Unset, int]):
        end_time (Union[Unset, str]): end_time of this punishment
        floor (Union[Unset, ModelsFloor]):
        floor_id (Union[Unset, int]): punished because of this floor
        id (Union[Unset, int]):
        made_by (Union[Unset, int]): admin user_id who made this punish
        reason (Union[Unset, str]): reason
        start_time (Union[Unset, str]): start from end_time of previous punishment (punishment accumulation of different
            floors)
            if no previous punishment or previous punishment end time less than time.Now() (synced), set start time
            time.Now()
        updated_at (Union[Unset, str]):
        user_id (Union[Unset, int]): user punished
    """

    created_at: Union[Unset, str] = UNSET
    day: Union[Unset, int] = UNSET
    division: Union[Unset, "ModelsDivision"] = UNSET
    division_id: Union[Unset, int] = UNSET
    duration: Union[Unset, int] = UNSET
    end_time: Union[Unset, str] = UNSET
    floor: Union[Unset, "ModelsFloor"] = UNSET
    floor_id: Union[Unset, int] = UNSET
    id: Union[Unset, int] = UNSET
    made_by: Union[Unset, int] = UNSET
    reason: Union[Unset, str] = UNSET
    start_time: Union[Unset, str] = UNSET
    updated_at: Union[Unset, str] = UNSET
    user_id: Union[Unset, int] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        created_at = self.created_at

        day = self.day

        division: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.division, Unset):
            division = self.division.to_dict()

        division_id = self.division_id

        duration = self.duration

        end_time = self.end_time

        floor: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.floor, Unset):
            floor = self.floor.to_dict()

        floor_id = self.floor_id

        id = self.id

        made_by = self.made_by

        reason = self.reason

        start_time = self.start_time

        updated_at = self.updated_at

        user_id = self.user_id

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if created_at is not UNSET:
            field_dict["created_at"] = created_at
        if day is not UNSET:
            field_dict["day"] = day
        if division is not UNSET:
            field_dict["division"] = division
        if division_id is not UNSET:
            field_dict["division_id"] = division_id
        if duration is not UNSET:
            field_dict["duration"] = duration
        if end_time is not UNSET:
            field_dict["end_time"] = end_time
        if floor is not UNSET:
            field_dict["floor"] = floor
        if floor_id is not UNSET:
            field_dict["floor_id"] = floor_id
        if id is not UNSET:
            field_dict["id"] = id
        if made_by is not UNSET:
            field_dict["made_by"] = made_by
        if reason is not UNSET:
            field_dict["reason"] = reason
        if start_time is not UNSET:
            field_dict["start_time"] = start_time
        if updated_at is not UNSET:
            field_dict["updated_at"] = updated_at
        if user_id is not UNSET:
            field_dict["user_id"] = user_id

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.models_division import ModelsDivision
        from ..models.models_floor import ModelsFloor

        d = dict(src_dict)
        created_at = d.pop("created_at", UNSET)

        day = d.pop("day", UNSET)

        _division = d.pop("division", UNSET)
        division: Union[Unset, ModelsDivision]
        if isinstance(_division, Unset):
            division = UNSET
        else:
            division = ModelsDivision.from_dict(_division)

        division_id = d.pop("division_id", UNSET)

        duration = d.pop("duration", UNSET)

        end_time = d.pop("end_time", UNSET)

        _floor = d.pop("floor", UNSET)
        floor: Union[Unset, ModelsFloor]
        if isinstance(_floor, Unset):
            floor = UNSET
        else:
            floor = ModelsFloor.from_dict(_floor)

        floor_id = d.pop("floor_id", UNSET)

        id = d.pop("id", UNSET)

        made_by = d.pop("made_by", UNSET)

        reason = d.pop("reason", UNSET)

        start_time = d.pop("start_time", UNSET)

        updated_at = d.pop("updated_at", UNSET)

        user_id = d.pop("user_id", UNSET)

        models_punishment = cls(
            created_at=created_at,
            day=day,
            division=division,
            division_id=division_id,
            duration=duration,
            end_time=end_time,
            floor=floor,
            floor_id=floor_id,
            id=id,
            made_by=made_by,
            reason=reason,
            start_time=start_time,
            updated_at=updated_at,
            user_id=user_id,
        )

        models_punishment.additional_properties = d
        return models_punishment

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
