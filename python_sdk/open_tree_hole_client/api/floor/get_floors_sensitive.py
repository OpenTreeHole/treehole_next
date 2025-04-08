from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.floor_sensitive_floor_response import FloorSensitiveFloorResponse
from ...models.get_floors_sensitive_order_by import GetFloorsSensitiveOrderBy
from ...models.models_message_model import ModelsMessageModel
from ...types import UNSET, Response, Unset


def _get_kwargs(
    *,
    all_: Union[Unset, bool] = UNSET,
    offset: Union[Unset, str] = UNSET,
    open_: Union[Unset, bool] = UNSET,
    order_by: Union[Unset, GetFloorsSensitiveOrderBy] = GetFloorsSensitiveOrderBy.TIME_CREATED,
    size: Union[Unset, int] = 10,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["all"] = all_

    params["offset"] = offset

    params["open"] = open_

    json_order_by: Union[Unset, str] = UNSET
    if not isinstance(order_by, Unset):
        json_order_by = order_by.value

    params["order_by"] = json_order_by

    params["size"] = size

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/floors/_sensitive",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[Union[ModelsMessageModel, list["FloorSensitiveFloorResponse"]]]:
    if response.status_code == 200:
        response_200 = []
        _response_200 = response.json()
        for response_200_item_data in _response_200:
            response_200_item = FloorSensitiveFloorResponse.from_dict(response_200_item_data)

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
) -> Response[Union[ModelsMessageModel, list["FloorSensitiveFloorResponse"]]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    all_: Union[Unset, bool] = UNSET,
    offset: Union[Unset, str] = UNSET,
    open_: Union[Unset, bool] = UNSET,
    order_by: Union[Unset, GetFloorsSensitiveOrderBy] = GetFloorsSensitiveOrderBy.TIME_CREATED,
    size: Union[Unset, int] = 10,
) -> Response[Union[ModelsMessageModel, list["FloorSensitiveFloorResponse"]]]:
    """List sensitive floors, admin only

    Args:
        all_ (Union[Unset, bool]):
        offset (Union[Unset, str]):
        open_ (Union[Unset, bool]):
        order_by (Union[Unset, GetFloorsSensitiveOrderBy]):  Default:
            GetFloorsSensitiveOrderBy.TIME_CREATED.
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[ModelsMessageModel, list['FloorSensitiveFloorResponse']]]
    """

    kwargs = _get_kwargs(
        all_=all_,
        offset=offset,
        open_=open_,
        order_by=order_by,
        size=size,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: Union[AuthenticatedClient, Client],
    all_: Union[Unset, bool] = UNSET,
    offset: Union[Unset, str] = UNSET,
    open_: Union[Unset, bool] = UNSET,
    order_by: Union[Unset, GetFloorsSensitiveOrderBy] = GetFloorsSensitiveOrderBy.TIME_CREATED,
    size: Union[Unset, int] = 10,
) -> Optional[Union[ModelsMessageModel, list["FloorSensitiveFloorResponse"]]]:
    """List sensitive floors, admin only

    Args:
        all_ (Union[Unset, bool]):
        offset (Union[Unset, str]):
        open_ (Union[Unset, bool]):
        order_by (Union[Unset, GetFloorsSensitiveOrderBy]):  Default:
            GetFloorsSensitiveOrderBy.TIME_CREATED.
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[ModelsMessageModel, list['FloorSensitiveFloorResponse']]
    """

    return sync_detailed(
        client=client,
        all_=all_,
        offset=offset,
        open_=open_,
        order_by=order_by,
        size=size,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    all_: Union[Unset, bool] = UNSET,
    offset: Union[Unset, str] = UNSET,
    open_: Union[Unset, bool] = UNSET,
    order_by: Union[Unset, GetFloorsSensitiveOrderBy] = GetFloorsSensitiveOrderBy.TIME_CREATED,
    size: Union[Unset, int] = 10,
) -> Response[Union[ModelsMessageModel, list["FloorSensitiveFloorResponse"]]]:
    """List sensitive floors, admin only

    Args:
        all_ (Union[Unset, bool]):
        offset (Union[Unset, str]):
        open_ (Union[Unset, bool]):
        order_by (Union[Unset, GetFloorsSensitiveOrderBy]):  Default:
            GetFloorsSensitiveOrderBy.TIME_CREATED.
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[ModelsMessageModel, list['FloorSensitiveFloorResponse']]]
    """

    kwargs = _get_kwargs(
        all_=all_,
        offset=offset,
        open_=open_,
        order_by=order_by,
        size=size,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: Union[AuthenticatedClient, Client],
    all_: Union[Unset, bool] = UNSET,
    offset: Union[Unset, str] = UNSET,
    open_: Union[Unset, bool] = UNSET,
    order_by: Union[Unset, GetFloorsSensitiveOrderBy] = GetFloorsSensitiveOrderBy.TIME_CREATED,
    size: Union[Unset, int] = 10,
) -> Optional[Union[ModelsMessageModel, list["FloorSensitiveFloorResponse"]]]:
    """List sensitive floors, admin only

    Args:
        all_ (Union[Unset, bool]):
        offset (Union[Unset, str]):
        open_ (Union[Unset, bool]):
        order_by (Union[Unset, GetFloorsSensitiveOrderBy]):  Default:
            GetFloorsSensitiveOrderBy.TIME_CREATED.
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[ModelsMessageModel, list['FloorSensitiveFloorResponse']]
    """

    return (
        await asyncio_detailed(
            client=client,
            all_=all_,
            offset=offset,
            open_=open_,
            order_by=order_by,
            size=size,
        )
    ).parsed
