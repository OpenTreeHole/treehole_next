from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.models_floor import ModelsFloor


T = TypeVar("T", bound="ModelsReport")


@_attrs_define
class ModelsReport:
    """
    Attributes:
        dealt (Union[Unset, bool]): the report has been dealt
        dealt_by (Union[Unset, int]): who dealt the report
        floor (Union[Unset, ModelsFloor]):
        floor_id (Union[Unset, int]):
        hole_id (Union[Unset, int]):
        id (Union[Unset, int]):
        reason (Union[Unset, str]):
        report_id (Union[Unset, int]):
        result (Union[Unset, str]): deal result
        time_created (Union[Unset, str]):
        time_updated (Union[Unset, str]):
    """

    dealt: Union[Unset, bool] = UNSET
    dealt_by: Union[Unset, int] = UNSET
    floor: Union[Unset, "ModelsFloor"] = UNSET
    floor_id: Union[Unset, int] = UNSET
    hole_id: Union[Unset, int] = UNSET
    id: Union[Unset, int] = UNSET
    reason: Union[Unset, str] = UNSET
    report_id: Union[Unset, int] = UNSET
    result: Union[Unset, str] = UNSET
    time_created: Union[Unset, str] = UNSET
    time_updated: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        dealt = self.dealt

        dealt_by = self.dealt_by

        floor: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.floor, Unset):
            floor = self.floor.to_dict()

        floor_id = self.floor_id

        hole_id = self.hole_id

        id = self.id

        reason = self.reason

        report_id = self.report_id

        result = self.result

        time_created = self.time_created

        time_updated = self.time_updated

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if dealt is not UNSET:
            field_dict["dealt"] = dealt
        if dealt_by is not UNSET:
            field_dict["dealt_by"] = dealt_by
        if floor is not UNSET:
            field_dict["floor"] = floor
        if floor_id is not UNSET:
            field_dict["floor_id"] = floor_id
        if hole_id is not UNSET:
            field_dict["hole_id"] = hole_id
        if id is not UNSET:
            field_dict["id"] = id
        if reason is not UNSET:
            field_dict["reason"] = reason
        if report_id is not UNSET:
            field_dict["report_id"] = report_id
        if result is not UNSET:
            field_dict["result"] = result
        if time_created is not UNSET:
            field_dict["time_created"] = time_created
        if time_updated is not UNSET:
            field_dict["time_updated"] = time_updated

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.models_floor import ModelsFloor

        d = dict(src_dict)
        dealt = d.pop("dealt", UNSET)

        dealt_by = d.pop("dealt_by", UNSET)

        _floor = d.pop("floor", UNSET)
        floor: Union[Unset, ModelsFloor]
        if isinstance(_floor, Unset):
            floor = UNSET
        else:
            floor = ModelsFloor.from_dict(_floor)

        floor_id = d.pop("floor_id", UNSET)

        hole_id = d.pop("hole_id", UNSET)

        id = d.pop("id", UNSET)

        reason = d.pop("reason", UNSET)

        report_id = d.pop("report_id", UNSET)

        result = d.pop("result", UNSET)

        time_created = d.pop("time_created", UNSET)

        time_updated = d.pop("time_updated", UNSET)

        models_report = cls(
            dealt=dealt,
            dealt_by=dealt_by,
            floor=floor,
            floor_id=floor_id,
            hole_id=hole_id,
            id=id,
            reason=reason,
            report_id=report_id,
            result=result,
            time_created=time_created,
            time_updated=time_updated,
        )

        models_report.additional_properties = d
        return models_report

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
