from collections.abc import Mapping
from typing import Any, TypeVar, Union, cast

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="ModelsUserConfig")


@_attrs_define
class ModelsUserConfig:
    """
    Attributes:
        notify (Union[Unset, list[str]]): used when notify
        show_folded (Union[Unset, str]): 对折叠内容的处理
            fold 折叠, hide 隐藏, show 展示
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

        models_user_config = cls(
            notify=notify,
            show_folded=show_folded,
        )

        models_user_config.additional_properties = d
        return models_user_config

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
