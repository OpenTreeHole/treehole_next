from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="FloorSensitiveFloorResponse")


@_attrs_define
class FloorSensitiveFloorResponse:
    """
    Attributes:
        content (Union[Unset, str]):
        deleted (Union[Unset, bool]):
        hole_id (Union[Unset, int]):
        id (Union[Unset, int]):
        is_actual_sensitive (Union[Unset, bool]):
        modified (Union[Unset, int]):
        sensitive_detail (Union[Unset, str]):
        time_created (Union[Unset, str]):
        time_updated (Union[Unset, str]):
    """

    content: Union[Unset, str] = UNSET
    deleted: Union[Unset, bool] = UNSET
    hole_id: Union[Unset, int] = UNSET
    id: Union[Unset, int] = UNSET
    is_actual_sensitive: Union[Unset, bool] = UNSET
    modified: Union[Unset, int] = UNSET
    sensitive_detail: Union[Unset, str] = UNSET
    time_created: Union[Unset, str] = UNSET
    time_updated: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        content = self.content

        deleted = self.deleted

        hole_id = self.hole_id

        id = self.id

        is_actual_sensitive = self.is_actual_sensitive

        modified = self.modified

        sensitive_detail = self.sensitive_detail

        time_created = self.time_created

        time_updated = self.time_updated

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if content is not UNSET:
            field_dict["content"] = content
        if deleted is not UNSET:
            field_dict["deleted"] = deleted
        if hole_id is not UNSET:
            field_dict["hole_id"] = hole_id
        if id is not UNSET:
            field_dict["id"] = id
        if is_actual_sensitive is not UNSET:
            field_dict["is_actual_sensitive"] = is_actual_sensitive
        if modified is not UNSET:
            field_dict["modified"] = modified
        if sensitive_detail is not UNSET:
            field_dict["sensitive_detail"] = sensitive_detail
        if time_created is not UNSET:
            field_dict["time_created"] = time_created
        if time_updated is not UNSET:
            field_dict["time_updated"] = time_updated

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        content = d.pop("content", UNSET)

        deleted = d.pop("deleted", UNSET)

        hole_id = d.pop("hole_id", UNSET)

        id = d.pop("id", UNSET)

        is_actual_sensitive = d.pop("is_actual_sensitive", UNSET)

        modified = d.pop("modified", UNSET)

        sensitive_detail = d.pop("sensitive_detail", UNSET)

        time_created = d.pop("time_created", UNSET)

        time_updated = d.pop("time_updated", UNSET)

        floor_sensitive_floor_response = cls(
            content=content,
            deleted=deleted,
            hole_id=hole_id,
            id=id,
            is_actual_sensitive=is_actual_sensitive,
            modified=modified,
            sensitive_detail=sensitive_detail,
            time_created=time_created,
            time_updated=time_updated,
        )

        floor_sensitive_floor_response.additional_properties = d
        return floor_sensitive_floor_response

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
