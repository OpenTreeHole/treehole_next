from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.get_user_favorite_groups_order import GetUserFavoriteGroupsOrder
from ...models.models_favorite_group import ModelsFavoriteGroup
from ...types import UNSET, Response, Unset


def _get_kwargs(
    *,
    order: Union[Unset, GetUserFavoriteGroupsOrder] = GetUserFavoriteGroupsOrder.TIME_CREATED,
    plain: Union[Unset, bool] = False,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    json_order: Union[Unset, str] = UNSET
    if not isinstance(order, Unset):
        json_order = order.value

    params["order"] = json_order

    params["plain"] = plain

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/user/favorite_groups",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[list["ModelsFavoriteGroup"]]:
    if response.status_code == 200:
        response_200 = []
        _response_200 = response.json()
        for response_200_item_data in _response_200:
            response_200_item = ModelsFavoriteGroup.from_dict(response_200_item_data)

            response_200.append(response_200_item)

        return response_200
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[list["ModelsFavoriteGroup"]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    order: Union[Unset, GetUserFavoriteGroupsOrder] = GetUserFavoriteGroupsOrder.TIME_CREATED,
    plain: Union[Unset, bool] = False,
) -> Response[list["ModelsFavoriteGroup"]]:
    """List User's Favorite Groups

    Args:
        order (Union[Unset, GetUserFavoriteGroupsOrder]):  Default:
            GetUserFavoriteGroupsOrder.TIME_CREATED.
        plain (Union[Unset, bool]):  Default: False.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['ModelsFavoriteGroup']]
    """

    kwargs = _get_kwargs(
        order=order,
        plain=plain,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: Union[AuthenticatedClient, Client],
    order: Union[Unset, GetUserFavoriteGroupsOrder] = GetUserFavoriteGroupsOrder.TIME_CREATED,
    plain: Union[Unset, bool] = False,
) -> Optional[list["ModelsFavoriteGroup"]]:
    """List User's Favorite Groups

    Args:
        order (Union[Unset, GetUserFavoriteGroupsOrder]):  Default:
            GetUserFavoriteGroupsOrder.TIME_CREATED.
        plain (Union[Unset, bool]):  Default: False.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['ModelsFavoriteGroup']
    """

    return sync_detailed(
        client=client,
        order=order,
        plain=plain,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    order: Union[Unset, GetUserFavoriteGroupsOrder] = GetUserFavoriteGroupsOrder.TIME_CREATED,
    plain: Union[Unset, bool] = False,
) -> Response[list["ModelsFavoriteGroup"]]:
    """List User's Favorite Groups

    Args:
        order (Union[Unset, GetUserFavoriteGroupsOrder]):  Default:
            GetUserFavoriteGroupsOrder.TIME_CREATED.
        plain (Union[Unset, bool]):  Default: False.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['ModelsFavoriteGroup']]
    """

    kwargs = _get_kwargs(
        order=order,
        plain=plain,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: Union[AuthenticatedClient, Client],
    order: Union[Unset, GetUserFavoriteGroupsOrder] = GetUserFavoriteGroupsOrder.TIME_CREATED,
    plain: Union[Unset, bool] = False,
) -> Optional[list["ModelsFavoriteGroup"]]:
    """List User's Favorite Groups

    Args:
        order (Union[Unset, GetUserFavoriteGroupsOrder]):  Default:
            GetUserFavoriteGroupsOrder.TIME_CREATED.
        plain (Union[Unset, bool]):  Default: False.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['ModelsFavoriteGroup']
    """

    return (
        await asyncio_detailed(
            client=client,
            order=order,
            plain=plain,
        )
    ).parsed
