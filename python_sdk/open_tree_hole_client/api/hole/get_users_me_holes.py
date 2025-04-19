from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.models_hole import ModelsHole
from ...types import UNSET, Response, Unset


def _get_kwargs(
    *,
    offset: Union[Unset, str] = UNSET,
    order: Union[Unset, str] = UNSET,
    size: Union[Unset, int] = 10,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["offset"] = offset

    params["order"] = order

    params["size"] = size

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/users/me/holes",
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
    offset: Union[Unset, str] = UNSET,
    order: Union[Unset, str] = UNSET,
    size: Union[Unset, int] = 10,
) -> Response[list["ModelsHole"]]:
    """List a Hole Created By User

    Args:
        offset (Union[Unset, str]):
        order (Union[Unset, str]):
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['ModelsHole']]
    """

    kwargs = _get_kwargs(
        offset=offset,
        order=order,
        size=size,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: Union[AuthenticatedClient, Client],
    offset: Union[Unset, str] = UNSET,
    order: Union[Unset, str] = UNSET,
    size: Union[Unset, int] = 10,
) -> Optional[list["ModelsHole"]]:
    """List a Hole Created By User

    Args:
        offset (Union[Unset, str]):
        order (Union[Unset, str]):
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['ModelsHole']
    """

    return sync_detailed(
        client=client,
        offset=offset,
        order=order,
        size=size,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    offset: Union[Unset, str] = UNSET,
    order: Union[Unset, str] = UNSET,
    size: Union[Unset, int] = 10,
) -> Response[list["ModelsHole"]]:
    """List a Hole Created By User

    Args:
        offset (Union[Unset, str]):
        order (Union[Unset, str]):
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['ModelsHole']]
    """

    kwargs = _get_kwargs(
        offset=offset,
        order=order,
        size=size,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: Union[AuthenticatedClient, Client],
    offset: Union[Unset, str] = UNSET,
    order: Union[Unset, str] = UNSET,
    size: Union[Unset, int] = 10,
) -> Optional[list["ModelsHole"]]:
    """List a Hole Created By User

    Args:
        offset (Union[Unset, str]):
        order (Union[Unset, str]):
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['ModelsHole']
    """

    return (
        await asyncio_detailed(
            client=client,
            offset=offset,
            order=order,
            size=size,
        )
    ).parsed
