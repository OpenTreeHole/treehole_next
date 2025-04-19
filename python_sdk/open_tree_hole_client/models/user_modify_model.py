from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.user_user_config_model import UserUserConfigModel


T = TypeVar("T", bound="UserModifyModel")


@_attrs_define
class UserModifyModel:
    """
    Attributes:
        config (Union[Unset, UserUserConfigModel]):
        nickname (Union[Unset, str]):
    """

    config: Union[Unset, "UserUserConfigModel"] = UNSET
    nickname: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        config: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.config, Unset):
            config = self.config.to_dict()

        nickname = self.nickname

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if config is not UNSET:
            field_dict["config"] = config
        if nickname is not UNSET:
            field_dict["nickname"] = nickname

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.user_user_config_model import UserUserConfigModel

        d = dict(src_dict)
        _config = d.pop("config", UNSET)
        config: Union[Unset, UserUserConfigModel]
        if isinstance(_config, Unset):
            config = UNSET
        else:
            config = UserUserConfigModel.from_dict(_config)

        nickname = d.pop("nickname", UNSET)

        user_modify_model = cls(
            config=config,
            nickname=nickname,
        )

        user_modify_model.additional_properties = d
        return user_modify_model

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
