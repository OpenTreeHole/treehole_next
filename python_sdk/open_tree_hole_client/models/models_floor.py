from collections.abc import Mapping
from typing import Any, TypeVar, Union, cast

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="ModelsFloor")


@_attrs_define
class ModelsFloor:
    """
    Attributes:
        anonyname (Union[Unset, str]): a random username
        content (Union[Unset, str]): content of the floor, no more than 15000, should be sensitive checked, no more than
            10000 in frontend
        deleted (Union[Unset, bool]): whether the floor is deleted
        dislike (Union[Unset, int]): dislike number
        disliked (Union[Unset, bool]): whether the user has disliked the floor
        floor_id (Union[Unset, int]): old version compatibility
        fold (Union[Unset, list[str]]): fold reason, for v1
        fold_v2 (Union[Unset, str]): fold reason
        hole_id (Union[Unset, int]): the hole it belongs to
        id (Union[Unset, int]): / saved fields
        is_actual_sensitive (Union[Unset, bool]): manual sensitive check
        is_me (Union[Unset, bool]): whether the user is the author of the floor
        is_sensitive (Union[Unset, bool]): auto sensitive check
        like (Union[Unset, int]): like number
        liked (Union[Unset, bool]): whether the user has liked the floor
        mention (Union[Unset, list['ModelsFloor']]): many to many mentions
        modified (Union[Unset, int]): the modification times of floor.content
        ranking (Union[Unset, int]): the ranking of this floor in the hole
        reply_to (Union[Unset, int]): floor_id that it replies to, for dialog mode, in the same hole
        sensitive_detail (Union[Unset, str]): auto sensitive check detail
        special_tag (Union[Unset, str]): additional info, like "树洞管理团队"
        time_created (Union[Unset, str]):
        time_updated (Union[Unset, str]):
    """

    anonyname: Union[Unset, str] = UNSET
    content: Union[Unset, str] = UNSET
    deleted: Union[Unset, bool] = UNSET
    dislike: Union[Unset, int] = UNSET
    disliked: Union[Unset, bool] = UNSET
    floor_id: Union[Unset, int] = UNSET
    fold: Union[Unset, list[str]] = UNSET
    fold_v2: Union[Unset, str] = UNSET
    hole_id: Union[Unset, int] = UNSET
    id: Union[Unset, int] = UNSET
    is_actual_sensitive: Union[Unset, bool] = UNSET
    is_me: Union[Unset, bool] = UNSET
    is_sensitive: Union[Unset, bool] = UNSET
    like: Union[Unset, int] = UNSET
    liked: Union[Unset, bool] = UNSET
    mention: Union[Unset, list["ModelsFloor"]] = UNSET
    modified: Union[Unset, int] = UNSET
    ranking: Union[Unset, int] = UNSET
    reply_to: Union[Unset, int] = UNSET
    sensitive_detail: Union[Unset, str] = UNSET
    special_tag: Union[Unset, str] = UNSET
    time_created: Union[Unset, str] = UNSET
    time_updated: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        anonyname = self.anonyname

        content = self.content

        deleted = self.deleted

        dislike = self.dislike

        disliked = self.disliked

        floor_id = self.floor_id

        fold: Union[Unset, list[str]] = UNSET
        if not isinstance(self.fold, Unset):
            fold = self.fold

        fold_v2 = self.fold_v2

        hole_id = self.hole_id

        id = self.id

        is_actual_sensitive = self.is_actual_sensitive

        is_me = self.is_me

        is_sensitive = self.is_sensitive

        like = self.like

        liked = self.liked

        mention: Union[Unset, list[dict[str, Any]]] = UNSET
        if not isinstance(self.mention, Unset):
            mention = []
            for mention_item_data in self.mention:
                mention_item = mention_item_data.to_dict()
                mention.append(mention_item)

        modified = self.modified

        ranking = self.ranking

        reply_to = self.reply_to

        sensitive_detail = self.sensitive_detail

        special_tag = self.special_tag

        time_created = self.time_created

        time_updated = self.time_updated

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if anonyname is not UNSET:
            field_dict["anonyname"] = anonyname
        if content is not UNSET:
            field_dict["content"] = content
        if deleted is not UNSET:
            field_dict["deleted"] = deleted
        if dislike is not UNSET:
            field_dict["dislike"] = dislike
        if disliked is not UNSET:
            field_dict["disliked"] = disliked
        if floor_id is not UNSET:
            field_dict["floor_id"] = floor_id
        if fold is not UNSET:
            field_dict["fold"] = fold
        if fold_v2 is not UNSET:
            field_dict["fold_v2"] = fold_v2
        if hole_id is not UNSET:
            field_dict["hole_id"] = hole_id
        if id is not UNSET:
            field_dict["id"] = id
        if is_actual_sensitive is not UNSET:
            field_dict["is_actual_sensitive"] = is_actual_sensitive
        if is_me is not UNSET:
            field_dict["is_me"] = is_me
        if is_sensitive is not UNSET:
            field_dict["is_sensitive"] = is_sensitive
        if like is not UNSET:
            field_dict["like"] = like
        if liked is not UNSET:
            field_dict["liked"] = liked
        if mention is not UNSET:
            field_dict["mention"] = mention
        if modified is not UNSET:
            field_dict["modified"] = modified
        if ranking is not UNSET:
            field_dict["ranking"] = ranking
        if reply_to is not UNSET:
            field_dict["reply_to"] = reply_to
        if sensitive_detail is not UNSET:
            field_dict["sensitive_detail"] = sensitive_detail
        if special_tag is not UNSET:
            field_dict["special_tag"] = special_tag
        if time_created is not UNSET:
            field_dict["time_created"] = time_created
        if time_updated is not UNSET:
            field_dict["time_updated"] = time_updated

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        anonyname = d.pop("anonyname", UNSET)

        content = d.pop("content", UNSET)

        deleted = d.pop("deleted", UNSET)

        dislike = d.pop("dislike", UNSET)

        disliked = d.pop("disliked", UNSET)

        floor_id = d.pop("floor_id", UNSET)

        fold = cast(list[str], d.pop("fold", UNSET))

        fold_v2 = d.pop("fold_v2", UNSET)

        hole_id = d.pop("hole_id", UNSET)

        id = d.pop("id", UNSET)

        is_actual_sensitive = d.pop("is_actual_sensitive", UNSET)

        is_me = d.pop("is_me", UNSET)

        is_sensitive = d.pop("is_sensitive", UNSET)

        like = d.pop("like", UNSET)

        liked = d.pop("liked", UNSET)

        mention = []
        _mention = d.pop("mention", UNSET)
        for mention_item_data in _mention or []:
            mention_item = ModelsFloor.from_dict(mention_item_data)

            mention.append(mention_item)

        modified = d.pop("modified", UNSET)

        ranking = d.pop("ranking", UNSET)

        reply_to = d.pop("reply_to", UNSET)

        sensitive_detail = d.pop("sensitive_detail", UNSET)

        special_tag = d.pop("special_tag", UNSET)

        time_created = d.pop("time_created", UNSET)

        time_updated = d.pop("time_updated", UNSET)

        models_floor = cls(
            anonyname=anonyname,
            content=content,
            deleted=deleted,
            dislike=dislike,
            disliked=disliked,
            floor_id=floor_id,
            fold=fold,
            fold_v2=fold_v2,
            hole_id=hole_id,
            id=id,
            is_actual_sensitive=is_actual_sensitive,
            is_me=is_me,
            is_sensitive=is_sensitive,
            like=like,
            liked=liked,
            mention=mention,
            modified=modified,
            ranking=ranking,
            reply_to=reply_to,
            sensitive_detail=sensitive_detail,
            special_tag=special_tag,
            time_created=time_created,
            time_updated=time_updated,
        )

        models_floor.additional_properties = d
        return models_floor

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
