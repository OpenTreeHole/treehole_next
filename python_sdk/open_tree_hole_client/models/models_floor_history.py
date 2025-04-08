from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="ModelsFloorHistory")


@_attrs_define
class ModelsFloorHistory:
    """
    Attributes:
        content (Union[Unset, str]):
        floor_id (Union[Unset, int]):
        id (Union[Unset, int]): / base info
        is_actual_sensitive (Union[Unset, bool]): manual sensitive check
        is_sensitive (Union[Unset, bool]): auto sensitive check
        reason (Union[Unset, str]):
        sensitive_detail (Union[Unset, str]):
        time_created (Union[Unset, str]):
        time_updated (Union[Unset, str]):
        user_id (Union[Unset, int]): The one who modified the floor
    """

    content: Union[Unset, str] = UNSET
    floor_id: Union[Unset, int] = UNSET
    id: Union[Unset, int] = UNSET
    is_actual_sensitive: Union[Unset, bool] = UNSET
    is_sensitive: Union[Unset, bool] = UNSET
    reason: Union[Unset, str] = UNSET
    sensitive_detail: Union[Unset, str] = UNSET
    time_created: Union[Unset, str] = UNSET
    time_updated: Union[Unset, str] = UNSET
    user_id: Union[Unset, int] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        content = self.content

        floor_id = self.floor_id

        id = self.id

        is_actual_sensitive = self.is_actual_sensitive

        is_sensitive = self.is_sensitive

        reason = self.reason

        sensitive_detail = self.sensitive_detail

        time_created = self.time_created

        time_updated = self.time_updated

        user_id = self.user_id

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if content is not UNSET:
            field_dict["content"] = content
        if floor_id is not UNSET:
            field_dict["floor_id"] = floor_id
        if id is not UNSET:
            field_dict["id"] = id
        if is_actual_sensitive is not UNSET:
            field_dict["is_actual_sensitive"] = is_actual_sensitive
        if is_sensitive is not UNSET:
            field_dict["is_sensitive"] = is_sensitive
        if reason is not UNSET:
            field_dict["reason"] = reason
        if sensitive_detail is not UNSET:
            field_dict["sensitive_detail"] = sensitive_detail
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
        content = d.pop("content", UNSET)

        floor_id = d.pop("floor_id", UNSET)

        id = d.pop("id", UNSET)

        is_actual_sensitive = d.pop("is_actual_sensitive", UNSET)

        is_sensitive = d.pop("is_sensitive", UNSET)

        reason = d.pop("reason", UNSET)

        sensitive_detail = d.pop("sensitive_detail", UNSET)

        time_created = d.pop("time_created", UNSET)

        time_updated = d.pop("time_updated", UNSET)

        user_id = d.pop("user_id", UNSET)

        models_floor_history = cls(
            content=content,
            floor_id=floor_id,
            id=id,
            is_actual_sensitive=is_actual_sensitive,
            is_sensitive=is_sensitive,
            reason=reason,
            sensitive_detail=sensitive_detail,
            time_created=time_created,
            time_updated=time_updated,
            user_id=user_id,
        )

        models_floor_history.additional_properties = d
        return models_floor_history

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
