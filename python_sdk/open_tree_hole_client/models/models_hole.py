from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.gorm_deleted_at import GormDeletedAt
    from ..models.models_hole_floors import ModelsHoleFloors
    from ..models.models_tag import ModelsTag


T = TypeVar("T", bound="ModelsHole")


@_attrs_define
class ModelsHole:
    """
    Attributes:
        division_id (Union[Unset, int]): 所属 division 的 id
        floors (Union[Unset, ModelsHoleFloors]): 返回给前端的楼层列表，包括首楼、尾楼和预加载的前 n 个楼层
        good (Union[Unset, bool]):
        hidden (Union[Unset, bool]): 是否隐藏，隐藏的洞用户不可见，管理员可见
        hole_id (Union[Unset, int]): 兼容旧版 id
        id (Union[Unset, int]): / saved fields
        locked (Union[Unset, bool]): 锁定帖子，如果锁定则非管理员无法发帖，也无法修改已有发帖
        no_purge (Union[Unset, bool]):
        reply (Union[Unset, int]): 回复量（即该洞下 floor 的数量 - 1）
        tags (Union[Unset, list['ModelsTag']]): tag 列表；不超过 10 个
        time_created (Union[Unset, str]):
        time_deleted (Union[Unset, GormDeletedAt]):
        time_updated (Union[Unset, str]):
        view (Union[Unset, int]): 浏览量
    """

    division_id: Union[Unset, int] = UNSET
    floors: Union[Unset, "ModelsHoleFloors"] = UNSET
    good: Union[Unset, bool] = UNSET
    hidden: Union[Unset, bool] = UNSET
    hole_id: Union[Unset, int] = UNSET
    id: Union[Unset, int] = UNSET
    locked: Union[Unset, bool] = UNSET
    no_purge: Union[Unset, bool] = UNSET
    reply: Union[Unset, int] = UNSET
    tags: Union[Unset, list["ModelsTag"]] = UNSET
    time_created: Union[Unset, str] = UNSET
    time_deleted: Union[Unset, "GormDeletedAt"] = UNSET
    time_updated: Union[Unset, str] = UNSET
    view: Union[Unset, int] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        division_id = self.division_id

        floors: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.floors, Unset):
            floors = self.floors.to_dict()

        good = self.good

        hidden = self.hidden

        hole_id = self.hole_id

        id = self.id

        locked = self.locked

        no_purge = self.no_purge

        reply = self.reply

        tags: Union[Unset, list[dict[str, Any]]] = UNSET
        if not isinstance(self.tags, Unset):
            tags = []
            for tags_item_data in self.tags:
                tags_item = tags_item_data.to_dict()
                tags.append(tags_item)

        time_created = self.time_created

        time_deleted: Union[Unset, dict[str, Any]] = UNSET
        if not isinstance(self.time_deleted, Unset):
            time_deleted = self.time_deleted.to_dict()

        time_updated = self.time_updated

        view = self.view

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if division_id is not UNSET:
            field_dict["division_id"] = division_id
        if floors is not UNSET:
            field_dict["floors"] = floors
        if good is not UNSET:
            field_dict["good"] = good
        if hidden is not UNSET:
            field_dict["hidden"] = hidden
        if hole_id is not UNSET:
            field_dict["hole_id"] = hole_id
        if id is not UNSET:
            field_dict["id"] = id
        if locked is not UNSET:
            field_dict["locked"] = locked
        if no_purge is not UNSET:
            field_dict["no_purge"] = no_purge
        if reply is not UNSET:
            field_dict["reply"] = reply
        if tags is not UNSET:
            field_dict["tags"] = tags
        if time_created is not UNSET:
            field_dict["time_created"] = time_created
        if time_deleted is not UNSET:
            field_dict["time_deleted"] = time_deleted
        if time_updated is not UNSET:
            field_dict["time_updated"] = time_updated
        if view is not UNSET:
            field_dict["view"] = view

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.gorm_deleted_at import GormDeletedAt
        from ..models.models_hole_floors import ModelsHoleFloors
        from ..models.models_tag import ModelsTag

        d = dict(src_dict)
        division_id = d.pop("division_id", UNSET)

        _floors = d.pop("floors", UNSET)
        floors: Union[Unset, ModelsHoleFloors]
        if isinstance(_floors, Unset):
            floors = UNSET
        else:
            floors = ModelsHoleFloors.from_dict(_floors)

        good = d.pop("good", UNSET)

        hidden = d.pop("hidden", UNSET)

        hole_id = d.pop("hole_id", UNSET)

        id = d.pop("id", UNSET)

        locked = d.pop("locked", UNSET)

        no_purge = d.pop("no_purge", UNSET)

        reply = d.pop("reply", UNSET)

        tags = []
        _tags = d.pop("tags", UNSET)
        for tags_item_data in _tags or []:
            tags_item = ModelsTag.from_dict(tags_item_data)

            tags.append(tags_item)

        time_created = d.pop("time_created", UNSET)

        _time_deleted = d.pop("time_deleted", UNSET)
        time_deleted: Union[Unset, GormDeletedAt]
        if isinstance(_time_deleted, Unset):
            time_deleted = UNSET
        else:
            time_deleted = GormDeletedAt.from_dict(_time_deleted)

        time_updated = d.pop("time_updated", UNSET)

        view = d.pop("view", UNSET)

        models_hole = cls(
            division_id=division_id,
            floors=floors,
            good=good,
            hidden=hidden,
            hole_id=hole_id,
            id=id,
            locked=locked,
            no_purge=no_purge,
            reply=reply,
            tags=tags,
            time_created=time_created,
            time_deleted=time_deleted,
            time_updated=time_updated,
            view=view,
        )

        models_hole.additional_properties = d
        return models_hole

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
