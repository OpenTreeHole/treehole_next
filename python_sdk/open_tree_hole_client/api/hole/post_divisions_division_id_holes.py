from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.hole_create_model import HoleCreateModel
from ...models.models_hole import ModelsHole
from ...types import Response


def _get_kwargs(
    division_id: int,
    *,
    body: HoleCreateModel,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "post",
        "url": f"/divisions/{division_id}/holes",
    }

    _body = body.to_dict()

    _kwargs["json"] = _body
    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(*, client: Union[AuthenticatedClient, Client], response: httpx.Response) -> Optional[ModelsHole]:
    if response.status_code == 201:
        response_201 = ModelsHole.from_dict(response.json())

        return response_201
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(*, client: Union[AuthenticatedClient, Client], response: httpx.Response) -> Response[ModelsHole]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    division_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: HoleCreateModel,
) -> Response[ModelsHole]:
    """Create A Hole

     Create a hole, create tags and floor binding to it and set the name mapping

    Args:
        division_id (int):
        body (HoleCreateModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[ModelsHole]
    """

    kwargs = _get_kwargs(
        division_id=division_id,
        body=body,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    division_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: HoleCreateModel,
) -> Optional[ModelsHole]:
    """Create A Hole

     Create a hole, create tags and floor binding to it and set the name mapping

    Args:
        division_id (int):
        body (HoleCreateModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        ModelsHole
    """

    return sync_detailed(
        division_id=division_id,
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    division_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: HoleCreateModel,
) -> Response[ModelsHole]:
    """Create A Hole

     Create a hole, create tags and floor binding to it and set the name mapping

    Args:
        division_id (int):
        body (HoleCreateModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[ModelsHole]
    """

    kwargs = _get_kwargs(
        division_id=division_id,
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    division_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: HoleCreateModel,
) -> Optional[ModelsHole]:
    """Create A Hole

     Create a hole, create tags and floor binding to it and set the name mapping

    Args:
        division_id (int):
        body (HoleCreateModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        ModelsHole
    """

    return (
        await asyncio_detailed(
            division_id=division_id,
            client=client,
            body=body,
        )
    ).parsed
