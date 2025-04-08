from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.models_user_permission_silent import ModelsUserPermissionSilent


T = TypeVar("T", bound="ModelsUserPermission")


@_attrs_define
class ModelsUserPermission:
    """
    Attributes:
        admin (Union[Unset, str]): 管理员权限到期时间
        offense_count (Union[Unset, int]):
        silent (Union[Unset, ModelsUserPermissionSilent]): key: division_id value: 对应分区禁言解除时间
    """

    admin: Union[Unset, str] = UNSET
    offense_count: Union[Unset, int] = UNSET
    silent: Union[Unset, "ModelsUserPermissionSilent"] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        admin = self.admin

        offense_count = self.offense_count

        silent: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.silent, Unset):
            silent = self.silent.to_dict()

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if admin is not UNSET:
            field_dict["admin"] = admin
        if offense_count is not UNSET:
            field_dict["offense_count"] = offense_count
        if silent is not UNSET:
            field_dict["silent"] = silent

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.models_user_permission_silent import ModelsUserPermissionSilent

        d = dict(src_dict)
        admin = d.pop("admin", UNSET)

        offense_count = d.pop("offense_count", UNSET)

        _silent = d.pop("silent", UNSET)
        silent: Union[Unset, ModelsUserPermissionSilent]
        if isinstance(_silent, Unset):
            silent = UNSET
        else:
            silent = ModelsUserPermissionSilent.from_dict(_silent)

        models_user_permission = cls(
            admin=admin,
            offense_count=offense_count,
            silent=silent,
        )

        models_user_permission.additional_properties = d
        return models_user_permission

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
