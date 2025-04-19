from collections.abc import Mapping
from typing import TYPE_CHECKING, Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.common_error_detail_element import CommonErrorDetailElement


T = TypeVar("T", bound="CommonHttpError")


@_attrs_define
class CommonHttpError:
    """
    Attributes:
        code (Union[Unset, int]):
        detail (Union[Unset, list['CommonErrorDetailElement']]):
        message (Union[Unset, str]):
    """

    code: Union[Unset, int] = UNSET
    detail: Union[Unset, list["CommonErrorDetailElement"]] = UNSET
    message: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        code = self.code

        detail: Union[Unset, list[dict[str, Any]]] = UNSET
        if not isinstance(self.detail, Unset):
            detail = []
            for detail_item_data in self.detail:
                detail_item = detail_item_data.to_dict()
                detail.append(detail_item)

        message = self.message

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if code is not UNSET:
            field_dict["code"] = code
        if detail is not UNSET:
            field_dict["detail"] = detail
        if message is not UNSET:
            field_dict["message"] = message

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        from ..models.common_error_detail_element import CommonErrorDetailElement

        d = dict(src_dict)
        code = d.pop("code", UNSET)

        detail = []
        _detail = d.pop("detail", UNSET)
        for detail_item_data in _detail or []:
            detail_item = CommonErrorDetailElement.from_dict(detail_item_data)

            detail.append(detail_item)

        message = d.pop("message", UNSET)

        common_http_error = cls(
            code=code,
            detail=detail,
            message=message,
        )

        common_http_error.additional_properties = d
        return common_http_error

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
