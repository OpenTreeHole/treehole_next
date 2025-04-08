from collections.abc import Mapping
from typing import Any, TypeVar, Union, cast

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..models.floor_modify_model_like import FloorModifyModelLike
from ..types import UNSET, Unset

T = TypeVar("T", bound="FloorModifyModel")


@_attrs_define
class FloorModifyModel:
    """
    Attributes:
        content (Union[Unset, str]): Owner or admin, the original content should be moved to  floor_history
        fold (Union[Unset, list[str]]): 仅管理员，留空则重置，低优先级
        fold_v2 (Union[Unset, str]): 仅管理员，留空则重置，高优先级
        like (Union[Unset, FloorModifyModelLike]): All user, deprecated, "add" is like, "cancel" is reset
        special_tag (Union[Unset, str]): Admin and Operator only
    """

    content: Union[Unset, str] = UNSET
    fold: Union[Unset, list[str]] = UNSET
    fold_v2: Union[Unset, str] = UNSET
    like: Union[Unset, FloorModifyModelLike] = UNSET
    special_tag: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        content = self.content

        fold: Union[Unset, list[str]] = UNSET
        if not isinstance(self.fold, Unset):
            fold = self.fold

        fold_v2 = self.fold_v2

        like: Union[Unset, str] = UNSET
        if not isinstance(self.like, Unset):
            like = self.like.value

        special_tag = self.special_tag

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if content is not UNSET:
            field_dict["content"] = content
        if fold is not UNSET:
            field_dict["fold"] = fold
        if fold_v2 is not UNSET:
            field_dict["fold_v2"] = fold_v2
        if like is not UNSET:
            field_dict["like"] = like
        if special_tag is not UNSET:
            field_dict["special_tag"] = special_tag

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        content = d.pop("content", UNSET)

        fold = cast(list[str], d.pop("fold", UNSET))

        fold_v2 = d.pop("fold_v2", UNSET)

        _like = d.pop("like", UNSET)
        like: Union[Unset, FloorModifyModelLike]
        if isinstance(_like, Unset):
            like = UNSET
        else:
            like = FloorModifyModelLike(_like)

        special_tag = d.pop("special_tag", UNSET)

        floor_modify_model = cls(
            content=content,
            fold=fold,
            fold_v2=fold_v2,
            like=like,
            special_tag=special_tag,
        )

        floor_modify_model.additional_properties = d
        return floor_modify_model

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
