from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union, cast

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.models_user_config import ModelsUserConfig
    from ..models.models_user_permission import ModelsUserPermission


T = TypeVar("T", bound="ModelsUser")


@_attrs_define
class ModelsUser:
    """
    Attributes:
        config (Union[Unset, ModelsUserConfig]):
        default_special_tag (Union[Unset, str]):
        favorite_group_count (Union[Unset, int]):
        has_answered_questions (Union[Unset, bool]):
        id (Union[Unset, int]): / base info
        is_admin (Union[Unset, bool]): get from jwt
        joined_time (Union[Unset, str]):
        nickname (Union[Unset, str]):
        permission (Union[Unset, ModelsUserPermission]):
        special_tags (Union[Unset, list[str]]):
        user_id (Union[Unset, int]):
    """

    config: Union[Unset, "ModelsUserConfig"] = UNSET
    default_special_tag: Union[Unset, str] = UNSET
    favorite_group_count: Union[Unset, int] = UNSET
    has_answered_questions: Union[Unset, bool] = UNSET
    id: Union[Unset, int] = UNSET
    is_admin: Union[Unset, bool] = UNSET
    joined_time: Union[Unset, str] = UNSET
    nickname: Union[Unset, str] = UNSET
    permission: Union[Unset, "ModelsUserPermission"] = UNSET
    special_tags: Union[Unset, list[str]] = UNSET
    user_id: Union[Unset, int] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        config: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.config, Unset):
            config = self.config.to_dict()

        default_special_tag = self.default_special_tag

        favorite_group_count = self.favorite_group_count

        has_answered_questions = self.has_answered_questions

        id = self.id

        is_admin = self.is_admin

        joined_time = self.joined_time

        nickname = self.nickname

        permission: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.permission, Unset):
            permission = self.permission.to_dict()

        special_tags: Union[Unset, list[str]] = UNSET
        if not isinstance(self.special_tags, Unset):
            special_tags = self.special_tags

        user_id = self.user_id

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if config is not UNSET:
            field_dict["config"] = config
        if default_special_tag is not UNSET:
            field_dict["default_special_tag"] = default_special_tag
        if favorite_group_count is not UNSET:
            field_dict["favorite_group_count"] = favorite_group_count
        if has_answered_questions is not UNSET:
            field_dict["has_answered_questions"] = has_answered_questions
        if id is not UNSET:
            field_dict["id"] = id
        if is_admin is not UNSET:
            field_dict["is_admin"] = is_admin
        if joined_time is not UNSET:
            field_dict["joined_time"] = joined_time
        if nickname is not UNSET:
            field_dict["nickname"] = nickname
        if permission is not UNSET:
            field_dict["permission"] = permission
        if special_tags is not UNSET:
            field_dict["special_tags"] = special_tags
        if user_id is not UNSET:
            field_dict["user_id"] = user_id

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.models_user_config import ModelsUserConfig
        from ..models.models_user_permission import ModelsUserPermission

        d = dict(src_dict)
        _config = d.pop("config", UNSET)
        config: Union[Unset, ModelsUserConfig]
        if isinstance(_config, Unset):
            config = UNSET
        else:
            config = ModelsUserConfig.from_dict(_config)

        default_special_tag = d.pop("default_special_tag", UNSET)

        favorite_group_count = d.pop("favorite_group_count", UNSET)

        has_answered_questions = d.pop("has_answered_questions", UNSET)

        id = d.pop("id", UNSET)

        is_admin = d.pop("is_admin", UNSET)

        joined_time = d.pop("joined_time", UNSET)

        nickname = d.pop("nickname", UNSET)

        _permission = d.pop("permission", UNSET)
        permission: Union[Unset, ModelsUserPermission]
        if isinstance(_permission, Unset):
            permission = UNSET
        else:
            permission = ModelsUserPermission.from_dict(_permission)

        special_tags = cast(list[str], d.pop("special_tags", UNSET))

        user_id = d.pop("user_id", UNSET)

        models_user = cls(
            config=config,
            default_special_tag=default_special_tag,
            favorite_group_count=favorite_group_count,
            has_answered_questions=has_answered_questions,
            id=id,
            is_admin=is_admin,
            joined_time=joined_time,
            nickname=nickname,
            permission=permission,
            special_tags=special_tags,
            user_id=user_id,
        )

        models_user.additional_properties = d
        return models_user

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
