from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.common_http_error import CommonHttpError
from ...models.floor_ban_division import FloorBanDivision
from ...types import Response


def _get_kwargs(
    id: int,
) -> dict[str, Any]:
    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": f"/floors/{id}/user_silence",
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[Union[CommonHttpError, FloorBanDivision]]:
    if response.status_code == 200:
        response_200 = FloorBanDivision.from_dict(response.json())

        return response_200
    if response.status_code == 403:
        response_403 = CommonHttpError.from_dict(response.json())

        return response_403
    if response.status_code == 404:
        response_404 = CommonHttpError.from_dict(response.json())

        return response_404
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[CommonHttpError, FloorBanDivision]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    id: int,
    *,
    client: Union[AuthenticatedClient, Client],
) -> Response[Union[CommonHttpError, FloorBanDivision]]:
    """Get A User's Silence Status According To FloorID, admin only

    Args:
        id (int):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[CommonHttpError, FloorBanDivision]]
    """

    kwargs = _get_kwargs(
        id=id,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    id: int,
    *,
    client: Union[AuthenticatedClient, Client],
) -> Optional[Union[CommonHttpError, FloorBanDivision]]:
    """Get A User's Silence Status According To FloorID, admin only

    Args:
        id (int):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[CommonHttpError, FloorBanDivision]
    """

    return sync_detailed(
        id=id,
        client=client,
    ).parsed


async def asyncio_detailed(
    id: int,
    *,
    client: Union[AuthenticatedClient, Client],
) -> Response[Union[CommonHttpError, FloorBanDivision]]:
    """Get A User's Silence Status According To FloorID, admin only

    Args:
        id (int):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[CommonHttpError, FloorBanDivision]]
    """

    kwargs = _get_kwargs(
        id=id,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    id: int,
    *,
    client: Union[AuthenticatedClient, Client],
) -> Optional[Union[CommonHttpError, FloorBanDivision]]:
    """Get A User's Silence Status According To FloorID, admin only

    Args:
        id (int):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[CommonHttpError, FloorBanDivision]
    """

    return (
        await asyncio_detailed(
            id=id,
            client=client,
        )
    ).parsed
