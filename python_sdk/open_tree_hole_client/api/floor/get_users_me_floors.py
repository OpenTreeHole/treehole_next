from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.get_users_me_floors_order_by import GetUsersMeFloorsOrderBy
from ...models.get_users_me_floors_sort import GetUsersMeFloorsSort
from ...models.models_floor import ModelsFloor
from ...models.models_message_model import ModelsMessageModel
from ...types import UNSET, Response, Unset


def _get_kwargs(
    *,
    offset: Union[Unset, int] = 0,
    order_by: Union[Unset, GetUsersMeFloorsOrderBy] = GetUsersMeFloorsOrderBy.ID,
    size: Union[Unset, int] = 30,
    sort: Union[Unset, GetUsersMeFloorsSort] = GetUsersMeFloorsSort.ASC,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["offset"] = offset

    json_order_by: Union[Unset, str] = UNSET
    if not isinstance(order_by, Unset):
        json_order_by = order_by.value

    params["order_by"] = json_order_by

    params["size"] = size

    json_sort: Union[Unset, str] = UNSET
    if not isinstance(sort, Unset):
        json_sort = sort.value

    params["sort"] = json_sort

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/users/me/floors",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[Union[ModelsMessageModel, list["ModelsFloor"]]]:
    if response.status_code == 200:
        response_200 = []
        _response_200 = response.json()
        for response_200_item_data in _response_200:
            response_200_item = ModelsFloor.from_dict(response_200_item_data)

            response_200.append(response_200_item)

        return response_200
    if response.status_code == 404:
        response_404 = ModelsMessageModel.from_dict(response.json())

        return response_404
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[ModelsMessageModel, list["ModelsFloor"]]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    offset: Union[Unset, int] = 0,
    order_by: Union[Unset, GetUsersMeFloorsOrderBy] = GetUsersMeFloorsOrderBy.ID,
    size: Union[Unset, int] = 30,
    sort: Union[Unset, GetUsersMeFloorsSort] = GetUsersMeFloorsSort.ASC,
) -> Response[Union[ModelsMessageModel, list["ModelsFloor"]]]:
    """List User's Reply Floors

    Args:
        offset (Union[Unset, int]):  Default: 0.
        order_by (Union[Unset, GetUsersMeFloorsOrderBy]):  Default: GetUsersMeFloorsOrderBy.ID.
        size (Union[Unset, int]):  Default: 30.
        sort (Union[Unset, GetUsersMeFloorsSort]):  Default: GetUsersMeFloorsSort.ASC.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[ModelsMessageModel, list['ModelsFloor']]]
    """

    kwargs = _get_kwargs(
        offset=offset,
        order_by=order_by,
        size=size,
        sort=sort,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: Union[AuthenticatedClient, Client],
    offset: Union[Unset, int] = 0,
    order_by: Union[Unset, GetUsersMeFloorsOrderBy] = GetUsersMeFloorsOrderBy.ID,
    size: Union[Unset, int] = 30,
    sort: Union[Unset, GetUsersMeFloorsSort] = GetUsersMeFloorsSort.ASC,
) -> Optional[Union[ModelsMessageModel, list["ModelsFloor"]]]:
    """List User's Reply Floors

    Args:
        offset (Union[Unset, int]):  Default: 0.
        order_by (Union[Unset, GetUsersMeFloorsOrderBy]):  Default: GetUsersMeFloorsOrderBy.ID.
        size (Union[Unset, int]):  Default: 30.
        sort (Union[Unset, GetUsersMeFloorsSort]):  Default: GetUsersMeFloorsSort.ASC.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[ModelsMessageModel, list['ModelsFloor']]
    """

    return sync_detailed(
        client=client,
        offset=offset,
        order_by=order_by,
        size=size,
        sort=sort,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    offset: Union[Unset, int] = 0,
    order_by: Union[Unset, GetUsersMeFloorsOrderBy] = GetUsersMeFloorsOrderBy.ID,
    size: Union[Unset, int] = 30,
    sort: Union[Unset, GetUsersMeFloorsSort] = GetUsersMeFloorsSort.ASC,
) -> Response[Union[ModelsMessageModel, list["ModelsFloor"]]]:
    """List User's Reply Floors

    Args:
        offset (Union[Unset, int]):  Default: 0.
        order_by (Union[Unset, GetUsersMeFloorsOrderBy]):  Default: GetUsersMeFloorsOrderBy.ID.
        size (Union[Unset, int]):  Default: 30.
        sort (Union[Unset, GetUsersMeFloorsSort]):  Default: GetUsersMeFloorsSort.ASC.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[ModelsMessageModel, list['ModelsFloor']]]
    """

    kwargs = _get_kwargs(
        offset=offset,
        order_by=order_by,
        size=size,
        sort=sort,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: Union[AuthenticatedClient, Client],
    offset: Union[Unset, int] = 0,
    order_by: Union[Unset, GetUsersMeFloorsOrderBy] = GetUsersMeFloorsOrderBy.ID,
    size: Union[Unset, int] = 30,
    sort: Union[Unset, GetUsersMeFloorsSort] = GetUsersMeFloorsSort.ASC,
) -> Optional[Union[ModelsMessageModel, list["ModelsFloor"]]]:
    """List User's Reply Floors

    Args:
        offset (Union[Unset, int]):  Default: 0.
        order_by (Union[Unset, GetUsersMeFloorsOrderBy]):  Default: GetUsersMeFloorsOrderBy.ID.
        size (Union[Unset, int]):  Default: 30.
        sort (Union[Unset, GetUsersMeFloorsSort]):  Default: GetUsersMeFloorsSort.ASC.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[ModelsMessageModel, list['ModelsFloor']]
    """

    return (
        await asyncio_detailed(
            client=client,
            offset=offset,
            order_by=order_by,
            size=size,
            sort=sort,
        )
    ).parsed
