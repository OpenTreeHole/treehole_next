from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.models_floor import ModelsFloor
from ...types import UNSET, Response, Unset


def _get_kwargs(
    *,
    hole_id: Union[Unset, int] = UNSET,
    length: Union[Unset, int] = UNSET,
    s: Union[Unset, str] = UNSET,
    start_floor: Union[Unset, int] = UNSET,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["hole_id"] = hole_id

    params["length"] = length

    params["s"] = s

    params["start_floor"] = start_floor

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/floors",
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
    hole_id: Union[Unset, int] = UNSET,
    length: Union[Unset, int] = UNSET,
    s: Union[Unset, str] = UNSET,
    start_floor: Union[Unset, int] = UNSET,
) -> Response[list["ModelsFloor"]]:
    """Old API for Listing Floors

    Args:
        hole_id (Union[Unset, int]):
        length (Union[Unset, int]):
        s (Union[Unset, str]):
        start_floor (Union[Unset, int]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['ModelsFloor']]
    """

    kwargs = _get_kwargs(
        hole_id=hole_id,
        length=length,
        s=s,
        start_floor=start_floor,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: Union[AuthenticatedClient, Client],
    hole_id: Union[Unset, int] = UNSET,
    length: Union[Unset, int] = UNSET,
    s: Union[Unset, str] = UNSET,
    start_floor: Union[Unset, int] = UNSET,
) -> Optional[list["ModelsFloor"]]:
    """Old API for Listing Floors

    Args:
        hole_id (Union[Unset, int]):
        length (Union[Unset, int]):
        s (Union[Unset, str]):
        start_floor (Union[Unset, int]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['ModelsFloor']
    """

    return sync_detailed(
        client=client,
        hole_id=hole_id,
        length=length,
        s=s,
        start_floor=start_floor,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    hole_id: Union[Unset, int] = UNSET,
    length: Union[Unset, int] = UNSET,
    s: Union[Unset, str] = UNSET,
    start_floor: Union[Unset, int] = UNSET,
) -> Response[list["ModelsFloor"]]:
    """Old API for Listing Floors

    Args:
        hole_id (Union[Unset, int]):
        length (Union[Unset, int]):
        s (Union[Unset, str]):
        start_floor (Union[Unset, int]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['ModelsFloor']]
    """

    kwargs = _get_kwargs(
        hole_id=hole_id,
        length=length,
        s=s,
        start_floor=start_floor,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: Union[AuthenticatedClient, Client],
    hole_id: Union[Unset, int] = UNSET,
    length: Union[Unset, int] = UNSET,
    s: Union[Unset, str] = UNSET,
    start_floor: Union[Unset, int] = UNSET,
) -> Optional[list["ModelsFloor"]]:
    """Old API for Listing Floors

    Args:
        hole_id (Union[Unset, int]):
        length (Union[Unset, int]):
        s (Union[Unset, str]):
        start_floor (Union[Unset, int]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['ModelsFloor']
    """

    return (
        await asyncio_detailed(
            client=client,
            hole_id=hole_id,
            length=length,
            s=s,
            start_floor=start_floor,
        )
    ).parsed
