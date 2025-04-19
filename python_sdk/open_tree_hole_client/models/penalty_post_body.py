from collections.abc import Mapping
from typing import Any, TypeVar, Union, cast

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="PenaltyPostBody")


@_attrs_define
class PenaltyPostBody:
    """
    Attributes:
        days (Union[Unset, int]): high priority
        divisions (Union[Unset, list[int]]): high priority
        penalty_level (Union[Unset, int]): low priority, deprecated
        reason (Union[Unset, str]): optional
    """

    days: Union[Unset, int] = UNSET
    divisions: Union[Unset, list[int]] = UNSET
    penalty_level: Union[Unset, int] = UNSET
    reason: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        days = self.days

        divisions: Union[Unset, list[int]] = UNSET
        if not isinstance(self.divisions, Unset):
            divisions = self.divisions

        penalty_level = self.penalty_level

        reason = self.reason

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if days is not UNSET:
            field_dict["days"] = days
        if divisions is not UNSET:
            field_dict["divisions"] = divisions
        if penalty_level is not UNSET:
            field_dict["penalty_level"] = penalty_level
        if reason is not UNSET:
            field_dict["reason"] = reason

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        days = d.pop("days", UNSET)

        divisions = cast(list[int], d.pop("divisions", UNSET))

        penalty_level = d.pop("penalty_level", UNSET)

        reason = d.pop("reason", UNSET)

        penalty_post_body = cls(
            days=days,
            divisions=divisions,
            penalty_level=penalty_level,
            reason=reason,
        )

        penalty_post_body.additional_properties = d
        return penalty_post_body

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
