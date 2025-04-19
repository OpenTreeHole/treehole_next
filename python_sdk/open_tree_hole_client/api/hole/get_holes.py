from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.models_hole import ModelsHole
from ...types import UNSET, Response, Unset


def _get_kwargs(
    *,
    division_id: Union[Unset, int] = UNSET,
    length: Union[Unset, int] = 10,
    order: Union[Unset, str] = UNSET,
    start_time: Union[Unset, str] = UNSET,
    tag: Union[Unset, str] = UNSET,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["division_id"] = division_id

    params["length"] = length

    params["order"] = order

    params["start_time"] = start_time

    params["tag"] = tag

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/holes",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[list["ModelsHole"]]:
    if response.status_code == 200:
        response_200 = []
        _response_200 = response.json()
        for response_200_item_data in _response_200:
            response_200_item = ModelsHole.from_dict(response_200_item_data)

            response_200.append(response_200_item)

        return response_200
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[list["ModelsHole"]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    division_id: Union[Unset, int] = UNSET,
    length: Union[Unset, int] = 10,
    order: Union[Unset, str] = UNSET,
    start_time: Union[Unset, str] = UNSET,
    tag: Union[Unset, str] = UNSET,
) -> Response[list["ModelsHole"]]:
    """Old API for Listing Holes

    Args:
        division_id (Union[Unset, int]):
        length (Union[Unset, int]):  Default: 10.
        order (Union[Unset, str]):
        start_time (Union[Unset, str]):
        tag (Union[Unset, str]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['ModelsHole']]
    """

    kwargs = _get_kwargs(
        division_id=division_id,
        length=length,
        order=order,
        start_time=start_time,
        tag=tag,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: Union[AuthenticatedClient, Client],
    division_id: Union[Unset, int] = UNSET,
    length: Union[Unset, int] = 10,
    order: Union[Unset, str] = UNSET,
    start_time: Union[Unset, str] = UNSET,
    tag: Union[Unset, str] = UNSET,
) -> Optional[list["ModelsHole"]]:
    """Old API for Listing Holes

    Args:
        division_id (Union[Unset, int]):
        length (Union[Unset, int]):  Default: 10.
        order (Union[Unset, str]):
        start_time (Union[Unset, str]):
        tag (Union[Unset, str]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['ModelsHole']
    """

    return sync_detailed(
        client=client,
        division_id=division_id,
        length=length,
        order=order,
        start_time=start_time,
        tag=tag,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    division_id: Union[Unset, int] = UNSET,
    length: Union[Unset, int] = 10,
    order: Union[Unset, str] = UNSET,
    start_time: Union[Unset, str] = UNSET,
    tag: Union[Unset, str] = UNSET,
) -> Response[list["ModelsHole"]]:
    """Old API for Listing Holes

    Args:
        division_id (Union[Unset, int]):
        length (Union[Unset, int]):  Default: 10.
        order (Union[Unset, str]):
        start_time (Union[Unset, str]):
        tag (Union[Unset, str]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['ModelsHole']]
    """

    kwargs = _get_kwargs(
        division_id=division_id,
        length=length,
        order=order,
        start_time=start_time,
        tag=tag,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: Union[AuthenticatedClient, Client],
    division_id: Union[Unset, int] = UNSET,
    length: Union[Unset, int] = 10,
    order: Union[Unset, str] = UNSET,
    start_time: Union[Unset, str] = UNSET,
    tag: Union[Unset, str] = UNSET,
) -> Optional[list["ModelsHole"]]:
    """Old API for Listing Holes

    Args:
        division_id (Union[Unset, int]):
        length (Union[Unset, int]):  Default: 10.
        order (Union[Unset, str]):
        start_time (Union[Unset, str]):
        tag (Union[Unset, str]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['ModelsHole']
    """

    return (
        await asyncio_detailed(
            client=client,
            division_id=division_id,
            length=length,
            order=order,
            start_time=start_time,
            tag=tag,
        )
    ).parsed
