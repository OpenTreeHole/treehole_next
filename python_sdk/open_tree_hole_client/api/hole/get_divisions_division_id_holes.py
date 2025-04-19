from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.models_hole import ModelsHole
from ...models.models_message_model import ModelsMessageModel
from ...types import UNSET, Response, Unset


def _get_kwargs(
    division_id: int,
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
        "url": f"/divisions/{division_id}/holes",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[Union[ModelsMessageModel, list["ModelsHole"]]]:
    if response.status_code == 200:
        response_200 = []
        _response_200 = response.json()
        for response_200_item_data in _response_200:
            response_200_item = ModelsHole.from_dict(response_200_item_data)

            response_200.append(response_200_item)

        return response_200
    if response.status_code == 404:
        response_404 = ModelsMessageModel.from_dict(response.json())

        return response_404
    if response.status_code == 500:
        response_500 = ModelsMessageModel.from_dict(response.json())

        return response_500
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[ModelsMessageModel, list["ModelsHole"]]]:
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
    offset: Union[Unset, str] = UNSET,
    order: Union[Unset, str] = UNSET,
    size: Union[Unset, int] = 10,
) -> Response[Union[ModelsMessageModel, list["ModelsHole"]]]:
    """List Holes In A Division

    Args:
        division_id (int):
        offset (Union[Unset, str]):
        order (Union[Unset, str]):
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[ModelsMessageModel, list['ModelsHole']]]
    """

    kwargs = _get_kwargs(
        division_id=division_id,
        offset=offset,
        order=order,
        size=size,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    division_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    offset: Union[Unset, str] = UNSET,
    order: Union[Unset, str] = UNSET,
    size: Union[Unset, int] = 10,
) -> Optional[Union[ModelsMessageModel, list["ModelsHole"]]]:
    """List Holes In A Division

    Args:
        division_id (int):
        offset (Union[Unset, str]):
        order (Union[Unset, str]):
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[ModelsMessageModel, list['ModelsHole']]
    """

    return sync_detailed(
        division_id=division_id,
        client=client,
        offset=offset,
        order=order,
        size=size,
    ).parsed


async def asyncio_detailed(
    division_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    offset: Union[Unset, str] = UNSET,
    order: Union[Unset, str] = UNSET,
    size: Union[Unset, int] = 10,
) -> Response[Union[ModelsMessageModel, list["ModelsHole"]]]:
    """List Holes In A Division

    Args:
        division_id (int):
        offset (Union[Unset, str]):
        order (Union[Unset, str]):
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[ModelsMessageModel, list['ModelsHole']]]
    """

    kwargs = _get_kwargs(
        division_id=division_id,
        offset=offset,
        order=order,
        size=size,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    division_id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    offset: Union[Unset, str] = UNSET,
    order: Union[Unset, str] = UNSET,
    size: Union[Unset, int] = 10,
) -> Optional[Union[ModelsMessageModel, list["ModelsHole"]]]:
    """List Holes In A Division

    Args:
        division_id (int):
        offset (Union[Unset, str]):
        order (Union[Unset, str]):
        size (Union[Unset, int]):  Default: 10.

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[ModelsMessageModel, list['ModelsHole']]
    """

    return (
        await asyncio_detailed(
            division_id=division_id,
            client=client,
            offset=offset,
            order=order,
            size=size,
        )
    ).parsed
