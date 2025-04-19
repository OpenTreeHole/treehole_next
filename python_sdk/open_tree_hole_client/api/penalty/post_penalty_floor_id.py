from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.models_user import ModelsUser
from ...models.penalty_post_body import PenaltyPostBody
from ...types import Response


def _get_kwargs(
    floor_id: int,
    *,
    body: PenaltyPostBody,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": f"/penalty/{floor_id}",
    }

    _body = body.to_dict()

    _kwargs["json"] = _body
    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(*, client: Union[AuthenticatedClient, Client], response: httpx.Response) -> Optional[ModelsUser]:
    if response.status_code == 201:
        response_201 = ModelsUser.from_dict(response.json())

        return response_201
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(*, client: Union[AuthenticatedClient, Client], response: httpx.Response) -> Response[ModelsUser]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    floor_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: PenaltyPostBody,
) -> Response[ModelsUser]:
    """Ban publisher of a floor

    Args:
        floor_id (int):
        body (PenaltyPostBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[ModelsUser]
    """

    kwargs = _get_kwargs(
        floor_id=floor_id,
        body=body,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    floor_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: PenaltyPostBody,
) -> Optional[ModelsUser]:
    """Ban publisher of a floor

    Args:
        floor_id (int):
        body (PenaltyPostBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        ModelsUser
    """

    return sync_detailed(
        floor_id=floor_id,
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    floor_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: PenaltyPostBody,
) -> Response[ModelsUser]:
    """Ban publisher of a floor

    Args:
        floor_id (int):
        body (PenaltyPostBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[ModelsUser]
    """

    kwargs = _get_kwargs(
        floor_id=floor_id,
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    floor_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: PenaltyPostBody,
) -> Optional[ModelsUser]:
    """Ban publisher of a floor

    Args:
        floor_id (int):
        body (PenaltyPostBody):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        ModelsUser
    """

    return (
        await asyncio_detailed(
            floor_id=floor_id,
            client=client,
            body=body,
        )
    ).parsed
