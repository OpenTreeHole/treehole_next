from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.models_floor import ModelsFloor
from ...types import UNSET, Response, Unset


def _get_kwargs(
    *,
    accurate: Union[Unset, bool] = False,
    end_time: Union[Unset, int] = UNSET,
    offset: Union[Unset, int] = 0,
    search: str,
    size: Union[Unset, int] = 10,
    start_time: Union[Unset, int] = UNSET,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["accurate"] = accurate

    params["end_time"] = end_time

    params["offset"] = offset

    params["search"] = search

    params["size"] = size

    params["start_time"] = start_time

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": "/floors/search",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[list["ModelsFloor"]]:
    if response.status_code == 200:
        response_200 = []
        _response_200 = response.json()
        for response_200_item_data in _response_200:
            response_200_item = ModelsFloor.from_dict(response_200_item_data)

            response_200.append(response_200_item)

        return response_200
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[list["ModelsFloor"]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    accurate: Union[Unset, bool] = False,
    end_time: Union[Unset, int] = UNSET,
    offset: Union[Unset, int] = 0,
    search: str,
    size: Union[Unset, int] = 10,
    start_time: Union[Unset, int] = UNSET,
) -> Response[list["ModelsFloor"]]:
    """SearchFloors In ElasticSearch

    Args:
        accurate (Union[Unset, bool]):  Default: False.
        end_time (Union[Unset, int]):
        offset (Union[Unset, int]):  Default: 0.
        search (str):
        size (Union[Unset, int]):  Default: 10.
        start_time (Union[Unset, int]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['ModelsFloor']]
    """

    kwargs = _get_kwargs(
        accurate=accurate,
        end_time=end_time,
        offset=offset,
        search=search,
        size=size,
        start_time=start_time,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: Union[AuthenticatedClient, Client],
    accurate: Union[Unset, bool] = False,
    end_time: Union[Unset, int] = UNSET,
    offset: Union[Unset, int] = 0,
    search: str,
    size: Union[Unset, int] = 10,
    start_time: Union[Unset, int] = UNSET,
) -> Optional[list["ModelsFloor"]]:
    """SearchFloors In ElasticSearch

    Args:
        accurate (Union[Unset, bool]):  Default: False.
        end_time (Union[Unset, int]):
        offset (Union[Unset, int]):  Default: 0.
        search (str):
        size (Union[Unset, int]):  Default: 10.
        start_time (Union[Unset, int]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['ModelsFloor']
    """

    return sync_detailed(
        client=client,
        accurate=accurate,
        end_time=end_time,
        offset=offset,
        search=search,
        size=size,
        start_time=start_time,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    accurate: Union[Unset, bool] = False,
    end_time: Union[Unset, int] = UNSET,
    offset: Union[Unset, int] = 0,
    search: str,
    size: Union[Unset, int] = 10,
    start_time: Union[Unset, int] = UNSET,
) -> Response[list["ModelsFloor"]]:
    """SearchFloors In ElasticSearch

    Args:
        accurate (Union[Unset, bool]):  Default: False.
        end_time (Union[Unset, int]):
        offset (Union[Unset, int]):  Default: 0.
        search (str):
        size (Union[Unset, int]):  Default: 10.
        start_time (Union[Unset, int]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['ModelsFloor']]
    """

    kwargs = _get_kwargs(
        accurate=accurate,
        end_time=end_time,
        offset=offset,
        search=search,
        size=size,
        start_time=start_time,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: Union[AuthenticatedClient, Client],
    accurate: Union[Unset, bool] = False,
    end_time: Union[Unset, int] = UNSET,
    offset: Union[Unset, int] = 0,
    search: str,
    size: Union[Unset, int] = 10,
    start_time: Union[Unset, int] = UNSET,
) -> Optional[list["ModelsFloor"]]:
    """SearchFloors In ElasticSearch

    Args:
        accurate (Union[Unset, bool]):  Default: False.
        end_time (Union[Unset, int]):
        offset (Union[Unset, int]):  Default: 0.
        search (str):
        size (Union[Unset, int]):  Default: 10.
        start_time (Union[Unset, int]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['ModelsFloor']
    """

    return (
        await asyncio_detailed(
            client=client,
            accurate=accurate,
            end_time=end_time,
            offset=offset,
            search=search,
            size=size,
            start_time=start_time,
        )
    ).parsed
