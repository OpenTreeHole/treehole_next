from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.floor_create_model import FloorCreateModel
from ...models.models_floor import ModelsFloor
from ...types import Response


def _get_kwargs(
    hole_id: int,
    *,
    body: FloorCreateModel,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": f"/holes/{hole_id}/floors",
    }

    _body = body.to_dict()

    _kwargs["json"] = _body
    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(*, client: Union[AuthenticatedClient, Client], response: httpx.Response) -> Optional[ModelsFloor]:
    if response.status_code == 201:
        response_201 = ModelsFloor.from_dict(response.json())

        return response_201
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(*, client: Union[AuthenticatedClient, Client], response: httpx.Response) -> Response[ModelsFloor]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    hole_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: FloorCreateModel,
) -> Response[ModelsFloor]:
    """Create A Floor

    Args:
        hole_id (int):
        body (FloorCreateModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[ModelsFloor]
    """

    kwargs = _get_kwargs(
        hole_id=hole_id,
        body=body,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    hole_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: FloorCreateModel,
) -> Optional[ModelsFloor]:
    """Create A Floor

    Args:
        hole_id (int):
        body (FloorCreateModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        ModelsFloor
    """

    return sync_detailed(
        hole_id=hole_id,
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    hole_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: FloorCreateModel,
) -> Response[ModelsFloor]:
    """Create A Floor

    Args:
        hole_id (int):
        body (FloorCreateModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[ModelsFloor]
    """

    kwargs = _get_kwargs(
        hole_id=hole_id,
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    hole_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: FloorCreateModel,
) -> Optional[ModelsFloor]:
    """Create A Floor

    Args:
        hole_id (int):
        body (FloorCreateModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        ModelsFloor
    """

    return (
        await asyncio_detailed(
            hole_id=hole_id,
            client=client,
            body=body,
        )
    ).parsed
