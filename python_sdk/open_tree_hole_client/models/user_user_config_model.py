from collections.abc import Mapping
from typing import Any, TypeVar, Union, cast

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="UserUserConfigModel")


@_attrs_define
class UserUserConfigModel:
    """
    Attributes:
        notify (Union[Unset, list[str]]):
        show_folded (Union[Unset, str]):
    """

    notify: Union[Unset, list[str]] = UNSET
    show_folded: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        notify: Union[Unset, list[str]] = UNSET
        if not isinstance(self.notify, Unset):
            notify = self.notify

        show_folded = self.show_folded

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if notify is not UNSET:
            field_dict["notify"] = notify
        if show_folded is not UNSET:
            field_dict["show_folded"] = show_folded

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        notify = cast(list[str], d.pop("notify", UNSET))

        show_folded = d.pop("show_folded", UNSET)

        user_user_config_model = cls(
            notify=notify,
            show_folded=show_folded,
        )

        user_user_config_model.additional_properties = d
        return user_user_config_model

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
