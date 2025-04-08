from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..models.models_message_type import ModelsMessageType
from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.models_message_data import ModelsMessageData


T = TypeVar("T", bound="ModelsMessage")


@_attrs_define
class ModelsMessage:
    """
    Attributes:
        code (Union[Unset, ModelsMessageType]):
        data (Union[Unset, ModelsMessageData]):
        description (Union[Unset, str]):
        has_read (Union[Unset, bool]): 兼容旧版, 永远为false，以MessageUser的HasRead为准
        id (Union[Unset, int]):
        message (Union[Unset, str]):
        message_id (Union[Unset, int]): 兼容旧版 id
        time_created (Union[Unset, str]):
        time_updated (Union[Unset, str]):
        url (Union[Unset, str]):
    """

    code: Union[Unset, ModelsMessageType] = UNSET
    data: Union[Unset, "ModelsMessageData"] = UNSET
    description: Union[Unset, str] = UNSET
    has_read: Union[Unset, bool] = UNSET
    id: Union[Unset, int] = UNSET
    message: Union[Unset, str] = UNSET
    message_id: Union[Unset, int] = UNSET
    time_created: Union[Unset, str] = UNSET
    time_updated: Union[Unset, str] = UNSET
    url: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        code: Union[Unset, str] = UNSET
        if not isinstance(self.code, Unset):
            code = self.code.value

        data: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.data, Unset):
            data = self.data.to_dict()

        description = self.description

        has_read = self.has_read

        id = self.id

        message = self.message

        message_id = self.message_id

        time_created = self.time_created

        time_updated = self.time_updated

        url = self.url

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if code is not UNSET:
            field_dict["code"] = code
        if data is not UNSET:
            field_dict["data"] = data
        if description is not UNSET:
            field_dict["description"] = description
        if has_read is not UNSET:
            field_dict["has_read"] = has_read
        if id is not UNSET:
            field_dict["id"] = id
        if message is not UNSET:
            field_dict["message"] = message
        if message_id is not UNSET:
            field_dict["message_id"] = message_id
        if time_created is not UNSET:
            field_dict["time_created"] = time_created
        if time_updated is not UNSET:
            field_dict["time_updated"] = time_updated
        if url is not UNSET:
            field_dict["url"] = url

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.models_message_data import ModelsMessageData

        d = dict(src_dict)
        _code = d.pop("code", UNSET)
        code: Union[Unset, ModelsMessageType]
        if isinstance(_code, Unset):
            code = UNSET
        else:
            code = ModelsMessageType(_code)

        _data = d.pop("data", UNSET)
        data: Union[Unset, ModelsMessageData]
        if isinstance(_data, Unset):
            data = UNSET
        else:
            data = ModelsMessageData.from_dict(_data)

        description = d.pop("description", UNSET)

        has_read = d.pop("has_read", UNSET)

        id = d.pop("id", UNSET)

        message = d.pop("message", UNSET)

        message_id = d.pop("message_id", UNSET)

        time_created = d.pop("time_created", UNSET)

        time_updated = d.pop("time_updated", UNSET)

        url = d.pop("url", UNSET)

        models_message = cls(
            code=code,
            data=data,
            description=description,
            has_read=has_read,
            id=id,
            message=message,
            message_id=message_id,
            time_created=time_created,
            time_updated=time_updated,
            url=url,
        )

        models_message.additional_properties = d
        return models_message

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
